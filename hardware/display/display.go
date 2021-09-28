package display

import (
	"github.com/xBlaz3kx/ChargePi-go/cache"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"log"
	"time"
)

const (
	DriverHD44780 = "hd44780"
)

type (
	// LCDMessage Object representing the message that will be displayed on the LCD.
	// Each array element in Messages represents a line being displayed on the 16x2 screen.
	LCDMessage struct {
		Messages        []string
		messageDuration time.Duration
	}

	// LCD is an abstraction layer for concrete implementation of a display.
	LCD interface {
		DisplayMessage(message LCDMessage)
		ListenForMessages()
		Cleanup()
		Clear()
		GetLcdChannel() chan LCDMessage
	}
)

// NewMessage creates a new message for the LCD.
func NewMessage(duration time.Duration, messages []string) LCDMessage {
	return LCDMessage{
		Messages:        messages,
		messageDuration: duration,
	}
}

// NewDisplay returns a concrete implementation of an LCD based on the drivers that are supported.
// The LCD is built with the settings from the settings file.
func NewDisplay() LCD {
	cacheSettings, isFound := cache.Cache.Get("settings")
	if !isFound {
		panic("settings not found")
	}
	config := cacheSettings.(*settings.Settings)
	lcdSettings := config.ChargePoint.Hardware.Lcd

	if lcdSettings.IsSupported {
		log.Println("Preparing LCD from config")
		switch lcdSettings.Driver {
		case DriverHD44780:
			lcdChannel := make(chan LCDMessage, 5)
			lcd, err := NewHD44780(lcdChannel, lcdSettings.I2CAddress, lcdSettings.I2CBus)
			if err != nil {
				log.Println("Could not create the LCD:", err)
				return nil
			}
			return lcd
		default:
			return nil
		}
	}
	return nil
}
