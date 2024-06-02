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
			color, err := symbolToColor(char)
			if err != nil {
				return "", fmt.Errorf("failed to convert screen to html: %w", err)
			}

			htmlScreen += "<div style=\"width: 1rem; height: 1rem; background-color:" + color + "\">" + "</div>"
		}

		htmlScreen += "</div>"
	}

	htmlScreen += "</div>"

	return template.HTML(htmlScreen), nil
}

func symbolToColor(symbol rune) (string, error) {
	switch symbol {
	case rune(' '):
		return "white", nil
	case domain.Border.Rune():
		return "black", nil
	case domain.Snake.Rune():
		return "green", nil
	case domain.Food.Rune():
		return "red", nil
	default:
		return "", fmt.Errorf("%w: %c", errUnknownSymbol, symbol)
	}
}
