package indicator

import (
	"github.com/xBlaz3kx/ChargePi-go/cache"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"log"
)

const (
	// color constants
	Off    = 0x0
	White  = 0xFFFFFF
	Red    = 0xff0000
	Green  = 0x00ff00
	Blue   = 0x000ff
	Yellow = 0xeeff00
	Orange = 0xff7b00

	// supported types
	TypeWS281x = "WS281x"
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
	cacheSettings, isFound := cache.Cache.Get("settings")
	if !isFound {
		panic("settings not found")
	}
	config := cacheSettings.(*settings.Settings)
	indicatorSettings := config.ChargePoint.Hardware.LedIndicator

	if indicatorSettings.Enabled {
		log.Println("Preparing Indicator from config: ", indicatorSettings.Type)
		switch indicatorSettings.Type {
		case TypeWS281x:
			if indicatorSettings.IndicateCardRead {
				stripLength++
			}

			ledStrip, ledError := NewWS281xStrip(stripLength, indicatorSettings.DataPin)
			if ledError != nil {
				log.Println(ledError)
				return nil
			}

			return ledStrip
		default:
			return nil
		}
	}
	return nil
}
