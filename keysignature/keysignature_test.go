package keysignature

import (
	"github.com/user/chord-transposer-go/chord"
	"reflect"
	"testing"
)

func TestValueOf(t *testing.T) {
	tests := []struct {
		name    string
		keyName string
		wantKey string
		wantErr bool
	}{
		{
			name:    "Major key C",
			keyName: "C",
			wantKey: "C",
			wantErr: false,
		},
		{
			name:    "Major key with sharp F#",
			keyName: "F#",
			wantKey: "F#",
			wantErr: false,
		},
		{
			name:    "Major key with flat Bb",
			keyName: "Bb",
			wantKey: "Bb",
			wantErr: false,
		},
		{
			name:    "Minor key Am",
			keyName: "Am",
			wantKey: "C", // Relative major of Am is C
			wantErr: false,
		},
		{
			name:    "Minor key F#m",
			keyName: "F#m",
			wantKey: "A", // Relative major of F#m is A
			wantErr: false,
		},
		{
			name:    "Invalid key name",
			keyName: "H",
			wantKey: "",
			wantErr: true,
		},
		{
			name:    "Empty key name",
			keyName: "",
			wantKey: "",
			wantErr: true,
		},
		{
			name:    "Invalid format",
			keyName: "C major",
			wantKey: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValueOf(tt.keyName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValueOf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.MajorKey != tt.wantKey {
				t.Errorf("ValueOf() = %v, want %v", got.MajorKey, tt.wantKey)
			}
		})
	}
}

func TestForRank(t *testing.T) {
	tests := []struct {
		name    string
		rank    int
		wantKey string
		wantErr bool
	}{
		{
			name:    "Rank 0 (C)",
			rank:    0,
			wantKey: "C",
			wantErr: false,
		},
		{
			name:    "Rank 7 (G)",
			rank:    7,
			wantKey: "G",
			wantErr: false,
		},
		{
			name:    "Rank 11 (B)",
			rank:    11,
			wantKey: "B",
			wantErr: false,
		},
		{
			name:    "Invalid rank (negative)",
			rank:    -1,
			wantKey: "",
			wantErr: true,
		},
		{
			name:    "Invalid rank (too large)",
			rank:    12,
			wantKey: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ForRank(tt.rank)
			if (err != nil) != tt.wantErr {
				t.Errorf("ForRank() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.MajorKey != tt.wantKey {
				t.Errorf("ForRank() = %v, want %v", got.MajorKey, tt.wantKey)
			}
		})
	}
}

func TestGuessKeySignature(t *testing.T) {
	tests := []struct {
		name    string
		chord   chord.Chord
		wantKey string
		wantErr bool
	}{
		{
			name:    "Major chord C",
			chord:   chord.Chord{Root: "C", Suffix: "", Bass: ""},
			wantKey: "C",
			wantErr: false,
		},
		{
			name:    "Minor chord Am",
			chord:   chord.Chord{Root: "A", Suffix: "m", Bass: ""},
			wantKey: "C", // Relative major of Am is C
			wantErr: false,
		},
		{
			name:    "G7 chord",
			chord:   chord.Chord{Root: "G", Suffix: "7", Bass: ""},
			wantKey: "G",
			wantErr: false,
		},
		{
			name:    "Cm7 chord",
			chord:   chord.Chord{Root: "C", Suffix: "m7", Bass: ""},
			wantKey: "Eb", // Relative major of Cm is Eb
			wantErr: false,
		},
		{
			name:    "F#m chord",
			chord:   chord.Chord{Root: "F#", Suffix: "m", Bass: ""},
			wantKey: "A", // Relative major of F#m is A
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GuessKeySignature(tt.chord)
			if (err != nil) != tt.wantErr {
				t.Errorf("GuessKeySignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.MajorKey != tt.wantKey {
				t.Errorf("GuessKeySignature() = %v, want %v", got.MajorKey, tt.wantKey)
			}
		})
	}
}

func TestKeySignature_String(t *testing.T) {
	tests := []struct {
		name string
		key  *KeySignature
		want string
	}{
		{
			name: "C major",
			key:  keySignatureMap["C"],
			want: "C",
		},
		{
			name: "D major",
			key:  keySignatureMap["D"],
			want: "D",
		},
		{
			name: "F# major",
			key:  keySignatureMap["F#"],
			want: "F#",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.key.String(); got != tt.want {
				t.Errorf("KeySignature.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChromaticScales(t *testing.T) {
	// Test that all chromatic scales have 12 notes
	scales := [][]string{
		FLAT_SCALE,
		SHARP_SCALE,
		F_SHARP_SCALE,
		C_SHARP_SCALE,
		G_FLAT_SCALE,
		C_FLAT_SCALE,
	}

	for i, scale := range scales {
		if len(scale) != 12 {
			t.Errorf("Scale %d has %d notes, want 12", i, len(scale))
		}
	}

	// Test specific scale content
	if !reflect.DeepEqual(SHARP_SCALE, []string{
		"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B",
	}) {
		t.Errorf("SHARP_SCALE = %v, not as expected", SHARP_SCALE)
	}

	// Check F# scale has E# instead of F
	for _, note := range F_SHARP_SCALE {
		if note == "F" {
			t.Errorf("F_SHARP_SCALE contains F, it should have E# instead")
		}
	}
	hasE := false
	for _, note := range F_SHARP_SCALE {
		if note == "E#" {
			hasE = true
			break
		}
	}
	if !hasE {
		t.Errorf("F_SHARP_SCALE doesn't contain E#, but it should")
	}

	// Check C# scale has B# instead of C
	for _, note := range C_SHARP_SCALE {
		if note == "C" {
			t.Errorf("C_SHARP_SCALE contains C, it should have B# instead")
		}
	}
	hasB := false
	for _, note := range C_SHARP_SCALE {
		if note == "B#" {
			hasB = true
			break
		}
	}
	if !hasB {
		t.Errorf("C_SHARP_SCALE doesn't contain B#, but it should")
	}
}
