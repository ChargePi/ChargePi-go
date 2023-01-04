package settings

type (
	I2C struct {
		Address string `json:"address,omitempty" yaml:"address,omitempty" mapstructure:"address,omitempty"`
		Bus     int    `json:"bus,omitempty" yaml:"bus,omitempty" mapstructure:"bus,omitempty"`
	}

	SPI struct {
		ChipSelect int `json:"chipSelect,omitempty" yaml:"chipSelect,omitempty" mapstructure:"chipSelect,omitempty"`
		Bus        int `json:"bus,omitempty" yaml:"bus,omitempty" mapstructure:"bus,omitempty"`
	}

	ModBus struct {
		DeviceAddress string `json:"deviceAddress,omitempty" yaml:"deviceAddress,omitempty" mapstructure:"deviceAddress,omitempty"`
		Protocol      string `json:"protocol,omitempty" yaml:"protocol,omitempty" mapstructure:"protocol,omitempty"`
	}

	Serial struct {
		DeviceAddress string `json:"deviceAddress,omitempty" yaml:"deviceAddress,omitempty" mapstructure:"deviceAddress,omitempty"`
		BaudRate      uint   `json:"baudRate,omitempty" yaml:"baudRate,omitempty" mapstructure:"baudRate,omitempty"`
		Parity        uint   `json:"parity,omitempty" yaml:"parity,omitempty" mapstructure:"parity,omitempty"`
		DataBits      uint8  `json:"dataBits,omitempty" yaml:"dataBits,omitempty" mapstructure:"dataBits,omitempty"`
		StopBits      uint8  `json:"stopBits,omitempty" yaml:"stopBits,omitempty" mapstructure:"stopBits,omitempty"`
	}
)
