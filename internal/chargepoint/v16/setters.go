package v16

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/reader"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
)

func (cp *ChargePoint) SetLogger(logger *log.Logger) {
	cp.logger = logger
}

func (cp *ChargePoint) SetReader(reader reader.Reader) error {
	if util.IsNilInterfaceOrPointer(reader) {
		return nil
	}

	cp.tagReader = reader
	return nil
}

func (cp *ChargePoint) SetDisplay(display display.Display) error {
	if util.IsNilInterfaceOrPointer(display) {
		return nil
	}

	cp.display = display
	cp.display.Clear()
	return nil
}

func (cp *ChargePoint) SetIndicator(indicator indicator.Indicator) error {
	if util.IsNilInterfaceOrPointer(indicator) {
		return nil
	}

	cp.indicator = indicator
	return nil
}

func (cp *ChargePoint) SetSettings(settings settings.Info) error {
	err := validator.New().StructCtx(context.Background(), settings)
	if err != nil {
		return err
	}

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

	cp.connectionSettings = settings
	return nil
}

func (cp *ChargePoint) GetSettings() settings.Info {
	return cp.info
}

func (cp *ChargePoint) GetIndicatorSettings() settings.IndicatorStatusMapping {
	return cp.indicatorMapping
}

func (cp *ChargePoint) SetIndicatorSettings(settings settings.IndicatorStatusMapping) error {
	err := validator.New().StructCtx(context.Background(), settings)
	if err != nil {
		return err
	}

	cp.indicatorMapping = settings
	return nil
}

func (cp *ChargePoint) SetAvailability(availabilityType core.AvailabilityType) error {
	// Check if there are ongoing transactions
	_, sessionErr := cp.sessionManager.GetSession(0, nil)
	switch sessionErr {
	case nil:
		cp.availability = availabilityType
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
