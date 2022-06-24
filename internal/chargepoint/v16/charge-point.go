package v16

import (
	"context"
	"github.com/go-co-op/gocron"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/reactivex/rxgo/v2"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/api"
	chargePointUtil "github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/util"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/auth"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/connector"
	connectorManager "github.com/xBlaz3kx/ChargePi-go/internal/components/connector-manager"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/hardware/reader"
	chargePoint "github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/util"
)

type (
	ChargePoint struct {
		chargePoint  ocpp16.ChargePoint
		availability core.AvailabilityType
		Settings     *settings.Settings
		// Hardware components
		TagReader reader.Reader
		Indicator indicator.Indicator
		LCD       display.LCD
		// Software components
		connectorManager   connectorManager.Manager
		connectorChannel   chan rxgo.Item
		meterValuesChannel chan rxgo.Item
		scheduler          *gocron.Scheduler
		authCache          *auth.Cache
		logger             *log.Logger
	}

	ChargePointV16 interface {
		chargePoint.ChargePoint
		notifyConnectorStatus(connector connector.Connector)
		startChargingConnector(connector connector.Connector, tagId string) error
		stopChargingConnector(connector connector.Connector, reason core.Reason) error
		displayConnectorStatus(connectorId int, status core.ChargePointStatus)
		restoreState()
	}

	Options func(point *ChargePoint)
)

// NewChargePoint creates a new ChargePoint for OCPP version 1.6.
func NewChargePoint(manager connectorManager.Manager, scheduler *gocron.Scheduler, cache *auth.Cache, opts ...Options) *ChargePoint {
	ch := make(chan rxgo.Item, 5)
	// Set the channel
	manager.SetNotificationChannel(ch)

	cp := &ChargePoint{
		availability:     core.AvailabilityTypeInoperative,
		connectorChannel: ch,
		scheduler:        scheduler,
		connectorManager: manager,
		authCache:        cache,
		logger:           log.StandardLogger(),
	}

	// Apply options
	for _, opt := range opts {
		opt(cp)
	}

	return cp
}

// Init initializes the charge point based on the settings and listens to the Reader.
func (cp *ChargePoint) Init(settings *settings.Settings) {
	if settings == nil {
		log.Fatal("no settings provided")
	}

	cp.Settings = settings

	var (
		info      = settings.ChargePoint.Info
		tlsConfig = settings.ChargePoint.TLS
		wsClient  = chargePointUtil.CreateClient(
			settings.ChargePoint.Info.BasicAuthUsername,
			settings.ChargePoint.Info.BasicAuthPassword,
			tlsConfig)
		logInfo = log.WithFields(log.Fields{
			"chargePointId": info.Id,
		})
	)

	logInfo.Debug("Creating charge point")
	cp.chargePoint = ocpp16.NewChargePoint(info.Id, nil, wsClient)

	// Set charging profiles
	chargePointUtil.SetProfilesFromConfig(cp.chargePoint, cp, cp, cp)

	cp.setMaxCachedTags()
}

// Connect to the central system and send a BootNotification
func (cp *ChargePoint) Connect(ctx context.Context, serverUrl string) {
	cp.logger.Infof("Trying to connect to the central system: %s", serverUrl)
	connectErr := cp.chargePoint.Start(serverUrl)

	// Check if the connection was successful
	if connectErr != nil {
		//cp.CleanUp(core.ReasonOther)
		cp.logger.WithError(connectErr).Fatalf("Cannot connect to the central system")
	}

	cp.logger.Infof("Successfully connected to: %s", serverUrl)
	cp.availability = core.AvailabilityTypeOperative

	go cp.ListenForConnectorStatusChange(ctx, cp.connectorChannel)
	cp.bootNotification()
}

// HandleChargingRequest Entry point for determining if the request is to start or stop charging. Trying to find a connector that has the tag stored in the Session; if such a connector exists,
// execute stopChargingConnector, otherwise startCharging.
func (cp *ChargePoint) HandleChargingRequest(tagId string) (*api.HandleChargingResponse, error) {
	var (
		err      error
		response = api.HandleChargingResponse{}
	)
	cp.logger.Infof("Handling request for tag %s", tagId)

	c := cp.connectorManager.FindConnectorWithTagId(tagId)
	if !util.IsNilInterfaceOrPointer(c) {
		err = cp.stopChargingConnector(c, core.ReasonLocal)
		if err != nil {
			cp.logger.WithError(err).Errorf("Error stopping charging the connector")
			response.ErrorMessage = err.Error()
		}

		response.ConnectorId = int32(c.GetConnectorId())

		return nil, err
	}

	err = cp.startCharging(tagId)
	if err != nil {
		cp.logger.WithError(err).Errorf("Cannot start charing the connector")
		return nil, err
	}

	return nil, err
}

// CleanUp When exiting the client, stop all the transactions, clean up all the peripherals and terminate the connection.
func (cp *ChargePoint) CleanUp(reason core.Reason) {
	cp.logger.Infof("Cleaning up ChargePoint, reason: %s", reason)

	switch reason {
	case core.ReasonRemote, core.ReasonLocal, core.ReasonHardReset, core.ReasonSoftReset:
		for _, c := range cp.connectorManager.GetConnectors() {
			// Stop charging the connectors
			err := cp.stopChargingConnector(c, reason)
			if err != nil {
				cp.logger.WithError(err).Errorf("Cannot stop the transaction at cleanup")
			}
		}
		break
	}

	cp.logger.Infof("Disconnecting the client..")
	cp.chargePoint.Stop()

	if !util.IsNilInterfaceOrPointer(cp.TagReader) {
		cp.logger.Info("Cleaning up the Tag Reader")
		cp.TagReader.Cleanup()
	}

	if !util.IsNilInterfaceOrPointer(cp.LCD) {
		cp.logger.Info("Cleaning up LCD")
		cp.LCD.Cleanup()
	}

	if !util.IsNilInterfaceOrPointer(cp.Indicator) {
		cp.logger.Info("Cleaning up Indicator")
		cp.Indicator.Cleanup()
	}

	close(cp.connectorChannel)
	cp.logger.Info("Clearing the scheduler...")
	cp.scheduler.Stop()
	cp.scheduler.Clear()

	cp.authCache.DumpTags()
}
