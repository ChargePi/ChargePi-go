//go:build linux && cgo

package reader

import (
	log "github.com/sirupsen/logrus"
)

type Settings struct {
	IsEnabled   bool           `json:"enabled,omitempty" yaml:"enabled,omitempty" mapstructure:"enabled,omitempty"`
	ReaderModel string         `json:"model,omitempty" yaml:"model,omitempty" mapstructure:"model,omitempty"`
	PN532       *PN532Settings `json:"pn532,omitempty" yaml:"pn532,omitempty" mapstructure:"pn532,omitempty"`
	DummyReader *DummySettings `json:"dummy,omitempty" yaml:"dummy,omitempty" mapstructure:"dummy,omitempty"`
}

// NewTagReader creates an instance of the Reader interface based on the provided configuration.
func NewTagReader(reader Settings) (Reader, error) {
	if reader.IsEnabled {
		log.Infof("Preparing tag reader from config: %s", reader.ReaderModel)
		switch reader.ReaderModel {
		case PN532, ACR122, PN533, BR500, R502:
			return NewReader(reader.PN532.Device, reader.ReaderModel, reader.PN532.ResetPin)
		case TypeDummy:
			return NewDummy(reader.DummyReader)
		default:
			return nil, ErrReaderUnsupported
		}
	}

	return nil, ErrReaderDisabled
}
