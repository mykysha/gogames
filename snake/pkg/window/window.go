package window

import (
	"bufio"
	"os"
	"strings"
	"time"
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
	for i, char := range text {
		if err := w.Set(byte(char), row, col+i); err != nil { // TODO: redo
			return err
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

func (w *Window) Display() {
	writer := bufio.NewWriter(os.Stdout) // TODO: well rework.

	for {
		time.Sleep(time.Second / 30)

		if w.border != nil {
			if _, err := writer.WriteString(strings.Repeat(string(*w.border), w.cols+2)); err != nil {
				panic(err) // TODO: replace with log.
			}
		}

		for i := range w.rows {
			row := w.data[i*w.cols : (i+1)*w.cols]

			if w.border != nil {
				row = append([]byte{*w.border}, row...)
				row = append(row, *w.border)
			}

			if _, err := writer.WriteString(string(row)); err != nil {
				panic(err) // TODO: replace with log.
			}
		}

		if w.border != nil {
			if _, err := writer.WriteString(strings.Repeat(string(*w.border), w.cols+2)); err != nil {
				panic(err) // TODO: replace with log.
			}
		}

		if err := writer.Flush(); err != nil {
			panic(err) // TODO: replace with log.
		}
	}
}
