package reader

import (
	"context"
	"errors"
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
