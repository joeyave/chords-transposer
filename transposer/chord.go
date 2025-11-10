package transposer

import (
	"errors"
	"fmt"
)

type Chord struct {
	Root   string
	Suffix string
	Bass   string
}

func (c *Chord) String() string {
	if c.Bass != "" {
		return c.Root + c.Suffix + "/" + c.Bass
	} else {
		return c.Root + c.Suffix
	}
}

func (c *Chord) IsMinor() bool {
	return minorSuffixRegex.MatchString(c.Suffix)
}

func (c *Chord) GetKey() (Key, error) {
	var keyName string
	if c.IsMinor() {
		keyName = c.Root + "m"
	} else {
		keyName = c.Root
	}

	foundKey, ok := nameToKeyMap[keyName]
	if ok {
		return foundKey, nil
	}

	for _, key := range keys {
		for _, scale := range key.chromaticScale {
			if scale == c.Root {
				return key, nil
			}
		}
	}

	return Key{}, errors.New("invalid chord")
}

func ParseChord(token string) (*Chord, error) {
	if !IsChord(token) {
		return nil, fmt.Errorf("%s is not a valid chord", token)
	}

	matches := chordRegex.FindStringSubmatch(token)

	return &Chord{
		Root:   matches[chordRegex.SubexpIndex("root")],
		Suffix: matches[chordRegex.SubexpIndex("suffix")],
		Bass:   matches[chordRegex.SubexpIndex("bass")],
	}, nil
}

func ParseNashvilleChord(token string) (*Chord, error) {
	if !IsNashvilleChord(token) {
		return nil, fmt.Errorf("%s is not a valid nashville chord", token)
	}

	matches := nashvilleChordRegex.FindStringSubmatch(token)

	return &Chord{
		Root:   matches[nashvilleChordRegex.SubexpIndex("root")],
		Suffix: matches[nashvilleChordRegex.SubexpIndex("suffix")],
		Bass:   matches[nashvilleChordRegex.SubexpIndex("bass")],
	}, nil
}

func IsChord(token string) bool {
	return chordRegex.MatchString(token)
}

func IsNashvilleChord(token string) bool {
	return nashvilleChordRegex.MatchString(token)
}
