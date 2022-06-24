package util

import (
	"fmt"
	"github.com/agrison/go-commons-lang/stringUtils"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ws"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/tls"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	v16 "github.com/xBlaz3kx/ocppManager-go/v16"
	"strings"
	"time"
)

// CreateConnectionUrl creates a connection url from the provided settings
func CreateConnectionUrl(point settings.ChargePoint) string {
	var (
		chargePointInfo = point.Info
		serverUrl       = fmt.Sprintf("ws://%s", chargePointInfo.ServerUri)
	)

	// Replace insecure Websockets
	if point.TLS.IsEnabled {
		serverUrl = strings.Replace(serverUrl, "ws", "wss", 1)
	}

	return serverUrl
}

// CreateClient creates a Websocket client based on the settings.
func CreateClient(basicAuthUser, basicAuthPass string, tlsConfig settings.TLS) *ws.Client {
	var (
		client            = ws.NewClient()
		clientConfig      = ws.NewClientTimeoutConfig()
		pingInterval, err = ocppConfigManager.GetConfigurationValue(v16.WebSocketPingInterval.String())
	)

	if err == nil {
		duration, err := time.ParseDuration(fmt.Sprintf("%ss", pingInterval))
		if err == nil {
			clientConfig.PingPeriod = duration
		}
	}

	// Check if the client has TLS
	if tlsConfig.IsEnabled {
		client = tls.GetTLSClient(tlsConfig.CACertificatePath, tlsConfig.ClientCertificatePath, tlsConfig.ClientKeyPath)
	}

	// If HTTP basic auth is provided, set it in the Websocket client
	if stringUtils.IsNoneEmpty(basicAuthUser, basicAuthPass) {
		client.SetBasicAuth(basicAuthUser, basicAuthPass)
	}

	client.SetTimeoutConfig(clientConfig)
	return client
}

// SetProfilesFromConfig based on the provided OCPP configuration, set the profiles
func SetProfilesFromConfig(
	chargePoint ocpp16.ChargePoint,
	coreHandler core.ChargePointHandler,
	reservationHandler reservation.ChargePointHandler,
	triggerHandler remotetrigger.ChargePointHandler,
) {
	// Set handlers based on configuration
	profiles, err := ocppConfigManager.GetConfigurationValue(v16.SupportedFeatureProfiles.String())
	if err != nil {
		log.WithError(err).Fatalf("No supported profiles specified")
	}

	for _, profile := range strings.Split(profiles, ", ") {
		switch strings.ToLower(profile) {
		case strings.ToLower(core.ProfileName):
			chargePoint.SetCoreHandler(coreHandler)
			log.Debug("Setting core handler")
			break
		case strings.ToLower(reservation.ProfileName):
			chargePoint.SetReservationHandler(reservationHandler)
			log.Debug("Setting reservation handler")
			break
		case strings.ToLower(smartcharging.ProfileName):
			//chargePoint.SetSmartChargingHandler(cp)
			break
		case strings.ToLower(localauth.ProfileName):
			//chargePoint.SetLocalAuthListHandler(cp)
			break
		case strings.ToLower(remotetrigger.ProfileName):
			log.Debug("Setting remote trigger handler")
			chargePoint.SetRemoteTriggerHandler(triggerHandler)
			break
		case strings.ToLower(firmware.ProfileName):
			//chargePoint.SetFirmwareManagementHandler(cp)
			break
		}
	}
}
