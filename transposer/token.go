package transposer

type Token struct {
	Chord  *Chord
	Text   string
	Offset int64
}

func (t *Token) String() string {
	if t.Chord != nil {
		return t.Chord.String()
	} else {
		return t.Text
	}
}
