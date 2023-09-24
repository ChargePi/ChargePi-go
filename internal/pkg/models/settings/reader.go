package settings

type (
	TagReader struct {
		IsEnabled   bool   `json:"enabled,omitempty" yaml:"enabled,omitempty" mapstructure:"enabled,omitempty"`
		ReaderModel string `json:"readerModel,omitempty" yaml:"readerModel,omitempty" mapstructure:"readerModel,omitempty"`
		PN532       *PN532 `json:"pn532,omitempty" yaml:"pn532,omitempty" mapstructure:"pn532,omitempty"`
	}

	PN532 struct {
		Device   string `json:"deviceAddress,omitempty" yaml:"deviceAddress,omitempty" mapstructure:"deviceAddress,omitempty"`
		ResetPin int    `json:"resetPin,omitempty" yaml:"resetPin,omitempty" mapstructure:"resetPin,omitempty"`
	}
)
