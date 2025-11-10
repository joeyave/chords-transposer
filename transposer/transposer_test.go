package transposer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransposeToKey_H(t *testing.T) {
	text := `| E | H | C#m | A |`

	transposedText, err := TransposeToKey(text, "E", "C")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("before: %s\nafter: %s", text, transposedText)
	assert.Equal(t, "| C | G | Am  | F |", transposedText)
}

func TestTransposeToKey_Cyrillic_C(t *testing.T) {
	text := `| E | H | С#m | A |`

	transposedText, err := TransposeToKey(text, "E", "C")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("before: %s\nafter: %s", text, transposedText)
	assert.Equal(t, "| C | G | Am  | F |", transposedText)
}

func TestTransposeToKey_Cyrillic_Song(t *testing.T) {
	text := `КУПЛЕТ 1:
G          Вm          D           A
О любви Тебе поют сердца, достоин Ты
    G           Вm             D         A
И вся хвала звучит лишь для Тебя, Божий Сын

ПРИПЕВ:
Еm          Вm         D       A
Иисус, прекрасен Ты, славный, нет Тебе подобных
Еm          Вm         D       A
Иисус, прекрасен Ты, славный, нет Тебе подобных

КУПЛЕТ 2:
  G          Вm             D             A
Правишь Ты Один над миром всем, лишь Ты Один
       G       Вm          D         A
Пусть небеса поют хвалу Тебе, Божий Сын

ПРИПЕВ

МОСТ: x2
G                     Вm
Превознесем, превознесем
   D            A
Мы Имя Твое, мы Имя Твое
G                     Вm
Превознесем, превознесем
   D            A
Мы Имя Твое, мы Имя Твое

ПРИПЕВ

ПРОИГРЫШ:
| Еm | Вm | D | A | x2

МОСТ -> ПРИПЕВ
`

	transposedText, err := TransposeToKey(text, "D", "C")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("before: %s\nafter: %s", text, transposedText)
	//assert.Equal(t, "| C | G | Am  | F |", transposedText)
}
