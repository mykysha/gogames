package window

import (
	"fmt"
	"strings"

	"github.com/mykysha/gogames/snake/domain"
)

type Window struct {
	empty byte
	rows  int
	cols  int
	data  []byte
}

func New(rows, cols int) *Window {
	empty := []byte(" ")[0]
	data := make([]byte, 0, rows*cols)

	for range rows * cols {
		data = append(data, empty)
	}

	return &Window{
		empty: empty,
		rows:  rows,
		cols:  cols,
		data:  data,
	}
}

func (w *Window) Set(data byte, row, col int) error {
	if row >= w.rows || row < 0 {
		return errInvalidRow
	}

	if col >= w.cols || col < 0 {
		return errInvalidCol
	}

	w.data[row*w.cols+col] = data

	return nil
}

func (w *Window) WriteText(text string, row, col int) error {
	for _, char := range text {
		col++
		if col >= w.cols {
			col = 0
			row++

			if row >= w.rows {
				return errTextOutOfScreen
			}
		}

		if err := w.Set(byte(char), row, col); err != nil {
			return fmt.Errorf("failed to write symbol: %w", err)
		}
	}

	return nil
}

func (w *Window) Remove(row, col int) error {
	return w.Set(w.empty, row, col)
}

func (w *Window) Clean() {
	for i := range w.data {
		w.data[i] = w.empty
	}
}

func (w *Window) GetSnapshot() []string {
	curScreenHeight := w.rows + 2

	currentScreen := make([]string, 0, curScreenHeight)

	currentScreen = append(currentScreen, strings.Repeat(string(domain.Border), w.cols+2))

	for i := range w.rows {
		row := w.data[i*w.cols : (i+1)*w.cols]

		row = append([]byte{byte(domain.Border)}, row...)
		row = append(row, byte(domain.Border))

		currentScreen = append(currentScreen, string(row))
	}

	currentScreen = append(currentScreen, strings.Repeat(string(domain.Border), w.cols+2))

	return currentScreen
}
