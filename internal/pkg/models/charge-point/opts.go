package chargePoint

import (
	"context"
	"errors"
	"go.uber.org/zap"

	"github.com/ChargePi/ChargePi-go/pkg/display"
	"github.com/ChargePi/ChargePi-go/pkg/indicator"
	"github.com/ChargePi/ChargePi-go/pkg/models/settings"
	"github.com/ChargePi/ChargePi-go/pkg/reader"
)

type Options func(point ChargePoint)

// WithLogger add logger to the ChargePoint
func WithLogger(logger *zap.Logger) Options {
	return func(point ChargePoint) {
		if logger != nil {
			point.SetLogger(logger)
		}
	}
}

// WithReaderFromSettings creates a TagReader based on the settings.
func WithReaderFromSettings(ctx context.Context, readerSettings settings.TagReader) Options {
	return func(point ChargePoint) {
		logger := zap.L()
		// Create reader based on settings
		tagReader, err := reader.NewTagReader(readerSettings)
		switch {
		case errors.Is(err, reader.ErrReaderDisabled):
			return
		case errors.Is(err, reader.ErrReaderUnsupported):
			logger.With(zap.Error(err)).Fatal("Error attaching a display")
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
		logger := zap.L()

		lcd, err := display.NewDisplay(lcdSettings)
		switch {
		case errors.Is(err, display.ErrDisplayDisabled):
			return
		case errors.Is(err, display.ErrDisplayUnsupported), errors.Is(err, display.ErrInvalidConnectionDetails):
			logger.With(zap.Error(err)).Fatal("Error attaching a display")
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
