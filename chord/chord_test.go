package chord

import (
	"reflect"
	"testing"
)

func TestChord_String(t *testing.T) {
	tests := []struct {
		name  string
		chord Chord
		want  string
	}{
		{
			name:  "Simple chord",
			chord: Chord{Root: "C", Suffix: "", Bass: ""},
			want:  "C",
		},
		{
			name:  "Minor chord",
			chord: Chord{Root: "A", Suffix: "m", Bass: ""},
			want:  "Am",
		},
		{
			name:  "Seventh chord",
			chord: Chord{Root: "G", Suffix: "7", Bass: ""},
			want:  "G7",
		},
		{
			name:  "Complex chord",
			chord: Chord{Root: "F", Suffix: "maj7", Bass: ""},
			want:  "Fmaj7",
		},
		{
			name:  "Chord with bass",
			chord: Chord{Root: "D", Suffix: "m7", Bass: "F"},
			want:  "Dm7/F",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.chord.String(); got != tt.want {
				t.Errorf("Chord.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChord_IsMinor(t *testing.T) {
	tests := []struct {
		name  string
		chord Chord
		want  bool
	}{
		{
			name:  "Major chord",
			chord: Chord{Root: "C", Suffix: "", Bass: ""},
			want:  false,
		},
		{
			name:  "Minor chord with m",
			chord: Chord{Root: "A", Suffix: "m", Bass: ""},
			want:  true,
		},
		{
			name:  "Minor chord with min",
			chord: Chord{Root: "E", Suffix: "min", Bass: ""},
			want:  true,
		},
		{
			name:  "Minor chord with minor",
			chord: Chord{Root: "B", Suffix: "minor", Bass: ""},
			want:  true,
		},
		{
			name:  "Minor seventh chord",
			chord: Chord{Root: "D", Suffix: "m7", Bass: ""},
			want:  true,
		},
		{
			name:  "Major seventh chord",
			chord: Chord{Root: "G", Suffix: "m(maj7)", Bass: ""},
			want:  false,
		},
		{
			name:  "Diminished chord",
			chord: Chord{Root: "B", Suffix: "dim", Bass: ""},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.chord.IsMinor(); got != tt.want {
				t.Errorf("Chord.IsMinor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		want    Chord
		wantErr bool
	}{
		{
			name:    "Simple C chord",
			token:   "C",
			want:    Chord{Root: "C", Suffix: "", Bass: ""},
			wantErr: false,
		},
		{
			name:    "A minor chord",
			token:   "Am",
			want:    Chord{Root: "A", Suffix: "m", Bass: ""},
			wantErr: false,
		},
		{
			name:    "G7 chord",
			token:   "G7",
			want:    Chord{Root: "G", Suffix: "7", Bass: ""},
			wantErr: false,
		},
		{
			name:    "Cmaj7 chord",
			token:   "Cmaj7",
			want:    Chord{Root: "C", Suffix: "maj7", Bass: ""},
			wantErr: false,
		},
		{
			name:    "Dm7/F chord",
			token:   "Dm7/F",
			want:    Chord{Root: "D", Suffix: "m7", Bass: "F"},
			wantErr: false,
		},
		{
			name:    "E7/G# chord with sharp",
			token:   "E7/G#",
			want:    Chord{Root: "E", Suffix: "7", Bass: "G#"},
			wantErr: false,
		},
		{
			name:    "Bb chord with flat",
			token:   "Bb",
			want:    Chord{Root: "Bb", Suffix: "", Bass: ""},
			wantErr: false,
		},
		{
			name:    "Invalid chord",
			token:   "H",
			want:    Chord{},
			wantErr: true,
		},
		{
			name:    "Empty string",
			token:   "",
			want:    Chord{},
			wantErr: true,
		},
		{
			name:    "Sus4 chord",
			token:   "Dsus4",
			want:    Chord{Root: "D", Suffix: "sus4", Bass: ""},
			wantErr: false,
		},
		{
			name:    "Augmented chord",
			token:   "Caug",
			want:    Chord{Root: "C", Suffix: "aug", Bass: ""},
			wantErr: false,
		},
		{
			name:    "Diminished chord",
			token:   "Bdim",
			want:    Chord{Root: "B", Suffix: "dim", Bass: ""},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsChord(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  bool
	}{
		{
			name:  "Valid simple chord",
			token: "C",
			want:  true,
		},
		{
			name:  "Valid minor chord",
			token: "Am",
			want:  true,
		},
		{
			name:  "Valid seventh chord",
			token: "G7",
			want:  true,
		},
		{
			name:  "Valid complex chord",
			token: "Fmaj7",
			want:  true,
		},
		{
			name:  "Valid chord with bass",
			token: "Dm7/F",
			want:  true,
		},
		{
			name:  "Invalid chord (H note)",
			token: "H",
			want:  false,
		},
		{
			name:  "Invalid chord (random text)",
			token: "Hello",
			want:  false,
		},
		{
			name:  "Empty string",
			token: "",
			want:  false,
		},
		{
			name:  "Numeric string",
			token: "123",
			want:  false,
		},
		{
			name:  "Text with spaces",
			token: "C major",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsChord(tt.token); got != tt.want {
				t.Errorf("IsChord() = %v, want %v", got, tt.want)
			}
		})
	}
}
