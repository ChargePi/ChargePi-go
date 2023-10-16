package settings

type EVCC struct {
	Type string `validate:"required" json:"type,omitempty" yaml:"type" mapstructure:"type"`
	// Based on the type, get the connection details
	Relay *Relay `json:"relay,omitempty" yaml:"relay,omitempty" mapstructure:"relay,omitempty"`

	Serial *Serial `json:"serial,omitempty" yaml:"serial,omitempty" mapstructure:"serial,omitempty"`

	Modbus *ModBus `json:"modbus,omitempty" yaml:"modbus,omitempty" mapstructure:"modbus,omitempty"`

	Dummy *EvccDummy `json:"dummy,omitempty" yaml:"dummy,omitempty" mapstructure:"dummy,omitempty"`
}

type Relay struct {
	RelayPin     int  `json:"relayPin,omitempty" yaml:"relayPin,omitempty" mapstructure:"relayPin,omitempty"`
	InverseLogic bool `json:"inverseLogic,omitempty" yaml:"inverseLogic,omitempty" mapstructure:"inverseLogic,omitempty"`
}

type EvccDummy struct {
}
