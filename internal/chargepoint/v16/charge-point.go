package v16

import (
	"context"
	"github.com/go-co-op/gocron"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/auth"
	connectorManager "github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/reader"
	chargePoint "github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
)

type (
	ChargePoint struct {
		chargePoint  ocpp16.ChargePoint
		availability core.AvailabilityType
		settings     *settings.Settings
		// Hardware components
		tagReader reader.Reader
		indicator indicator.Indicator
		display   display.Display
		// Software components
		connectorManager   connectorManager.Manager
		connectorChannel   chan chargePoint.StatusNotification
		meterValuesChannel chan chargePoint.MeterValueNotification
		scheduler          *gocron.Scheduler
		tagManager         auth.TagManager
		logger             *log.Logger
	}
)

// NewChargePoint creates a new ChargePoint for OCPP version 1.6.
func NewChargePoint(manager connectorManager.Manager, scheduler *gocron.Scheduler, cache auth.TagManager, opts ...chargePoint.Options) *ChargePoint {
	ch := make(chan chargePoint.StatusNotification, 5)
	// Set the channel
	manager.SetNotificationChannel(ch)

	cp := &ChargePoint{
		availability:     core.AvailabilityTypeInoperative,
		connectorChannel: ch,
		scheduler:        scheduler,
		connectorManager: manager,
		tagManager:       cache,
		logger:           log.StandardLogger(),
	}

	cp.ApplyOpts(opts...)

	return cp
}

// Connect to the central system and send a BootNotification
func (cp *ChargePoint) Connect(ctx context.Context, serverUrl string) {
	var (
		connectionSettings = cp.settings.ChargePoint.ConnectionSettings
		tlsConfig          = connectionSettings.TLS
		wsClient           = util.CreateClient(
			connectionSettings.BasicAuthUsername,
			connectionSettings.BasicAuthPassword,
			tlsConfig)
		logInfo = log.WithFields(log.Fields{
			"chargePointId": connectionSettings.Id,
		})
	)

	logInfo.Debug("Creating charge point")
	cp.chargePoint = ocpp16.NewChargePoint(connectionSettings.Id, nil, wsClient)

	// Set charging profiles
	util.SetProfilesFromConfig(cp.chargePoint, cp, cp, cp)

	cp.logger.Infof("Trying to connect to the central system: %s", serverUrl)
	connectErr := cp.chargePoint.Start(serverUrl)
	if connectErr != nil {
		//cp.CleanUp(core.ReasonOther)
		cp.logger.WithError(connectErr).Fatalf("Cannot connect to the central system")
	}

	cp.logger.Infof("Successfully connected to: %s", serverUrl)
	cp.availability = core.AvailabilityTypeOperative

	go cp.ListenForConnectorStatusChange(ctx, cp.connectorChannel)
	cp.bootNotification()
}

// CleanUp When exiting the client, stop all the transactions, clean up all the peripherals and terminate the connection.
func (cp *ChargePoint) CleanUp(reason core.Reason) {
	cp.logger.Infof("Cleaning up ChargePoint, reason: %s", reason)

	switch reason {
	case core.ReasonRemote, core.ReasonLocal, core.ReasonHardReset, core.ReasonSoftReset:
		for _, c := range cp.connectorManager.GetEVSEs() {
			// Stop charging the connectors
			err := cp.stopChargingConnector(c, reason)
			if err != nil {
				cp.logger.WithError(err).Errorf("Cannot stop the transaction at cleanup")
			}
		}
	}

	if !util.IsNilInterfaceOrPointer(cp.tagReader) {
		cp.logger.Info("Cleaning up the Tag Reader")
		cp.tagReader.Cleanup()
	}

	if !util.IsNilInterfaceOrPointer(cp.display) {
		cp.logger.Info("Cleaning up display")
		cp.display.Cleanup()
	}

	if !util.IsNilInterfaceOrPointer(cp.indicator) {
		cp.logger.Info("Cleaning up Indicator")
		cp.indicator.Cleanup()
	}

	close(cp.connectorChannel)
	cp.logger.Info("Clearing the scheduler...")
	cp.scheduler.Stop()
	cp.scheduler.Clear()

	_ = cp.tagManager.WriteLocalAuthList()

	cp.logger.Infof("Disconnecting the client..")
	cp.chargePoint.Stop()
}
