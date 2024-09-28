//go:build linux

package evcc

import (
	"errors"

	"github.com/ChargePi/ChargePi-go/pkg/hardware"
)

type Settings struct {
	Type string `validate:"required" json:"type,omitempty" yaml:"type" mapstructure:"type"`
	// Based on the type, get the connection details
	Relay *RelaySettings `json:"relay,omitempty" yaml:"relay,omitempty" mapstructure:"relay,omitempty"`

	Serial *hardware.Serial `json:"serial,omitempty" yaml:"serial,omitempty" mapstructure:"serial,omitempty"`

	Modbus *hardware.ModBus `json:"modbus,omitempty" yaml:"modbus,omitempty" mapstructure:"modbus,omitempty"`

	Dummy *DummySettings `json:"dummy,omitempty" yaml:"dummy,omitempty" mapstructure:"dummy,omitempty"`
}

// NewEVCCFromType creates a new EVCC instance based on the provided type.
func NewEVCCFromType(evccSettings Settings) (EVCC, error) {
	switch evccSettings.Type {
	case Relay:
		if evccSettings.Relay == nil {
			return nil, errors.New("missing relay settings")
		}

		return NewRelay(*evccSettings.Relay)
	case TypeDummy:
		return NewDummy(evccSettings.Dummy)
	default:
		return nil, ErrInvalidType
	}
}
