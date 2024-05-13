package snaker

import (
	"fmt"
	"sync"
	"time"

	"github.com/mykysha/gogames/snake/domain"
)

type (
	Direction int
)

const (
	DirectionUp Direction = iota
	DirectionRight
	DirectionDown
	DirectionLeft
)

type Snake struct {
	updateFlag   chan struct{}
	updateTime   time.Duration
	nextDir      *Direction
	dir          Direction
	body         []domain.Coordinate
	becameBigger bool
	rows         int
	cols         int
	rwMux        *sync.RWMutex
}

func NewSnake(
	updateFlag chan struct{},
	startDir Direction,
	startBody []domain.Coordinate,
	rows, cols int,
) *Snake {
	return &Snake{
		updateFlag:   updateFlag,
		updateTime:   time.Second / 5,
		nextDir:      nil,
		dir:          startDir,
		body:         startBody,
		becameBigger: false,
		rows:         rows,
		cols:         cols,
		rwMux:        new(sync.RWMutex),
	}
}

func (s *Snake) SetDirection(dir Direction) error {
	reverseDir, err := getReverseDir(s.dir)
	if err != nil {
		return fmt.Errorf("failed to get reverse direction: %w", err)
	}

	if dir == reverseDir {
		return nil
	}

	s.nextDir = &dir

	return nil
}

func (s *Snake) IncreaseSpeed() {
	s.updateTime = s.updateTime - time.Millisecond*10
}

func getReverseDir(dir Direction) (Direction, error) {
	switch dir {
	case DirectionUp:
		return DirectionDown, nil
	case DirectionRight:
		return DirectionLeft, nil
	case DirectionDown:
		return DirectionUp, nil
	case DirectionLeft:
		return DirectionRight, nil
	}

	return 0, errUnknownDirection
}

func (s *Snake) MakeBigger() {
	s.becameBigger = true
}

func (s *Snake) GetLocation() []domain.Coordinate {
	s.rwMux.RLock()
	defer s.rwMux.RUnlock()

	return s.body
}

func (s *Snake) Live() {
	for {
		time.Sleep(s.updateTime)

		if s.nextDir != nil {
			s.dir = *s.nextDir
			s.nextDir = nil
		}

		newHead := move(s.body[len(s.body)-1], s.dir, s.rows, s.cols)

		s.rwMux.Lock()

		s.body = append(s.body, newHead)

		if !s.becameBigger {
			s.body = s.body[1:]
		} else {
			s.becameBigger = false
		}

		s.rwMux.Unlock()

		s.updateFlag <- struct{}{}
	}
}

func move(cur domain.Coordinate, dir Direction, rows, cols int) domain.Coordinate {
	newRow := cur.Row
	newCol := cur.Col

	switch dir {
	case DirectionUp:
		newRow = updateCoordinate(newRow, -1, rows)
	case DirectionRight:
		newCol = updateCoordinate(newCol, 1, cols)
	case DirectionDown:
		newRow = updateCoordinate(newRow, 1, rows)
	case DirectionLeft:
		newCol = updateCoordinate(newCol, -1, cols)
	}

	return domain.Coordinate{
		Row: newRow,
		Col: newCol,
	}
}

func updateCoordinate(current, addendant, max int) int {
	newCoordinate := current + addendant

	if newCoordinate < 0 {
		return max - 1
	}

	if newCoordinate == max {
		return 0
	}

	return newCoordinate
}
