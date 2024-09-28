package v16

import (
	"context"
	"go.uber.org/zap"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ChargePi/ChargePi-go/internal/auth"
	chargePoint "github.com/ChargePi/ChargePi-go/internal/chargepoint"
	"github.com/ChargePi/ChargePi-go/internal/diagnostics"
	"github.com/ChargePi/ChargePi-go/internal/evse/manager"
	"github.com/ChargePi/ChargePi-go/internal/networking"
	settings "github.com/ChargePi/ChargePi-go/internal/pkg/configuration/manager"
	"github.com/ChargePi/ChargePi-go/internal/pkg/notifications"
	"github.com/ChargePi/ChargePi-go/internal/pkg/scheduler"
	"github.com/ChargePi/ChargePi-go/internal/sessions"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/display"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/indicator"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/reader"
	"github.com/ChargePi/ChargePi-go/pkg/util"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/avast/retry-go"
	"github.com/go-co-op/gocron"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ChargePoint struct {
	chargePoint ocpp16.ChargePoint
	// State
	state State

	// Settings management
	settingsManager    settings.Manager
	info               chargePoint.Info
	connectionSettings chargePoint.ConnectionSettings
	modem              networking.Modem

	// Hardware components
	tagReader        reader.Reader
	indicator        indicator.Indicator
	display          display.Display
	indicatorMapping indicator.StatusMapping

	// Software components
	evseManager        manager.Manager
	sessionService     sessions.Service
	tagAuthService     auth.Service
	diagnosticsService diagnostics.Service

	meterValuesChannel chan notifications.MeterValueNotification
	scheduler          *gocron.Scheduler
	logger             *zap.Logger
}

// NewChargePoint creates a new ChargePoint for OCPP version 1.6.
func NewChargePoint(logger *zap.Logger, evseManager manager.Manager, settingsManager settings.Manager, tagManager auth.Service, sessionManager sessions.Service, diagnosticsManager diagnostics.Service, opts ...chargePoint.Options) (*ChargePoint, error) {
	cp := &ChargePoint{
		state: State{
			availability: core.AvailabilityTypeInoperative,
			isConnected:  false,
		},
		scheduler:          scheduler.NewScheduler(),
		evseManager:        evseManager,
		tagAuthService:     tagManager,
		sessionService:     sessionManager,
		settingsManager:    settingsManager,
		diagnosticsService: diagnosticsManager,
		logger:             logger.Named("charge_point_v16"),
	}

	err := cp.ApplyOpts(opts...)
	if err != nil {
		return nil, err
	}

	return cp, nil
}

// ApplyOpts applies a list of options to the charge point. Returns an error if any of the options fails to apply.
func (cp *ChargePoint) ApplyOpts(opts ...chargePoint.Options) error {
	// Apply options
	for _, opt := range opts {
		err := opt(cp)
		if err != nil {
			cp.logger.WithError(err).Error("Error applying option")
			return err
		}
	}

	return nil
}

// Connect to the central system and send a BootNotification
func (cp *ChargePoint) Connect(ctx context.Context, serverUrl string) error {
	// Disconnect if already connected to a central system, might be a reconnection or a new URL
	if cp.state.IsConnected() {
		cp.chargePoint.Stop()
	}

	logger := cp.logger.With(zap.String("serverUrl", serverUrl), zap.String("id", cp.connectionSettings.Id))

	// Get the ping interval from the configuration
	pingInterval, _ := cp.settingsManager.GetConfigurationValue(ocpp_v16.WebSocketPingInterval)

	// Create a new websocket client
	wsClient, err := chargePoint.CreateClient(cp.logger, cp.connectionSettings, pingInterval)
	if err != nil {
		return errors.Wrap(err, "cannot create a new websocket client")
	}

	// Create a new OCPP charge point handler
	cp.chargePoint = ocpp16.NewChargePoint(cp.connectionSettings.Id, nil, wsClient)

	// Set profiles
	err = cp.setProfilesFromConfig()
	if err != nil {
		return err
	}

	cp.logger.Info("Trying to connect to the central system")
	connectErr := cp.chargePoint.Start(serverUrl)
	if connectErr != nil {
		// cp.Cleanup(core.ReasonOther)
		cp.state.SetConnected(false)
		return connectErr
	}

	cp.logger.Info("Successfully connected to backend")
	cp.bootNotification()
	return nil
}

// Cleanup When exiting the client, stop all the transactions, clean up all the peripherals and terminate the connection.
func (cp *ChargePoint) Cleanup(reason core.Reason) error {
	cp.logger.With(zap.Any("reason", reason)).Info("Cleaning up ChargePoint")

	if !util.IsNilInterfaceOrPointer(cp.tagReader) {
		cp.tagReader.Cleanup()
	}

	if !util.IsNilInterfaceOrPointer(cp.display) {
		cp.logger.Debug("Cleaning up display")
		err := cp.display.Cleanup(nil)
		if err != nil {
			cp.logger.With(zap.Error(err)).Error("Error cleaning up display")
		}
	}

	if !util.IsNilInterfaceOrPointer(cp.indicator) {
		cp.logger.Debug("Cleaning up settings")
		cp.indicator.Cleanup()
	}

	cp.logger.Debug("Clearing the scheduler")
	cp.scheduler.Stop()
	cp.scheduler.Clear()

	cp.logger.Info("Disconnecting the client..")
	cp.chargePoint.Stop()

	close(cp.meterValuesChannel)
	return nil
}

// GetVersion get the firmware version of the charge point
func (cp *ChargePoint) GetVersion() string {
	return chargePoint.FirmwareVersion
}

// GetStatus of the charge point
func (cp *ChargePoint) GetStatus() string {
	return string(cp.state.GetAvailability())
}

// IsConnected of the charge point
func (cp *ChargePoint) IsConnected() bool {
	return cp.state.IsConnected()
}

func (cp *ChargePoint) SetLogger(logger *zap.Logger) {
	cp.logger = logger
}

// Reset the charge point.
func (cp *ChargePoint) Reset(resetType string) error {
	cp.logger.Info("Resetting the charge point")

	// Todo check if conditions are met
	var err error
	switch resetType {
	case string(core.ResetTypeHard):
		// Stop the client gracefully first
		_, err = cp.scheduler.Every(3).Seconds().LimitRunsTo(1).Do(cp.Cleanup, core.ReasonHardReset)
		if err != nil {
			return err
		}

		if !util.IsRunningInContainer() {
			// Schedule a reboot
			_, err = cp.scheduler.Every(10).Seconds().LimitRunsTo(1).Do(exec.Command, "sudo reboot")
		}
	case string(core.ResetTypeSoft):
		_, err = cp.scheduler.Every(3).Seconds().LimitRunsTo(1).Do(cp.Cleanup, core.ReasonSoftReset)
	}

	return err
}

func (cp *ChargePoint) OnReset(request *core.ResetRequest) (confirmation *core.ResetConfirmation, err error) {
	cp.logger.Sugar().Infof("Received request %s", request.GetFeatureName())
	var response = core.ResetStatusRejected
	var retries = 3

	resetRetries, _ := cp.settingsManager.GetConfigurationValue(ocpp_v16.ResetRetries)
	if resetRetries != nil {
		aRetries, err := strconv.Atoi(*resetRetries)
		if err == nil {
			retries = aRetries
		}
	}

	resetErr := retry.Do(
		func() error {
			return cp.Reset(string(request.Type))
		},
		retry.Attempts(uint(retries)),
		retry.Delay(time.Second*10),
	)
	if resetErr == nil {
		response = core.ResetStatusAccepted
	}

	return core.NewResetConfirmation(response), nil
}

func (cp *ChargePoint) OnChangeAvailability(request *core.ChangeAvailabilityRequest) (confirmation *core.ChangeAvailabilityConfirmation, err error) {
	cp.logger.Sugar().Infof("Received request %s", request.GetFeatureName())
	response := core.AvailabilityStatusRejected

	// This would mean a request to change the availability of the whole charge point
	if request.ConnectorId == 0 {

		// Try to set the availability of the charge point, if it fails, schedule the change
		err := cp.SetAvailability(request.Type)
		if err != nil {
			_, err := cp.scheduler.Every(3).Seconds().LimitRunsTo(1).Do(cp.SetAvailability, request.Type)
			if err != nil {
				return core.NewChangeAvailabilityConfirmation(core.AvailabilityStatusScheduled), nil
			}
		}

		return core.NewChangeAvailabilityConfirmation(core.AvailabilityStatusAccepted), nil
	}

	// Check if there are ongoing transactions, schedule the change if there are
	_, sessionErr := cp.sessionService.GetSession(request.ConnectorId, nil)
	switch sessionErr {
	case nil:
		response = core.AvailabilityStatusScheduled
		// todo set evse availability
		// cp.evseManager.Get
	default:
		cp.logger.With(zap.Error(sessionErr)).Error("Error checking for ongoing transactions")
	}

	return core.NewChangeAvailabilityConfirmation(response), nil
}

// SetAvailability sets the availability status of the charge point.
func (cp *ChargePoint) SetAvailability(availabilityType core.AvailabilityType) error {
	cp.logger.With(zap.String("availability", string(availabilityType))).Debug("Setting availability")

	// Check if there are ongoing transactions
	_, sessionErr := cp.sessionService.GetSession(0, nil)
	if sessionErr != nil {
		return errors.Wrap(sessionErr, "error checking for ongoing transactions")
	}

	// Set availability status and notify the backend
	cp.state.SetAvailability(availabilityType)

	switch availabilityType {
	case core.AvailabilityTypeInoperative:
		cp.notifyStatus(0, core.ChargePointStatusUnavailable, core.NoError)
	case core.AvailabilityTypeOperative:
		cp.notifyStatus(0, core.ChargePointStatusAvailable, core.NoError)
	}
	return nil
}

func (cp *ChargePoint) Pass() bool {
	// todo check EVSEs status
	// cp.evseManager
	// todo check components' status
	cp.sessionService.Pass()
	cp.diagnosticsService.Pass()

	return true
}

func (cp *ChargePoint) Name() string {
	return "charge-point-v16"
}

// sendRequest is a middleware function that implements a retry mechanism for sending requests. If the max attempts is reached, return an error
func (cp *ChargePoint) sendRequest(request ocpp.Request, callback func(confirmation ocpp.Response, protoError error)) error {
	maxMessageAttempts, attemptErr := cp.settingsManager.GetConfigurationValue(ocpp_v16.TransactionMessageAttempts)
	maxRetries, convError := strconv.Atoi(*maxMessageAttempts)
	if attemptErr != nil || convError != nil {
		maxRetries = 5
	}

	retryIntervalValue, intervalErr := cp.settingsManager.GetConfigurationValue(ocpp_v16.TransactionMessageRetryInterval)
	retryInterval, convError2 := strconv.Atoi(*retryIntervalValue)
	if intervalErr != nil || convError2 != nil {
		retryInterval = 30
	}

	return retry.Do(
		func() error {
			return cp.chargePoint.SendRequestAsync(request, callback)
		},
		retry.Attempts(uint(maxRetries)),
		retry.Delay(time.Duration(retryInterval)),
	)
}

// setProfilesFromConfig based on the provided OCPP configuration, set the profiles.
func (cp *ChargePoint) setProfilesFromConfig() error {
	cp.chargePoint.SetCoreHandler(cp)

	// Setup core configuration validation
	err := cp.setupCoreConfigurationValidation()
	if err != nil {
		return errors.Wrap(err, "unable to setup configuration validation")
	}

	// Setup custom configuration variables validation
	err = cp.setupCustomConfigurationValidation()
	if err != nil {
		return errors.Wrap(err, "unable to setup configuration validation")
	}

	// Set handlers based on configuration
	profiles, err := cp.settingsManager.GetConfigurationValue(ocpp_v16.SupportedFeatureProfiles)
	if err != nil {
		return errors.Wrap(err, "unable to retrieve supported profiles")
	}

	logInfo := log.WithField("profiles", profiles)

	for _, profile := range strings.Split(*profiles, ", ") {
		switch profile {
		case reservation.ProfileName:
			cp.chargePoint.SetReservationHandler(cp)
			logInfo.Debug("Setting reservation handler")
		case smartcharging.ProfileName:
			logInfo.Debug("Setting smart charging handler")
			err = cp.setupSmartChargingConfigurationValidation()
			// cp.chargePoint.SetSmartChargingHandler(cp)
		case localauth.ProfileName:
			logInfo.Debug("Setting local auth handler")
			cp.chargePoint.SetLocalAuthListHandler(cp)

			err = cp.setupCoreConfigurationValidation()
		case remotetrigger.ProfileName:
			logInfo.Debug("Setting remote trigger handler")
			cp.chargePoint.SetRemoteTriggerHandler(cp)
		case firmware.ProfileName:
			logInfo.Debug("Setting firmware handler")
			cp.chargePoint.SetFirmwareManagementHandler(cp)
		default:
			return errors.New("unsupported profile")
		}

		if err != nil {
			return errors.Wrap(err, "unable to setup profile configuration validation")
		}
	}

	return nil
}

func (cp *ChargePoint) handleRequestErr(err error, text string) {
	if err != nil {
		cp.logger.WithError(err).Errorf(text)
	}
}
