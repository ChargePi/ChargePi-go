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
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
	settings2 "github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"os/exec"
)

type (
	ChargePoint struct {
		chargePoint        ocpp16.ChargePoint
		availability       core.AvailabilityType
		info               settings2.Info
		connectionSettings settings2.ConnectionSettings
		// Hardware components
		tagReader        reader.Reader
		indicator        indicator.Indicator
		display          display.Display
		indicatorMapping settings2.IndicatorStatusMapping
		// Software components
		connectorManager connectorManager.Manager

		meterValuesChannel chan notifications.MeterValueNotification
		scheduler          *gocron.Scheduler
		tagManager         auth.TagManager
		logger             *log.Logger
	}
)

// NewChargePoint creates a new ChargePoint for OCPP version 1.6.
func NewChargePoint(manager connectorManager.Manager, scheduler *gocron.Scheduler, cache auth.TagManager, opts ...chargePoint.Options) *ChargePoint {
	cp := &ChargePoint{
		availability:     core.AvailabilityTypeInoperative,
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
		connectionSettings = cp.connectionSettings
		tlsConfig          = connectionSettings.TLS
		wsClient           = util.CreateClient(
			connectionSettings.BasicAuthUsername,
			connectionSettings.BasicAuthPassword,
			tlsConfig)
		logInfo = log.WithFields(log.Fields{
			"chargePointId": connectionSettings.Id,
		})
	)

	logInfo.Debug("Creating an ocpp connection")
	cp.chargePoint = ocpp16.NewChargePoint(connectionSettings.Id, nil, wsClient)

	// Set charging profiles
	util.SetProfilesFromConfig(cp.chargePoint, cp, cp, cp)

	cp.logger.Infof("Trying to connect to the central system: %s", serverUrl)
	connectErr := cp.chargePoint.Start(serverUrl)
	if connectErr != nil {
		// cp.CleanUp(core.ReasonOther)
		cp.logger.WithError(connectErr).Fatalf("Cannot connect to the central system")
	}

	cp.logger.Infof("Successfully connected to: %s", serverUrl)
	cp.availability = core.AvailabilityTypeOperative

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

	cp.logger.Info("Clearing the scheduler...")
	cp.scheduler.Stop()
	cp.scheduler.Clear()

	// Persist tags
	_ = cp.tagManager.WriteLocalAuthList()

	cp.logger.Infof("Disconnecting the client..")
	cp.chargePoint.Stop()
}

// Reset the charge point.
func (cp *ChargePoint) Reset(resetType string) error {
	cp.logger.Infof("Resetting the charge point")

	// Todo check if conditions are met

	switch resetType {
	case string(core.ResetTypeHard):
		_, err := cp.scheduler.Every(3).Seconds().LimitRunsTo(1).Do(cp.CleanUp, core.ReasonHardReset)
		if err != nil {
			return err
		}

		_, err = cp.scheduler.Every(30).Seconds().LimitRunsTo(1).Do(exec.Command, "sudo reboot")
	case string(core.ResetTypeSoft):
		_, err := cp.scheduler.Every(3).Seconds().LimitRunsTo(1).Do(cp.CleanUp, core.ReasonSoftReset)
		if err != nil {
			return err
		}
	}

	return nil
}
