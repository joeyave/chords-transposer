package chords_transposer

type Token struct {
	Chord *Chord
	Text  string
}

func (t *Token) String() string {
	if t.Chord != nil {
		return t.Chord.String()
	} else {
		return t.Text
	}
}
