package v16

import (
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ocppManager-go/ocpp_v16"
)

// sendRequest is a middleware function that implements a retry mechanism for sending requests. If the max attempts is reached, return an error
func (cp *ChargePoint) sendRequest(request ocpp.Request, callback func(confirmation ocpp.Response, protoError error)) error {
	variableManager := cp.settingsManager.GetOcppV16Manager()
	var (
		maxMessageAttempts, attemptErr  = variableManager.GetConfigurationValue(ocpp_v16.TransactionMessageAttempts)
		retryIntervalValue, intervalErr = variableManager.GetConfigurationValue(ocpp_v16.TransactionMessageRetryInterval)
		maxRetries, convError           = strconv.Atoi(*maxMessageAttempts)
		retryInterval, convError2       = strconv.Atoi(*retryIntervalValue)
	)

	if attemptErr != nil || convError != nil {
		maxRetries = 5
	}

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

// SetProfilesFromConfig based on the provided OCPP configuration, set the profiles
func (cp *ChargePoint) SetProfilesFromConfig() {
	cp.chargePoint.SetCoreHandler(cp)

	// Set handlers based on configuration
	profiles, err := cp.settingsManager.GetOcppV16Manager().GetConfigurationValue(ocpp_v16.SupportedFeatureProfiles)
	if err != nil {
		log.WithError(err).Panic("No supported profiles specified")
	}

	logInfo := log.WithField("profiles", profiles)

	for _, profile := range strings.Split(*profiles, ", ") {
		switch strings.ToLower(profile) {
		case strings.ToLower(reservation.ProfileName):
			cp.chargePoint.SetReservationHandler(cp)
			logInfo.Debug("Setting reservation handler")
		case strings.ToLower(smartcharging.ProfileName):
			logInfo.Debug("Setting local auth handler")
			// chargePoint.SetSmartChargingHandler(cp)
		case strings.ToLower(localauth.ProfileName):
			logInfo.Debug("Setting local auth handler")
			cp.chargePoint.SetLocalAuthListHandler(cp)
		case strings.ToLower(remotetrigger.ProfileName):
			logInfo.Debug("Setting remote trigger handler")
			cp.chargePoint.SetRemoteTriggerHandler(cp)
		case strings.ToLower(firmware.ProfileName):
			cp.chargePoint.SetFirmwareManagementHandler(cp)
		}
	}
}

func (cp *ChargePoint) handleRequestErr(err error, text string) {
	if err != nil {
		cp.logger.WithError(err).Errorf(text)
	}
}
