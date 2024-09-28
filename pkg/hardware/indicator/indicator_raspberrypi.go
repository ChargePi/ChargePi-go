//go:build raspberrypi4

package indicator

import (
	"github.com/ChargePi/ChargePi-go/pkg/util"
	log "github.com/sirupsen/logrus"
)

type Settings struct {
	// Enable or disable the indicator
	Enabled bool `json:"enabled" yaml:"enabled" mapstructure:"enabled"`

	// Indicate card read
	IndicateCardRead bool `json:"indicateCardRead,omitempty" yaml:"indicateCardRead,omitempty" mapstructure:"indicateCardRead,omitempty"`

	// Type of indicator
	Type string `json:"type,omitempty" yaml:"type" mapstructure:"type"`

	// Statuses
	IndicatorMappings *StatusMapping `json:"statuses,omitempty" yaml:"statuses,omitempty" mapstructure:"statuses,omitempty"`

	// Based on the type, get the connection details
	WS281x *WS281xSettings `json:"ws281x,omitempty" yaml:"ws281x,omitempty" mapstructure:"ws281x,omitempty"`

	// Dummy indicator
	IndicatorDummy *DummySettings `json:"dummy,omitempty" yaml:"dummy,omitempty" mapstructure:"dummy,omitempty"`
}

// NewIndicator constructs the Settings based on the type provided by the settings file.
func NewIndicator(stripLength int, indicator Settings) Indicator {
	if indicator.Enabled {

		// Last LED is used to indicate card read
		if indicator.IndicateCardRead {
			stripLength++
		}

		log.Infof("Preparing Settings from config: %s", indicator.Type)
		switch indicator.Type {
		case TypeWS281x:
			if util.IsNilInterfaceOrPointer(indicator.WS281x) {
				return nil
			}

			//ledStrip, ledError := NewWS281xStrip(stripLength, indicator.WS281x.DataPin)
			//if ledError != nil {
			//	log.WithError(ledError).Errorf("Error creating indicator")
			//	return nil
			//}
			//return ledStrip
			return nil
		case TypeDummy:
			return NewDummy(*indicator.IndicatorDummy)
		default:
			return nil
		}
	}

	return nil
}
