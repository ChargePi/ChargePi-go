package v16

import (
	"context"
	"os/exec"

	"github.com/go-co-op/gocron"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/auth"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/reader"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/scheduler"
	settings2 "github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"github.com/xBlaz3kx/ChargePi-go/internal/sessions/service/session"
)

type ChargePoint struct {
	chargePoint        ocpp16.ChargePoint
	availability       core.AvailabilityType
	isConnected        bool
	settingsManager    settings2.Manager
	info               settings.Info
	connectionSettings settings.ConnectionSettings
	// Hardware components
	tagReader        reader.Reader
	indicator        indicator.Indicator
	display          display.Display
	indicatorMapping settings.IndicatorStatusMapping
	// Software components
	evseManager        evse.Manager
	sessionManager     session.Manager
	meterValuesChannel chan notifications.MeterValueNotification
	scheduler          *gocron.Scheduler
	tagManager         auth.TagManager
	logger             *log.Logger
}

// NewChargePoint creates a new ChargePoint for OCPP version 1.6.
func NewChargePoint(manager evse.Manager, tagManager auth.TagManager, sessionManager session.Manager, opts ...chargePoint.Options) *ChargePoint {
	cp := &ChargePoint{
		availability:    core.AvailabilityTypeInoperative,
		scheduler:       scheduler.NewScheduler(),
		evseManager:     manager,
		tagManager:      tagManager,
		sessionManager:  sessionManager,
		settingsManager: settings2.GetManager(),
		logger:          log.StandardLogger(),
	}

	// Set profiles
	util.SetProfilesFromConfig(cp.chargePoint, cp, cp, cp, cp)

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

	cp.logger.Infof("Trying to connect to the central system: %s", serverUrl)
	connectErr := cp.chargePoint.Start(serverUrl)
	if connectErr != nil {
		// cp.CleanUp(core.ReasonOther)
		cp.isConnected = false
		cp.logger.WithError(connectErr).Panic("Cannot connect to the central system")
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
		for _, c := range cp.evseManager.GetEVSEs() {
			// Stop charging the connectors
			err := cp.stopChargingConnector(c, reason)
			if err != nil {
				cp.logger.WithError(err).Errorf("Cannot stop the transaction at cleanup")
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

	cp.logger.Infof("Disconnecting the client..")
	cp.chargePoint.Stop()
}

// Reset the charge point.
func (cp *ChargePoint) Reset(resetType string) error {
	cp.logger.Infof("Resetting the charge point")

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
