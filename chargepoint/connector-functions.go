package chargepoint

import (
	"fmt"
	"github.com/kr/pretty"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	types2 "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/reactivex/rxgo/v2"
	"github.com/xBlaz3kx/ChargePi-go/cache"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"github.com/xBlaz3kx/ChargePi-go/hardware"
	"github.com/xBlaz3kx/ChargePi-go/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/hardware/power-meter"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"log"
	"time"
)

// AddConnectors Add the Connectors from the connectors.json file to the handler. Create and add all their components and initialize the struct.
func (handler *ChargePointHandler) AddConnectors() {
	connectors := settings.GetConnectors()
	log.Println("Adding connectors")
	handler.Connectors = []*Connector{}
	for _, connector := range connectors {

		// Create a power meter from settings
		powerMeter, err := power_meter.NewPowerMeter(connector)
		if err != nil {
			log.Printf("Cannot instantiate power meter: %s", err)
		}

		// Create a connector object
		connectorObj, err := NewConnector(
			connector.EvseId,
			connector.ConnectorId,
			connector.Type,
			hardware.NewRelay(
				connector.Relay.RelayPin,
				connector.Relay.InverseLogic,
			),
			powerMeter,
			connector.PowerMeter.Enabled,
			handler.Settings.ChargePoint.Info.MaxChargingTime,
		)
		if err != nil {
			log.Println("Error while creating a connector:", err)
			continue
		}

		handler.Connectors = append(handler.Connectors, connectorObj)
		pretty.Print("Added a connector", connectorObj)
	}
	// Add an indicator with the length of valid connectors
	handler.Indicator = indicator.NewIndicator(len(handler.Connectors))
}

// restoreState After connecting to the central system, try to restore the previous state of each Connector and notify the system about its state.
// If the ConnectorStatus was "Preparing" or "Charging", try to resume or start charging. If the charging fails, change the connector status and notify the central system.
func (handler *ChargePointHandler) restoreState() {
	var err error
	for _, connector := range handler.Connectors {
		// load connector configuration from cache
		connectorSettings, isFound := cache.Cache.Get(fmt.Sprintf("connectorEvse%dId%dConfiguration", connector.EvseId, connector.ConnectorId))
		if !isFound {
			continue
		}
		cachedConnector := connectorSettings.(*settings.Connector)

		if cachedConnector != nil {
			//set the previous status to determine what action to do
			connector.SetStatus(core.ChargePointStatus(cachedConnector.Status), core.NoError)
			switch core.ChargePointStatus(cachedConnector.Status) {
			case core.ChargePointStatusPreparing:
				log.Println("Attempting to start charging at connector", connector.ConnectorId)
				err = handler.startCharging(cachedConnector.Session.TagId)
				if err != nil {
					log.Println(err)
					connector.SetStatus(core.ChargePointStatusAvailable, core.InternalError)
					continue
				}
				break
			case core.ChargePointStatusCharging:
				err = handler.attemptToResumeChargingAtConnector(connector, data.Session(cachedConnector.Session))
				if err != nil {
					log.Printf("Resume charging failed at %d %d, reason: %v", connector.EvseId, connector.ConnectorId, err)
					//attempt to stop charging
					err = handler.stopChargingConnector(connector, core.ReasonDeAuthorized)
					if err != nil {
						log.Println("Stopping the charging returned", err)
						connector.SetStatus(core.ChargePointStatusFaulted, core.InternalError)
					}
					continue
				}
				log.Println("Charging continued at ", connector.EvseId, connector.ConnectorId)
				break
			case core.ChargePointStatusFaulted:
				break
			}
		}
	}
}

// attemptToResumeChargingAtConnector try to resume or stop charging at a connector based on the status in the connector persistence file.
func (handler *ChargePointHandler) attemptToResumeChargingAtConnector(connector *Connector, session data.Session) error {
	log.Println("Attempt to resume charging at charging at", connector.ConnectorId)
	parse, err := time.Parse(time.RFC3339, session.Started)
	if err != nil {
		return err
	}
	chargingTimeElapsed := int(time.Now().Sub(parse).Minutes())
	if connector.MaxChargingTime < chargingTimeElapsed {
		//set the transaction id so connector is able to stop the transaction
		connector.session.TransactionId = session.TransactionId
		return fmt.Errorf("session time limit exceeded")
	}
	err = connector.ResumeCharging(session)
	if err != nil {
		return fmt.Errorf("charging session is unable to be resumed")
	}
	_, err = scheduler.GetScheduler().Every(connector.MaxChargingTime-chargingTimeElapsed).Minutes().LimitRunsTo(1).
		Tag(fmt.Sprintf("connector%dTimer", connector.ConnectorId)).Do(handler.stopChargingConnector, connector, core.ReasonLocal)
	return nil
}

// notifyConnectorStatus Notify the central system about the connector's status and updates the LED indicator.
func (handler *ChargePointHandler) notifyConnectorStatus(connector *Connector) {
	if connector != nil {
		request := core.StatusNotificationRequest{
			ConnectorId: connector.ConnectorId,
			Status:      connector.ConnectorStatus,
			ErrorCode:   connector.ErrorCode,
			Timestamp:   &types2.DateTime{Time: time.Now()},
		}
		callback := func(confirmation ocpp.Response, protoError error) {
			log.Printf("Notified status of the connector %d: %s", connector.ConnectorId, connector.ConnectorStatus)
		}
		err := handler.SendRequest(request, callback)
		if err != nil {
			log.Println("Cannot send status notification of connector: ", err)
			return
		}
	}
}

// listenForConnectorStatusChange listen for change in connector and notify the central system about the state
func (handler *ChargePointHandler) listenForConnectorStatusChange() {
	handler.connectorChannel = make(chan rxgo.Item)
	observableConnectors := rxgo.FromChannel(handler.connectorChannel)

	if observableConnectors != nil {
		// Set the communication channel before listening
		for _, connector := range handler.Connectors {
			connector.connectorNotificationChannel = handler.connectorChannel
		}
		// Start observing the connector for changes in status
		for item := range observableConnectors.Observe() {
			connector := item.V.(*Connector)
			connectorIndex := connector.ConnectorId - 1
			handler.displayLEDStatus(connectorIndex, connector.ConnectorStatus)
			handler.displayConnectorStatus(connector.ConnectorId, connector.ConnectorStatus)
			handler.notifyConnectorStatus(connector)
		}
	}
}

func (handler *ChargePointHandler) displayConnectorStatus(connectorId int, status core.ChargePointStatus) {
	switch status {
	case core.ChargePointStatusAvailable:
		handler.sendToLCD(fmt.Sprintf("Connector %d", connectorId), "available")
		break
	case core.ChargePointStatusFinishing:
		handler.sendToLCD("Stopped charging", fmt.Sprintf("at %d", connectorId))
		break
	case core.ChargePointStatusCharging:
		handler.sendToLCD("Started charging", fmt.Sprintf("at %d", connectorId))
		break
	case core.ChargePointStatusFaulted:
		handler.sendToLCD(fmt.Sprintf("Connector %d", connectorId), "faulted")
		break
	}
}
