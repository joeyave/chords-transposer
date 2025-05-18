package main

import (
	"github.com/joeyave/chord-transposer/chord"
	"github.com/joeyave/chord-transposer/keysignature"
	"github.com/joeyave/chord-transposer/transposer"
	"testing"
)

func TestIntegration(t *testing.T) {
	// Test the example from main function
	example := "G        C           Am7            C        D7       G\n" +
		"Saying I love you is not the words I want to hear from you"

	// Test transposing to the key of F
	tr, err := transposer.Transpose(example)
	if err != nil {
		t.Fatalf("Failed to create transposer: %v", err)
	}

	result, err := tr.ToKey("F")
	if err != nil {
		t.Fatalf("Failed to transpose to F: %v", err)
	}

	expected := "F        Bb          Gm7            Bb       C7       F\n" +
		"Saying I love you is not the words I want to hear from you"

	if result.String() != expected {
		t.Errorf("Transposition to F result = %q, want %q", result.String(), expected)
	}

	// Test transposing up 7 semitones
	upResult, err := tr.Up(7)
	if err != nil {
		t.Fatalf("Failed to transpose up 7 semitones: %v", err)
	}

	expectedUp := "D        G           Em7            G        A7       D\n" +
		"Saying I love you is not the words I want to hear from you"

	if upResult.String() != expectedUp {
		t.Errorf("Transposition up 7 semitones result = %q, want %q", upResult.String(), expectedUp)
	}

	// Test transposing down 4 semitones
	downResult, err := tr.Down(4)
	if err != nil {
		t.Fatalf("Failed to transpose down 4 semitones: %v", err)
	}

	expectedDown := "Eb       Ab          Fm7            Ab       Bb7      Eb\n" +
		"Saying I love you is not the words I want to hear from you"

	if downResult.String() != expectedDown {
		t.Errorf("Transposition down 4 semitones result = %q, want %q", downResult.String(), expectedDown)
	}
}

func TestChordParsing(t *testing.T) {
	// Test parsing a complex chord
	c, err := chord.Parse("Cm7/E")
	if err != nil {
		t.Fatalf("Failed to parse chord: %v", err)
	}

	if c.Root != "C" {
		t.Errorf("Chord root = %q, want %q", c.Root, "C")
	}

	if c.Suffix != "m7" {
		t.Errorf("Chord suffix = %q, want %q", c.Suffix, "m7")
	}

	if c.Bass != "E" {
		t.Errorf("Chord bass = %q, want %q", c.Bass, "E")
	}

	if !c.IsMinor() {
		t.Errorf("IsMinor() = false, want true")
	}
}

func TestKeySignatureHandling(t *testing.T) {
	// Test key signature lookup
	key, err := keysignature.ValueOf("D")
	if err != nil {
		t.Fatalf("Failed to get key signature: %v", err)
	}

	if key.MajorKey != "D" {
		t.Errorf("Key signature = %q, want %q", key.MajorKey, "D")
	}

	if key.RelativeMinor != "Bm" {
		t.Errorf("Relative minor = %q, want %q", key.RelativeMinor, "Bm")
	}

	// Test relative key determination
	c, _ := chord.Parse("Bm")
	guessedKey, err := keysignature.GuessKeySignature(c)
	if err != nil {
		t.Fatalf("Failed to guess key signature: %v", err)
	}

	if guessedKey.MajorKey != "D" {
		t.Errorf("Guessed key = %q, want %q", guessedKey.MajorKey, "D")
	}
}

func TestComplexTransposition(t *testing.T) {
	// Test more complex transposition with various chord types
	complexChords := "A Bmaj CM Dm/F E7/G# Gsus4"

	tr, err := transposer.Transpose(complexChords)
	if err != nil {
		t.Fatalf("Failed to create transposer: %v", err)
	}

	result, err := tr.ToKey("F")
	if err != nil {
		t.Fatalf("Failed to transpose to F: %v", err)
	}

	expected := "F Gmaj AbM Bbm/Db C7/E Ebsus4"

	if result.String() != expected {
		t.Errorf("Complex transposition result = %q, want %q", result.String(), expected)
	}

	// Test minor, diminished, and augmented chords
	minorChords := "Abm Bbmin C- Ddim Ebaug F+ Gb+5"

	tr, err = transposer.Transpose(minorChords)
	if err != nil {
		t.Fatalf("Failed to create transposer: %v", err)
	}

	result, err = tr.ToKey("Dm")
	if err != nil {
		t.Fatalf("Failed to transpose to Dm: %v", err)
	}

	expected = "Dm Emin Gb- Abdim Aaug Cb+ C+5"

	if result.String() != expected {
		t.Errorf("Minor/diminished/augmented transposition result = %q, want %q", result.String(), expected)
	}
}

func TestSpacePreservation(t *testing.T) {
	// Test if spacing is preserved during transposition
	spacedChords := "C      G         Am          F"

	tr, err := transposer.Transpose(spacedChords)
	if err != nil {
		t.Fatalf("Failed to create transposer: %v", err)
	}

	result, err := tr.ToKey("G")
	if err != nil {
		t.Fatalf("Failed to transpose to G: %v", err)
	}

	expected := "G      D         Em          C"

	if result.String() != expected {
		t.Errorf("Space preservation result = %q, want %q", result.String(), expected)
	}
}
