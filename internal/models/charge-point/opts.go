package chargePoint

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/reader"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
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
		if !readerSettings.IsEnabled {
			return
		}

		// Create reader based on settings
		tagReader, err := reader.NewTagReader(readerSettings)
		if err != nil {
			return
		}

		point.SetReader(tagReader)

		// Listen for incoming tags
		go tagReader.ListenForTags(ctx)
		go point.ListenForTag(ctx, tagReader.GetTagChannel())
	}
}

// WithReader adds the reader to the charge point and starts listening to the Reader.
func WithReader(ctx context.Context, tagReader reader.Reader) Options {
	return func(point ChargePoint) {
		if util.IsNilInterfaceOrPointer(tagReader) {
			return
		}
		point.SetReader(tagReader)

		// Listen for incoming tags
		go tagReader.ListenForTags(ctx)
		go point.ListenForTag(ctx, tagReader.GetTagChannel())
	}
}

// WithDisplayFromSettings create a Display based on the provided settings.
func WithDisplayFromSettings(lcdSettings settings.Display) Options {
	return func(point ChargePoint) {
		if !lcdSettings.IsEnabled {
			return
		}

		lcd, err := display.NewDisplay(lcdSettings)
		if err != nil {
			return
		}

		point.SetDisplay(lcd)

	}
}

// WithDisplay add the provided Display to the ChargePoint.
func WithDisplay(display display.Display) Options {
	return func(point ChargePoint) {
		if util.IsNilInterfaceOrPointer(display) {
			return
		}
		point.SetDisplay(display)
	}
}
