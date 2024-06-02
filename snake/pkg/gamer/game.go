package gamer

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"strings"

	"github.com/mykysha/gogames/snake/domain"
	"github.com/mykysha/gogames/snake/pkg/displayer"
	"github.com/mykysha/gogames/snake/pkg/interceptor"
	"github.com/mykysha/gogames/snake/pkg/log"
	"github.com/mykysha/gogames/snake/pkg/snaker"
	"github.com/mykysha/gogames/snake/pkg/window"
)

type Game struct {
	logger     log.Logger
	screenChan chan string

	rows      int
	cols      int
	screen    Window
	displayer Displayer
	keys      chan byte

	score int
	food  domain.Coordinate
	snake Snake
}

func NewGame(
	logger log.Logger,
	screenChan chan string,
	startDir snaker.Direction,
	startBody []domain.Coordinate,
	rows, cols int,
) (*Game, error) {
	screen := window.New(rows, cols)

	display := displayer.NewDisplayer(bufio.NewWriter(os.Stdout))

	if err := display.InitialDisplay(screen.GetSnapshot()); err != nil {
		return nil, fmt.Errorf("failed to display initial screen: %w", err)
	}

	return &Game{
		logger:     logger,
		screenChan: screenChan,
		rows:       rows,
		cols:       cols,
		screen:     screen,
		displayer:  display,
		score:      0,
		keys:       nil,
		food:       generateNewFood(rows, cols, startBody),
		snake:      snaker.NewSnake(startDir, startBody, rows, cols),
	}, nil
}

func (g *Game) Run() {
	g.keys = g.startInterceptingKeystrokes()

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

		if err := g.displayer.DisplayScreen(snapshot); err != nil {
			g.logger.Error("failed to display screen", "error", err)
		}
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

	middleRow := g.rows / 2
	if err := g.screen.WriteText("Game Over", middleRow-1, 1); err != nil {
		return fmt.Errorf("failed to write game over text: %w", err)
	}

	if err := g.screen.WriteText(fmt.Sprintf("Score: %d", g.score), middleRow+1, 1); err != nil {
		return fmt.Errorf("failed to write score: %w", err)
	}

	if err := g.displayer.DisplayScreen(g.screen.GetSnapshot()); err != nil {
		return fmt.Errorf("failed to display game over screen: %w", err)
	}

	return nil
}
