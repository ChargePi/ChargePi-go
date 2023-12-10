package v16

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"github.com/xBlaz3kx/ChargePi-go/pkg/display"
	"github.com/xBlaz3kx/ChargePi-go/pkg/indicator"
	settings2 "github.com/xBlaz3kx/ChargePi-go/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/reader"
)

func (cp *ChargePoint) SetLogger(logger log.FieldLogger) {
	cp.logger = logger
}

func (cp *ChargePoint) SetReader(reader reader.Reader) error {
	if util.IsNilInterfaceOrPointer(reader) {
		return nil
	}

	cp.logger.Debugf("Setting reader")
	cp.tagReader = reader
	return nil
}

func (cp *ChargePoint) SetDisplay(display display.Display) error {
	if util.IsNilInterfaceOrPointer(display) {
		return nil
	}

	cp.logger.Debug("Setting display")
	cp.display = display
	cp.display.Clear()
	return nil
}

func (cp *ChargePoint) SetIndicator(indicator indicator.Indicator) error {
	if util.IsNilInterfaceOrPointer(indicator) {
		return nil
	}

	cp.logger.Debug("Setting indicator")
	cp.indicator = indicator
	return nil
}

func (cp *ChargePoint) SetSettings(settings settings.Info) error {
	err := validator.New().StructCtx(context.Background(), settings)
	if err != nil {
		return err
	}

	cp.logger.Debug("Setting charge point settings")
	cp.info = settings
	return nil
}

func (cp *ChargePoint) GetConnectionSettings() settings.ConnectionSettings {
	return cp.connectionSettings
}

func (cp *ChargePoint) SetConnectionSettings(settings settings.ConnectionSettings) error {
	err := validator.New().StructCtx(context.Background(), settings)
	if err != nil {
		return err
	}

	cp.logger.Debug("Setting backend connection settings")
	cp.connectionSettings = settings
	return nil
}

func (cp *ChargePoint) GetSettings() settings.Info {
	return cp.info
}

func (cp *ChargePoint) GetIndicatorSettings() settings2.IndicatorStatusMapping {
	return cp.indicatorMapping
}

func (cp *ChargePoint) SetIndicatorSettings(settings settings2.IndicatorStatusMapping) error {
	err := validator.New().StructCtx(context.Background(), settings)
	if err != nil {
		return err
	}

	cp.logger.Debug("Setting indicator settings")
	cp.indicatorMapping = settings
	return nil
}

func (cp *ChargePoint) SetAvailability(availabilityType core.AvailabilityType) error {
	cp.logger.WithField("availability", availabilityType).Debug("Setting availability")

	// Check if there are ongoing transactions
	_, sessionErr := cp.sessionManager.GetSession(0, nil)
	switch sessionErr {
	case nil:
		// Set availability status and notify the backend
		cp.availability = availabilityType

		switch availabilityType {
		case core.AvailabilityTypeInoperative:
			cp.notifyConnectorStatus(0, core.ChargePointStatusUnavailable, core.NoError)
		case core.AvailabilityTypeOperative:
			cp.notifyConnectorStatus(0, core.ChargePointStatusAvailable, core.NoError)
		}

	default:
		return errors.New("error checking for ongoing transactions")
	}

	return nil
}

func (cp *ChargePoint) ApplyOpts(opts ...chargePoint.Options) {
	// Apply options
	for _, opt := range opts {
		opt(cp)
	}
}
