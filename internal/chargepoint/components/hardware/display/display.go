package display

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/models"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
)

const (
	DriverHD44780 = "hd44780"
)

var (
	ErrDisplayUnsupported = errors.New("display type unsupported")
	ErrDisplayDisabled    = errors.New("display disabled")
)

type (
	// Display is an abstraction layer for concrete implementation of a display.
	Display interface {
		DisplayMessage(message models.Message)
		Cleanup()
		Clear()
	}
)

// NewDisplay returns a concrete implementation of an Display based on the drivers that are supported.
// The Display is built with the settings from the settings file.
func NewDisplay(lcdSettings settings.Display) (Display, error) {
	if lcdSettings.IsEnabled {
		log.Info("Preparing display from config")

		switch lcdSettings.Driver {
		case DriverHD44780:
			lcd, err := NewHD44780(lcdSettings.I2CAddress, lcdSettings.I2CBus)
			if err != nil {
				return nil, err
			}

			return lcd, nil
		default:
			return nil, ErrDisplayUnsupported
		}
	}

	return nil, ErrDisplayDisabled
}
