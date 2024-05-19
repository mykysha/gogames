package main

import (
	"github.com/mykysha/gogames/snake/domain"
	"github.com/mykysha/gogames/snake/pkg/gamer"
	"github.com/mykysha/gogames/snake/pkg/snaker"
)

func main() {
	rows := 20
	cols := 20

	borderSprite := []byte("#")[0]
	snakeSprite := []byte("*")[0]
	foodSprite := []byte("@")[0]

	dir := snaker.DirectionRight

	snakeRow := 15
	snakeColTail := 15
	snakeColMiddle := 16
	snakeColHead := 17

	startBody := []domain.Coordinate{
		{
			Row: snakeRow,
			Col: snakeColTail,
		},
		{
			Row: snakeRow,
			Col: snakeColMiddle,
		},
		{
			Row: snakeRow,
			Col: snakeColHead,
		},
	}

	game, err := gamer.NewGame(dir, startBody, rows, cols, snakeSprite, foodSprite, &borderSprite)
	if err != nil {
		panic(err)
	}

	game.Run()
}
