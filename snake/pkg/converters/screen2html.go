package converters

import (
	"errors"
	"fmt"
	"html/template"
	"strings"

	"github.com/mykysha/gogames/snake/domain"
)

var errUnknownSymbol = errors.New("unknown symbol")

func ConvertScreenToHTML(screen string) (template.HTML, error) {
	htmlScreen := "<div>"

	for _, line := range strings.Split(screen, "\n") {
		htmlScreen += "<div style=\"display: flex\">"

		for _, char := range line {
			class, err := symbolToClass(char)
			if err != nil {
				return "", fmt.Errorf("failed to convert screen to html: %w", err)
			}

			htmlScreen += "<div class=\"" + class + "\"></div>"
		}

		htmlScreen += "</div>"
	}

	htmlScreen += "</div>"

	return template.HTML(htmlScreen), nil
}

func symbolToClass(symbol rune) (string, error) {
	switch symbol {
	case rune(' '):
		return "space", nil
	case domain.Border.Rune():
		return "wall", nil
	case domain.Snake.Rune():
		return "snake", nil
	case domain.Food.Rune():
		return "food", nil
	default:
		return "", fmt.Errorf("%w: %c", errUnknownSymbol, symbol)
	}
}
