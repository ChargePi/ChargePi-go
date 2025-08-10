package v16

import (
	"context"
	"github.com/ChargePi/ChargePi-go/internal/evse/manager"
	"go.uber.org/zap"
	"os/exec"

	"github.com/ChargePi/ChargePi-go/internal/auth"
	"github.com/ChargePi/ChargePi-go/internal/diagnostics"
	"github.com/ChargePi/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/ChargePi/ChargePi-go/internal/pkg/models/notifications"
	"github.com/ChargePi/ChargePi-go/internal/pkg/models/settings"
	"github.com/ChargePi/ChargePi-go/internal/pkg/scheduler"
	settings2 "github.com/ChargePi/ChargePi-go/internal/pkg/settings"
	"github.com/ChargePi/ChargePi-go/internal/pkg/util"
	"github.com/ChargePi/ChargePi-go/internal/sessions/service/session"
	"github.com/ChargePi/ChargePi-go/pkg/display"
	"github.com/ChargePi/ChargePi-go/pkg/indicator"
	hardwareSettings "github.com/ChargePi/ChargePi-go/pkg/models/settings"
	"github.com/ChargePi/ChargePi-go/pkg/reader"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/go-co-op/gocron"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)

type ChargePoint struct {
	// OCPP related
	chargePoint  ocpp16.ChargePoint
	availability core.AvailabilityType
	isConnected  bool

	// Settings management
	settingsManager    settings2.Manager
	info               settings.Info
	connectionSettings settings.ConnectionSettings

	// Hardware components
	tagReader        reader.Reader
	indicator        indicator.Indicator
	display          display.Display
	indicatorMapping hardwareSettings.IndicatorStatusMapping

	// Software components
	evseManager        manager.Manager
	sessionManager     session.Manager
	diagnosticsManager diagnostics.Manager
	meterValuesChannel chan notifications.MeterValueNotification
	scheduler          *gocron.Scheduler
	tagManager         auth.Manager
	logger             *zap.Logger
}

// NewChargePoint creates a new ChargePoint for OCPP version 1.6.
func NewChargePoint(logger *zap.Logger, manager manager.Manager, tagManager auth.Manager, sessionManager session.Manager, diagnosticsManager diagnostics.Manager, opts ...chargePoint.Options) *ChargePoint {
	cp := &ChargePoint{
		availability:       core.AvailabilityTypeOperative,
		scheduler:          scheduler.NewScheduler(),
		evseManager:        manager,
		tagManager:         tagManager,
		sessionManager:     sessionManager,
		settingsManager:    settings2.GetManager(),
		diagnosticsManager: diagnosticsManager,
		logger:             logger.Named("charge_point_v16"),
	}

	cp.ApplyOpts(opts...)

	return cp
}

// Connect to the central system and send a BootNotification
func (cp *ChargePoint) Connect(ctx context.Context, serverUrl string) {
	// Disconnect if already connected to a central system, might be a reconnection or a new URL
	if cp.isConnected {
		cp.chargePoint.Stop()
	}

	logger := cp.logger.With(zap.String("serverUrl", serverUrl), zap.String("id", cp.connectionSettings.Id))

	// Get the ping interval from the configuration
	pingInterval, _ := cp.settingsManager.GetOcppV16Manager().GetConfigurationValue(ocpp_v16.WebSocketPingInterval)

	// Create a new websocket client
	wsClient, err := util.CreateClient(cp.logger, cp.connectionSettings, pingInterval)
	if err != nil {
		logger.With(zap.Error(err)).Panic("Cannot create a new websocket client")
	}

	// Create a new OCPP charge point handler
	cp.chargePoint = ocpp16.NewChargePoint(cp.connectionSettings.Id, nil, wsClient)

	// Set profiles
	cp.SetProfilesFromConfig()

	cp.logger.Info("Trying to connect to the central system")
	connectErr := cp.chargePoint.Start(serverUrl)
	if connectErr != nil {
		// cp.CleanUp(core.ReasonOther)
		cp.isConnected = false
		cp.logger.With(zap.Error(err)).Error("Cannot connect to the central system")
		return
	}

	cp.logger.Info("Successfully connected to backend")
	cp.bootNotification()
}

// CleanUp When exiting the client, stop all the transactions, clean up all the peripherals and terminate the connection.
func (cp *ChargePoint) CleanUp(reason core.Reason) {
	cp.logger.With(zap.Any("reason", reason)).Info("Cleaning up ChargePoint")

	switch reason {
	case core.ReasonRemote, core.ReasonLocal, core.ReasonHardReset, core.ReasonSoftReset:
		for _, c := range cp.evseManager.GetEVSEs() {
			// Stop charging the connectors
			err := cp.stopChargingConnector(c, reason)
			if err != nil {
				cp.logger.With(zap.Error(err)).Error("Cannot stop the transaction at cleanup")
			}
		}
	}

	if !util.IsNilInterfaceOrPointer(cp.tagReader) {
		cp.logger.Debug("Cleaning up the Tag Reader")
		cp.tagReader.Cleanup()
	}

	if !util.IsNilInterfaceOrPointer(cp.display) {
		cp.logger.Debug("Cleaning up display")
		cp.display.Cleanup()
	}

	if !util.IsNilInterfaceOrPointer(cp.indicator) {
		cp.logger.Debug("Cleaning up Indicator")
		cp.indicator.Cleanup()
	}

	cp.logger.Debug("Clearing the scheduler...")
	cp.scheduler.Stop()
	cp.scheduler.Clear()

	cp.logger.Info("Disconnecting the client..")
	cp.chargePoint.Stop()
}

// Reset the charge point.
func (cp *ChargePoint) Reset(resetType string) error {
	cp.logger.Info("Resetting the charge point")

	// Todo check if conditions are met
	var err error
	switch resetType {
	case string(core.ResetTypeHard):
		_, err = cp.scheduler.Every(3).Seconds().LimitRunsTo(1).Do(cp.CleanUp, core.ReasonHardReset)
		if err != nil {
			return err
		}

		_, err = cp.scheduler.Every(30).Seconds().LimitRunsTo(1).Do(exec.Command, "sudo reboot")
	case string(core.ResetTypeSoft):
		_, err = cp.scheduler.Every(3).Seconds().LimitRunsTo(1).Do(cp.CleanUp, core.ReasonSoftReset)
	}

	return err
}

// GetVersion get the firmware version of the charge point
func (cp *ChargePoint) GetVersion() string {
	return chargePoint.FirmwareVersion
}

// GetStatus of the charge point
func (cp *ChargePoint) GetStatus() string {
	return string(cp.availability)
}

// IsConnected of the charge point
func (cp *ChargePoint) IsConnected() bool {
	return cp.isConnected
}
