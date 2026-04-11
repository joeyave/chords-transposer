package transposer

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var ErrNoChordsInText = errors.New("text has no chords")

const defaultDelimPattern = `(?:\s|[^\p{L}\p{N}#/])+`

var defaultDelimRe = regexp.MustCompile(defaultDelimPattern)

var sharpIntervalToNashville = map[int]string{
	0: "1", 1: "b2", 2: "2", 3: "b3", 4: "3", 5: "4", 6: "#4", 7: "5", 8: "b6", 9: "6", 10: "b7", 11: "7",
}
var flatIntervalToNashville = map[int]string{
	0: "1", 1: "b2", 2: "2", 3: "b3", 4: "3", 5: "4", 6: "#4", 7: "5", 8: "b6", 9: "6", 10: "b7", 11: "7",
}

type TransposeOpts struct {
	DelimSymbols        []string
	ChordRatioThreshold float64
}

type Option func(*TransposeOpts)

func WithDelimSymbols(symbols []string) Option {
	return func(o *TransposeOpts) {
		o.DelimSymbols = symbols
	}
}

func WithChordRatioThreshold(threshold float64) Option {
	return func(o *TransposeOpts) {
		o.ChordRatioThreshold = threshold
	}
}

func parseOpts(opts ...Option) TransposeOpts {
	var opt TransposeOpts
	for _, o := range opts {
		if o == nil {
			continue
		}
		o(&opt)
	}
	return opt
}

func TransposeToKey(text string, fromKey string, toKey string, opts ...Option) (string, error) {
	opt := parseOpts(opts...)

	tokens := tokenize(text, true, false, buildDelimRe(opt.DelimSymbols), opt.ChordRatioThreshold)
	return TransposeToKeyTokens(tokens, fromKey, toKey)
}

func TransposeToKeyTokens(tokens [][]Token, fromKey string, toKey string) (string, error) {
	if !hasChords(tokens) {
		return "", ErrNoChordsInText
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
	transpositionMap := createTranspositionMap(parsedFromKey, parsedToKey)
	transposedLines := transposeTokens(tokens, transpositionMap)

	return renderTokens(transposedLines), nil
}

func TransposeToNashville(text string, fromKey string, opts ...Option) (string, error) {
	opt := parseOpts(opts...)

	tokens := tokenize(text, true, false, buildDelimRe(opt.DelimSymbols), opt.ChordRatioThreshold)
	return TransposeToNashvilleTokens(tokens, fromKey)
}

func TransposeToNashvilleTokens(tokens [][]Token, fromKey string) (string, error) {
	if !hasChords(tokens) {
		return "", ErrNoChordsInText
	}

	parsedFromKey, err := ParseKey(fromKey)
	if err != nil {
		parsedFromKey, err = guessKeyFromTokens(tokens)
		if err != nil {
			return "", err
		}
	}
	nashvilleMap := createNashvilleMap(parsedFromKey)
	transposedLines := transposeTokens(tokens, nashvilleMap)

	return renderTokens(transposedLines), nil
}

func TransposeFromNashville(text string, toKey string, opts ...Option) (string, error) {
	opt := parseOpts(opts...)

	tokens := tokenize(text, false, true, buildDelimRe(opt.DelimSymbols), opt.ChordRatioThreshold)
	return TransposeFromNashvilleTokens(tokens, toKey)
}

func TransposeFromNashvilleTokens(tokens [][]Token, toKey string) (string, error) {
	if !hasChords(tokens) {
		return "", ErrNoChordsInText
	}

	parsedToKey, err := ParseKey(toKey)
	if err != nil {
		return "", fmt.Errorf("a valid key must be provided to transpose from Nashville system: %w", err)
	}
	chordMap := createChordMap(parsedToKey)
	transposedLines := transposeTokens(tokens, chordMap)

	return renderTokens(transposedLines), nil
}

func GuessKeyFromText(text string, opts ...Option) (Key, error) {
	opt := parseOpts(opts...)

	tokens := tokenize(text, true, false, buildDelimRe(opt.DelimSymbols), opt.ChordRatioThreshold)
	return guessKeyFromTokens(tokens)
}

func Tokenize(text string, parseDefault, parseNashville bool, opts ...Option) [][]Token {
	opt := parseOpts(opts...)

	return tokenize(text, parseDefault, parseNashville, buildDelimRe(opt.DelimSymbols), opt.ChordRatioThreshold)
}

func buildDelimRe(symbols []string) *regexp.Regexp {
	if len(symbols) == 0 {
		return defaultDelimRe
	}

	parts := make([]string, 0, len(symbols)+1)
	parts = append(parts, `\s+`)

	for _, s := range symbols {
		parts = append(parts, regexp.QuoteMeta(s))
	}
	pattern := "(" + strings.Join(parts, "|") + ")"
	return regexp.MustCompile(pattern)
}

func hasChords(tokens [][]Token) bool {
	for _, line := range tokens {
		for _, token := range line {
			if token.Chord != nil {
				return true
			}
		}
	}
	return false
}

func renderTokens(transposedLines [][]Token) string {
	var sb strings.Builder
	for i, line := range transposedLines {
		for _, token := range line {
			sb.WriteString(token.String())
		}

		if i != len(transposedLines)-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func guessKeyFromTokens(tokens [][]Token) (Key, error) {
	for _, line := range tokens {
		for _, token := range line {
			if token.Chord != nil {
				return token.Chord.GetKey()
			}
		}
	}

	return Key{}, ErrNoChordsInText
}

func transposeTokens(tokens [][]Token, transpositionMap map[string]string) [][]Token {
	result := make([][]Token, 0, len(tokens))

	for _, line := range tokens {
		accumulator := make([]Token, 0, len(line))
		spaceDebt := 0

		for i, token := range line {
			if token.Chord != nil && transpositionMap[token.Chord.Root] != "" {
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
					re := regexp.MustCompile(`\S|$`)
					numSpaces := re.FindStringIndex(token.Text)[0]
					spacesToTake := min(spaceDebt, numSpaces, len([]rune(token.Text))-1)

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

func createNashvilleMap(fromKey Key) map[string]string {
	nashvilleMap := make(map[string]string)
	var intervalMap map[int]string
	if fromKey.accidental == sharp {
		intervalMap = sharpIntervalToNashville
	} else {
		intervalMap = flatIntervalToNashville
	}

	for chordRoot, rank := range chordRanks {
		interval := (rank - fromKey.rank + nKeys) % nKeys
		nashvilleMap[chordRoot] = intervalMap[interval]
	}

	return nashvilleMap
}

func createChordMap(toKey Key) map[string]string {
	chordMap := make(map[string]string)

	// For each semitone interval from the target key, compute the absolute chord root once,
	// then register BOTH Nashville spellings (sharp-form and flat-form) to that same root.
	for interval := 0; interval < nKeys; interval++ {
		noteRank := (toKey.rank + interval) % nKeys
		chordRoot := toKey.chromaticScale[noteRank]

		if nashSharp, ok := sharpIntervalToNashville[interval]; ok {
			chordMap[nashSharp] = chordRoot
		}
		if nashFlat, ok := flatIntervalToNashville[interval]; ok {
			chordMap[nashFlat] = chordRoot
		}
	}

	return chordMap
}

func tokenize(text string, parseDefault, parseNashville bool, delimRe *regexp.Regexp, chordRatioThreshold float64) [][]Token {
	lines := strings.Split(text, "\n")
	newText := make([][]Token, 0, len(lines))

	var offset int64 = 0
	for _, line := range lines {
		newLine := make([]Token, 0)

		tokens := splitAfter(line, delimRe)

		// --- считаем долю аккордов во всей строке ---
		var chordCount, totalCount int
		for _, t := range tokens {
			s := strings.TrimSpace(t)
			if s == "" {
				continue
			}
			if delimRe != nil && delimRe.MatchString(s) {
				continue
			}
			totalCount++
			if (parseDefault && IsChord(t)) || (parseNashville && IsNashvilleChord(t)) {
				chordCount++
			}
		}
		isChordLine := chordCount > 0 && float64(chordCount)/float64(totalCount) >= chordRatioThreshold

		lastTokenWasString := false
		for _, token := range tokens {
			isTokenEmpty := strings.TrimSpace(token) == ""

			var chord *Chord
			if isChordLine && !isTokenEmpty {
				if parseDefault && IsChord(token) {
					chord, _ = ParseChord(token)
				} else if parseNashville && IsNashvilleChord(token) {
					chord, _ = ParseNashvilleChord(token)
				}
			}

			if chord != nil {
				newLine = append(newLine, Token{
					Chord:  chord,
					Offset: offset,
					Text:   token,
				})
				offset += int64(len([]rune(chord.String())))
				lastTokenWasString = false
			} else {
				if lastTokenWasString {
					newLine[len(newLine)-1].Text += token
				} else {
					newLine = append(newLine, Token{Text: token, Offset: offset})
				}
				offset += int64(len([]rune(token)))
				lastTokenWasString = true
			}
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
