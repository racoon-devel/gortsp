package rtsp

import (
	"errors"
)

var (
	ErrInvalidURL      = errors.New("URL must be rtsp://host:port/path")
	ErrMethodMustBeSet = errors.New("method must be set")
)
