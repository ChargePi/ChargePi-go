//go:build cgo && raspberrypi4

package display

import (
	"github.com/ChargePi/ChargePi-go/pkg/util"
	log "github.com/sirupsen/logrus"
)

type Settings struct {
	// Enable or disable the display from the configuration
	IsEnabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty" mapstructure:"enabled,omitempty"`

	// Enable the display to be controller remotely through OCPP
	RemoteEnabled bool `json:"remote,omitempty" yaml:"remote,omitempty" mapstructure:"remote,omitempty"`

	// Display driver/type of the display - can be a direct implementation of the driver or an HTTP server
	Driver string `json:"driver,omitempty" yaml:"driver,omitempty" mapstructure:"driver,omitempty"`

	// Default display language
	Language string `json:"language,omitempty" yaml:"language,omitempty" mapstructure:"language,omitempty"`

	// Hitachi HD44780 display configuration details
	HD44780 *HD44780Settings `json:"hd44780,omitempty" yaml:"hd44780,omitempty" mapstructure:"hd44780,omitempty"`

	// Dummy display configuration details
	DisplayDummy *DummySettings `json:"dummy,omitempty" yaml:"dummy,omitempty" mapstructure:"dummy,omitempty"`
}

// NewDisplay returns a concrete implementation of a Display based on the drivers that are supported.
// The Display is built with the settings from the settings file.
func NewDisplay(lcdSettings Settings) (Display, error) {
	if lcdSettings.IsEnabled {
		log.Info("Preparing display from config")

		switch lcdSettings.Driver {
		case DriverHD44780:
			if util.IsNilInterfaceOrPointer(lcdSettings.HD44780) {
				return nil, ErrInvalidConnectionDetails
			}

			lcd, err := NewHD44780(*lcdSettings.HD44780)
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
