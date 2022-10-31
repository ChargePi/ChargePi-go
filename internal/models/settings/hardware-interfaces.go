package settings

type (
	I2C struct {
		Address string `fig:"address" json:"address,omitempty" yaml:"address" mapstructure:"address"`
		Bus     int    `fig:"bus" json:"bus,omitempty" yaml:"bus" mapstructure:"bus"`
	}

	SPI struct {
		ChipSelect int `fig:"chipSelect" json:"chipSelect,omitempty" yaml:"chipSelect" mapstructure:"chipSelect"`
		Bus        int `fig:"bus" json:"bus,omitempty" yaml:"bus" mapstructure:"bus"`
	}

	ModBus struct {
		DeviceAddress string `fig:"deviceAddress" json:"deviceAddress,omitempty" yaml:"deviceAddress" mapstructure:"deviceAddress"`
		Protocol      string `fig:"protocol" json:"protocol,omitempty" yaml:"protocol" mapstructure:"protocol"`
	}

	Serial struct {
		DeviceAddress string `fig:"deviceAddress" json:"deviceAddress,omitempty" yaml:"deviceAddress" mapstructure:"deviceAddress"`
		BaudRate      uint   `fig:"baudRate" json:"baudRate,omitempty" yaml:"baudRate" mapstructure:"baudRate"`
		Parity        uint   `fig:"parity" json:"parity,omitempty" yaml:"parity" mapstructure:"parity"`
		DataBits      uint8  `fig:"dataBits" json:"dataBits,omitempty" yaml:"dataBits" mapstructure:"dataBits"`
		StopBits      uint8  `fig:"stopBits" json:"stopBits,omitempty" yaml:"stopBits" mapstructure:"stopBits"`
	}
)
