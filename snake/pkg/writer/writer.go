package writer

import (
	"fmt"
	"strings"
)

type Writer interface {
	WriteString(string) (int, error)
	Flush() error
}

func DisplayScreen(data []string, writer Writer) {
	if _, err := writer.WriteString(fmt.Sprintf("\x1b[%dF", len(data))); err != nil {
		panic(err) // TODO: yeah. Fix.
	}

	for _, col := range data {
		if strings.Contains(col, "\n") {
			panic("Newline in column") // TODO: ok, have to stop putting panic everywhere at some point
		}

		if _, err := writer.WriteString(col + "\x1b[K\n"); err != nil {
			panic(err) // TODO: no comments.
		}
	}

	if err := writer.Flush(); err != nil {
		panic(err) // TODO todo
	}
}
