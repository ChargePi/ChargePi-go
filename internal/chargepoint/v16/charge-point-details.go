package v16

import (
	data "github.com/ChargePi/ChargePi-go/pkg/models/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)

// sendHeartBeat Send a setHeartbeat to the central system.
func (cp *ChargePoint) sendChargePointInfo() error {
	cp.logger.Info("Sending charge point information to the central system..")
	dataTransfer := core.NewDataTransferRequest(cp.info.OCPPDetails.Vendor)
	dataTransfer.Data = data.NewChargePointInfo(cp.info.Type, cp.info.MaxPower)
	// dataTransfer.MessageId

	return cp.sendRequest(dataTransfer, func(confirmation ocpp.Response, protoError error) {
		if protoError != nil {
			cp.logger.WithError(protoError).Warn("Error sending data")
			return
		}

		resp := confirmation.(*core.DataTransferConfirmation)
		if resp.Status == core.DataTransferStatusAccepted {
			cp.logger.Info("Sent additional charge point information")
		}
	})
}

// sendEvses Send the EVSEs' configuration to the central system.
func (cp *ChargePoint) sendEvses() {
	for _, evse := range cp.evseManager.GetEVSEs() {
		var connectors []data.Connector
		for _, connector := range evse.GetConnectors() {
			connectors = append(connectors, data.NewConnector(connector.ConnectorId, connector.Type))
		}

		cp.SendEVSEsDetails(evse.GetEvseId(), float32(evse.GetMaxChargingPower()), connectors...)
	}
}

func (cp *ChargePoint) SendEVSEsDetails(evseId int, maxPower float32, connectors ...data.Connector) {
	logInfo := cp.logger.WithField("evseId", evseId)
	logInfo.Info("Sending EVSE details to the central system")

	dataTransfer := core.NewDataTransferRequest(cp.info.OCPPDetails.Vendor)
	dataTransfer.Data = data.NewEvseInfo(evseId, maxPower, connectors...)

	err := cp.sendRequest(dataTransfer, func(confirmation ocpp.Response, protoError error) {
		if protoError != nil {
			logInfo.WithError(protoError).Warn("Error sending data")
			return
		}

		resp := confirmation.(*core.DataTransferConfirmation)
		if resp.Status == core.DataTransferStatusAccepted {
			logInfo.Info("Sent additional charge point information")
		}
	})
	if err != nil {
		logInfo.WithError(err).Warn("Error sending data")
	}
}
