//go:build demo

package reader

import (
	"go.uber.org/zap"
)

// NewTagReader creates an instance of the Reader interface based on the provided configuration.
func NewTagReader(reader Settings) (Reader, error) {
	if reader.IsEnabled {
		logger := zap.L()
		logger.Sugar().Infof("Preparing tag reader from config: %s", reader.ReaderModel)
		switch reader.ReaderModel {
		case TypeDummy:
			return NewDummy(logger, reader.DummyReader)
		default:
			return nil, ErrReaderUnsupported
		}
	}

	return nil, ErrReaderDisabled
}
