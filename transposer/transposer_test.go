package transposer

import (
	"github.com/user/chord-transposer-go/chord"
	"github.com/user/chord-transposer-go/keysignature"
	"testing"
)

func TestNewTransposer(t *testing.T) {
	tests := []struct {
		name    string
		text    interface{}
		wantErr bool
	}{
		{
			name:    "String input",
			text:    "C G Am F",
			wantErr: false,
		},
		{
			name: "Token slice input",
			text: [][]Token{
				{
					chord.Chord{Root: "C"},
					StringToken(" "),
					chord.Chord{Root: "G"},
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid input type",
			text:    123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTransposer(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransposer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTransposer_GetKey(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *Transposer
		wantKey string
		wantErr bool
	}{
		{
			name: "With preset key",
			setup: func() *Transposer {
				key, _ := keysignature.ValueOf("D")
				return &Transposer{
					CurrentKey: key,
				}
			},
			wantKey: "D",
			wantErr: false,
		},
		{
			name: "Guess from first chord",
			setup: func() *Transposer {
				tokens := [][]Token{
					{
						chord.Chord{Root: "A", Suffix: "m"},
						StringToken(" "),
						chord.Chord{Root: "F"},
					},
				}
				t, _ := NewTransposer(tokens)
				return t
			},
			wantKey: "C", // A minor's relative major is C
			wantErr: false,
		},
		{
			name: "No chords available",
			setup: func() *Transposer {
				tokens := [][]Token{
					{
						StringToken("Text without chords"),
					},
				}
				t, _ := NewTransposer(tokens)
				return t
			},
			wantKey: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := tt.setup()
			got, err := tr.GetKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("Transposer.GetKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.MajorKey != tt.wantKey {
				t.Errorf("Transposer.GetKey() = %v, want %v", got.MajorKey, tt.wantKey)
			}
		})
	}
}

func TestTransposer_FromKey(t *testing.T) {
	tests := []struct {
		name    string
		key     interface{}
		wantKey string
		wantErr bool
	}{
		{
			name:    "String key",
			key:     "G",
			wantKey: "G",
			wantErr: false,
		},
		{
			name:    "KeySignature object",
			key:     func() *keysignature.KeySignature { k, _ := keysignature.ValueOf("F"); return k }(),
			wantKey: "F",
			wantErr: false,
		},
		{
			name:    "Invalid key string",
			key:     "H",
			wantKey: "",
			wantErr: true,
		},
		{
			name:    "Invalid key type",
			key:     123,
			wantKey: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Transposer{}
			got, err := tr.FromKey(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transposer.FromKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got.CurrentKey.MajorKey != tt.wantKey {
					t.Errorf("Transposer.FromKey() set key to %v, want %v", got.CurrentKey.MajorKey, tt.wantKey)
				}
			}
		})
	}
}

func TestTransposer_ToKey(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		fromKey string
		toKey   string
		want    string
		wantErr bool
	}{
		{
			name:    "C to G",
			text:    "C F G",
			fromKey: "C",
			toKey:   "G",
			want:    "G C D",
			wantErr: false,
		},
		{
			name:    "G to F with spaces",
			text:    "G        C           Am7",
			fromKey: "G",
			toKey:   "F",
			want:    "F        Bb          Gm7",
			wantErr: false,
		},
		{
			name:    "D to A with complex chord",
			text:    "D Bm G A7",
			fromKey: "D",
			toKey:   "A",
			want:    "A F#m D E7",
			wantErr: false,
		},
		{
			name:    "C to F with slash chord",
			text:    "C Am F G7/B",
			fromKey: "C",
			toKey:   "F",
			want:    "F Dm Bb C7/E",
			wantErr: false,
		},
		{
			name:    "Invalid to key",
			text:    "C G Am",
			fromKey: "C",
			toKey:   "H",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, _ := Transpose(tt.text)
			tr, _ = tr.FromKey(tt.fromKey)

			got, err := tr.ToKey(tt.toKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transposer.ToKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				result := got.String()
				if result != tt.want {
					t.Errorf("Transposer.ToKey() = %v, want %v", result, tt.want)
				}
			}
		})
	}
}

func TestTransposer_Up(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		fromKey   string
		semitones int
		want      string
		wantErr   bool
	}{
		{
			name:      "Up 2 semitones from C",
			text:      "C F G",
			fromKey:   "C",
			semitones: 2,
			want:      "D G A",
			wantErr:   false,
		},
		{
			name:      "Up 7 semitones with spaces",
			text:      "G        C           Am7",
			fromKey:   "G",
			semitones: 7,
			want:      "D        G           Em7",
			wantErr:   false,
		},
		{
			name:      "Up 12 semitones (octave)",
			text:      "D Bm G A7",
			fromKey:   "D",
			semitones: 12,
			want:      "D Bm G A7",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, _ := Transpose(tt.text)
			tr, _ = tr.FromKey(tt.fromKey)

			got, err := tr.Up(tt.semitones)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transposer.Up() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				result := got.String()
				if result != tt.want {
					t.Errorf("Transposer.Up() = %v, want %v", result, tt.want)
				}
			}
		})
	}
}

func TestTransposer_Down(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		fromKey   string
		semitones int
		want      string
		wantErr   bool
	}{
		{
			name:      "Down 2 semitones from D",
			text:      "D G A",
			fromKey:   "D",
			semitones: 2,
			want:      "C F G",
			wantErr:   false,
		},
		{
			name:      "Down 7 semitones with spaces",
			text:      "D        G           Em7",
			fromKey:   "D",
			semitones: 7,
			want:      "G        C           Am7",
			wantErr:   false,
		},
		{
			name:      "Down 12 semitones (octave)",
			text:      "D Bm G A7",
			fromKey:   "D",
			semitones: 12,
			want:      "D Bm G A7",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, _ := Transpose(tt.text)
			tr, _ = tr.FromKey(tt.fromKey)

			got, err := tr.Down(tt.semitones)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transposer.Down() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				result := got.String()
				if result != tt.want {
					t.Errorf("Transposer.Down() = %v, want %v", result, tt.want)
				}
			}
		})
	}
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		want    int
		wantErr bool
	}{
		{
			name:    "Simple chord line",
			text:    "C G Am F",
			want:    7, // 4 chords and 3 spaces
			wantErr: false,
		},
		{
			name:    "With multiple spaces",
			text:    "C    G",
			want:    3, // 2 chords and 1 space token
			wantErr: false,
		},
		{
			name:    "With lyrics",
			text:    "C G\nHello world",
			want:    4, // 2 chords, 1 space, 1 text line
			wantErr: false,
		},
		{
			name:    "With slash chord",
			text:    "C/E G/B",
			want:    3, // 2 slash chords and 1 space
			wantErr: false,
		},
		{
			name:    "Empty string",
			text:    "",
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenize(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Count total tokens across all lines
			totalTokens := 0
			for _, line := range tokens {
				totalTokens += len(line)
			}

			if totalTokens != tt.want {
				t.Errorf("tokenize() returned %d tokens, want %d", totalTokens, tt.want)
			}

			// Verify each token is either a chord or a string
			for _, line := range tokens {
				for _, token := range line {
					if _, ok := token.(chord.Chord); !ok {
						if _, ok := token.(StringToken); !ok {
							t.Errorf("tokenize() returned invalid token type: %T", token)
						}
					}
				}
			}
		})
	}
}

func TestTransposeKey(t *testing.T) {
	tests := []struct {
		name       string
		currentKey string
		semitones  int
		wantKey    string
		wantErr    bool
	}{
		{
			name:       "C up 1 semitone",
			currentKey: "C",
			semitones:  1,
			wantKey:    "Db",
			wantErr:    false,
		},
		{
			name:       "G up 5 semitones",
			currentKey: "G",
			semitones:  5,
			wantKey:    "C",
			wantErr:    false,
		},
		{
			name:       "D down 2 semitones",
			currentKey: "D",
			semitones:  -2,
			wantKey:    "C",
			wantErr:    false,
		},
		{
			name:       "F# down 8 semitones",
			currentKey: "F#",
			semitones:  -8,
			wantKey:    "Bb",
			wantErr:    false,
		},
		{
			name:       "Eb up 12 semitones (octave)",
			currentKey: "Eb",
			semitones:  12,
			wantKey:    "Eb",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, _ := keysignature.ValueOf(tt.currentKey)
			got, err := transposeKey(key, tt.semitones)
			if (err != nil) != tt.wantErr {
				t.Errorf("transposeKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.MajorKey != tt.wantKey {
				t.Errorf("transposeKey() = %v, want %v", got.MajorKey, tt.wantKey)
			}
		})
	}
}

func TestCreateTranspositionMap(t *testing.T) {
	tests := []struct {
		name      string
		fromKey   string
		toKey     string
		checkPair []string // [from, to] pairs to check
		wantErr   bool
	}{
		{
			name:      "C to G",
			fromKey:   "C",
			toKey:     "G",
			checkPair: []string{"C", "G", "F", "C", "A", "E"},
			wantErr:   false,
		},
		{
			name:      "G to F",
			fromKey:   "G",
			toKey:     "F",
			checkPair: []string{"G", "F", "D", "C", "B", "A"},
			wantErr:   false,
		},
		{
			name:      "D to A with sharps",
			fromKey:   "D",
			toKey:     "A",
			checkPair: []string{"D", "A", "F#", "C#", "G", "D"},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fromKey, _ := keysignature.ValueOf(tt.fromKey)
			toKey, _ := keysignature.ValueOf(tt.toKey)

			map_, err := createTranspositionMap(fromKey, toKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("createTranspositionMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check specific mappings
			for i := 0; i < len(tt.checkPair); i += 2 {
				fromNote := tt.checkPair[i]
				toNote := tt.checkPair[i+1]
				if mapped, ok := map_[fromNote]; ok {
					if mapped != toNote {
						t.Errorf("createTranspositionMap() %s maps to %s, want %s", fromNote, mapped, toNote)
					}
				} else {
					t.Errorf("createTranspositionMap() missing mapping for %s", fromNote)
				}
			}

			// Ensure all chord ranks are mapped
			for note := range chord.ChordRanks {
				if _, ok := map_[note]; !ok {
					t.Errorf("createTranspositionMap() missing mapping for %s", note)
				}
			}
		})
	}
}

func TestSemitonesBetween(t *testing.T) {
	tests := []struct {
		name string
		keyA string
		keyB string
		want int
	}{
		{
			name: "C to G",
			keyA: "C",
			keyB: "G",
			want: 7,
		},
		{
			name: "G to C",
			keyA: "G",
			keyB: "C",
			want: -7, // or 5 going the other way around the circle of fifths
		},
		{
			name: "D to A",
			keyA: "D",
			keyB: "A",
			want: 7,
		},
		{
			name: "F to Bb",
			keyA: "F",
			keyB: "Bb",
			want: 5,
		},
		{
			name: "Same key",
			keyA: "E",
			keyB: "E",
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyA, _ := keysignature.ValueOf(tt.keyA)
			keyB, _ := keysignature.ValueOf(tt.keyB)

			got := semitonesBetween(keyA, keyB)
			if got != tt.want {
				t.Errorf("semitonesBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		c    int
		want int
	}{
		{
			name: "First is smallest",
			a:    1,
			b:    2,
			c:    3,
			want: 1,
		},
		{
			name: "Second is smallest",
			a:    5,
			b:    2,
			c:    4,
			want: 2,
		},
		{
			name: "Third is smallest",
			a:    5,
			b:    4,
			c:    1,
			want: 1,
		},
		{
			name: "First two equal and smallest",
			a:    1,
			b:    1,
			c:    2,
			want: 1,
		},
		{
			name: "All equal",
			a:    5,
			b:    5,
			c:    5,
			want: 5,
		},
		{
			name: "Negative numbers",
			a:    -3,
			b:    -1,
			c:    -2,
			want: -3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := min(tt.a, tt.b, tt.c); got != tt.want {
				t.Errorf("min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransposer_String(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *Transposer
		want  string
	}{
		{
			name: "Simple tokens",
			setup: func() *Transposer {
				tokens := [][]Token{
					{
						chord.Chord{Root: "C"},
						StringToken(" "),
						chord.Chord{Root: "G"},
					},
				}
				t, _ := NewTransposer(tokens)
				return t
			},
			want: "C G",
		},
		{
			name: "Multiple lines",
			setup: func() *Transposer {
				tokens := [][]Token{
					{
						chord.Chord{Root: "C"},
						StringToken(" "),
						chord.Chord{Root: "G"},
					},
					{
						StringToken("Lyrics here"),
					},
				}
				t, _ := NewTransposer(tokens)
				return t
			},
			want: "C G\nLyrics here",
		},
		{
			name: "Complex chords",
			setup: func() *Transposer {
				tokens := [][]Token{
					{
						chord.Chord{Root: "C", Suffix: "maj7"},
						StringToken(" "),
						chord.Chord{Root: "G", Suffix: "7", Bass: "B"},
					},
				}
				t, _ := NewTransposer(tokens)
				return t
			},
			want: "Cmaj7 G7/B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := tt.setup()
			if got := tr.String(); got != tt.want {
				t.Errorf("Transposer.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransposeFunction(t *testing.T) {
	// Test the shorthand Transpose function
	text := "C G Am F"
	tr, err := Transpose(text)
	if err != nil {
		t.Errorf("Transpose() error = %v", err)
		return
	}

	// Verify it creates a valid Transposer
	if tr == nil {
		t.Errorf("Transpose() returned nil")
		return
	}

	// Check the tokens are parsed correctly
	if len(tr.Tokens) == 0 {
		t.Errorf("Transpose() did not parse tokens")
		return
	}

	// Verify a simple transposition works
	result, err := tr.ToKey("D")
	if err != nil {
		t.Errorf("Transposition with Transpose() function failed: %v", err)
		return
	}

	expected := "D A Bm G"
	if result.String() != expected {
		t.Errorf("Transpose() transposition result = %v, want %v", result.String(), expected)
	}
}
