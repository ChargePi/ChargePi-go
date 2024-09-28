//go:build demo

package reader

import log "github.com/sirupsen/logrus"

// NewTagReader creates an instance of the Reader interface based on the provided configuration.
func NewTagReader(reader Settings) (Reader, error) {
	if reader.IsEnabled {
		log.Infof("Preparing tag reader from config: %s", reader.ReaderModel)
		switch reader.ReaderModel {
		case TypeDummy:
			return NewDummy(reader.DummyReader)
		default:
			return nil, ErrReaderUnsupported
		}
	}

	return nil, ErrReaderDisabled
}
