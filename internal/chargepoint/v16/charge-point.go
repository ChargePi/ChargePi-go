package v16

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/reactivex/rxgo/v2"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/auth"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/connector"
	connectorManager "github.com/xBlaz3kx/ChargePi-go/internal/components/connector-manager"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/hardware/reader"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/tls"
	"github.com/xBlaz3kx/ChargePi-go/pkg/util"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	v16 "github.com/xBlaz3kx/ocppManager-go/v16"
	"strings"
	"time"
)

type (
	ChargePoint struct {
		chargePoint      ocpp16.ChargePoint
		availability     core.AvailabilityType
		connectorManager connectorManager.Manager
		connectorChannel chan rxgo.Item
		Settings         *settings.Settings
		TagReader        reader.Reader
		Indicator        indicator.Indicator
		LCD              display.LCD
		scheduler        *gocron.Scheduler
		authCache        *auth.Cache
	}

	ChargePointV16 interface {
		chargepoint.ChargePoint
		notifyConnectorStatus(connector connector.Connector)
		startChargingConnector(connector connector.Connector, tagId string) error
		stopChargingConnector(connector connector.Connector, reason core.Reason) error
		displayConnectorStatus(connectorId int, status core.ChargePointStatus)
		restoreState()
	}
)

// NewChargePoint creates a new ChargePoint for OCPP version 1.6.
func NewChargePoint(tagReader reader.Reader, lcd display.LCD, manager connectorManager.Manager, scheduler *gocron.Scheduler, cache *auth.Cache) *ChargePoint {
	ch := make(chan rxgo.Item, 5)
	// Set the channel
	manager.SetNotificationChannel(ch)

	return &ChargePoint{
		TagReader:        tagReader,
		LCD:              lcd,
		availability:     core.AvailabilityTypeInoperative,
		connectorChannel: ch,
		scheduler:        scheduler,
		connectorManager: manager,
		authCache:        cache,
	}
}

// createClientConfiguration creates default configuration for websocket client
func createClientConfiguration() ws.ClientTimeoutConfig {
	var (
		clientConfig      = ws.NewClientTimeoutConfig()
		pingInterval, err = ocppConfigManager.GetConfigurationValue(v16.WebSocketPingInterval.String())
	)

	if err == nil {
		duration, err := time.ParseDuration(fmt.Sprintf("%ss", pingInterval))
		if err == nil {
			clientConfig.PingPeriod = duration
		}
	}

	return clientConfig
}

// Init initializes the charge point based on the settings and listens to the Reader.
func (cp *ChargePoint) Init(ctx context.Context, settings *settings.Settings) {
	if settings == nil {
		log.Fatal("no settings provided")
	}

	cp.Settings = settings

	var (
		client    ws.WsClient = ws.NewClient()
		info                  = settings.ChargePoint.Info
		tlsConfig             = settings.ChargePoint.TLS
		logInfo               = log.WithFields(log.Fields{
			"chargePointId": info.Id,
		})
	)

	// Check if the client has TLS
	if tlsConfig.IsEnabled {
		client = tls.GetTLSClient(tlsConfig.CACertificatePath, tlsConfig.ClientCertificatePath, tlsConfig.ClientKeyPath)
	}

	client.SetTimeoutConfig(createClientConfiguration())

	logInfo.Debug("Creating charge point")
	cp.chargePoint = ocpp16.NewChargePoint(info.Id, nil, client)

	// Set handlers based on configuration
	profiles, err := ocppConfigManager.GetConfigurationValue(v16.SupportedFeatureProfiles.String())
	if err != nil {
		log.WithError(err).Fatalf("No supported profiles specified")
	}

	for _, profile := range strings.Split(profiles, " ,") {
		switch strings.ToLower(profile) {
		case core.ProfileName:
			cp.chargePoint.SetCoreHandler(cp)
			break
		case reservation.ProfileName:
			cp.chargePoint.SetReservationHandler(cp)
			break
		case smartcharging.ProfileName:
			//cp.chargePoint.SetSmartChargingHandler(cp)
			break
		case localauth.ProfileName:
			//cp.chargePoint.SetLocalAuthListHandler(cp)
			break
		case remotetrigger.ProfileName:
			cp.chargePoint.SetRemoteTriggerHandler(cp)
			break
		case firmware.ProfileName:
			//cp.chargePoint.SetFirmwareManagementHandler(cp)
			break
		}
	}

	cp.setMaxCachedTags()

	// Start listening for tags from reader
	if cp.TagReader != nil && cp.Settings.ChargePoint.Hardware.TagReader.IsEnabled {
		go cp.ListenForTag(ctx, cp.TagReader.GetTagChannel())
	}
}

// Connect to the central system and send a BootNotification
func (cp *ChargePoint) Connect(ctx context.Context, serverUrl string) {
	log.Infof("Trying to connect to the central system: %s", serverUrl)
	connectErr := cp.chargePoint.Start(serverUrl)

	// Check if the connection was successful
	if connectErr != nil {
		//cp.CleanUp(core.ReasonOther)
		log.WithError(connectErr).Fatalf("Cannot connect to the central system")
	}

	log.Infof("Successfully connected to: %s", serverUrl)
	cp.availability = core.AvailabilityTypeOperative

	go cp.ListenForConnectorStatusChange(ctx, cp.connectorChannel)
	cp.bootNotification()
}

// HandleChargingRequest Entry point for determining if the request is to start or stop charging. Trying to find a connector that has the tag stored in the Session; if such a connector exists,
// execute stopChargingConnector, otherwise startCharging.
func (cp *ChargePoint) HandleChargingRequest(tagId string) {
	log.Infof("Handling request for tag %s", tagId)

	c := cp.connectorManager.FindConnectorWithTagId(tagId)
	if !util.IsNilInterfaceOrPointer(c) {
		err := cp.stopChargingConnector(c, core.ReasonLocal)
		if err != nil {
			log.WithError(err).Errorf("Error stopping charing the connector")
		}
		return
	}

	err := cp.startCharging(tagId)
	if err != nil {
		log.WithError(err).Errorf("Cannot start charing the connector")
	}
}

// CleanUp When exiting the client, stop all the transactions, clean up all the peripherals and terminate the connection.
func (cp *ChargePoint) CleanUp(reason core.Reason) {
	log.Infof("Cleaning up ChargePoint, reason: %s", reason)

	switch reason {
	case core.ReasonRemote, core.ReasonLocal, core.ReasonHardReset, core.ReasonSoftReset:
		for _, c := range cp.connectorManager.GetConnectors() {
			// Stop charging the connectors
			err := cp.stopChargingConnector(c, reason)
			if err != nil {
				log.WithError(err).Errorf("Cannot stop the transaction at cleanup")
			}
		}
		break
	}

	log.Infof("Disconnecting the client..")
	cp.chargePoint.Stop()

	if cp.TagReader != nil {
		log.Info("Cleaning up the Tag Reader")
		cp.TagReader.Cleanup()
	}

	if cp.LCD != nil {
		log.Info("Cleaning up LCD")
		cp.LCD.Cleanup()
	}

	if cp.Indicator != nil {
		log.Info("Cleaning up Indicator")
		cp.Indicator.Cleanup()
	}

	close(cp.connectorChannel)
	log.Info("Clearing the scheduler...")
	cp.scheduler.Stop()
	cp.scheduler.Clear()

	cp.authCache.DumpTags()
}
