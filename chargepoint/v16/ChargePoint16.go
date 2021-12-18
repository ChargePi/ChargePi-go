package v16

import (
	"context"
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/reactivex/rxgo/v2"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint/tls"
	"github.com/xBlaz3kx/ChargePi-go/components/connector"
	connector_manager "github.com/xBlaz3kx/ChargePi-go/components/connector-manager"
	"github.com/xBlaz3kx/ChargePi-go/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/components/hardware/reader"
	"github.com/xBlaz3kx/ChargePi-go/components/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/data/auth"
	"github.com/xBlaz3kx/ChargePi-go/data/settings"
	"log"
)

type (
	ChargePointHandler struct {
		chargePoint      ocpp16.ChargePoint
		availability     core.AvailabilityType
		connectorManager connector_manager.Manager
		connectorChannel chan rxgo.Item
		Settings         *settings.Settings
		TagReader        reader.Reader
		Indicator        indicator.Indicator
		LCD              display.LCD
	}

	ChargePoint interface {
		Run(ctx context.Context)
		connect(ctx context.Context)
		ListenForTag(ctx context.Context)
		HandleChargingRequest(tagId string)
		CleanUp(reason core.Reason)
		restoreState()
		notifyConnectorStatus(connector connector.Connector)
		startChargingConnector(connector connector.Connector, tagId string) error
		displayConnectorStatus(connectorId int, status core.ChargePointStatus)
	}
)

func NewChargePoint(tagReader reader.Reader, lcd display.LCD) *ChargePointHandler {
	ch := make(chan rxgo.Item, 5)
	return &ChargePointHandler{
		TagReader:        tagReader,
		LCD:              lcd,
		availability:     core.AvailabilityTypeInoperative,
		connectorChannel: ch,
		connectorManager: connector_manager.NewManager(ch),
	}
}

func (handler *ChargePointHandler) Run(ctx context.Context, settings *settings.Settings) {
	if settings == nil {
		log.Fatal("no settings provided")
	}

	var (
		client    ws.WsClient = nil
		info                  = settings.ChargePoint.Info
		tlsConfig             = settings.ChargePoint.TLS
	)

	handler.Settings = settings

	// Check if the client has TLS
	if tlsConfig.IsEnabled {
		client = tls.GetTLSClient(tlsConfig.CACertificatePath, tlsConfig.ClientCertificatePath, tlsConfig.ClientKeyPath)
		handler.chargePoint = ocpp16.NewChargePoint(info.Id, nil, client)
	} else {
		handler.chargePoint = ocpp16.NewChargePoint(info.Id, nil, nil)
	}

	// Start listening for tags from reader
	if handler.TagReader != nil {
		go handler.ListenForTag(ctx)
	}

	handler.addConnectors()
	handler.setMaxCachedTags()

	// Set handlers for Core, Reservation and RemoteTrigger
	handler.chargePoint.SetCoreHandler(handler)
	handler.chargePoint.SetReservationHandler(handler)
	handler.chargePoint.SetRemoteTriggerHandler(handler)

	handler.connect(ctx, fmt.Sprintf("ws://%s/%s", info.ServerUri, info.Id))
}

// connect to the central system and attempt to boot
func (handler *ChargePointHandler) connect(ctx context.Context, serverUrl string) {
	log.Println("Trying to connect to the central system: ", serverUrl)
	connectErr := handler.chargePoint.Start(serverUrl)

	// Check if the connection was successful
	if connectErr != nil {
		handler.chargePoint.Stop()
		handler.CleanUp(core.ReasonOther)
		log.Fatalf("Error connecting to the central system: %s \n", connectErr)
	}

	log.Printf("connected to central server: %s", serverUrl)
	handler.availability = core.AvailabilityTypeOperative

	go handler.listenForConnectorStatusChange(ctx)
	handler.bootNotification()
}

// HandleChargingRequest Entry point for determining if the request is to start or stop charging. Trying to find a connector that has the tag stored in the Session; if such a connector exists,
// execute stopChargingConnector, otherwise startCharging.
func (handler *ChargePointHandler) HandleChargingRequest(tagId string) {
	log.Printf("Handling request for tag %s", tagId)

	c := handler.connectorManager.FindConnectorWithTagId(tagId)
	if c != nil {
		err := handler.stopChargingConnector(c, core.ReasonLocal)
		if err != nil {
			log.Printf("Error stopping charing the connector: %s", err)
		}
	} else {
		err := handler.startCharging(tagId)
		if err != nil {
			log.Printf("Error started charing the connector: %s \n", err)
		}
	}
}

// CleanUp When exiting the client, stop all the transactions, clean up all the peripherals and terminate the connection.
func (handler *ChargePointHandler) CleanUp(reason core.Reason) {
	log.Println("Cleaning up ChargePoint, reason:", reason)

	handler.connectorManager.StopAllConnectors(reason)
	for _, c := range handler.connectorManager.GetConnectors() {
		if c.IsCharging() {
			log.Println("Stopping a transaction at connector: ", c.GetConnectorId())
			err := handler.stopChargingConnector(c, reason)
			if err != nil {
				log.Printf("error while stopping the transaction at cleanup: %v", err)
			}
		}
	}

	log.Println("Disconnecting the client..")
	handler.chargePoint.Stop()

	if handler.TagReader != nil {
		log.Println("Cleaning up the Tag Reader")
		handler.TagReader.Cleanup()
	}

	if handler.LCD != nil {
		log.Println("Cleaning up LCD")
		handler.LCD.Cleanup()
	}

	if handler.Indicator != nil {
		log.Println("Cleaning up Indicator")
		handler.Indicator.Cleanup()
	}

	close(handler.connectorChannel)
	log.Println("Clearing the scheduler...")
	scheduler.GetScheduler().Stop()
	scheduler.GetScheduler().Clear()

	auth.DumpTags()
}
