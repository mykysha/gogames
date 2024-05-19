package displayer

import (
	"fmt"
	"strings"
)

type OutputWriter interface {
	WriteString(s string) (int, error)
	Flush() error
}

type Displayer struct {
	writer OutputWriter
}

func NewDisplayer(writer OutputWriter) *Displayer {
	return &Displayer{
		writer: writer,
	}
}

func (d Displayer) InitialDisplay(data []string) error {
	for _, row := range data {
		if strings.Contains(row, "\n") {
			return errNewLineInRow
		}

		if _, err := d.writer.WriteString(row + "\x1b[K\n"); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	if err := d.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush: %w", err)
	}

	return nil
}

func (d Displayer) DisplayScreen(data []string) error {
	if _, err := d.writer.WriteString(fmt.Sprintf("\x1b[%dF", len(data))); err != nil {
		return fmt.Errorf("failed to clear screen: %w", err)
	}

	if err := d.InitialDisplay(data); err != nil {
		return fmt.Errorf("failed to display screen: %w", err)
	}

	return nil
}
