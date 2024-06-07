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

	screenChan := make(chan [][]domain.Cell)
	keyChan := make(chan string)

	go setupGame(logger, screenChan, keyChan)
	setupServer(logger, screenChan, keyChan)
}

func setupServer(logger log.Logger, screenChan chan [][]domain.Cell, keyChan chan string) {
	handlers := api.NewAPI(logger, screenChan, keyChan)

	logger.Info("Server started at :8080")

	if err := http.ListenAndServe(":8080", handlers); err != nil {
		panic(err)
	}
}

func setupGame(logger log.Logger, screenChan chan [][]domain.Cell, keyChan chan string) {
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

	game := gamer.NewGame(logger, screenChan, keyChan, dir, startBody, rows, cols)

	game.Run()
}
