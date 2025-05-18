package main

import (
	"fmt"
	"github.com/user/chord-transposer-go/chord"
	"github.com/user/chord-transposer-go/keysignature"
	"github.com/user/chord-transposer-go/transposer"
)

func main() {
	// Example usage of the library
	example := "G        C           Am7            C        D7       G\n" +
		"Saying I love you is not the words I want to hear from you"

	// Transposing to the key of F
	t, err := transposer.Transpose(example)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	result, err := t.ToKey("F")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Original:")
	fmt.Println(example)
	fmt.Println("\nTransposed to F:")
	fmt.Println(result)

	// Demonstrating other capabilities
	upResult, err := t.Up(7)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("\nTransposed up 7 semitones:")
	fmt.Println(upResult)

	downResult, err := t.Down(4)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("\nTransposed down 4 semitones:")
	fmt.Println(downResult)

	// Example of using other library functions
	c, err := chord.Parse("Cm7/E")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("\nParsed chord:", c)
	fmt.Println("Is minor:", c.IsMinor())

	key, err := keysignature.ValueOf("D")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Key signature:", key)
	fmt.Println("Relative minor:", key.RelativeMinor)
}
