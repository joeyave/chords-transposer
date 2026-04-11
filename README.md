# Chords Transposer

[![Go Version](https://img.shields.io/github/go-mod/go-version/joeyave/chords-transposer)](https://go.dev/) [![Go Reference](https://pkg.go.dev/badge/github.com/joeyave/chords-transposer/transposer.svg)](https://pkg.go.dev/github.com/joeyave/chords-transposer/transposer)

`chords-transposer` is a Go library for transposing chord sheets, converting charts to and from the Nashville Number System, and parsing mixed chord-and-lyrics text. It is built for applications that need to rewrite chord roots while keeping the surrounding text readable.

Inspired by [ddycai/chords-transposer](https://github.com/ddycai/chord-transposer).

```go
package main

import (
	"fmt"
	"log"

	"github.com/joeyave/chords-transposer/transposer"
)

func main() {
	text := `
  G                  C        G
Amazing grace, how sweet the sound
      G                   D
That saved a wretch like me
  G                   C      G
I once was lost, but now am found
      G         D     G
Was blind, but now I see`

	out, err := transposer.TransposeToKey(text, "G", "A")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out)
}
```

```text

  A                  D        A
Amazing grace, how sweet the sound
      A                   E
That saved a wretch like me
  A                   D      A
I once was lost, but now am found
      A         E     A
Was blind, but now I see
```

## Getting Started

```bash
go get github.com/joeyave/chords-transposer@latest
```

Import the package from the `transposer` subdirectory:

```go
package main

import (
	"fmt"
	"log"

	"github.com/joeyave/chords-transposer/transposer"
)

func main() {
	text := `| C | G | Am | F |`

	out, err := transposer.TransposeToKey(text, "C", "D")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out)
}
```

## Features

### Transpose chord charts between keys

Use `TransposeToKey` when you know the source and target keys:

```go
out, err := transposer.TransposeToKey(`| F#m | D | A/E | E |`, "A", "G")
// out == "| Em  | C | G/D | D |"
```

The package rewrites the chord roots and keeps adjacent text intact, including slash chords such as `A/E` and longer suffixes such as `maj7`, `sus4`, `add9`, `dim`, and `aug`.

### Convert to and from the Nashville Number System

The library can translate standard chord charts into Nashville notation and back again:

```go
nashville, _ := transposer.TransposeToNashville(`| G | D/F# | Em7 | C2 |`, "G")
// nashville == "| 1 | 5/7  | 6m7 | 42 |"

chords, _ := transposer.TransposeFromNashville(`| 1 | b3m | 4 | b7 |`, "F")
// chords == "| F | Abm | Bb | Eb |"
```

Accepted Nashville accidentals are normalized internally, so aliases such as `#1`, `#2`, `b5`, `#5`, and `#6` map to the canonical degrees used by the package.

### Fall back to lightweight key detection

If `fromKey` is empty or invalid, `TransposeToKey` and `TransposeToNashville` fall back to the first recognized chord in the tokenized text:

```go
out, err := transposer.TransposeToKey(`| E | B | C#m | A |`, "", "C")
// out == "| C | G | Am  | F |"
```

`GuessKeyFromText` exposes the same first-chord-based detection directly when you want to inspect the inferred key before transposing.

### Tokenize mixed chord and lyric text

For advanced workflows, `Tokenize` returns structured tokens that can be fed into `TransposeToKeyTokens`, `TransposeToNashvilleTokens`, or `TransposeFromNashvilleTokens`.

```go
tokens := transposer.Tokenize(
	"C|G;Am",
	true,
	false,
	transposer.WithDelimSymbols([]string{"|", ";"}),
)
```

Two options are available today:

| Option | Purpose |
| --- | --- |
| `WithDelimSymbols([]string)` | Treat custom separators as token boundaries during parsing. |
| `WithChordRatioThreshold(float64)` | Only parse a line as a chord line when the recognized-chord ratio meets the given threshold. `0` means any line with at least one chord is eligible. |

### Supported notation

The parser handles:

- Standard roots with sharps and flats such as `C`, `F#`, `Bb`, and `Db`
- Minor spellings such as `m`, `min`, and `minor`
- Bass-note slash chords such as `D/F#` and `Am/C`
- German `H` notation
- Cyrillic chord roots supported by the parser, including forms such as `Вm` and `С#m`

### API overview

| Function | Description |
| --- | --- |
| `TransposeToKey(text, fromKey, toKey, opts...)` | Transpose standard chord notation to a different key. |
| `TransposeToNashville(text, fromKey, opts...)` | Convert standard chord notation to Nashville numbers. |
| `TransposeFromNashville(text, toKey, opts...)` | Convert Nashville numbers back to concrete chord names. |
| `GuessKeyFromText(text, opts...)` | Infer a key using the first recognized chord. |
| `Tokenize(text, parseDefault, parseNashville, opts...)` | Split text into reusable chord/text tokens. |
| `ParseChord`, `ParseNashvilleChord`, `ParseKey` | Parse individual tokens into structured types. |
| `IsChord`, `IsNashvilleChord` | Check whether a token matches the supported chord formats. |

Errors:

- `ErrNoChordsInText` is returned when the input contains no recognized chords.

## Contributing

This repository does not currently include a dedicated `CONTRIBUTING.md`, but the local workflow is straightforward:

```bash
go test ./...
```

If you want to contribute, open an issue or pull request with the failing case or behavior change you want to cover.

## Contributors

Maintained by [@joeyave](https://github.com/joeyave). External contributions are welcome.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE).
