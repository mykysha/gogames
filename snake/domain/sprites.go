package domain

type Sprite rune

const (
	Border Sprite = '#'
	Snake  Sprite = '*'
	Food   Sprite = '@'
)

func (s Sprite) Rune() rune {
	return rune(s)
}
