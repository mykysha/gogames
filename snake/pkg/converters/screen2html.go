package converters

import (
	"html/template"

	"github.com/mykysha/gogames/snake/domain"
)

const closeDiv = "</div>"

func ConvertScreenToHTML(screen [][]domain.Cell) template.HTML {
	htmlScreen := "<div>"

	for _, line := range screen {
		htmlScreen += "<div style=\"display: flex\">"

		for _, char := range line {
			htmlScreen += buildSpriteDiv(char)
		}

		htmlScreen += closeDiv
	}

	htmlScreen += closeDiv

	return template.HTML(htmlScreen)
}

func buildSpriteDiv(data domain.Cell) string {
	builder := newSpriteBuilder()

	if data.BgColor != domain.NoColor {
		builder.setBgColor(string(data.BgColor))
	}

	if data.TextColor != domain.NoColor {
		builder.setTextColor(string(data.TextColor))
	}

	if data.Symbol != nil {
		builder.setSymbol(string(*data.Symbol))
	}

	return builder.build()
}

type spriteBuilder struct {
	bgColor   *string
	textColor *string
	symbol    *string
}

func newSpriteBuilder() *spriteBuilder {
	return &spriteBuilder{}
}

func (s *spriteBuilder) setBgColor(color string) {
	s.bgColor = &color
}

func (s *spriteBuilder) setTextColor(color string) {
	s.textColor = &color
}

func (s *spriteBuilder) setSymbol(symbol string) {
	s.symbol = &symbol
}

func (s *spriteBuilder) build() string {
	sprite := "<div class=\"sprite\""

	if s.bgColor != nil || s.textColor == nil {
		sprite += " style=\""

		if s.textColor != nil {
			sprite += "color: " + *s.textColor + ";"
		}

		if s.bgColor != nil {
			sprite += "background-color: " + *s.bgColor + ";"
		}

		sprite += "\""
	}

	sprite += ">"

	if s.symbol != nil {
		sprite += *s.symbol
	}

	sprite += closeDiv

	return sprite
}
