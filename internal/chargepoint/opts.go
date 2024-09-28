package chargepoint

import (
	"context"
	"errors"
	"go.uber.org/zap"

	"github.com/ChargePi/ChargePi-go/pkg/hardware/display"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/indicator"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/reader"
)

type Options func(point ChargePoint) error

// WithLogger add logger to the ChargePoint
func WithLogger(logger *zap.Logger) Options {
	return func(point ChargePoint) error {
		if logger != nil {
			point.SetLogger(logger)
		}

		return nil
	}
}

// WithReaderFromSettings creates a TagReader based on the settings.
func WithReaderFromSettings(ctx context.Context, readerSettings reader.Settings) Options {
	return func(point ChargePoint) error {
		// Create reader based on settings
		tagReader, err := reader.NewTagReader(readerSettings)
		switch {
		case errors.Is(err, reader.ErrReaderDisabled):
			return nil
		case errors.Is(err, reader.ErrReaderUnsupported):
			return err
		}

		return point.SetReader(tagReader)
	}
}

// WithReader adds the reader to the charge point and starts listening to the Reader.
func WithReader(ctx context.Context, tagReader reader.Reader) Options {
	return func(point ChargePoint) error {
		return point.SetReader(tagReader)
	}
}

// WithDisplayFromSettings create a Display based on the provided settings.
func WithDisplayFromSettings(lcdSettings display.Settings) Options {
	return func(point ChargePoint) error {
		logger := zap.L()
		lcd, err := display.NewDisplay(lcdSettings)
		switch {
		case errors.Is(err, display.ErrDisplayDisabled):
			return nil
		case errors.Is(err, display.ErrDisplayUnsupported), errors.Is(err, display.ErrInvalidConnectionDetails):
			logger.With(zap.Error(err)).Error("Error attaching a display")
			return err
		}

		return point.SetDisplay(lcd)
	}
}

// WithDisplay add the provided Display to the ChargePoint.
func WithDisplay(display display.Display) Options {
	return func(point ChargePoint) error {
		return point.SetDisplay(display)
	}
}

// WithIndicator add an indicator
func WithIndicator(indicator indicator.Indicator) Options {
	return func(point ChargePoint) error {
		return point.SetIndicator(indicator)
	}
}
