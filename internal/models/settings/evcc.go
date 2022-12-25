package settings

type (
	EVCC struct {
		Type string
		// Based on the type, get the connection details
		RelayPin     int     `fig:"RelayPin" json:"RelayPin,omitempty" yaml:"RelayPin" mapstructure:"RelayPin"`
		InverseLogic bool    `fig:"InverseLogic" json:"InverseLogic,omitempty" yaml:"InverseLogic" mapstructure:"InverseLogic"`
		Serial       *Serial `fig:"serial" json:"serial,omitempty" yaml:"serial" mapstructure:"serial"`
		Modbus       *ModBus `fig:"modbus" json:"modbus,omitempty" yaml:"modbus" mapstructure:"modbus"`
	}
)
