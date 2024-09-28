package v16

import (
	"fmt"

	chargePoint "github.com/ChargePi/ChargePi-go/internal/chargepoint"
	ocpp2 "github.com/ChargePi/ChargePi-go/pkg/ocpp"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)

// bootNotification After connecting to the central system, send a BootNotification in order to get the
// Charge point status.
//
// If the central system does not accept the charge point, exit the client.
func (cp *ChargePoint) bootNotification() {
	ocppInfo := cp.info.OCPPDetails
	// _, _ = cp.modem.GetInfo()
	request := core.BootNotificationRequest{
		ChargePointVendor:       ocppInfo.Vendor,
		ChargePointModel:        ocppInfo.Model,
		ChargePointSerialNumber: ocppInfo.ChargePointSerialNumber,
		ChargeBoxSerialNumber:   ocppInfo.ChargeBoxSerialNumber,
		FirmwareVersion:         chargePoint.FirmwareVersion,
		// Todo fetch from 4G/LTE/SIM module
		// Iccid: ocppInfo.Iccid,
		// Imsi:  ocppInfo.Imsi,
	}

	callback := func(confirmation ocpp.Response, protoError error) {
		bootConf := confirmation.(*core.BootNotificationConfirmation)

		cp.logger.Infof("Registration status: %s", bootConf.Status)

		switch bootConf.Status {
		case core.RegistrationStatusAccepted:
			cp.logger.Info("Accepted by the central system")

			cp.state.SetConnected(true)
			cp.state.SetAvailability(core.AvailabilityTypeOperative)

			cp.setHeartbeat(bootConf.Interval)

			go func() {
				// Send details about the charge point and its EVSEs concurrently
				err := cp.sendChargePointInfo()
				if err != nil {
					cp.logger.WithError(err).Warn("Error sending charge point information")
				}
				cp.sendEvses()
			}()

			// Notify charge point status
			cp.notifyStatus(0, core.ChargePointStatusAvailable, core.NoError)

			// Restore the state of the charge point
			/*	err := cp.evseManager.RestoreEVSEs()
				if err != nil {
					cp.logger.WithError(err).Warn("Unable to restore states")
				}*/
		case core.RegistrationStatusPending, core.RegistrationStatusRejected:
			// Reschedule the boot notification in 1 minute
			cp.rescheduleBootNotification(bootConf.Interval)
		}
	}

	cp.logger.Info("Sending a BootNotification")
	err := cp.sendRequest(request, callback)
	cp.handleRequestErr(err, "Error sending BootNotification")
}

func (cp *ChargePoint) rescheduleBootNotification(interval int) {
	// Schedule a new boot notification in 1 minute
	_, err := cp.scheduler.Every(interval).Seconds().LimitRunsTo(1).Tag("bootNotification").Do(cp.bootNotification)
	if err != nil {
		cp.logger.WithError(err).Fatal("Error rescheduling BootNotification")
	}
}

// setHeartbeat Set the heartbeat interval for the charge point based on the OCPP configuration
func (cp *ChargePoint) setHeartbeat(interval int) {
	cp.logger.Debug("Setting a heartbeat schedule")

	if interval > 0 {
		// Default to the configuration value
		heartBeatInterval, err := cp.settingsManager.GetConfigurationValue(ocpp_v16.HeartbeatInterval)
		if err != nil {
			cp.logger.WithError(err).Fatal("Error getting heartbeat interval from configuration")
		}

		_, err = cp.scheduler.Every(fmt.Sprintf("%ss", *heartBeatInterval)).Tag("heartbeat").Do(cp.sendHeartBeat)
		if err != nil {
			cp.logger.WithError(err).Fatal("Error scheduling heartbeat")
		}
		return
	}

	_, err := cp.scheduler.Every(fmt.Sprintf("%ds", interval)).Tag("heartbeat").Do(cp.sendHeartBeat)
	if err != nil {
		cp.logger.WithError(err).Fatal("Error scheduling heartbeat")
	}
}

// sendHeartBeat Send a setHeartbeat to the central system.
func (cp *ChargePoint) sendHeartBeat() error {
	cp.logger.Info("Sending a heartbeat")

	return cp.sendRequest(core.NewHeartbeatRequest(), func(confirmation ocpp.Response, protoError error) {
		if protoError != nil {
			cp.logger.Info("Heartbeat failed")
			return
		}

		cp.logger.Info("Sent heartbeat")
	})
}

// sendChargePointInfo sends information regarding the charge point to the central system.
func (cp *ChargePoint) sendChargePointInfo() error {
	cp.logger.Info("Sending charge point information to the central system..")
	dataTransfer := core.NewDataTransferRequest(cp.info.OCPPDetails.Vendor)
	dataTransfer.Data = ocpp2.NewChargePointInfo(cp.info.Type, cp.info.MaxPower)

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
		var connectors []ocpp2.Connector
		for _, connector := range evse.GetConnectors() {
			connectors = append(connectors, ocpp2.NewConnector(connector.ConnectorId, connector.Type))
		}

		cp.SendEVSEsDetails(evse.GetEvseId(), float32(evse.GetMaxChargingPower()), connectors...)
	}
}

// SendEVSEsDetails Send the EVSE's configuration to the central system.
func (cp *ChargePoint) SendEVSEsDetails(evseId int, maxPower float32, connectors ...ocpp2.Connector) {
	logInfo := cp.logger.WithField("evseId", evseId)
	logInfo.Info("Sending EVSE details to the central system")

	dataTransfer := core.NewDataTransferRequest(cp.info.OCPPDetails.Vendor)
	dataTransfer.Data = ocpp2.NewEvseInfo(evseId, maxPower, connectors...)

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
	cp.handleRequestErr(err, "Error sending data")
}
