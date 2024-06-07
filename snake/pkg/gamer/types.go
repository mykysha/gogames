package gamer

import (
	"github.com/mykysha/gogames/snake/domain"
	"github.com/mykysha/gogames/snake/pkg/snaker"
)

type Window interface {
	Set(data domain.Cell, row, col int) error
	WriteText(text string, row, col int) error
	Clean()
	GetSnapshot() [][]domain.Cell
}

type Snake interface {
	SetDirection(dir snaker.Direction) error
	MakeBigger()
	IncreaseSpeed()
	Move() []domain.Coordinate
}
