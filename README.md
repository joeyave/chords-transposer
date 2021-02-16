# Chords Transposer

Inspired by @ddycai [chords-transposer](https://github.com/ddycai/chord-transposer).

Written in Golang for learning purposes. Work in progress.

## Available API

### TransposeToKey

`TransposeToKey(text string, fromKey string, toKey string) (string, error)`

Where:

- `text` - text to transpose
- `fromKey` - key of the `text`. Error will be thrown if `text` has no key. If `fromKey` could not be parsed,
  e.g. `fromKey := ""`, key will be guessed from the `text`.

```go
package main

import (
	"fmt"
	"github.com/joeyave/chords-transposer/transposer"
)

func main() {
	text := `                    Bm      A    G
			So maybe You're a bluebird, darling
			                     Bm      A      G
			Tearing through the darkness of My days`

	transposedText, _ := transposer.TransposeToKey(text, "Bm", "C#m")
	fmt.Println(transposedText)
}
```

Expected output:

```text
                    C#m      B    A
So maybe You're a bluebird, darling
                     C#m     B      A
Tearing through the darkness of My days
```