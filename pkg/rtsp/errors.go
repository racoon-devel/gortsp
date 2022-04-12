package rtsp

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidURL = errors.New("URL must be rtsp://host:port/path")
)

type ErrInvalidMethod struct {
	Method Method
}

func (e ErrInvalidMethod) Error() string {
	return fmt.Sprintf("invalid method: %s", e.Method)
}
