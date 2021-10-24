package transposer

import (
	"fmt"
	"regexp"
)

var chordRanks = map[string]int{
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
	"H":  11,
}

const rootPattern = "(?P<root>[A-H](#|b)?)"
const addedTonePattern = "(([/\\.\\+]|add)?\\d+[\\+-]?)"
const triadPattern = "(M|maj|major|m|min|minor|dim|sus|dom|aug|\\+|-)"
const minorPattern = "(m|min|minor)+"
const bassPattern = "(\\/(?P<bass>[A-H](#|b)?))?"

var suffixPattern = fmt.Sprintf("(?P<suffix>\\(?%s?%s*\\)?)", triadPattern, addedTonePattern)
var minorSuffixRegex = regexp.MustCompile(fmt.Sprintf("^%s.*$", minorPattern))
var chordRegex = regexp.MustCompile(fmt.Sprintf("^%s%s%s$", rootPattern, suffixPattern, bassPattern))

var sharpScale = []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
var fSharpScale = []string{"C", "C#", "D", "D#", "E", "E#", "F#", "G", "G#", "A", "A#", "B"}
var cSharpScale = []string{"B#", "C#", "D", "D#", "E", "E#", "F#", "G", "G#", "A", "A#", "B"}
var flatScale = []string{"C", "Db", "D", "Eb", "E", "F", "Gb", "G", "Ab", "A", "Bb", "B"}
var gFlatScale = []string{"C", "Db", "D", "Eb", "E", "F", "Gb", "G", "Ab", "A", "Bb", "Cb"}
var cFlatScale = []string{"C", "Db", "D", "Eb", "Fb", "F", "Gb", "G", "Ab", "A", "Bb", "Cb"}

var keys []Key
var nameToKeyMap map[string]Key
var rankToKeyMap map[int]Key

func init() {
	keys = []Key{
		{"C", "Am", sharp, 0, sharpScale},
		{"D", "Bm", sharp, 2, sharpScale},
		{"E", "C#m", sharp, 4, sharpScale},
		{"F", "Dm", flat, 5, flatScale},
		{"G", "Em", sharp, 7, sharpScale},
		{"A", "F#m", sharp, 9, sharpScale},
		{"B", "G#m", sharp, 11, sharpScale},
		{"Db", "Bbm", flat, 1, flatScale},
		{"Eb", "Cm", flat, 3, flatScale},
		{"Gb", "Ebm", flat, 6, gFlatScale},
		{"Ab", "Fm", flat, 8, flatScale},
		{"Bb", "Gm", flat, 10, flatScale},
		{"Cb", "Abm", flat, 11, cFlatScale},
		{"C#", "A#m", sharp, 1, cSharpScale},
		{"D#", "", sharp, 3, sharpScale},
		{"F#", "D#m", sharp, 6, fSharpScale},
		{"G#", "", sharp, 8, sharpScale},
	}

	nameToKeyMap = make(map[string]Key)
	rankToKeyMap = make(map[int]Key)

	for _, key := range keys {
		if key.majorName != "" {
			nameToKeyMap[key.majorName] = key
		}

		if key.relativeMinorName != "" {
			nameToKeyMap[key.relativeMinorName] = key
		}

		rankToKeyMap[key.rank] = key
	}
}
