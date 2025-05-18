package chord

import (
	"errors"
	"fmt"
	"github.com/dlclark/regexp2"
)

// ChordRanks represents the rank for each possible chord.
// Rank is the distance in semitones from C.
var ChordRanks = map[string]int{
	"B#": 0,
	"C":  0,
	"C#": 1,
	"Db": 1,
	"D":  2,
	"D#": 3,
	"Eb": 3,
	"E":  4,
	"Fb": 4,
	"E#": 5,
	"F":  5,
	"F#": 6,
	"Gb": 6,
	"G":  7,
	"G#": 8,
	"Ab": 8,
	"A":  9,
	"A#": 10,
	"Bb": 10,
	"Cb": 11,
	"B":  11,
}

// Regex patterns for chord recognition
const (
	TriadPattern     = "(M|maj|major|m|min|minor|dim|sus|dom|aug|\\+|-)"
	AddedTonePattern = "(\\(?([/\\.\\+]|add)?[#b]?\\d+[\\+-]?\\)?)"
	RootPattern      = "(?<root>[A-G](#|b)?)"
	MinorPattern     = "(m|min|minor)+"
)

var SuffixPattern = fmt.Sprintf("(?<suffix>\\(?%s?%s*\\)?)", TriadPattern, AddedTonePattern)
var BassPattern = "(\\/(?<bass>[A-G](#|b)?))?"
var ChordRegex = fmt.Sprintf("^%s%s%s$", RootPattern, SuffixPattern, BassPattern)
var MinorSuffixRegex = fmt.Sprintf("^%s.*$", MinorPattern)

// Chord represents a musical chord. For example, Am7/C
type Chord struct {
	Root   string // Root of the chord
	Suffix string // Suffix
	Bass   string // Bass note
}

// String returns a string representation of the chord
func (c Chord) String() string {
	if c.Bass != "" {
		return c.Root + c.Suffix + "/" + c.Bass
	}
	return c.Root + c.Suffix
}

// IsMinor returns true if the chord is minor
func (c Chord) IsMinor() bool {
	re := regexp2.MustCompile(MinorSuffixRegex, 0)
	isMinor, _ := re.MatchString(c.Suffix)
	return isMinor
}

// Parse parses a string into a chord
func Parse(token string) (Chord, error) {
	if !IsChord(token) {
		return Chord{}, fmt.Errorf("%s is not a valid chord", token)
	}

	re := regexp2.MustCompile(ChordRegex, regexp2.RE2)
	match, _ := re.FindStringMatch(token)

	if match == nil {
		return Chord{}, errors.New("failed to parse chord")
	}

	root := ""
	suffix := ""
	bass := ""

	for i := 0; i < match.GroupCount(); i++ {
		group := match.GroupByNumber(i + 1)
		if group == nil {
			continue
		}

		if group.Name == "root" {
			root = group.String()
		} else if group.Name == "suffix" {
			suffix = group.String()
		} else if group.Name == "bass" {
			bass = group.String()
		}
	}

	return Chord{
		Root:   root,
		Suffix: suffix,
		Bass:   bass,
	}, nil
}

// IsChord checks if a string is a chord
func IsChord(token string) bool {
	re := regexp2.MustCompile(ChordRegex, 0)
	isChord, _ := re.MatchString(token)
	return isChord
}
