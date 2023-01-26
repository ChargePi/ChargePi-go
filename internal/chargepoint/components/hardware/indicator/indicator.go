package indicator

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
)

// color constants
const (
	Off    = Color("Off")
	White  = Color("White")
	Red    = Color("Red")
	Green  = Color("Green")
	Blue   = Color("Blue")
	Yellow = Color("Yellow")
	Orange = Color("Orange")
)

// Supported types
const (
	TypeWS281x = "WS281x"
)

var (
	ErrInvalidIndex        = errors.New("invalid index")
	ErrInvalidPin          = errors.New("invalid data pin number")
	ErrInvalidNumberOfLeds = errors.New("number of leds must be greater than zero")
)

type (
	Color string

	// Indicator is an abstraction layer for connector status indication, usually an RGB LED strip.
	Indicator interface {
		DisplayColor(index int, color Color) error
		Blink(index int, times int, color Color) error
		Cleanup()
		GetType() string
	}
)

// NewIndicator constructs the Indicator based on the type provided by the settings file.
func NewIndicator(stripLength int, indicator settings.LedIndicator) Indicator {
	if indicator.Enabled {
		if indicator.IndicateCardRead {
			stripLength++
		}

		log.Infof("Preparing Indicator from config: %s", indicator.Type)
		switch indicator.Type {
		case TypeWS281x:
			ledStrip, ledError := NewWS281xStrip(stripLength, *indicator.DataPin)
			if ledError != nil {
				log.WithError(ledError).Errorf("Error creating indicator")
				return nil
			}

			return ledStrip
		default:
			return nil
		}
	}

	return nil
}
