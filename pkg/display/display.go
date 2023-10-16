package display

import (
	"errors"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/settings"
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
func NewDisplay(lcdSettings settings.Display) (Display, error) {
	if lcdSettings.IsEnabled {
		log.Info("Preparing display from config")

		switch lcdSettings.Driver {
		case DriverHD44780:
			if util.IsNilInterfaceOrPointer(lcdSettings.HD44780) {
				return nil, ErrInvalidConnectionDetails
			}

			lcd, err := NewHD44780(lcdSettings.HD44780.I2C.Address, lcdSettings.HD44780.I2C.Bus)
			if err != nil {
				return nil, err
			}

			return lcd, nil
		case TypeDummy:
			return NewDummy(lcdSettings.DisplayDummy)
		default:
			return nil, ErrDisplayUnsupported
		}
	}

	return nil, ErrDisplayDisabled
}
