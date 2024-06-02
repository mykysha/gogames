package main

import (
	"log/slog"
	"net/http"

	"github.com/mykysha/gogames/snake/api"
	"github.com/mykysha/gogames/snake/domain"
	"github.com/mykysha/gogames/snake/pkg/gamer"
	"github.com/mykysha/gogames/snake/pkg/log"
	"github.com/mykysha/gogames/snake/pkg/snaker"
)

func main() {
	logger := slog.Default()

	screenChan := make(chan string)

	go setupServer(logger, screenChan)
	setupGame(logger, screenChan)
}

func setupServer(logger log.Logger, screenChan chan string) {
	handlers := api.NewAPI(logger, screenChan)

	logger.Info("Server started at :8080")

	if err := http.ListenAndServe(":8080", handlers); err != nil {
		panic(err)
	}
}

func setupGame(logger log.Logger, screenChan chan string) {
	rows := 20
	cols := 20

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

	game, err := gamer.NewGame(logger, screenChan, dir, startBody, rows, cols)
	if err != nil {
		panic(err)
	}

	game.Run()
}
