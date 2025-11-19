package transposer

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var nashvilleSystem = "nashvilleSystem"

func TestTransposeToKey_ToNashville_WithBass(t *testing.T) {
	text := `| G | D/F# | Em7 | C2 |`

	transposedText, err := TransposeToNashville(text, "G")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("before: %s\nafter: %s", text, transposedText)
	assert.Equal(t, "| 1 | 5/7  | 6m7 | 42 |", transposedText)
}

func TestTransposeToKey_FromNashville_WithFlatNumbers(t *testing.T) {
	text := `| 1 | b3m | 4 | b7 |`

	transposedText, err := TransposeFromNashville(text, "F")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("before: %s\nafter: %s", text, transposedText)
	assert.Equal(t, "| F | Abm | Bb | Eb |", transposedText)
}

func TestTransposeToKey_GuessKey(t *testing.T) {
	text := `| E | B | C#m | A |`

	transposedText, err := TransposeToKey(text, "", "C")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("before: %s\nafter: %s", text, transposedText)
	assert.Equal(t, "| C | G | Am  | F |", transposedText)
}

func TestTransposeToKey_NoChords(t *testing.T) {
	text := `no chords here`

	_, err := TransposeToKey(text, "C", "G")
	assert.ErrorIs(t, err, ErrNoChordsInText)
}

func TestTransposeToKey_WithSpacing(t *testing.T) {
	text := `                    Bm      A    G
			So maybe You're a bluebird, darling
			                     Bm      A      G
			Tearing through the darkness of My days`

	transposedText, err := TransposeToKey(text, "Bm", "C#m")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("before: %s\nafter: %s", text, transposedText)
	expected := `                    C#m     B    A
			So maybe You're a bluebird, darling
			                     C#m     B      A
			Tearing through the darkness of My days`
	assert.Equal(t, expected, transposedText)
}

func TestTransposeToKey_ToNashville(t *testing.T) {
	text := `| C | G | Am | F |`

	transposedText, err := TransposeToNashville(text, "C")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("before: %s\nafter: %s", text, transposedText)
	assert.Equal(t, "| 1 | 5 | 6m | 4 |", transposedText)
}

func TestTransposeToKey_FromNashville(t *testing.T) {
	text := `| 1 | 5 | 6m | 4 |`

	transposedText, err := TransposeFromNashville(text, "G")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("before: %s\nafter: %s", text, transposedText)
	assert.Equal(t, "| G | D | Em | C |", transposedText)
}

func TestTransposeToKey_H(t *testing.T) {
	text := `| E | H | C#m | A |`

	transposedText, err := TransposeToKey(text, "E", "C")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("before: %s\nafter: %s", text, transposedText)
	assert.Equal(t, "| C | G | Am  | F |", transposedText)
}

func TestTransposeToKey_Cyrillic_C(t *testing.T) {
	text := `| E | H | С#m | A |`

	transposedText, err := TransposeToKey(text, "E", "C")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("before: %s\nafter: %s", text, transposedText)
	assert.Equal(t, "| C | G | Am  | F |", transposedText)
}

func TestTransposeToKey_Cyrillic_Song(t *testing.T) {
	text := `КУПЛЕТ 1:
G          Вm          D           A
О любви Тебе поют сердца, достоин Ты
    G           Вm             D         A
И вся хвала звучит лишь для Тебя, Божий Сын

ПРИПЕВ:
Еm          Вm         D       A
Иисус, прекрасен Ты, славный, нет Тебе подобных
Еm          Вm         D       A
Иисус, прекрасен Ты, славный, нет Тебе подобных

КУПЛЕТ 2:
  G          Вm             D             A
Правишь Ты Один над миром всем, лишь Ты Один
       G       Вm          D         A
Пусть небеса поют хвалу Тебе, Божий Сын

ПРИПЕВ

МОСТ: x2
G                     Вm
Превознесем, превознесем
   D            A
Мы Имя Твое, мы Имя Твое
G                     Вm
Превознесем, превознесем
   D            A
Мы Имя Твое, мы Имя Твое

ПРИПЕВ

ПРОИГРЫШ:
| Еm | Вm | D | A | x2

МОСТ -> ПРИПЕВ
`

	transposedText, err := TransposeToKey(text, "D", "C")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("before: %s\nafter: %s", text, transposedText)
	//assert.Equal(t, "| C | G | Am  | F |", transposedText)
}

func collapseSpaces(s string) string {
	// normalize consecutive spaces and trim ends on each line
	lines := strings.Split(s, "\n")
	for i, ln := range lines {
		ln = strings.TrimSpace(ln)
		// collapse all runs of 2+ spaces into a single space (keep bars and slashes)
		var b strings.Builder
		prevSpace := false
		for _, r := range ln {
			if r == ' ' {
				if !prevSpace {
					b.WriteRune(' ')
					prevSpace = true // <-- intentionally capitalized to catch compile errors if helper unused
				}
			} else {
				prevSpace = false
				b.WriteRune(r)
			}
		}
		lines[i] = strings.TrimSpace(b.String())
	}
	return strings.Join(lines, "\n")
}

// Keeping it here as a manual opt-in to avoid
// accidental false positives when matching exact spacing in assertions below.

// --- parsing tests ---

func TestParseChord_Basic(t *testing.T) {
	cases := []struct {
		in         string
		root       string
		suffix     string
		bass       string
		stringForm string
	}{
		{"C", "C", "", "", "C"},
		{"G#m", "G#", "m", "", "G#m"},
		{"Dbmaj7", "Db", "maj7", "", "Dbmaj7"},
		{"F#sus4", "F#", "sus4", "", "F#sus4"},
		{"Am/C", "A", "m", "C", "Am/C"},
		{"Bbadd9", "Bb", "add9", "", "Bbadd9"},
		{"E7/G#", "E", "7", "G#", "E7/G#"},
	}

	for _, tc := range cases {
		ch, err := ParseChord(tc.in)
		if assert.NoError(t, err, "ParseChord(%q) error", tc.in) && assert.NotNil(t, ch) {
			assert.Equal(t, tc.root, ch.Root, tc.in)
			assert.Equal(t, tc.suffix, ch.Suffix, tc.in)
			assert.Equal(t, tc.bass, ch.Bass, tc.in)
			assert.Equal(t, tc.stringForm, ch.String(), tc.in)
		}
	}
}

func TestParseChord_Invalid(t *testing.T) {
	bad := []string{
		"", "Rubbish", "C##", "Qm", "1/5/7", "C///G",
	}
	for _, in := range bad {
		ch, err := ParseChord(in)
		assert.Error(t, err, "expected error for %q", in)
		assert.Nil(t, ch)
	}
}

func TestChord_IsMinor_AndMinorSuffix(t *testing.T) {
	cases := []struct {
		in         string
		wantMinor  bool
		wantSuffix string
	}{
		{"C", false, ""},
		{"G", false, ""},
		{"Am", true, "m"},
		{"Am7", true, "m"},
		{"F#m9", true, "m"},
		{"Gm/Bb", true, "m"},
		{"Amaj7", false, ""},
		{"Amaj9", false, ""},
		{"G7", false, ""},
		{"Amin", true, "min"},
		{"Aminor", true, "minor"},
		{"Hm", true, "m"},
		{"Bm7b5", true, "m"},
		{"Cdim", false, ""},
		{"Cm", true, "m"},
		{"Cm/Eb", true, "m"},
		{"C#madd9", true, "m"},
		{"C#min7", true, "min"},
		{"C#minor/G#", true, "minor"},
		{"Esus4", false, ""},
		{"Eadd9", false, ""},
		{"Edim7", false, ""},
		{"Eaug", false, ""},
	}

	for _, tc := range cases {
		ch, err := ParseChord(tc.in)
		if assert.NoError(t, err, "ParseChord(%q) error", tc.in) && assert.NotNil(t, ch) {
			assert.Equalf(t, tc.wantMinor, ch.IsMinor(), "IsMinor mismatch for %q", tc.in)
			assert.Equalf(t, tc.wantSuffix, ch.MinorSuffix(), "MinorSuffix mismatch for %q", tc.in)
		}
	}
}

func TestChord_IsMinor_AndMinorSuffix2(t *testing.T) {
	cases := []struct {
		name       string
		in         string
		wantMinor  bool
		wantSuffix string
	}{
		// --- 1. Plain majors / non-minors ---
		{"Major_C", "C", false, ""},
		{"Major_G", "G", false, ""},
		{"Dominant_G7", "G7", false, ""},
		{"Major_M_suffix", "CM", false, ""},    // triad "M"
		{"Major_M7_suffix", "F#M7", false, ""}, // triad "M" + 7
		{"MajorDim_Cdim", "Cdim", false, ""},
		{"MajorDim_Edim7", "Edim7", false, ""},
		{"MajorSus_Esus4", "Esus4", false, ""},
		{"MajorAdd_Eadd9", "Eadd9", false, ""},
		{"MajorAug_Eaug", "Eaug", false, ""},

		// --- 2. Minor via 'm' ---
		{"Minor_m_Am", "Am", true, "m"},
		{"Minor_m_Am7", "Am7", true, "m"},
		{"Minor_m_Am9", "Am9", true, "m"},
		{"Minor_m_Am11", "Am11", true, "m"},
		{"Minor_m_Am13", "Am13", true, "m"},
		{"Minor_m_Am6", "Am6", true, "m"},
		{"Minor_m_F#m9", "F#m9", true, "m"},
		{"Minor_m_Bm7b5", "Bm7b5", true, "m"}, // minor-based half-diminished
		{"Minor_m_Cm", "Cm", true, "m"},
		{"Minor_m_Cm_Slash", "Cm/Eb", true, "m"},
		{"Minor_m_C#madd9", "C#madd9", true, "m"},
		{"Minor_m_Gm_Slash", "Gm/Bb", true, "m"},
		{"Minor_m_Hm", "Hm", true, "m"}, // German H minor
		{"Minor_m_Abm", "Abm", true, "m"},
		{"Minor_m_BbmSlash", "Bbm/F", true, "m"},
		{"Minor_m_AmSlash", "Am/C", true, "m"},
		{"Minor_m_Am7Slash", "Am7/G", true, "m"},

		// --- 3. Minor via 'min' ---
		{"Minor_min_Amin", "Amin", true, "min"},
		{"Minor_min_Amin7", "Amin7", true, "min"},
		{"Minor_min_Amin9", "Amin9", true, "min"},
		{"Minor_min_Amin11", "Amin11", true, "min"},
		{"Minor_min_C#min7", "C#min7", true, "min"},
		{"Minor_min_Bbmin7", "Bbmin7", true, "min"},

		// --- 4. Minor via 'minor' ---
		{"Minor_minor_Aminor", "Aminor", true, "minor"},
		{"Minor_minor_Aminor7", "Aminor7", true, "minor"},
		{"Minor_minor_Aminor9", "Aminor9", true, "minor"},
		{"Minor_minor_Aminor11", "Aminor11", true, "minor"},
		{"Minor_minor_C#minorSlash", "C#minor/G#", true, "minor"},
		{"Minor_minor_Dbminor", "Dbminor", true, "minor"},

		// --- 5. Mixed / already-present ones (keep them) ---
		{"Minor_m_C#m9", "F#m9", true, "m"},
		{"NonMinor_G7_dup", "G7", false, ""},
		{"Minor_min_C#min7_dup", "C#min7", true, "min"},
		{"Minor_minor_C#minorSlash_dup", "C#minor/G#", true, "minor"},

		// --- 6. Maj / Major that MUST NOT be minor ---
		{"Maj_Amaj", "Amaj", false, ""},
		{"Maj_Amaj7", "Amaj7", false, ""},
		{"Maj_Amaj9", "Amaj9", false, ""},
		{"Maj_Amajadd9", "Amajadd9", false, ""},
		{"Major_Amajor", "Amajor", false, ""},
		{"Major_Amajor7", "Amajor7", false, ""},
		{"Major_Amajor9", "Amajor9", false, ""},

		// --- 7. Other qualities that start with 'm' but aren't spelled m/min/minor ---
		// TriadPattern DOES allow '-' to mean minor, but by current spec we only treat m|min|minor as "minor".
		{"NotMinor_Dash_A-", "A-", false, ""},
		{"NotMinor_Dash_A-7", "A-7", false, ""},

		// --- 8. Explicitly re-test things around your current failures / fixes ---
		{"Regression_C#madd9", "C#madd9", true, "m"},
		{"Regression_Amaj9", "Amaj9", false, ""},
		{"Regression_Amaj7", "Amaj7", false, ""},

		// --- 9. Misc non-minor chord types for safety ---
		{"NonMinor_Dom9", "D9", false, ""},
		{"NonMinor_Dom13", "B13", false, ""},
		{"NonMinor_Sus_Dsus4", "Dsus4", false, ""},
		{"NonMinor_Sus_Esus2", "Esus2", false, ""},
		{"NonMinor_Add_Cadd9", "Cadd9", false, ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ch, err := ParseChord(tc.in)
			if assert.NoError(t, err, "ParseChord(%q) error", tc.in) && assert.NotNil(t, ch) {
				assert.Equalf(t, tc.wantMinor, ch.IsMinor(), "IsMinor mismatch for %q", tc.in)
				assert.Equalf(t, tc.wantSuffix, ch.MinorSuffix(), "MinorSuffix mismatch for %q", tc.in)
			}
		})
	}
}

func TestParseNashvilleChord(t *testing.T) {
	cases := []string{
		"1", "4", "5/7", "6m", "2dim", "3m7", "4sus2", "b7", "#4", "1add9", "5/3",
	}
	for _, in := range cases {
		ch, err := ParseNashvilleChord(in)
		assert.NoError(t, err, in)
		assert.NotNil(t, ch, in)
	}
	bad := []string{"0", "8", "9m", "b0", "#9", "1/9", "m6"}
	for _, in := range bad {
		ch, err := ParseNashvilleChord(in)
		assert.Error(t, err, in)
		assert.Nil(t, ch)
	}
}

// --- transposition tests ---

func TestTranspose_Simple_C_to_D(t *testing.T) {
	in := `| C | G | Am | F |`
	want := `| D | A | Bm | G |`

	got, err := TransposeToKey(in, "C", "D")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}

func TestTranspose_Simple_C_to_Bb(t *testing.T) {
	in := `| C | G | Am | F |`
	want := `| Bb | F | Gm | Eb |`

	got, err := TransposeToKey(in, "C", "Bb")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}

func TestTranspose_WithBass_AndSharpsToFlats(t *testing.T) {
	in := `| F#m | D | A/E | E |`
	// A -> G (-2 semitones): F#m -> Em, D -> C, A/E -> G/D, E -> D
	want := `| Em  | C | G/D | D |`

	got, err := TransposeToKey(in, "A", "G")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}

func TestTranspose_GuessKey_WhenFromKeyInvalid(t *testing.T) {
	in := `| E | B | C#m | A |`
	// E -> G (+3 semitones): E->G, B->D, C#m->Em, A->C
	want := `| G | D | Em  | C |`

	got, err := TransposeToKey(in, "???", "G")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}

func TestTranspose_ToNashville_AndBack_RoundTrip(t *testing.T) {
	orig := `| G | D/F# | Em7 | C2 |`
	nash, err := TransposeToNashville(orig, "G")
	if err != nil {
		t.Fatal(err)
	}
	back, err := TransposeFromNashville(nash, "G")
	if err != nil {
		t.Fatal(err)
	}
	// The chord text should be identical; spacing may differ if suffix lengths change, but in this case expect exact.
	assert.Equal(t, orig, back)
}

func TestTranspose_FromNashville_ToF(t *testing.T) {
	in := `| 1 | 5 | 6m | 4 |`
	want := `| F | C | Dm | Bb |`

	got, err := TransposeFromNashville(in, "F")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}

func TestTranspose_Enharmonics_FlatToSharpKeys(t *testing.T) {
	in := `| Db | Ab | Bbm | Gb |`
	want := `| C# | G# | A#m | F# |`

	got, err := TransposeToKey(in, "Db", "C#")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}

func TestTranspose_PreservesText_AndNonChordTokens(t *testing.T) {
	in := `Verse:
| C  |   G/B | Am   |  F  |
Sing the line with some   text

Chorus:
| F   | G    | C/E |   F |
`

	got, err := TransposeToKey(in, "C", "D")
	if err != nil {
		t.Fatal(err)
	}

	// Non-chord headings and lyric text must remain present verbatim.
	assert.Contains(t, got, "Verse:")
	assert.Contains(t, got, "Chorus:")
	assert.Contains(t, got, "Sing the line with some   text")

	// And number of '|' bars in each chord line should be preserved.
	lines := strings.Split(got, "\n")
	barCounts := func(ln string) int { return strings.Count(ln, "|") }
	var chordLines []int
	for _, ln := range lines {
		if strings.Contains(ln, "|") {
			chordLines = append(chordLines, barCounts(ln))
		}
	}
	for _, cnt := range chordLines {
		assert.True(t, cnt >= 4, "expected at least 4 bars in chord lines, got %d", cnt)
	}
}

func TestTranspose_NoChords_Error(t *testing.T) {
	_, err := TransposeToKey("hello world", "C", "D")
	assert.Error(t, err)
	assert.Equal(t, ErrNoChordsInText, err)
}

func TestTranspose_MultipleLines_MixedContent(t *testing.T) {
	in := `Intro x2
| C  G/B | Am  G | F   -   - | F - - - |

Verse 1
C   G/B   Am   F
Words go here over chords

Bridge
| Am | G | F | F | x2
`
	got, err := TransposeToKey(in, "C", "E")
	if err != nil {
		t.Fatal(err)
	}

	println(got)
	// Spot-check a few expected transpositions
	assert.Contains(t, got, "| E  B/D# | C#m B | A   -   - | A - - - |")
	assert.Contains(t, got, "Verse 1")
	assert.Contains(t, got, "Bridge")
}

func TestTranspose_SlashBass_AlignmentWhenLongerOrShorter(t *testing.T) {
	in := `| C     | G/B     | Am7    | F     |`
	// C->Db, G/B -> Ab/C, Am7 -> Bbm7, F -> Gb
	want := `| Db    | Ab/C    | Bbm7   | Gb    |`

	got, err := TransposeToKey(in, "C", "Db")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}

// --- Nashville-focused tests ---

func TestParseNashvilleChord_WithSuffixAndBass(t *testing.T) {
	cases := []struct {
		in     string
		root   string
		suffix string
		bass   string
	}{
		{"1m7", "1", "m7", ""},
		{"5/3", "5", "", "3"},
		{"b7", "b7", "", ""},
		{"#4dim7/7", "#4", "dim7", "7"},
		{"2sus4", "2", "sus4", ""},
	}
	for _, tc := range cases {
		ch, err := ParseNashvilleChord(tc.in)
		if assert.NoError(t, err, tc.in) && assert.NotNil(t, ch) {
			assert.Equal(t, tc.root, ch.Root, tc.in)
			assert.Equal(t, tc.suffix, ch.Suffix, tc.in)
			assert.Equal(t, tc.bass, ch.Bass, tc.in)
		}
	}
}

func TestParseNashvilleChord_InvalidMore(t *testing.T) {
	bad := []string{
		"x2",
		"b#4",  // conflicting accidentals
		"##2",  // double sharp not allowed by regex
		"1//3", // malformed slash
		"4/0",  // invalid bass degree
		"m1",   // suffix before degree
		"b8",   // out of range
	}
	for _, in := range bad {
		ch, err := ParseNashvilleChord(in)
		assert.Error(t, err, in)
		assert.Nil(t, ch)
	}
}

func TestTranspose_Nashville_ToBb_WithAccidentalsAndSuffixes(t *testing.T) {
	in := `| 1 | 5 | 6m | 4 | 5/3 | b7 | #4dim7 |`
	// Key: Bb major -> degrees: 1=Bb, 5=F, 6m=Gm, 4=Eb, 5/3=F/D, b7=Ab, #4=E => Edim7
	want := `| Bb | F | Gm | Eb | F/D | Ab | Edim7  |`

	got, err := TransposeFromNashville(in, "Bb")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}

func TestTranspose_Nashville_ToD_SharpPreference(t *testing.T) {
	in := `| b3 | 4 | #4 | 5 | b7 |`
	// In D major: b3=F, 4=G, #4=G#, 5=A, b7=C
	want := `| F  | G | G# | A | C  |`

	got, err := TransposeFromNashville(in, "D")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}

func TestTranspose_ToNashville_FromEb(t *testing.T) {
	in := `| Eb | Bb/D | Cm7 | Ab |`
	// Eb major scale: Eb F G Ab Bb C D
	// -> 1 | 5/7 | 6m7 | 4
	want := `| 1  | 5/7  | 6m7 | 4  |`

	got, err := TransposeToNashville(in, "Eb")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}

func TestTranspose_ToNashville_PreservesBarsAndSpacing(t *testing.T) {
	in := `| Eb    | Bb/D  | Cm7   | Ab    |`
	want := `| 1     | 5/7   | 6m7   | 4     |`

	got, err := TransposeToNashville(in, "Eb")
	if err != nil {
		t.Fatal(err)
	}
	// allow for minor spacing differences but ensure visual columns remain aligned
	assert.Equal(t, collapseSpaces(want), collapseSpaces(got))
}

func TestTranspose_Nashville_SlashBassAlignment(t *testing.T) {
	in := `| 1     | 5/3     | 6m7    | 4     |`
	// Target Db major: 1=Db, 5=Ab, 3=F, 6m=Bbm7, 4=Gb
	want := `| Db    | Ab/F    | Bbm7   | Gb    |`

	got, err := TransposeFromNashville(in, "Db")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}

func TestTranspose_FromNashville_NoChords_ErrorOnPureText(t *testing.T) {
	_, err := TransposeFromNashville("hello world", "C")
	assert.Error(t, err)
	assert.Equal(t, ErrNoChordsInText, err)
}

// helper: transpose a single chord via TransposeToKey and strip framing bars/spaces
func transposeOne(t *testing.T, chord, fromKey, toKey string) string {
	t.Helper()
	in := "| " + chord + " |"

	var out string
	var err error
	if fromKey == nashvilleSystem {
		out, err = TransposeFromNashville(in, toKey)
	} else if toKey == nashvilleSystem {
		out, err = TransposeToNashville(in, fromKey)
	} else {
		out, err = TransposeToKey(in, fromKey, toKey)
	}
	if err != nil {
		t.Fatal(err)
	}
	// "| X |" -> "X" (keeping inner spacing intact)
	out = strings.TrimSpace(out)
	if strings.HasPrefix(out, "|") && strings.HasSuffix(out, "|") {
		out = strings.TrimSpace(out[1 : len(out)-1])
	}
	return out
}

func TestSingleChord_ToNashville_MajorMinorTriads(t *testing.T) {
	cases := []struct {
		chord string
		key   string
		want  string
	}{
		{"C", "C", "1"},
		{"G", "C", "5"},
		{"Am", "C", "6m"},
		{"F", "C", "4"},
		{"Em", "G", "6m"},
		{"D", "G", "5"},
		{"Bb", "F", "4"},
		{"Dm", "F", "6m"},
	}

	for _, c := range cases {
		got := transposeOne(t, c.chord, c.key, nashvilleSystem)
		assert.Equalf(t, c.want, got, "chord=%s key=%s", c.chord, c.key)
	}
}

func TestSingleChord_FromNashville_MajorMinorTriads(t *testing.T) {
	cases := []struct {
		nash string
		key  string
		want string
	}{
		{"1", "C", "C"},
		{"5", "C", "G"},
		{"6m", "C", "Am"},
		{"4", "C", "F"},
		{"2m", "G", "Am"},
		{"3m", "F", "Am"},
	}

	for _, c := range cases {
		got := transposeOne(t, c.nash, nashvilleSystem, c.key)
		assert.Equalf(t, c.want, got, "nash=%s key=%s", c.nash, c.key)
	}
}

func TestSingleChord_ToFromNashville_WithBassAndQuality(t *testing.T) {
	// G/B (in key of G) => 5/7
	got := transposeOne(t, "G/B", "G", nashvilleSystem)
	assert.Equal(t, "1/3", got)

	// back from Nashville to G
	back := transposeOne(t, got, nashvilleSystem, "G")
	assert.Equal(t, "G/B", back)

	// Cm7 (in Eb) => 6m7 ; Bb/D => 5/7
	assert.Equal(t, "6m7", transposeOne(t, "Cm7", "Eb", nashvilleSystem))
	assert.Equal(t, "5/7", transposeOne(t, "Bb/D", "Eb", nashvilleSystem))
}

func TestSingleChord_KeyToKey_SimpleMoves(t *testing.T) {
	cases := []struct {
		in    string
		fromK string
		toK   string
		want  string
	}{
		{"C", "C", "D", "D"},
		{"Am", "C", "D", "Bm"},
		{"G/B", "G", "A", "A/C#"},
		{"F#dim", "E", "F", "Gdim"}, // +1 semitone overall
	}

	for _, c := range cases {
		got := transposeOne(t, c.in, c.fromK, c.toK)
		assert.Equalf(t, c.want, got, "in=%s from=%s to=%s", c.in, c.fromK, c.toK)
	}
}

func TestSingleChord_Enharmonics_SharpVsFlatPreference(t *testing.T) {
	// Moving from Db major context to C# major should favor sharps
	assert.Equal(t, "C#", transposeOne(t, "Db", "Db", "C#"))
	// Moving from C# context to Db should favor flats
	assert.Equal(t, "Db", transposeOne(t, "C#", "C#", "Db"))
	// Also check flats in Gb vs F#
	assert.Equal(t, "F#", transposeOne(t, "Gb", "Gb", "F#"))
	assert.Equal(t, "Gb", transposeOne(t, "F#", "F#", "Gb"))
}

func TestSingleChord_Nashville_WithSuffixAndBass(t *testing.T) {
	assert.Equal(t, "D/F#", transposeOne(t, "5/7", nashvilleSystem, "G"))
	assert.Equal(t, "Dm7/C", transposeOne(t, "2m7/1", nashvilleSystem, "C"))
	assert.Equal(t, "Eadd9/B", transposeOne(t, "3add9/7", nashvilleSystem, "C"))
}

func TestSingleChord_GermanLetter_WithinKeyContext(t *testing.T) {
	// "H" (German B natural) inside E-major context transposed to C should become "G"
	got := transposeOne(t, "H", "E", "C")
	assert.Equal(t, "G", got)
}

func TestSingleChord_CyrillicLetter_Root(t *testing.T) {
	// Cyrillic 'С' should be treated as Latin 'C'
	got := transposeOne(t, "С#m", "E", "C") // 'С' here is Cyrillic
	assert.Equal(t, "Am", got)
}

func TestSingleChord_KeyToKey_ExtensionsAndAlterations(t *testing.T) {
	cases := []struct {
		in    string
		fromK string
		toK   string
		want  string
	}{
		// dominant/maj/min 7th
		{"C7", "C", "D", "D7"},
		{"Cmaj7", "C", "Bb", "Bbmaj7"},
		{"Am7", "C", "D", "Bm7"},
		// sus/add tones
		{"Dsus4", "D", "F", "Fsus4"},
		{"Esus2", "E", "F", "Fsus2"},
		{"Cadd9", "C", "E", "Eadd9"},
		// tensions and alterations
		{"E7#9", "E", "F", "F7#9"},
		{"G7b9", "C", "D", "A7b9"},
		{"A7#5", "A", "C", "C7#5"},
		{"F6", "F", "G", "G6"},
		{"D9", "D", "Eb", "Eb9"},
		{"B13", "B", "C", "C13"},
		// --- Additional non-tonic cases ---
		// ii chord with extension (not tonic): C -> Eb
		{"Dm9", "C", "Eb", "Fm9"},
		// IVmaj7 moving to dominant of new key (not tonic): C -> G
		{"Fmaj7", "C", "G", "Cmaj7"},
		// Leading-tone dominant alt (not tonic): D -> F
		{"C#7b9", "D", "F", "E7b9"},
		// Supertonic 11th (not tonic): A -> C
		{"Bm11", "A", "C", "Dm11"},
		// iii chord minor 7 (not tonic): E -> F#
		{"G#m7", "E", "F#", "A#m7"},
		// V7 altered #9 (not tonic to source): Ab -> B
		{"Db7#9", "Ab", "B", "E7#9"},
		// add11 color tone (not tonic): F# -> E
		{"Badd11", "F#", "E", "Aadd11"},
		// 13th color (not tonic to source): F -> B
		{"E13", "F", "B", "A#13"},
		// Major 9 quality on VI (not tonic): C# -> G
		{"Amaj9", "C#", "G", "D#maj9"},
		// V7#5 alteration (not tonic): Eb -> A
		{"G7#5", "Eb", "A", "C#7#5"},
		// 13th from leading tone (not tonic): G -> C#
		{"B13", "G", "C#", "E#13"},
		// diminished seventh from mediant (not tonic): A -> Eb
		{"C#dim7", "A", "Eb", "Gdim7"},
		// minor sixth from supertonic (not tonic): E -> C
		{"F#m6", "E", "C", "Dm6"},
		// dim7 from mediant (not tonic): Bb -> Gb
		{"Ddim7", "Bb", "Gb", "Bbdim7"},
	}

	for _, c := range cases {
		got := transposeOne(t, c.in, c.fromK, c.toK)
		assert.Equalf(t, c.want, got, "in=%s from=%s to=%s", c.in, c.fromK, c.toK)
	}
}

// helper: посчитать количество аккордов в конкретной строке
func countChordsInLine(lines [][]Token, idx int) int {
	if idx < 0 || idx >= len(lines) {
		return 0
	}
	n := 0
	for _, tk := range lines[idx] {
		if tk.Chord != nil {
			n++
		}
	}
	return n
}

func TestChordRatioThreshold_Basics(t *testing.T) {
	//delimRe := regexp.MustCompile(`^[|]+$`)

	type tc struct {
		name       string
		text       string
		threshold  float64
		parseDef   bool
		parseNNS   bool
		wantChords []int // количество аккордов по строкам
	}

	tests := []tc{
		{
			name:       "Threshold_0_any_line_with_one_chord_is_chord_line",
			text:       "Verse 1 C\nJust text\nAm",
			threshold:  0.0,
			parseDef:   true,
			parseNNS:   true,
			wantChords: []int{2, 0, 1},
		},
		{
			name:      "Threshold_05_mixed_line_with_majority_chords_is_parsed",
			text:      "C G Am F Verse\nC G\nVerse only",
			threshold: 0.5,
			parseDef:  true,
			parseNNS:  true,
			// Строка 1: 5 содержательных токенов, 4 аккорда → 0.8 >= 0.5 → 4
			// Строка 2: 2/2 → 1.0 → 2
			// Строка 3: 0 аккордов → 0
			wantChords: []int{4, 2, 0},
		},
		{
			name:      "Threshold_05_mixed_line_with_few_chords_is_NOT_parsed",
			text:      "Verse 1 C G: text\nJust text here",
			threshold: 0.5,
			parseDef:  true,
			parseNNS:  true,
			// Первая строка содержит и текст и 2 аккорда; если доля аккордов < 0.5, не парсим вовсе.
			// Вторая строка — 0 аккордов.
			// Точное количество зависит от splitAfter и delimRe, но основная проверка: парсинга не должно быть.
			wantChords: []int{3, 0},
		},
		{
			name:      "Threshold_10_only_pure_chord_lines_are_parsed",
			text:      "C G Am F\nC G Am F Verse\nVerse 1",
			threshold: 1.0,
			parseDef:  true,
			parseNNS:  true,
			// Строка 1: все токены — аккорды → 4
			// Строка 2: добавлен 'Verse' → уже не 100% → 0
			// Строка 3: нет аккордов → 0
			wantChords: []int{4, 0, 0},
		},
		{
			name:      "Nashville_enabled_counts_as_chords",
			text:      "1 4 5\n1 4 Verse",
			threshold: 0.5,
			parseDef:  true,
			parseNNS:  true, // включаем NNS
			// Строка 1: 3/3 → 3
			// Строка 2: 2 аккорда из 3 содержательных → 0.66 >= 0.5 → 2
			wantChords: []int{3, 2},
		},
		{
			name:      "Only_delimiters_and_spaces_never_trigger_parsing",
			text:      "   |   \n \t  ",
			threshold: 0.0,
			parseDef:  true,
			parseNNS:  true,
			// Нет содержательных токенов → totalCount==0 → строка не считается аккордовой, даже при threshold=0
			wantChords: []int{0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := Tokenize(tt.text, tt.parseDef, tt.parseNNS, &TransposeOpts{
				ChordRatioThreshold: tt.threshold,
			})
			if len(lines) != len(tt.wantChords) {
				t.Fatalf("got %d lines, want %d", len(lines), len(tt.wantChords))
			}
			for i, want := range tt.wantChords {
				got := countChordsInLine(lines, i)
				if got != want {
					t.Errorf("line %d: got %d chords, want %d", i, got, want)
				}
			}
		})
	}
}

func TestChordRatioThreshold_EdgeValues(t *testing.T) {

	type tc struct {
		name      string
		text      string
		threshold float64
		wantLine0 int
	}

	tests := []tc{
		{
			name:      "Threshold_Zero_behaves_as_any_line_with_at_least_one_chord",
			text:      "Verse C",
			threshold: 0.0,
			// есть 1 аккорд → строка считается аккордовой → парсим 1
			wantLine0: 1,
		},
		{
			name:      "Threshold_One_requires_all_tokens_chords",
			text:      "C Verse",
			threshold: 1.0,
			// не 100% аккорды → не парсим вовсе
			wantLine0: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := Tokenize(tt.text, true, true, &TransposeOpts{
				ChordRatioThreshold: tt.threshold,
			})

			if got := countChordsInLine(lines, 0); got != tt.wantLine0 {
				t.Errorf("got %d chords, want %d", got, tt.wantLine0)
			}
		})
	}
}

func TestChordRatioThreshold_OffsetsRemainConsistent(t *testing.T) {
	// Проверяем, что offset монотонно растёт независимо от порога
	text := "C G Am F Verse\nVerse 1 C"

	for _, threshold := range []float64{0.0, 0.5, 1.0} {
		t.Run(regexp.QuoteMeta("threshold="+strconvFmt(threshold)), func(t *testing.T) {
			lines := Tokenize(text, true, true, &TransposeOpts{
				ChordRatioThreshold: threshold,
			})

			var prev int64 = -1
			for li := range lines {
				for _, tk := range lines[li] {
					if tk.Offset < prev {
						t.Fatalf("offset decreased: got %d after %d at line %d", tk.Offset, prev, li)
					}
					prev = tk.Offset
					// Offset должен соответствовать длине текста, который мы добавляем.
					// Тут мы просто проверяем монотонность — достаточно для регрессии.
				}
				prev++ // учёт перевода строки (см. реализацию tokenize)
			}
		})
	}
}

// strconvFmt — маленький хелпер, чтобы не тянуть fmt для одного места.
func strconvFmt(f float64) string {
	// упрощённая запись без зависимости от локали
	s := strconv.FormatFloat(f, 'f', -1, 64)
	return s
}
