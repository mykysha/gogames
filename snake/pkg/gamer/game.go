package gamer

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"

	"github.com/mykysha/gogames/snake/domain"
	"github.com/mykysha/gogames/snake/pkg/log"
	"github.com/mykysha/gogames/snake/pkg/snaker"
	"github.com/mykysha/gogames/snake/pkg/window"
)

type Game struct {
	logger     log.Logger
	screenChan chan string

	rows   int
	cols   int
	screen Window
	keys   chan string

	score int
	food  domain.Coordinate
	snake Snake
}

func NewGame(
	logger log.Logger,
	screenChan, keyChan chan string,
	startDir snaker.Direction,
	startBody []domain.Coordinate,
	rows, cols int,
) *Game {
	screen := window.New(rows, cols)

	screenChan <- strings.Join(screen.GetSnapshot(), "\n")

	return &Game{
		logger:     logger,
		screenChan: screenChan,
		rows:       rows,
		cols:       cols,
		screen:     screen,
		score:      0,
		keys:       keyChan,
		food:       generateNewFood(rows, cols, startBody),
		snake:      snaker.NewSnake(startDir, startBody, rows, cols),
	}
}

func (g *Game) Run() {
	for {
		gameOver := g.gameCycle()
		if gameOver {
			break
		}
	}
}

func (g *Game) gameCycle() bool {
	g.handleMovement()
	newSnakeLocation := g.snake.Move()

	if newSnakeLocation != nil {
		gameOver := g.handleSnakeMovement(newSnakeLocation)
		if gameOver {
			if err := g.displayGameOver(); err != nil {
				g.logger.Error("failed to display game over", "error", err)
			}

			return true
		}

		snapshot := g.screen.GetSnapshot()
		singularScreen := strings.Join(snapshot, "\n")
		g.screenChan <- singularScreen
	}

	return false
}

func (g *Game) handleSnakeMovement(newSnakeLocation []domain.Coordinate) bool {
	g.screen.Clean()

	for ind, coordinate := range newSnakeLocation {
		if slices.Contains(newSnakeLocation[ind+1:], coordinate) {
			return true
		}

		if coordinate == g.food {
			g.handleEatenFood(newSnakeLocation)
		}

		if err := g.screen.Set(byte(domain.Snake), coordinate.Row, coordinate.Col); err != nil {
			g.logger.Error("failed to set snake sprite", "error", err)
		}
	}

	if err := g.screen.Set(byte(domain.Food), g.food.Row, g.food.Col); err != nil {
		g.logger.Error("failed to set food sprite", "error", err)
	}

	return false
}

func (g *Game) handleEatenFood(newSnakeLocation []domain.Coordinate) {
	g.snake.MakeBigger()
	g.snake.IncreaseSpeed()

	g.score++

	g.food = generateNewFood(g.rows, g.cols, newSnakeLocation)
}

func generateNewFood(rows, cols int, snake []domain.Coordinate) domain.Coordinate {
	for {
		food := domain.Coordinate{Row: rand.IntN(rows), Col: rand.IntN(cols)}

		if !slices.Contains(snake, food) {
			return food
		}
	}
}

func (g *Game) handleMovement() {
	select {
	case key := <-g.keys:
		g.setSnakeDirection(key)
	default:
	}
}

func (g *Game) setSnakeDirection(key string) {
	switch key {
	case "up":
		if err := g.snake.SetDirection(snaker.DirectionUp); err != nil {
			g.logger.Error("failed to set direction up", "error", err)
		}
	case "right":
		if err := g.snake.SetDirection(snaker.DirectionRight); err != nil {
			g.logger.Error("failed to set direction right", "error", err)
		}
	case "down":
		if err := g.snake.SetDirection(snaker.DirectionDown); err != nil {
			g.logger.Error("failed to set direction down", "error", err)
		}
	case "left":
		if err := g.snake.SetDirection(snaker.DirectionLeft); err != nil {
			g.logger.Error("failed to set direction left", "error", err)
		}
	}
}

func (g *Game) displayGameOver() error {
	g.screen.Clean()

	middleRow := g.rows / 2
	if err := g.screen.WriteText("Game Over", middleRow-1, 1); err != nil {
		return fmt.Errorf("failed to write game over text: %w", err)
	}

	if err := g.screen.WriteText(fmt.Sprintf("Score: %d", g.score), middleRow+1, 1); err != nil {
		return fmt.Errorf("failed to write score: %w", err)
	}
	// TODO: Display the score in htmx somehow.
	return nil
}
