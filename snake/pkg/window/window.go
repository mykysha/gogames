package window

import (
	"fmt"

	"github.com/mykysha/gogames/snake/domain"
)

type Window struct {
	empty domain.Cell
	rows  int
	cols  int
	data  []domain.Cell
}

func New(rows, cols int) *Window {
	empty := domain.Cell{
		BgColor: domain.White,
	}

	data := make([]domain.Cell, 0, rows*cols)

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

func (w *Window) Set(data domain.Cell, row, col int) error {
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

		byteSymbol := byte(char)

		if err := w.Set(domain.Cell{
			BgColor:   domain.White,
			TextColor: domain.Black,
			Symbol:    &byteSymbol,
		}, row, col); err != nil {
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

func (w *Window) GetSnapshot() [][]domain.Cell {
	curScreenHeight := w.rows + 2

	currentScreen := make([][]domain.Cell, 0, curScreenHeight)

	currentScreen = append(
		currentScreen,
		repeatCell(
			domain.Cell{
				BgColor: domain.Black,
			},
			w.cols+2,
		),
	)

	for i := range w.rows {
		row := w.data[i*w.cols : (i+1)*w.cols]

		row = append(
			[]domain.Cell{{
				BgColor: domain.Black,
			}},
			row...)
		row = append(
			row,
			domain.Cell{
				BgColor: domain.Black,
			},
		)

		currentScreen = append(currentScreen, row)
	}

	currentScreen = append(
		currentScreen,
		repeatCell(
			domain.Cell{
				BgColor: domain.Black,
			},
			w.cols+2,
		),
	)

	return currentScreen
}

func repeatCell(cell domain.Cell, times int) []domain.Cell {
	result := make([]domain.Cell, 0, times)

	for range times {
		result = append(result, cell)
	}

	return result
}
