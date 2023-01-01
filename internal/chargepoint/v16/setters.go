package v16

import (
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/reader"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	settings2 "github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
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
	return nil
}

func (cp *ChargePoint) SetIndicator(indicator indicator.Indicator) error {
	if util.IsNilInterfaceOrPointer(indicator) {
		return nil
	}

	cp.indicator = indicator
	return nil
}

func (cp *ChargePoint) SetSettings(settings settings2.Info) {
	cp.info = settings
}

func (cp *ChargePoint) SetIndicatorSettings(settings settings2.IndicatorStatusMapping) {
	cp.indicatorMapping = settings
}

func (cp *ChargePoint) SetConnectionSettings(settings settings2.ConnectionSettings) {
	cp.connectionSettings = settings
}

func (cp *ChargePoint) GetSettings() settings2.Info {
	return cp.info
}

func (cp *ChargePoint) GetIndicatorSettings() settings2.IndicatorStatusMapping {
	return cp.indicatorMapping
}

func (cp *ChargePoint) GetConnectionSettings() settings2.ConnectionSettings {
	return cp.connectionSettings
}

func (cp *ChargePoint) ApplyOpts(opts ...chargePoint.Options) {
	// Apply options
	for _, opt := range opts {
		opt(cp)
	}
}
