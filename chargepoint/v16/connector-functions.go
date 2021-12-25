package v16

import (
	"context"
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/reactivex/rxgo/v2"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/components/cache"
	connector2 "github.com/xBlaz3kx/ChargePi-go/components/connector"
	"github.com/xBlaz3kx/ChargePi-go/components/hardware/display/i18n"
	"github.com/xBlaz3kx/ChargePi-go/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/components/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/components/settings/settings-manager"
	"github.com/xBlaz3kx/ChargePi-go/data"
	settingsData "github.com/xBlaz3kx/ChargePi-go/data/settings"
	"time"
)

// addConnectors Add the Connectors from the connectors.json file to the handler. Create and add all their components and initialize the struct.
func (handler *ChargePointHandler) addConnectors() {
	connectors := settings_manager.GetConnectors()
	if data.IsNilInterfaceOrPointer(connectors) {
		log.Fatal("no connectors configured")
	}

	log.Info("Adding connectors")
	err := handler.connectorManager.AddConnectorsFromConfiguration(handler.Settings.ChargePoint.Info.MaxChargingTime, connectors)
	if err != nil {
		log.Fatal(err)
	}

	// Add an indicator with the length of valid connectors
	handler.Indicator = indicator.NewIndicator(len(handler.connectorManager.GetConnectors()))
}

// restoreState After connecting to the central system, try to restore the previous state of each ConnectorImpl and notify the system about its state.
// If the ConnectorStatus was "Preparing" or "Charging", try to resume or start charging. If the charging fails, change the connector status and notify the central system.
func (handler *ChargePointHandler) restoreState() {
	log.Info("Restoring connectors' state")
	var err error

	for _, connector := range handler.connectorManager.GetConnectors() {
		// Load connector configuration from cache
		connectorSettings, isFound := cache.Cache.Get(fmt.Sprintf("connectorEvse%dId%dConfiguration", connector.GetEvseId(), connector.GetConnectorId()))
		if !isFound {
			continue
		}
		cachedConnector := connectorSettings.(*settingsData.Connector)

		err = handler.connectorManager.RestoreConnectorStatus(cachedConnector)
		switch err {
		case nil:
			_, err = scheduler.GetScheduler().Every(connector.GetMaxChargingTime()-0).Minutes().LimitRunsTo(1).
				Tag(fmt.Sprintf("connector%dTimer", connector.GetConnectorId())).Do(handler.stopChargingConnector, connector, core.ReasonLocal)
			break
		default:
			// Attempt to stop charging
			err = handler.stopChargingConnector(connector, core.ReasonDeAuthorized)
			if err != nil {
				log.Debugf("Stopping the charging returned %v", err)
				connector.SetStatus(core.ChargePointStatusFaulted, core.InternalError)
			}
		}
	}
}

// notifyConnectorStatus Notify the central system about the connector's status and updates the LED indicator.
func (handler *ChargePointHandler) notifyConnectorStatus(connector connector2.Connector) {
	if !data.IsNilInterfaceOrPointer(connector) {
		var (
			status, errorCode = connector.GetStatus()
			connectorId       = connector.GetConnectorId()
			request           = core.NewStatusNotificationRequest(connectorId, errorCode, status)
		)

		request.Timestamp = types.NewDateTime(time.Now())

		callback := func(confirmation ocpp.Response, protoError error) {
			log.Infof("Notified status of the connector %d: %s", connectorId, status)
		}

		err := handler.SendRequest(request, callback)
		if err != nil {
			log.Errorf("Cannot send status notification of connector: %v", err)
		}
	}
}

// listenForConnectorStatusChange listen for change in connector and notify the central system about the state
func (handler *ChargePointHandler) listenForConnectorStatusChange(ctx context.Context) {
	log.Infof("Starting to listen for connector status change")
	observableConnectors := rxgo.FromChannel(handler.connectorChannel)

	if observableConnectors != nil {
	Listener:
		for {
			select {
			// Start observing the connector for changes in status
			case item := <-observableConnectors.Observe():
				connector, canCast := item.V.(*connector2.ConnectorImpl)
				if canCast {
					// Connector starts with index 1,
					connectorIndex := connector.ConnectorId - 1

					handler.displayLEDStatus(connectorIndex, connector.ConnectorStatus)
					go handler.displayConnectorStatus(connector.ConnectorId, connector.ConnectorStatus)
					handler.notifyConnectorStatus(connector)
				}
				break
			case <-ctx.Done():
				break Listener
			default:
			}
		}
	}
}

func (handler *ChargePointHandler) displayConnectorStatus(connectorId int, status core.ChargePointStatus) {
	var (
		language = handler.Settings.ChargePoint.Hardware.Lcd.Language
		message  = []string{}
		err      error
	)

	switch status {
	case core.ChargePointStatusAvailable:
		message, err = i18n.TranslateConnectorAvailableMessage(language, connectorId)
		break
	case core.ChargePointStatusFinishing:
		message, err = i18n.TranslateConnectorFinishingMessage(language, connectorId)
		break
	case core.ChargePointStatusCharging:
		message, err = i18n.TranslateConnectorChargingMessage(language, connectorId)
		break
	case core.ChargePointStatusFaulted:
		message, err = i18n.TranslateConnectorFaultedMessage(language, connectorId)
		break
	default:
		return
	}

	if err != nil {
		log.Println(err)
		return
	}

	handler.sendToLCD(message...)
}
