package v16

import (
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/reader"
	chargePoint "github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
)

func (cp *ChargePoint) SetLogger(logger *log.Logger) {
	cp.logger = logger
}

func (cp *ChargePoint) SetReader(reader reader.Reader) {
	cp.tagReader = reader
}

func (cp *ChargePoint) SetDisplay(display display.Display) {
	cp.display = display
}

func (cp *ChargePoint) SetIndicator(indicator indicator.Indicator) {
	cp.indicator = indicator
}

func (cp *ChargePoint) SetSettings(settings *settings.Settings) {
	cp.settings = settings
}

func (cp *ChargePoint) ApplyOpts(opts ...chargePoint.Options) {
	// Apply options
	for _, opt := range opts {
		opt(cp)
	}
}
