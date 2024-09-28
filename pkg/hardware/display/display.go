package display

import (
	"context"
	"errors"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
)

var (
	ErrDisplayUnsupported       = errors.New("display type unsupported")
	ErrInvalidConnectionDetails = errors.New("connection details invalid or empty")
	ErrDisplayDisabled          = errors.New("display disabled")
)

// Display is an abstraction layer for concrete implementation of a display.
type Display interface {
	DisplayMessage(message display.MessageInfo) error
	Clear() error
	Cleanup(ctx context.Context) error
	GetType() string
}

type Message struct {
	Content  string
	Language string
}
