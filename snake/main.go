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

	startBody := []domain.Coordinate{
		{
			Row: 15,
			Col: 15,
		},
		{
			Row: 15,
			Col: 16,
		},
		{
			Row: 15,
			Col: 17,
		},
	}

	game, err := gamer.NewGame(dir, startBody, rows, cols, snakeSprite, foodSprite, &borderSprite)
	if err != nil {
		panic(err)
	}

	game.Run()
}
