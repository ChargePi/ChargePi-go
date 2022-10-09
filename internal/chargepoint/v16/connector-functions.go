package v16

import (
	"context"
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/display/i18n"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
	settingsData "github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	data "github.com/xBlaz3kx/ChargePi-go/pkg/models/ocpp"
	"time"
)

// AddEVSEs Add the Connectors from the connectors.json file to the handler. Create and add all their components and initialize the struct.
func (cp *ChargePoint) AddEVSEs(connectors []*settingsData.EVSE) {
	if util.IsNilInterfaceOrPointer(connectors) {
		cp.logger.Fatal("no evses configured")
	}

	cp.logger.Info("Adding evses")
	err := cp.connectorManager.AddEVSEsFromSettings(cp.settings.ChargePoint.Info.MaxChargingTime, connectors)
	if err != nil {
		cp.logger.WithError(err).Fatalf("Unable to add evses from configuration")
	}

	// Add an indicator with the length of valid connectors
	cp.indicator = indicator.NewIndicator(len(cp.connectorManager.GetEVSEs()))

	cp.sendInfo(connectors)
}

func (cp *ChargePoint) sendInfo(evses []*settingsData.EVSE) {
	for _, evse := range evses {
		cp.sendConnectorInfo(evse)
	}
}

// sendConnectorInfo sends the connector type and maximum output power information to the backend.
func (cp *ChargePoint) sendConnectorInfo(evse *settingsData.EVSE) {
	dataTransfer := core.NewDataTransferRequest(cp.settings.ChargePoint.Info.OCPPInfo.Vendor)

	var connectors []data.Connector
	for _, connector := range evse.Connectors {
		connectors = append(connectors, data.NewConnector(connector.ConnectorId, connector.Type))
	}

	dataTransfer.Data = data.NewEvseInfo(evse.EvseId, evse.MaxPower, connectors...)

	_ = util.SendRequest(cp.chargePoint,
		dataTransfer,
		func(confirmation ocpp.Response, protoError error) {
			if protoError != nil {
				cp.logger.Info("Error sending data")
				return
			}

			resp := confirmation.(*core.DataTransferConfirmation)
			if resp.Status == core.DataTransferStatusAccepted {
				cp.logger.Info("Sent additional charge point information")
			}
		})
}

// restoreState After connecting to the central system, try to restore the previous state of each EVSE and notify
// the system about its state.
//
// If the ConnectorStatus was "Preparing" or "Charging", try to resume or start charging.
// If the charging fails, change the connector status and notify the central system.
func (cp *ChargePoint) restoreState() {
	cp.logger.Info("Restoring evses' state")

	for _, c := range cp.connectorManager.GetEVSEs() {
		var (
			cacheKey = fmt.Sprintf("evse%d", c.GetEvseId())
			conn     settingsData.EVSE
		)

		// Fetch the viper configuration
		connectorCfg, isFound := settings.EVSESettings.Load(cacheKey)
		if !isFound {
			continue
		}
		cfg := connectorCfg.(*viper.Viper)

		// Unmarshall
		err := cfg.Unmarshal(&conn)
		if err != nil {
			continue
		}

		err = cp.connectorManager.RestoreEVSEStatus(&conn)
		switch err {
		case nil:
			_, err = cp.scheduler.Every(c.GetMaxChargingTime()-0).Minutes().LimitRunsTo(1).
				Tag(fmt.Sprintf("evse%dTimer", c.GetEvseId())).Do(cp.stopChargingConnector, c, core.ReasonLocal)
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
func (cp *ChargePoint) notifyConnectorStatus(evseId int, status core.ChargePointStatus, errCode core.ChargePointErrorCode) {
	var (
		request = core.NewStatusNotificationRequest(evseId, errCode, status)
	)

	request.Timestamp = types.NewDateTime(time.Now())

	callback := func(confirmation ocpp.Response, protoError error) {
		cp.logger.Infof("Notified status of the connector %d: %s", evseId, status)
	}

	err := util.SendRequest(cp.chargePoint, request, callback)
	util.HandleRequestErr(err, "Cannot send status of connector")
}

// ListenForConnectorStatusChange listen for change in connector and notify the central system about the state
func (cp *ChargePoint) ListenForConnectorStatusChange(ctx context.Context, ch <-chan chargePoint.StatusNotification) {
	cp.logger.Debug("Starting to listen for connector status change")

Listener:
	for {
		select {
		// Start observing the connector for changes in status
		case c := <-ch:
			// Connector starts with index 1
			connectorIndex := c.EvseId - 1
			status := core.ChargePointStatus(c.Status)
			errCode := core.ChargePointErrorCode(c.Status)

			go cp.displayLEDStatus(connectorIndex, status)
			go cp.displayConnectorStatus(c.EvseId, status)

			// Send a status notification to the Central System
			cp.notifyConnectorStatus(c.EvseId, status, errCode)

		case meterVal := <-cp.meterValuesChannel:
			// Send a meter value notification to the Central System
			values := core.NewMeterValuesRequest(*meterVal.ConnectorId, meterVal.MeterValues)
			err := util.SendRequest(cp.chargePoint, values, func(confirmation ocpp.Response, protoError error) {})
			if err != nil {
				cp.logger.WithError(err).Errorf("Cannot send meter values")
			}

		case <-ctx.Done():
			break Listener
		}
	}
}

func (cp *ChargePoint) displayConnectorStatus(connectorId int, status core.ChargePointStatus) {
	var (
		language = cp.settings.ChargePoint.Hardware.Display.Language
		message  []string
		err      error
	)

	switch status {
	case core.ChargePointStatusAvailable:
		message, err = i18n.TranslateConnectorAvailableMessage(language, connectorId)
	case core.ChargePointStatusFinishing:
		message, err = i18n.TranslateConnectorFinishingMessage(language, connectorId)
	case core.ChargePointStatusCharging:
		message, err = i18n.TranslateConnectorChargingMessage(language, connectorId)
	case core.ChargePointStatusFaulted:
		message, err = i18n.TranslateConnectorFaultedMessage(language, connectorId)
	default:
		return
	}

	if err != nil {
		cp.logger.WithError(err).Errorf("Error displaying status")
		return
	}

	cp.sendToLCD(message...)
}
