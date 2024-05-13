package window

import "errors"

var (
	errInvalidRow      = errors.New("invalid row")
	errInvalidCol      = errors.New("invalid col")
	errTextOutOfScreen = errors.New("text out of screen")
)
