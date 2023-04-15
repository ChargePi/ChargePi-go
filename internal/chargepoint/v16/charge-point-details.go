package v16

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	data "github.com/xBlaz3kx/ChargePi-go/pkg/models/ocpp"
)

// sendHeartBeat Send a setHeartbeat to the central system.
func (cp *ChargePoint) sendChargePointInfo() error {
	cp.logger.Info("Sending charge point information to the central system..")
	dataTransfer := core.NewDataTransferRequest(cp.info.OCPPDetails.Vendor)
	dataTransfer.Data = data.NewChargePointInfo(cp.info.Type, cp.info.MaxPower)
	// dataTransfer.MessageId

	return util.SendRequest(cp.chargePoint, dataTransfer,
		func(confirmation ocpp.Response, protoError error) {
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

func (cp *ChargePoint) SendEVSEsDetails(evses ...evse.EVSE) {
	for _, e := range evses {
		cp.sendEvseDetails(e)
	}
}

// sendEvseDetails sends the connector type and maximum output power information to the backend.
func (cp *ChargePoint) sendEvseDetails(evse evse.EVSE) {
	logInfo := cp.logger.WithField("evseId", evse.GetEvseId())
	logInfo.Info("Sending EVSE details to the central system")

	dataTransfer := core.NewDataTransferRequest(cp.info.OCPPDetails.Vendor)

	var connectors []data.Connector
	for _, connector := range evse.GetConnectors() {
		connectors = append(connectors, data.NewConnector(connector.ConnectorId, connector.Type))
	}

	dataTransfer.Data = data.NewEvseInfo(evse.GetEvseId(), float32(evse.GetMaxChargingPower()), connectors...)

	err := util.SendRequest(cp.chargePoint, dataTransfer,
		func(confirmation ocpp.Response, protoError error) {
			if protoError != nil {
				logInfo.WithError(protoError).Warn("Error sending data")
				return
			}

			resp := confirmation.(*core.DataTransferConfirmation)
			if resp.Status == core.DataTransferStatusAccepted {
				logInfo.Info("Sent additional charge point information")
			}
		})
	util.HandleRequestErr(err, "Cannot send EVSE details")
}
