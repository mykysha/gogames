package gamer

import (
	"bufio"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"slices"
	"time"

	"github.com/mykysha/gogames/snake/domain"
	"github.com/mykysha/gogames/snake/pkg/displayer"
	"github.com/mykysha/gogames/snake/pkg/interceptor"
	"github.com/mykysha/gogames/snake/pkg/log"
	"github.com/mykysha/gogames/snake/pkg/snaker"
	"github.com/mykysha/gogames/snake/pkg/window"
)

const thirtyFPS = 30

type Window interface {
	Set(data byte, row, col int) error
	WriteText(text string, row, col int) error
	Clean()
	GetSnapshot() []string
}

type Snake interface {
	SetDirection(dir snaker.Direction) error
	GetLocation() []domain.Coordinate
	MakeBigger()
	IncreaseSpeed()
	Live()
}

type Displayer interface {
	DisplayScreen(data []string) error
	InitialDisplay(data []string) error
}

type Game struct {
	logger     log.Logger
	framerate  int
	updateFlag chan struct{}

	rows        int
	cols        int
	snakeSprite byte
	foodSprite  byte
	screen      Window
	displayer   Displayer

	score int
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

	updateFlag := make(chan struct{})

	return &Game{
		logger:      slog.Default(),
		framerate:   thirtyFPS,
		updateFlag:  updateFlag,
		rows:        rows,
		cols:        cols,
		snakeSprite: snakeSprite,
		foodSprite:  foodSprite,
		screen:      screen,
		displayer:   display,
		score:       0,
		snake:       snaker.NewSnake(updateFlag, startDir, startBody, rows, cols),
	}, nil
}

func (g *Game) Run() {
	food := domain.Coordinate{Row: 10, Col: 10}

	go g.snake.Live()
	go g.handleMovements()
	go g.drawFrames()

	for {
		g.screen.Clean()

		snakeLocation := g.snake.GetLocation()

		for ind, coordinate := range snakeLocation {
			if slices.Contains(snakeLocation[ind+1:], coordinate) {
				if err := g.displayGameOver(); err != nil {
					g.logger.Error("failed to display game over", "error", err)
				}

				return
			}

			if coordinate == food {
				g.snake.MakeBigger()
				g.snake.IncreaseSpeed()
				g.score++
				food = generateNewFood(g.rows, g.cols, snakeLocation)
			}

			g.screen.Set(g.snakeSprite, coordinate.Row, coordinate.Col)
		}

		g.screen.Set(g.foodSprite, food.Row, food.Col)

		<-g.updateFlag
	}
}

func generateNewFood(rows, cols int, snake []domain.Coordinate) domain.Coordinate {
	for {
		food := domain.Coordinate{Row: rand.IntN(rows), Col: rand.IntN(cols)}

		if !slices.Contains(snake, food) {
			return food
		}
	}
}

func (g *Game) handleMovements() {
	keys := make(chan rune)

	go func() {
		for {
			if err := interceptor.InterceptKeystrokes(keys); err != nil {
				g.logger.Error("failed to intercept keystrokes", "error", err)
			}
		}
	}()

	for {
		key := <-keys

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
}

func (g *Game) drawFrames() {
	for {
		time.Sleep(time.Second / time.Duration(g.framerate))

		if err := g.displayer.DisplayScreen(g.screen.GetSnapshot()); err != nil {
			g.logger.Error("failed to display screen", "error", err)
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
