package main

import (
	"math/rand/v2"
	"slices"

	"github.com/mykysha/gogames/snake/domain"
	"github.com/mykysha/gogames/snake/pkg/interceptor"
	"github.com/mykysha/gogames/snake/pkg/snake"
	"github.com/mykysha/gogames/snake/pkg/window"
)

func main() {
	symbol := []byte("*")[0]

	rows := 20
	cols := 20

	gameScreen := window.NewWithBorder(rows, cols, []byte("#")[0])

	updateFlag := make(chan struct{})
	dir := snake.DirectionRight
	startBody := []domain.Coordinate{{Row: 15, Col: 15}, {Row: 15, Col: 16}, {Row: 15, Col: 17}}

	snakeInstance := snake.NewSnake(updateFlag, dir, startBody, rows, cols)

	go func() {
		food := domain.Coordinate{Row: 10, Col: 10}

		for {
			gameScreen.Clean()

			snakeLocation := snakeInstance.GetLocation()

			for ind, coordinate := range snakeLocation {
				if slices.Contains(snakeLocation[ind+1:], coordinate) {
					gameScreen.Clean()
					gameScreen.WriteText("Game Over", 5, 5)
					gameScreen.Display()

					return
				}

				if coordinate == food {
					snakeInstance.MakeBigger()
					food = generateNewFood(rows, cols, snakeLocation)
				}

				gameScreen.Set(symbol, coordinate.Row, coordinate.Col)
			}

			gameScreen.Set([]byte("@")[0], food.Row, food.Col)

			<-updateFlag
		}
	}()

	go snakeInstance.Live()

	go func() {
		keys := make(chan rune)

		go interceptor.InterceptKeystrokes(keys)

		for {
			key := <-keys

			switch key {
			case 'w':
				snakeInstance.SetDirection(snake.DirectionUp)
			case 'd':
				snakeInstance.SetDirection(snake.DirectionRight)
			case 's':
				snakeInstance.SetDirection(snake.DirectionDown)
			case 'a':
				snakeInstance.SetDirection(snake.DirectionLeft)
			}
		}
	}()

	gameScreen.Display()
}

func generateNewFood(rows, cols int, snake []domain.Coordinate) domain.Coordinate {
	for {
		food := domain.Coordinate{Row: rand.IntN(rows), Col: rand.IntN(cols)}

		if !slices.Contains(snake, food) {
			return food
		}
	}
}
