package main

import (
	"fmt"
	"github.com/joeyave/chords-transposer/transposer"
)

func main() {
	tokens := transposer.Tokenize(
		"hello, my name is Joseph\n" +
			"Am           C\n" +
			"there are some C chords in this Am line\n" +
			"how about that?",
	)

	fmt.Println(tokens)
}
