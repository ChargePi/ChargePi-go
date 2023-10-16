package chargePoint

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/pkg/display"
	"github.com/xBlaz3kx/ChargePi-go/pkg/indicator"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/reader"
)

type Options func(point ChargePoint)

// WithLogger add logger to the ChargePoint
func WithLogger(logger *log.Logger) Options {
	return func(point ChargePoint) {
		if logger != nil {
			point.SetLogger(logger)
		}
	}
}

// WithReaderFromSettings creates a TagReader based on the settings.
func WithReaderFromSettings(ctx context.Context, readerSettings settings.TagReader) Options {
	return func(point ChargePoint) {
		// Create reader based on settings
		tagReader, err := reader.NewTagReader(readerSettings)
		switch {
		case errors.Is(err, reader.ErrReaderDisabled):
			return
		case errors.Is(err, reader.ErrReaderUnsupported):
			log.WithError(err).Fatal("Error attaching a display")
		}

		point.SetReader(tagReader)
	}
}

// WithReader adds the reader to the charge point and starts listening to the Reader.
func WithReader(ctx context.Context, tagReader reader.Reader) Options {
	return func(point ChargePoint) {
		point.SetReader(tagReader)
	}
}

// WithDisplayFromSettings create a Display based on the provided settings.
func WithDisplayFromSettings(lcdSettings settings.Display) Options {
	return func(point ChargePoint) {
		lcd, err := display.NewDisplay(lcdSettings)
		switch {
		case errors.Is(err, display.ErrDisplayDisabled):
			return
		case errors.Is(err, display.ErrDisplayUnsupported), errors.Is(err, display.ErrInvalidConnectionDetails):
			log.WithError(err).Fatal("Error attaching a display")
		}

		point.SetDisplay(lcd)
	}
}

// WithDisplay add the provided Display to the ChargePoint.
func WithDisplay(display display.Display) Options {
	return func(point ChargePoint) {
		point.SetDisplay(display)
	}
}

// WithIndicator add an indicator
func WithIndicator(indicator indicator.Indicator) Options {
	return func(point ChargePoint) {
		point.SetIndicator(indicator)
	}
}
