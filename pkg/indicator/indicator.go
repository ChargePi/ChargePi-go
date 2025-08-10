package indicator

import (
	"errors"
	"go.uber.org/zap"

	"github.com/ChargePi/ChargePi-go/internal/pkg/util"
	"github.com/ChargePi/ChargePi-go/pkg/models/settings"
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
	TypeDummy  = "dummy"
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
		ChangeColor(index int, color Color) error
		Blink(index int, times int, color Color) error
		SetBrightness(brightness int) error
		GetBrightness() int
		Cleanup()
		GetType() string
	}
)

// NewIndicator constructs the Indicator based on the type provided by the settings file.
func NewIndicator(stripLength int, indicator settings.Indicator) Indicator {
	logger := zap.L()

	if indicator.Enabled {
		// Last LED is used to indicate card read
		if indicator.IndicateCardRead {
			stripLength++
		}

		logger.With(zap.String("indicator_type", indicator.Type)).Info("Preparing Indicator from config")
		switch indicator.Type {
		case TypeWS281x:
			if util.IsNilInterfaceOrPointer(indicator.WS281x) {
				return nil
			}

			ledStrip, ledError := NewWS281xStrip(logger, stripLength, indicator.WS281x.DataPin)
			if ledError != nil {
				logger.With(zap.Error(ledError)).Error("Error creating indicator")
				return nil
			}

			return ledStrip
		case TypeDummy:
			return NewDummy(logger, indicator.IndicatorDummy)
		default:
			return nil
		}
	}

	return nil
}
