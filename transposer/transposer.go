package transposer

import (
	"fmt"
	"github.com/user/chord-transposer-go/chord"
	"github.com/user/chord-transposer-go/keysignature"
	"regexp"
	"strings"
)

// Token can be either a chord or regular text
type Token interface {
	String() string
}

// StringToken represents regular text
type StringToken string

func (s StringToken) String() string {
	return string(s)
}

// Ensure Chord implements the Token interface
var _ Token = chord.Chord{}
var _ Token = StringToken("")

// Constants
const NKeys = 12

// Transposer provides an API for transposing text containing chords
type Transposer struct {
	Tokens     [][]Token
	CurrentKey *keysignature.KeySignature
}

// NewTransposer creates a new Transposer instance
func NewTransposer(text interface{}) (*Transposer, error) {
	switch v := text.(type) {
	case string:
		tokens, err := tokenize(v)
		if err != nil {
			return nil, err
		}
		return &Transposer{Tokens: tokens}, nil
	case [][]Token:
		return &Transposer{Tokens: v}, nil
	default:
		return nil, fmt.Errorf("invalid argument (must be text or parsed text)")
	}
}

// GetKey returns the key of the text
// If not explicitly set, it will be guessed from the first chord
func (t *Transposer) GetKey() (*keysignature.KeySignature, error) {
	if t.CurrentKey != nil {
		return t.CurrentKey, nil
	}

	// Find the first chord to determine the key
	for _, line := range t.Tokens {
		for _, token := range line {
			if c, ok := token.(chord.Chord); ok {
				return keysignature.GuessKeySignature(c)
			}
		}
	}
	return nil, fmt.Errorf("given text has no chords")
}

// FromKey sets the source key
func (t *Transposer) FromKey(key interface{}) (*Transposer, error) {
	switch v := key.(type) {
	case string:
		k, err := keysignature.ValueOf(v)
		if err != nil {
			return nil, err
		}
		t.CurrentKey = k
	case *keysignature.KeySignature:
		t.CurrentKey = v
	default:
		return nil, fmt.Errorf("invalid key type")
	}
	return t, nil
}

// Up transposes up by the specified number of semitones
func (t *Transposer) Up(semitones int) (*Transposer, error) {
	key, err := t.GetKey()
	if err != nil {
		return nil, err
	}

	newKey, err := transposeKey(key, semitones)
	if err != nil {
		return nil, err
	}

	tokens, err := transposeTokens(t.Tokens, key, newKey)
	if err != nil {
		return nil, err
	}

	newTransposer, err := NewTransposer(tokens)
	if err != nil {
		return nil, err
	}

	return newTransposer.FromKey(newKey)
}

// Down transposes down by the specified number of semitones
func (t *Transposer) Down(semitones int) (*Transposer, error) {
	return t.Up(-semitones)
}

// ToKey transposes to the specified key
func (t *Transposer) ToKey(toKey string) (*Transposer, error) {
	key, err := t.GetKey()
	if err != nil {
		return nil, err
	}

	newKey, err := keysignature.ValueOf(toKey)
	if err != nil {
		return nil, err
	}

	tokens, err := transposeTokens(t.Tokens, key, newKey)
	if err != nil {
		return nil, err
	}

	newTransposer, err := NewTransposer(tokens)
	if err != nil {
		return nil, err
	}

	return newTransposer.FromKey(newKey)
}

// String returns a string representation of the text
func (t *Transposer) String() string {
	lines := make([]string, len(t.Tokens))
	for i, line := range t.Tokens {
		var sb strings.Builder
		for _, token := range line {
			sb.WriteString(token.String())
		}
		lines[i] = sb.String()
	}
	return strings.Join(lines, "\n")
}

// transposeKey finds the key that is a specified number of semitones above/below the current key
func transposeKey(currentKey *keysignature.KeySignature, semitones int) (*keysignature.KeySignature, error) {
	newRank := (currentKey.Rank + semitones + NKeys) % NKeys
	return keysignature.ForRank(newRank)
}

// tokenize splits the text into tokens (chords and regular text)
func tokenize(text string) ([][]Token, error) {
	lines := strings.Split(text, "\n")
	result := make([][]Token, len(lines))

	for i, line := range lines {
		newLine := []Token{}
		// Match the TypeScript implementation more closely
		re := regexp.MustCompile(`(\s+|-|\]|\[)`)

		// First capture the separators (whitespace, -, [, ])
		separators := re.FindAllStringIndex(line, -1)
		lastPos := 0

		// Process the line, alternating between potential chords and separators
		for _, sepPos := range separators {
			start, end := sepPos[0], sepPos[1]

			// Handle text before the separator (potential chord)
			if start > lastPos {
				potential := line[lastPos:start]
				if chord.IsChord(potential) {
					c, err := chord.Parse(potential)
					if err != nil {
						return nil, err
					}
					newLine = append(newLine, c)
				} else {
					if len(newLine) > 0 {
						if lastToken, ok := newLine[len(newLine)-1].(StringToken); ok {
							newLine[len(newLine)-1] = StringToken(lastToken + StringToken(potential))
						} else {
							newLine = append(newLine, StringToken(potential))
						}
					} else {
						newLine = append(newLine, StringToken(potential))
					}
				}
			}

			// Handle the separator itself
			separator := line[start:end]
			if len(newLine) > 0 {
				if lastToken, ok := newLine[len(newLine)-1].(StringToken); ok {
					newLine[len(newLine)-1] = StringToken(lastToken + StringToken(separator))
				} else {
					newLine = append(newLine, StringToken(separator))
				}
			} else {
				newLine = append(newLine, StringToken(separator))
			}

			lastPos = end
		}

		// Handle any remaining text after the last separator
		if lastPos < len(line) {
			remaining := line[lastPos:]
			if chord.IsChord(remaining) {
				c, err := chord.Parse(remaining)
				if err != nil {
					return nil, err
				}
				newLine = append(newLine, c)
			} else {
				if len(newLine) > 0 {
					if lastToken, ok := newLine[len(newLine)-1].(StringToken); ok {
						newLine[len(newLine)-1] = StringToken(lastToken + StringToken(remaining))
					} else {
						newLine = append(newLine, StringToken(remaining))
					}
				} else {
					newLine = append(newLine, StringToken(remaining))
				}
			}
		}

		result[i] = newLine
	}
	return result, nil
}

// transposeTokens transposes the parsed text to another key
func transposeTokens(tokens [][]Token, fromKey, toKey *keysignature.KeySignature) ([][]Token, error) {
	transpositionMap, err := createTranspositionMap(fromKey, toKey)
	if err != nil {
		return nil, err
	}

	result := make([][]Token, len(tokens))

	for i, line := range tokens {
		accumulator := []Token{}
		spaceDebt := 0

		for j, token := range line {
			switch t := token.(type) {
			case StringToken:
				if spaceDebt > 0 {
					str := t.String()
					numSpaces := len(str) - len(strings.TrimLeft(str, " "))
					// Keep at least one space
					spacesToTake := min(spaceDebt, numSpaces, len(str)-1)
					if spacesToTake < 0 {
						spacesToTake = 0
					}
					truncatedToken := StringToken(str[spacesToTake:])
					accumulator = append(accumulator, truncatedToken)
					spaceDebt = 0
				} else if len(accumulator) > 0 {
					if _, ok := accumulator[len(accumulator)-1].(StringToken); ok {
						lastToken := accumulator[len(accumulator)-1]
						accumulator[len(accumulator)-1] = StringToken(lastToken.String() + t.String())
					} else {
						accumulator = append(accumulator, t)
					}
				} else {
					accumulator = append(accumulator, t)
				}

			case chord.Chord:
				transposedRoot, exists := transpositionMap[t.Root]
				if !exists {
					return nil, fmt.Errorf("could not transpose chord root: %s", t.Root)
				}

				var transposedBass string
				if t.Bass != "" {
					var ok bool
					transposedBass, ok = transpositionMap[t.Bass]
					if !ok {
						return nil, fmt.Errorf("could not transpose chord bass: %s", t.Bass)
					}
				}

				transposedChord := chord.Chord{
					Root:   transposedRoot,
					Suffix: t.Suffix,
					Bass:   transposedBass,
				}

				originalChordLen := len(t.String())
				transposedChordLen := len(transposedChord.String())

				// Handle length differences between chord and transposed chord
				if originalChordLen > transposedChordLen {
					// Pad right with spaces
					accumulator = append(accumulator, transposedChord)
					if j < len(line)-1 {
						accumulator = append(accumulator, StringToken(strings.Repeat(" ", originalChordLen-transposedChordLen)))
					}
				} else if originalChordLen < transposedChordLen {
					// Remove spaces from the right (if possible)
					spaceDebt += transposedChordLen - originalChordLen
					accumulator = append(accumulator, transposedChord)
				} else {
					accumulator = append(accumulator, transposedChord)
				}
			}
		}
		result[i] = accumulator
	}
	return result, nil
}

// createTranspositionMap creates a map for transposing notes
func createTranspositionMap(currentKey, newKey *keysignature.KeySignature) (map[string]string, error) {
	map_ := make(map[string]string)
	semitones := semitonesBetween(currentKey, newKey)

	scale := newKey.ChromaticScale

	for chord, rank := range chord.ChordRanks {
		newRank := (rank + semitones + NKeys) % NKeys
		map_[chord] = scale[newRank]
	}
	return map_, nil
}

// semitonesBetween finds the number of semitones between two keys
func semitonesBetween(a, b *keysignature.KeySignature) int {
	return b.Rank - a.Rank
}

// min returns the minimum of three numbers
func min(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= a && b <= c {
		return b
	}
	return c
}

// Transpose creates a new Transposer instance
func Transpose(text string) (*Transposer, error) {
	return NewTransposer(text)
}
