package v16

import (
	"context"
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/reactivex/rxgo/v2"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/connector"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/hardware/display/i18n"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/models"
	settingsData "github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/util"
	"time"
)

// AddConnectors Add the Connectors from the connectors.json file to the handler. Create and add all their components and initialize the struct.
func (cp *ChargePoint) AddConnectors(connectors []*settingsData.Connector) {
	if util.IsNilInterfaceOrPointer(connectors) {
		cp.logger.Fatal("no connectors configured")
	}

	cp.logger.Debugf("Adding connectors")
	err := cp.connectorManager.AddConnectorsFromConfiguration(cp.Settings.ChargePoint.Info.MaxChargingTime, connectors)
	if err != nil {
		cp.logger.WithError(err).Fatalf("Unable to add connectors from configuration")
	}

	// Add an indicator with the length of valid connectors
	cp.Indicator = indicator.NewIndicator(len(cp.connectorManager.GetConnectors()))
}

// restoreState After connecting to the central system, try to restore the previous state of each ConnectorImpl and notify the system about its state.
// If the ConnectorStatus was "Preparing" or "Charging", try to resume or start charging. If the charging fails, change the connector status and notify the central system.
func (cp *ChargePoint) restoreState() {
	cp.logger.Debugf("Restoring connectors' state")

	for _, c := range cp.connectorManager.GetConnectors() {
		var (
			cacheKey = fmt.Sprintf("connectorEvse%dId%d", c.GetEvseId(), c.GetConnectorId())
			conn     settingsData.Connector
		)

		// Fetch the viper configuration
		connectorCfg, isFound := settings.ConnectorSettings.Load(cacheKey)
		if !isFound {
			continue
		}
		cfg := connectorCfg.(*viper.Viper)

		// Unmarshall
		err := cfg.Unmarshal(&conn)
		if err != nil {
			continue
		}

		err = cp.connectorManager.RestoreConnectorStatus(&conn)
		switch err {
		case nil:
			_, err = cp.scheduler.Every(c.GetMaxChargingTime()-0).Minutes().LimitRunsTo(1).
				Tag(fmt.Sprintf("connector%dTimer", c.GetConnectorId())).Do(cp.stopChargingConnector, c, core.ReasonLocal)
			break
		default:
			// Attempt to stop charging
			err = cp.stopChargingConnector(c, core.ReasonDeAuthorized)
			if err != nil {
				cp.logger.Debugf("Stopping the charging returned %v", err)
				c.SetStatus(core.ChargePointStatusFaulted, core.InternalError)
			}
		}
	}
}

// notifyConnectorStatus Notify the central system about the connector's status and updates the LED indicator.
func (cp *ChargePoint) notifyConnectorStatus(connector connector.Connector) {
	if util.IsNilInterfaceOrPointer(connector) {
		return
	}

	var (
		status, errorCode = connector.GetStatus()
		connectorId       = connector.GetConnectorId()
		request           = core.NewStatusNotificationRequest(connectorId, errorCode, status)
	)

	request.Timestamp = types.NewDateTime(time.Now())

	callback := func(confirmation ocpp.Response, protoError error) {
		cp.logger.Infof("Notified status of the connector %d: %s", connectorId, status)
	}

	err := util.SendRequest(cp.chargePoint, request, callback)
	util.HandleRequestErr(err, "Cannot send status of connector")
}

// ListenForConnectorStatusChange listen for change in connector and notify the central system about the state
func (cp *ChargePoint) ListenForConnectorStatusChange(ctx context.Context, ch <-chan rxgo.Item) {
	cp.logger.Debug("Starting to listen for connector status change")
	observableConnectors := rxgo.FromChannel(ch)

	if observableConnectors != nil {
	Listener:
		for {
			select {
			// Start observing the connector for changes in status
			case item := <-observableConnectors.Observe():
				c, canCast := item.V.(connector.Connector)
				if canCast {
					// Connector starts with index 1
					connectorIndex := c.GetConnectorId() - 1
					status, _ := c.GetStatus()

					cp.displayLEDStatus(connectorIndex, status)
					go cp.displayConnectorStatus(c.GetConnectorId(), status)
					cp.notifyConnectorStatus(c)
				}
				break
			case meterVal := <-cp.meterValuesChannel:
				meterValues, canCast := meterVal.V.(models.MeterValueNotification)
				if canCast {
					values := core.NewMeterValuesRequest(meterValues.ConnectorId, meterValues.MeterValues)
					err := util.SendRequest(cp.chargePoint, values, func(confirmation ocpp.Response, protoError error) {})
					if err != nil {
						cp.logger.WithError(err).Errorf("Cannot send meter values")
					}
				}
				break
			case <-ctx.Done():
				break Listener
			default:
			}
		}
	}
}

func (cp *ChargePoint) displayConnectorStatus(connectorId int, status core.ChargePointStatus) {
	var (
		language = cp.Settings.ChargePoint.Hardware.Lcd.Language
		message  []string
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
		cp.logger.WithError(err).Errorf("Error displaying status")
		return
	}

	cp.sendToLCD(message...)
}
