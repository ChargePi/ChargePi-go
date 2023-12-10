package settings

type (
	TagReader struct {
		IsEnabled   bool         `json:"enabled,omitempty" yaml:"enabled,omitempty" mapstructure:"enabled,omitempty"`
		ReaderModel string       `json:"model,omitempty" yaml:"model,omitempty" mapstructure:"model,omitempty"`
		PN532       *PN532       `json:"pn532,omitempty" yaml:"pn532,omitempty" mapstructure:"pn532,omitempty"`
		DummyReader *DummyReader `json:"dummy,omitempty" yaml:"dummy,omitempty" mapstructure:"dummy,omitempty"`
	}

	PN532 struct {
		Device   string `json:"deviceAddress,omitempty" yaml:"deviceAddress,omitempty" mapstructure:"deviceAddress,omitempty"`
		ResetPin int    `json:"resetPin,omitempty" yaml:"resetPin,omitempty" mapstructure:"resetPin,omitempty"`
	}

	DummyReader struct {
		TagIds []string `json:"tagIds,omitempty" yaml:"tagIds,omitempty" mapstructure:"tagIds,omitempty"`
	}
)
