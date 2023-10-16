package reader

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/settings"
)

// Supported readers - by libnfc
const (
	PN532     = "PN532"
	ACR122    = "ACR122"
	PN533     = "PN533"
	BR500     = "BR500"
	R502      = "R502"
	TypeDummy = "dummy"
)

var (
	ErrReaderUnsupported = errors.New("reader type unsupported")
	ErrReaderDisabled    = errors.New("reader disabled")
)

// Reader is an abstraction for an RFID/NFC tag reader.
type Reader interface {
	ListenForTags(ctx context.Context)
	Cleanup()
	Reset()
	GetTagChannel() <-chan string
	GetType() string
}

// NewTagReader creates an instance of the Reader interface based on the provided configuration.
func NewTagReader(reader settings.TagReader) (Reader, error) {
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