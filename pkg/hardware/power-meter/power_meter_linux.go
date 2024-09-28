//go:build linux && cgo

package powerMeter

import (
	"github.com/ChargePi/ChargePi-go/pkg/hardware"
	"github.com/ChargePi/ChargePi-go/pkg/util"
	log "github.com/sirupsen/logrus"
)

type Settings struct {
	Enabled bool `json:"enabled" yaml:"enabled" mapstructure:"enabled"`

	Type string `json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type,omitempty"`

	// Based on the type, get the connection details
	// For smarter energy meters, using Modbus RTU or similar based on the device type
	ModBus *hardware.ModBus `json:"modbus,omitempty" yaml:"modbus,omitempty" mapstructure:"modbus,omitempty"`

	// CS5460 specific details
	CS5460 *CS5460Settings `json:"cs5460,omitempty" yaml:"cs5460,omitempty" mapstructure:"cs5460,omitempty"`

	// Dummy power meter for testing
	PowerMeterDummy *DummySettings `json:"dummy,omitempty" yaml:"dummy,omitempty" mapstructure:"dummy,omitempty"`
}

// NewPowerMeter creates a new power meter based on the connector settings.
func NewPowerMeter(meterSettings Settings) (PowerMeter, error) {
	if meterSettings.Enabled {
		log.Infof("Creating a new power meter: %s", meterSettings.Type)

		switch meterSettings.Type {
		case TypeC5460A:
			if util.IsNilInterfaceOrPointer(meterSettings.CS5460) {
				return nil, ErrInvalidConnectionSettings
			}

			powerMeter, err := NewCS5460PowerMeter(*meterSettings.CS5460)
			if err != nil {
				return nil, err
			}

			return powerMeter, nil
		case TypeDummy:
			return NewDummy(meterSettings.PowerMeterDummy)
		default:
			return nil, ErrPowerMeterUnsupported
		}
	}

	return nil, ErrPowerMeterDisabled
}
