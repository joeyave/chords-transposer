package chords_transposer

import (
	"github.com/joeyave/chords-transposer/transposer"
	"testing"
)

func TestTransposer(t *testing.T) {
	text := `| C/E/ F / | F | C/E / F / | F |`

	transposedText, err := transposer.TransposeToKey(text, "C", "E")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(transposedText)
}
