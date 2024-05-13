package window

import (
	"fmt"
	"strings"
)

type Window struct {
	border *byte
	empty  byte
	rows   int
	cols   int
	data   []byte
}

func New(rows, cols int) *Window {
	empty := []byte(" ")[0]
	data := make([]byte, 0, rows*cols)

	for i := 0; i < rows*cols; i++ {
		data = append(data, empty)
	}

	return &Window{
		border: nil,
		empty:  empty,
		rows:   rows,
		cols:   cols,
		data:   data,
	}
}

func NewWithBorder(rows, cols int, border byte) *Window {
	window := New(rows, cols)
	window.border = &border

	return window
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

	if w.border != nil {
		currentScreen = append(currentScreen, strings.Repeat(string(*w.border), w.cols+2))
	}

	for i := range w.rows {
		row := w.data[i*w.cols : (i+1)*w.cols]

		if w.border != nil {
			row = append([]byte{*w.border}, row...)
			row = append(row, *w.border)
		}

		currentScreen = append(currentScreen, string(row))
	}

	if w.border != nil {
		currentScreen = append(currentScreen, strings.Repeat(string(*w.border), w.cols+2))
	}

	return currentScreen
}
