package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/avast/retry-go"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ws"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/util/tls"
	ocppManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

// CreateConnectionUrl creates a connection url from the provided settings
func CreateConnectionUrl(connectionSettings settings.ConnectionSettings) string {
	var (
		serverUrl = fmt.Sprintf("ws://%s", connectionSettings.ServerUri)
	)

	// Replace insecure Websockets
	if connectionSettings.TLS.IsEnabled {
		serverUrl = strings.Replace(serverUrl, "ws", "wss", 1)
	}

	return serverUrl
}

// CreateClient creates a Websocket client based on the settings.
func CreateClient(basicAuthUser, basicAuthPass string, tlsConfig settings.TLS) (*ws.Client, error) {
	client := ws.NewClient()
	clientConfig := ws.NewClientTimeoutConfig()

	// Set the ping interval
	pingInterval, err := ocppManager.GetConfigurationValue(configuration.WebSocketPingInterval.String())
	if err == nil {
		duration, err := time.ParseDuration(fmt.Sprintf("%ss", *pingInterval))
		if err == nil {
			clientConfig.PingPeriod = duration
		}
	}

	// Check if the TLS is enabled for the client
	if tlsConfig.IsEnabled {

		client, err = tls.CreateWssClient(tlsConfig.CACertificatePath, tlsConfig.ClientCertificatePath, tlsConfig.PrivateKeyPath)
		if err != nil {
			log.WithError(err).Error("Couldn't create TLS client")
			return nil, err
		}
	}

	// If HTTP basic auth is provided, set it in the Websocket client
	if stringUtils.IsNoneEmpty(basicAuthUser, basicAuthPass) {
		client.SetBasicAuth(basicAuthUser, basicAuthPass)
	}

	client.SetTimeoutConfig(clientConfig)
	return client, nil
}

// SetProfilesFromConfig based on the provided OCPP configuration, set the profiles
func SetProfilesFromConfig(
	chargePoint ocpp16.ChargePoint,
	coreHandler core.ChargePointHandler,
	reservationHandler reservation.ChargePointHandler,
	triggerHandler remotetrigger.ChargePointHandler,
	localAuth localauth.ChargePointHandler,
) {
	chargePoint.SetCoreHandler(coreHandler)

	// Set handlers based on configuration
	profiles, err := ocppManager.GetConfigurationValue(configuration.SupportedFeatureProfiles.String())
	if err != nil {
		log.WithError(err).Panic("No supported profiles specified")
	}

	logInfo := log.WithField("profiles", profiles)

	for _, profile := range strings.Split(*profiles, ", ") {
		switch strings.ToLower(profile) {
		case strings.ToLower(reservation.ProfileName):
			chargePoint.SetReservationHandler(reservationHandler)
			logInfo.Debug("Setting reservation handler")
		case strings.ToLower(smartcharging.ProfileName):
			logInfo.Debug("Setting local auth handler")
			// chargePoint.SetSmartChargingHandler(cp)
		case strings.ToLower(localauth.ProfileName):
			logInfo.Debug("Setting local auth handler")
			chargePoint.SetLocalAuthListHandler(localAuth)
		case strings.ToLower(remotetrigger.ProfileName):
			logInfo.Debug("Setting remote trigger handler")
			chargePoint.SetRemoteTriggerHandler(triggerHandler)
		case strings.ToLower(firmware.ProfileName):
			// chargePoint.SetFirmwareManagementHandler(cp)
		}
	}
}

// GetTypesToSample get the measurands to sample from the OCPP configuration.
func GetTypesToSample() []types.Measurand {
	var (
		measurands []types.Measurand
		// Get the types to sample
		measurandsString, err = ocppManager.GetConfigurationValue(configuration.MeterValuesSampledData.String())
	)

	if err != nil {
		return measurands
	}

	for _, measurand := range strings.Split(*measurandsString, ",") {
		measurands = append(measurands, types.Measurand(measurand))
	}

	return measurands
}

// SendRequest is a middleware function that implements a retry mechanism for sending requests. If the max attempts is reached, return an error
func SendRequest(chargePoint ocpp16.ChargePoint, request ocpp.Request, callback func(confirmation ocpp.Response, protoError error)) error {
	var (
		maxMessageAttempts, attemptErr  = ocppManager.GetConfigurationValue(configuration.TransactionMessageAttempts.String())
		retryIntervalValue, intervalErr = ocppManager.GetConfigurationValue(configuration.TransactionMessageRetryInterval.String())
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
			return chargePoint.SendRequestAsync(request, callback)
		},
		retry.Attempts(uint(maxRetries)),
		retry.Delay(time.Duration(retryInterval)),
	)
}
