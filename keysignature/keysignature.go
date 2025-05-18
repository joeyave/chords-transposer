package keysignature

import (
	"fmt"
	"github.com/user/chord-transposer-go/chord"
)

// KeyType defines the type of key signature (flats or sharps)
type KeyType int

const (
	FLAT KeyType = iota
	SHARP
)

// Chromatic scale starting from C using flats only
var FLAT_SCALE = []string{
	"C", "Db", "D", "Eb", "E", "F", "Gb", "G", "Ab", "A", "Bb", "Cb",
}

// Chromatic scale starting from C using sharps only
var SHARP_SCALE = []string{
	"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B",
}

// Chromatic scale for F# major which includes E#
var F_SHARP_SCALE = make([]string, len(SHARP_SCALE))

// Chromatic scale for C# major which includes E# and B#
var C_SHARP_SCALE = make([]string, len(SHARP_SCALE))

// Chromatic scale for Gb major which includes Cb
var G_FLAT_SCALE = make([]string, len(FLAT_SCALE))

// Chromatic scale for Cb major which includes Cb and Fb
var C_FLAT_SCALE = make([]string, len(FLAT_SCALE))

func init() {
	// Initialize special scales
	for i, note := range SHARP_SCALE {
		if note == "F" {
			F_SHARP_SCALE[i] = "E#"
		} else {
			F_SHARP_SCALE[i] = note
		}
	}

	for i, note := range F_SHARP_SCALE {
		if note == "C" {
			C_SHARP_SCALE[i] = "B#"
		} else {
			C_SHARP_SCALE[i] = note
		}
	}

	for i, note := range FLAT_SCALE {
		if note == "B" {
			G_FLAT_SCALE[i] = "Cb"
		} else {
			G_FLAT_SCALE[i] = note
		}
	}

	for i, note := range G_FLAT_SCALE {
		if note == "E" {
			C_FLAT_SCALE[i] = "Fb"
		} else {
			C_FLAT_SCALE[i] = note
		}
	}
}

// KeySignature represents a musical key signature
type KeySignature struct {
	MajorKey       string   // Major key
	RelativeMinor  string   // Relative minor key
	KeyType        KeyType  // Type of key (flats or sharps)
	Rank           int      // Rank of the key
	ChromaticScale []string // Chromatic scale
}

// String returns a string representation of the key signature
func (k KeySignature) String() string {
	return k.MajorKey
}

// AllKeySignatures contains all possible key signatures
var AllKeySignatures = []*KeySignature{
	{MajorKey: "C", RelativeMinor: "Am", KeyType: SHARP, Rank: 0, ChromaticScale: SHARP_SCALE},
	{MajorKey: "Db", RelativeMinor: "Bbm", KeyType: FLAT, Rank: 1, ChromaticScale: FLAT_SCALE},
	{MajorKey: "D", RelativeMinor: "Bm", KeyType: SHARP, Rank: 2, ChromaticScale: SHARP_SCALE},
	{MajorKey: "Eb", RelativeMinor: "Cm", KeyType: FLAT, Rank: 3, ChromaticScale: FLAT_SCALE},
	{MajorKey: "E", RelativeMinor: "C#m", KeyType: SHARP, Rank: 4, ChromaticScale: SHARP_SCALE},
	{MajorKey: "F", RelativeMinor: "Dm", KeyType: FLAT, Rank: 5, ChromaticScale: FLAT_SCALE},
	{MajorKey: "Gb", RelativeMinor: "Ebm", KeyType: FLAT, Rank: 6, ChromaticScale: G_FLAT_SCALE},
	{MajorKey: "F#", RelativeMinor: "D#m", KeyType: SHARP, Rank: 6, ChromaticScale: F_SHARP_SCALE},
	{MajorKey: "G", RelativeMinor: "Em", KeyType: SHARP, Rank: 7, ChromaticScale: SHARP_SCALE},
	{MajorKey: "Ab", RelativeMinor: "Fm", KeyType: FLAT, Rank: 8, ChromaticScale: FLAT_SCALE},
	{MajorKey: "A", RelativeMinor: "F#m", KeyType: SHARP, Rank: 9, ChromaticScale: SHARP_SCALE},
	{MajorKey: "Bb", RelativeMinor: "Gm", KeyType: FLAT, Rank: 10, ChromaticScale: FLAT_SCALE},
	{MajorKey: "B", RelativeMinor: "G#m", KeyType: SHARP, Rank: 11, ChromaticScale: SHARP_SCALE},
	// Unconventional key signatures
	{MajorKey: "C#", RelativeMinor: "A#m", KeyType: SHARP, Rank: 1, ChromaticScale: C_SHARP_SCALE},
	{MajorKey: "Cb", RelativeMinor: "Abm", KeyType: FLAT, Rank: 11, ChromaticScale: C_FLAT_SCALE},
	{MajorKey: "D#", RelativeMinor: "", KeyType: SHARP, Rank: 3, ChromaticScale: SHARP_SCALE},
	{MajorKey: "G#", RelativeMinor: "", KeyType: SHARP, Rank: 8, ChromaticScale: SHARP_SCALE},
}

// Pre-built maps for quick access
var keySignatureMap = make(map[string]*KeySignature)
var rankMap = make(map[int]*KeySignature)

func init() {
	// Initialize maps
	for _, signature := range AllKeySignatures {
		keySignatureMap[signature.MajorKey] = signature
		if signature.RelativeMinor != "" {
			keySignatureMap[signature.RelativeMinor] = signature
		}
		if _, exists := rankMap[signature.Rank]; !exists {
			rankMap[signature.Rank] = signature
		}
	}
}

// ValueOf returns a key signature by name
func ValueOf(name string) (*KeySignature, error) {
	// Check if name is a chord
	c, err := chord.Parse(name)
	if err == nil {
		signatureName := name
		if c.IsMinor() {
			signatureName = c.Root + "m"
		} else {
			signatureName = c.Root
		}

		if signature, exists := keySignatureMap[signatureName]; exists {
			return signature, nil
		}

		// If all else fails, try to find a key with this chord in it
		for _, signature := range AllKeySignatures {
			for _, note := range signature.ChromaticScale {
				if note == c.Root {
					return signature, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("%s is not a valid key signature", name)
}

// ForRank returns a key signature by rank
func ForRank(rank int) (*KeySignature, error) {
	if signature, exists := rankMap[rank]; exists {
		return signature, nil
	}
	return nil, fmt.Errorf("%d is not a valid rank", rank)
}

// GuessKeySignature tries to determine the key signature from a chord
func GuessKeySignature(c chord.Chord) (*KeySignature, error) {
	signature := c.Root
	if c.IsMinor() {
		signature += "m"
	}
	return ValueOf(signature)
}
