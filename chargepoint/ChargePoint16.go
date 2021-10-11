package chargepoint

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/reactivex/rxgo/v2"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"github.com/xBlaz3kx/ChargePi-go/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/hardware/reader"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"log"
	"sync"
)

type ChargePointHandler struct {
	mu               sync.Mutex
	chargePoint      ocpp16.ChargePoint
	IsAvailable      bool
	Connectors       []*Connector
	Settings         *settings.Settings
	TagReader        reader.Reader
	Indicator        indicator.Indicator
	LCD              display.LCD
	connectorChannel chan rxgo.Item
}

func (handler *ChargePointHandler) Run() {
	var (
		client ws.WsClient = nil
		info               = handler.Settings.ChargePoint.Info
		tls                = info.TLS
	)

	// Check if the client has TLS
	if tls.IsEnabled {
		client = GetTLSClient(tls.CACertificatePath, tls.ClientCertificatePath, tls.ClientKeyPath)
		handler.chargePoint = ocpp16.NewChargePoint(info.Id, nil, client)
	} else {
		handler.chargePoint = ocpp16.NewChargePoint(info.Id, nil, nil)
	}

	// Set handlers for Core, Reservation and RemoteTrigger
	handler.chargePoint.SetCoreHandler(handler)
	handler.chargePoint.SetReservationHandler(handler)
	handler.chargePoint.SetRemoteTriggerHandler(handler)

	handler.setMaxCachedTags()
	handler.connect()
}

// connect to the central system and attempt to boot
func (handler *ChargePointHandler) connect() {
	var (
		info = handler.Settings.ChargePoint.Info
	)

	serverUrl := fmt.Sprintf("ws://%s/%s", info.ServerUri, info.Id)
	log.Println("Trying to connect to the central system: ", serverUrl)
	connectErr := handler.chargePoint.Start(serverUrl)

	go handler.listenForTag()

	// Check if the connection was successful
	if connectErr != nil {
		log.Printf("Error connecting to the central system: %s", connectErr)
		handler.CleanUp(core.ReasonOther)
		handler.chargePoint.Stop()
	} else {
		log.Printf("connected to central server: %s with ID: %s", serverUrl, info.Id)
		handler.IsAvailable = true

		go handler.listenForConnectorStatusChange()
		handler.bootNotification()
	}
}

// FindAvailableConnector Find first Connector with the status "Available" from the handler.
func (handler *ChargePointHandler) FindAvailableConnector() *Connector {
	for _, connector := range handler.Connectors {
		if connector.IsAvailable() {
			return connector
		}
	}
	return nil
}

// FindConnectorWithId Find the Connector with the specified connectorID.
func (handler *ChargePointHandler) FindConnectorWithId(connectorID int) *Connector {
	for _, connector := range handler.Connectors {
		if connector.ConnectorId == connectorID {
			return connector
		}
	}
	return nil
}

// FindConnectorWithTagId Find the Connector that has the same tagId as the session of the connector.
func (handler *ChargePointHandler) FindConnectorWithTagId(tagId string) *Connector {
	for _, connector := range handler.Connectors {
		if connector.GetTagId() == tagId {
			return connector
		}
	}
	return nil
}

// FindConnectorWithTransactionId Find the Connector that contains the transactionId in the session of the connector.
func (handler *ChargePointHandler) FindConnectorWithTransactionId(transactionId string) *Connector {
	for _, connector := range handler.Connectors {
		if connector.GetTransactionId() == transactionId {
			return connector
		}
	}
	return nil
}

// FindConnectorWithReservationId Find the Connector that contains the reservationId.
func (handler *ChargePointHandler) FindConnectorWithReservationId(reservationId int) *Connector {
	for _, connector := range handler.Connectors {
		if connector.GetReservationId() == reservationId {
			return connector
		}
	}
	return nil
}

// HandleChargingRequest Entry point for determining if the request is to start or stop charging. Trying to find a connector that has the tag stored in the Session; if such a connector exists,
// execute stopChargingConnector, otherwise startCharging.
func (handler *ChargePointHandler) HandleChargingRequest(tagId string) {
	log.Printf("Handling request for tag %s", tagId)
	var connector = handler.FindConnectorWithTagId(tagId)
	if connector != nil {
		err := handler.stopChargingConnector(connector, core.ReasonLocal)
		if err != nil {
			log.Printf("Error stopping charing the connector: %s", err)
			return
		}
	} else {
		err := handler.startCharging(tagId)
		if err != nil {
			log.Printf("Error started charing the connector: %s", err)
			return
		}
	}
}

// CleanUp When exiting the client, stop all the transactions, clean up all the peripherals and terminate the connection.
func (handler *ChargePointHandler) CleanUp(reason core.Reason) {
	log.Println("Cleaning up ChargePoint, reason:", reason)
	for _, connector := range handler.Connectors {
		if connector.IsCharging() {
			log.Println("Stopping a transaction at connector: ", connector.ConnectorId)
			err := handler.stopChargingConnector(connector, reason)
			if err != nil {
				log.Printf("error while stopping the transaction at cleanup: %v", err)
			}
		}
	}

	log.Println("Disconnecting the client..")
	handler.chargePoint.Stop()

	if handler.TagReader != nil {
		log.Println("Cleaning up Tag Reader")
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

	data.DumpTags()
}
