package v16

import (
	"fmt"
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

func createClient(tlsConfig settings.TLS) *ws.Client {
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

	client.SetTimeoutConfig(clientConfig)
	return client
}

func setProfilesFromConfig(chargePoint ocpp16.ChargePoint, coreHandler core.ChargePointHandler,
	reservationHandler reservation.ChargePointHandler, triggerHandler remotetrigger.ChargePointHandler) {
	// Set handlers based on configuration
	profiles, err := ocppConfigManager.GetConfigurationValue(v16.SupportedFeatureProfiles.String())
	if err != nil {
		log.WithError(err).Fatalf("No supported profiles specified")
	}

	for _, profile := range strings.Split(profiles, " ,") {
		switch strings.ToLower(profile) {
		case core.ProfileName:
			chargePoint.SetCoreHandler(coreHandler)
			break
		case reservation.ProfileName:
			chargePoint.SetReservationHandler(reservationHandler)
			break
		case smartcharging.ProfileName:
			//cp.chargePoint.SetSmartChargingHandler(cp)
			break
		case localauth.ProfileName:
			//cp.chargePoint.SetLocalAuthListHandler(cp)
			break
		case remotetrigger.ProfileName:
			chargePoint.SetRemoteTriggerHandler(triggerHandler)
			break
		case firmware.ProfileName:
			//cp.chargePoint.SetFirmwareManagementHandler(cp)
			break
		}
	}
}
