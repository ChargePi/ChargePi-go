package display

import (
	"errors"
	"go.uber.org/zap"

	"github.com/ChargePi/ChargePi-go/internal/pkg/util"
	"github.com/ChargePi/ChargePi-go/pkg/models/settings"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
)

const (
	DriverHD44780 = "hd44780"
	TypeDummy     = "dummy"
)

var (
	ErrDisplayUnsupported       = errors.New("display type unsupported")
	ErrInvalidConnectionDetails = errors.New("connection details invalid or empty")
	ErrDisplayDisabled          = errors.New("display disabled")
)

// Display is an abstraction layer for concrete implementation of a display.
type Display interface {
	DisplayMessage(message display.MessageInfo)
	//	GetCurrentMessage(display.MessageInfo) (display.MessageStatus, error)
	// GetMessages(reqId int) ([]display.GetDisplayMessagesResponse, error)
	Clear()
	Cleanup()
	GetType() string
}

// NewDisplay returns a concrete implementation of a Display based on the drivers that are supported.
// The Display is built with the settings from the settings file.
func NewDisplay(displaySettings settings.Display) (Display, error) {
	logger := zap.L()
	if displaySettings.IsEnabled {
		logger.With(zap.String("display_type", displaySettings.Driver)).Info("Preparing display from config")

		switch displaySettings.Driver {
		case DriverHD44780:
			if util.IsNilInterfaceOrPointer(displaySettings.HD44780) {
				return nil, ErrInvalidConnectionDetails
			}

			lcd, err := NewHD44780(logger, *displaySettings.HD44780)
			if err != nil {
				return nil, err
			}

			return lcd, nil
		case TypeDummy:
			return NewDummy(logger, displaySettings.DisplayDummy)
		default:
			return nil, ErrDisplayUnsupported
		}
	}

	return nil, ErrDisplayDisabled
}
