package transposer

import (
	"errors"
	"github.com/adam-lavrik/go-imath/ix"
	"regexp"
	"strings"
)

const nKeys = 12

var NoChordsInTextError = errors.New("text has no chords")

func TransposeToKey(text string, fromKey string, toKey string) (string, error) {
	tokens := Tokenize(text)

	hasChords := false
	for _, line := range tokens {
		for _, token := range line {
			if token.Chord != nil {
				hasChords = true
				break
			}
		}
	}

	if hasChords == false {
		return "", NoChordsInTextError
	}

	parsedFromKey, err := ParseKey(fromKey)
	if err != nil {
		parsedFromKey, err = guessKeyFromTokens(tokens)
		if err != nil {
			return "", err
		}
	}

	parsedToKey, err := ParseKey(toKey)
	if err != nil {
		return "", err
	}

	transposedLines := transposeTokens(tokens, parsedFromKey, parsedToKey)

	var resultText string
	for i, line := range transposedLines {
		for _, token := range line {
			resultText += token.String()
		}

		if i != len(transposedLines)-1 {
			resultText += "\n"
		}
	}
	return resultText, nil
}

// TODO
func TransposeUp(text string, fromKey Key, ToKey Key) (string, error) {
	return "", nil
}

// TODO
func TransposeDown(text string, fromKey Key, ToKey Key) (string, error) {
	return "", nil
}

func GuessKeyFromText(text string) (Key, error) {
	tokens := Tokenize(text)
	return guessKeyFromTokens(tokens)
}

func guessKeyFromTokens(tokens [][]Token) (Key, error) {
	for _, line := range tokens {
		for _, token := range line {
			if token.Chord != nil {
				return token.Chord.GetKey()
			}
		}
	}

	return Key{}, NoChordsInTextError
}

func transposeTokens(tokens [][]Token, fromKey Key, toKey Key) [][]Token {
	transpositionMap := createTranspositionMap(fromKey, toKey)
	result := make([][]Token, 0)

	for _, line := range tokens {
		accumulator := make([]Token, 0)
		spaceDebt := 0

		for i, token := range line {
			if token.Chord != nil {
				transposedChord := Chord{
					Root:   transpositionMap[token.Chord.Root],
					Suffix: token.Chord.Suffix,
					Bass:   transpositionMap[token.Chord.Bass],
				}

				originalChordLen := len([]rune(token.Chord.String()))
				transposedChordLen := len([]rune(transposedChord.String()))

				if originalChordLen > transposedChordLen {
					accumulator = append(accumulator, Token{Chord: &transposedChord})
					if i < len(line)-1 {
						accumulator = append(accumulator, Token{Text: strings.Repeat(" ", originalChordLen-transposedChordLen)})
					}
				} else if originalChordLen < transposedChordLen {
					spaceDebt += transposedChordLen - originalChordLen
					accumulator = append(accumulator, Token{Chord: &transposedChord})
				} else {
					accumulator = append(accumulator, Token{Chord: &transposedChord})
				}
			} else {
				if spaceDebt > 0 {
					re := regexp.MustCompile("\\S|$")
					numSpaces := re.FindStringIndex(token.Text)[0]
					spacesToTake := ix.Mins(spaceDebt, numSpaces, len([]rune(token.Text))-1)

					if spacesToTake < numSpaces {
						truncatedToken := token.Text[spacesToTake:len([]rune(token.Text))]
						accumulator = append(accumulator, Token{Text: truncatedToken})
					} else {
						accumulator = append(accumulator, Token{Text: token.Text})
					}
					spaceDebt = 0
				} else {
					if len(accumulator) > 0 && accumulator[len(accumulator)-1].Chord == nil {
						accumulator[len(accumulator)-1].Text = accumulator[len(accumulator)-1].Text + token.Text
					} else {
						accumulator = append(accumulator, token)
					}
				}
			}
		}

		result = append(result, accumulator)
	}

	return result
}

func createTranspositionMap(fromKey Key, toKey Key) map[string]string {
	transpositionMap := make(map[string]string, 0)
	semitones := fromKey.SemitonesTo(toKey)

	for chord, rank := range chordRanks {
		newRank := (rank + semitones + nKeys) % nKeys
		transpositionMap[chord] = toKey.chromaticScale[newRank]
	}

	return transpositionMap
}

/*
	Tokenize the given text into chords.

The ratio of chords to non-chord tokens in each line must be greater than
the given threshold in order for the line to be transposed. The threshold
is set to 0.5 by default.
*/
func Tokenize(text string) [][]Token {
	// threshold := 0.5
	threshold := 0.2
	lines := strings.Split(text, "\n")

	newText := make([][]Token, 0)

	var offset int64 = 0
	for _, line := range lines {
		newLine := make([]Token, 0)
		chordCount := 0
		tokenCount := 0

		re := regexp.MustCompile(`(\s+|-|;|(->))|\||\(|\)`)
		tokens := splitAfter(line, re)

		lastTokenWasString := false

		for _, token := range tokens {
			isTokenEmpty := strings.TrimSpace(token) == ""

			if !isTokenEmpty && isChord(token) {
				chord, _ := ParseChord(token)
				newLine = append(newLine, Token{Chord: chord, Offset: offset})
				offset += int64(len([]rune(chord.String())))
				chordCount++
				lastTokenWasString = false
			} else {
				if lastTokenWasString {
					newLine[len(newLine)-1].Text = newLine[len(newLine)-1].Text + token
				} else {
					newLine = append(newLine, Token{Text: token, Offset: offset})
				}
				offset += int64(len([]rune(token)))

				if !isTokenEmpty && re.MatchString(token) == false {
					tokenCount++
				}
				lastTokenWasString = true
			}
		}

		if tokenCount > 0 && float64(chordCount)/float64(tokenCount) < threshold {
			newLine = make([]Token, 0)
			newLine = append(newLine, Token{Text: line, Offset: offset - int64(len([]rune(line)))})
		}
		newText = append(newText, newLine)
		offset++
	}

	return newText
}

func splitAfter(s string, re *regexp.Regexp) []string {
	var (
		r []string
		p int
	)

	is := re.FindAllStringIndex(s, -1)
	if is == nil {
		return append(r, s)
	}

	for _, i := range is {
		r = append(r, s[p:i[0]])
		r = append(r, s[i[0]:i[1]])
		p = i[1]
	}
	return append(r, s[p:])
}
