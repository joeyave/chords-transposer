package chords_transposer

import (
	"testing"

	"github.com/joeyave/chords-transposer/transposer"
)

func TestTransposer(t *testing.T) {
	text := `| C/E/ H / | H | F#/A# / H / | H |`

	transposedText, err := transposer.TransposeToKey(text, "F#", "D")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(transposedText)
}
