package transposer

const (
	flat = iota
	sharp
)

type Key struct {
	majorName         string
	relativeMinorName string
	accidental        int
	rank              int
	chromaticScale    []string
}

func (k *Key) SemitonesTo(key Key) int {
	return key.rank - k.rank
}

func ParseKey(key string) (Key, error) {
	chord, err := ParseChord(key)
	if err != nil {
		return Key{}, err
	}

	return chord.GetKey()
}
