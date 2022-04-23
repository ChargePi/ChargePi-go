package indicator

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// color constants
const (
	Off    = 0x0
	White  = 0xFFFFFF
	Red    = 0xff0000
	Green  = 0x00ff00
	Blue   = 0x000ff
	Yellow = 0xeeff00
	Orange = 0xff7b00
)

// Supported types
const (
	TypeWS281x = "WS281x"
)

var (
	ErrInvalidIndex        = errors.New("invalid index")
	ErrInvalidPin          = errors.New("invalid data pin #")
	ErrInvalidNumberOfLeds = errors.New("number of leds must be greater than zero")
)

type (
	// Indicator is an abstraction layer for connector status indication, usually an RGB LED strip.
	Indicator interface {
		DisplayColor(index int, colorHex uint32) error
		Blink(index int, times int, colorHex uint32) error
		Cleanup()
	}
)

// NewIndicator constructs the Indicator based on the type provided by the settings file.
func NewIndicator(stripLength int) Indicator {
	var (
		indicatorEnabled = viper.GetBool("chargepoint.hardware.ledIndicator.enabled")
		indicatorType    = viper.GetString("chargepoint.hardware.ledIndicator.type")
		indicateCardRead = viper.GetBool("chargepoint.hardware.ledIndicator.indicateCardRead")
		indicatorDataPin = viper.GetInt("chargepoint.hardware.ledIndicator.dataPin")
	)

	if indicatorEnabled {
		if indicateCardRead {
			stripLength++
		}

		log.Infof("Preparing Indicator from config: %s", indicatorType)
		switch indicatorType {
		case TypeWS281x:
			ledStrip, ledError := NewWS281xStrip(stripLength, indicatorDataPin)
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
