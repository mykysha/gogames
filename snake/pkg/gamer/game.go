package gamer

import (
	"bufio"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"slices"

	"github.com/mykysha/gogames/snake/domain"
	"github.com/mykysha/gogames/snake/pkg/displayer"
	"github.com/mykysha/gogames/snake/pkg/interceptor"
	"github.com/mykysha/gogames/snake/pkg/log"
	"github.com/mykysha/gogames/snake/pkg/snaker"
	"github.com/mykysha/gogames/snake/pkg/window"
)

type Window interface {
	Set(data byte, row, col int) error
	WriteText(text string, row, col int) error
	Clean()
	GetSnapshot() []string
}

type Snake interface {
	SetDirection(dir snaker.Direction) error
	MakeBigger()
	IncreaseSpeed()
	Move() []domain.Coordinate
}

type Displayer interface {
	DisplayScreen(data []string) error
	InitialDisplay(data []string) error
}

type Game struct {
	logger log.Logger

	rows        int
	cols        int
	snakeSprite byte
	foodSprite  byte
	screen      Window
	displayer   Displayer
	keys        chan byte

	score int
	food  domain.Coordinate
	snake Snake
}

func NewGame(
	startDir snaker.Direction,
	startBody []domain.Coordinate,
	rows, cols int,
	snakeSprite, foodSprite byte,
	borderSprite *byte,
) (*Game, error) {
	var screen *window.Window

	if borderSprite != nil {
		screen = window.NewWithBorder(rows, cols, *borderSprite)
	} else {
		screen = window.New(rows, cols)
	}

	display := displayer.NewDisplayer(bufio.NewWriter(os.Stdout))

	if err := display.InitialDisplay(screen.GetSnapshot()); err != nil {
		return nil, fmt.Errorf("failed to display initial screen: %w", err)
	}

	return &Game{
		logger:      slog.Default(),
		rows:        rows,
		cols:        cols,
		snakeSprite: snakeSprite,
		foodSprite:  foodSprite,
		screen:      screen,
		displayer:   display,
		score:       0,
		keys:        nil,
		food:        domain.Coordinate{Row: 10, Col: 10},
		snake:       snaker.NewSnake(startDir, startBody, rows, cols),
	}, nil
}

func (g *Game) Run() {
	g.keys = g.startInterceptingKeystrokes()

	for {
		gameOver := g.gameCycle()
		if gameOver {
			return
		}
	}
}

func (g *Game) gameCycle() bool {
	g.handleMovement()
	newSnakeLocation := g.snake.Move()

	if newSnakeLocation != nil {
		g.screen.Clean()

		for ind, coordinate := range newSnakeLocation {
			if slices.Contains(newSnakeLocation[ind+1:], coordinate) {
				if err := g.displayGameOver(); err != nil {
					g.logger.Error("failed to display game over", "error", err)
				}

				return true
			}

			if coordinate == g.food {
				g.snake.MakeBigger()
				g.snake.IncreaseSpeed()
				g.score++
				g.food = generateNewFood(g.rows, g.cols, newSnakeLocation)
			}

			g.screen.Set(g.snakeSprite, coordinate.Row, coordinate.Col)
		}

		g.screen.Set(g.foodSprite, g.food.Row, g.food.Col)

		if err := g.displayer.DisplayScreen(g.screen.GetSnapshot()); err != nil {
			g.logger.Error("failed to display screen", "error", err)
		}
	}

	return false
}

func generateNewFood(rows, cols int, snake []domain.Coordinate) domain.Coordinate {
	for {
		food := domain.Coordinate{Row: rand.IntN(rows), Col: rand.IntN(cols)}

		if !slices.Contains(snake, food) {
			return food
		}
	}
}

func (g *Game) startInterceptingKeystrokes() chan byte {
	keys := make(chan byte)

	go func() {
		for {
			if err := interceptor.InterceptKeystrokes(keys); err != nil {
				g.logger.Error("failed to intercept keystrokes", "error", err)
			}
		}
	}()

	return keys
}

func (g *Game) handleMovement() {
	select {
	case key := <-g.keys:
		g.setSnakeDirection(key)
	default:
	}
}

func (g *Game) setSnakeDirection(key byte) {
	switch key {
	case 'w':
		if err := g.snake.SetDirection(snaker.DirectionUp); err != nil {
			g.logger.Error("failed to set direction up", "error", err)
		}
	case 'd':
		if err := g.snake.SetDirection(snaker.DirectionRight); err != nil {
			g.logger.Error("failed to set direction right", "error", err)
		}
	case 's':
		if err := g.snake.SetDirection(snaker.DirectionDown); err != nil {
			g.logger.Error("failed to set direction down", "error", err)
		}
	case 'a':
		if err := g.snake.SetDirection(snaker.DirectionLeft); err != nil {
			g.logger.Error("failed to set direction left", "error", err)
		}
	}
}

func (g *Game) displayGameOver() error {
	g.screen.Clean()
	if err := g.screen.WriteText("Game Over", 5, 5); err != nil {
		return fmt.Errorf("failed to write game over text: %w", err)
	}

	if err := g.screen.WriteText(fmt.Sprintf("Score: %d", g.score), 7, 5); err != nil {
		return fmt.Errorf("failed to write score: %w", err)
	}

	if err := g.displayer.DisplayScreen(g.screen.GetSnapshot()); err != nil {
		return fmt.Errorf("failed to display game over screen: %w", err)
	}

	return nil
}
