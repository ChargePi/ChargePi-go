package v16

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	configManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

// bootNotification Notify the central system that the charging point is online. Set the setHeartbeat interval and call restoreState.
// If the central system does not accept the charge point, exit the client.
func (cp *ChargePoint) bootNotification() {
	var (
		ocppInfo = cp.info.OCPPDetails
		request  = core.BootNotificationRequest{
			ChargePointVendor:       ocppInfo.Vendor,
			ChargePointModel:        ocppInfo.Model,
			ChargePointSerialNumber: ocppInfo.ChargePointSerialNumber,
			ChargeBoxSerialNumber:   ocppInfo.ChargeBoxSerialNumber,
			FirmwareVersion:         chargePoint.FirmwareVersion,
			// Todo fetch from 4G/LTE/SIM module
			Iccid: ocppInfo.Iccid,
			Imsi:  ocppInfo.Imsi,
		}
	)

	callback := func(confirmation ocpp.Response, protoError error) {
		bootConf := confirmation.(*core.BootNotificationConfirmation)

		cp.logger.Infof("Registration status: %s", bootConf.Status)

		switch bootConf.Status {
		case core.RegistrationStatusAccepted:
			cp.logger.Info("Accepted by the central system")
			cp.isConnected = true
			cp.availability = core.AvailabilityTypeOperative
			cp.setHeartbeat(bootConf.Interval)

			// Send details about the charge point and its EVSes
			_ = cp.sendChargePointInfo()
			cp.sendEvses()

			// Notify charge point status
			cp.notifyConnectorStatus(0, core.ChargePointStatusAvailable, core.NoError)

			// Restore the state of the charge point
			cp.restoreState()
		case core.RegistrationStatusPending, core.RegistrationStatusRejected:
			// Schedule a new boot notification in 1 minute
			_, err := cp.scheduler.Every(bootConf.Interval).Seconds().LimitRunsTo(1).Tag("bootNotification").Do(cp.bootNotification)
			if err != nil {
				cp.logger.WithError(err).Errorf("Error rescheduling BootNotification")
			}
		}
	}

	cp.logger.Info("Sending a BootNotification")
	err := cp.sendRequest(request, callback)
	cp.handleRequestErr(err, "Error sending BootNotification")
}

func (cp *ChargePoint) setHeartbeat(interval int) {
	cp.logger.Debug("Setting a heartbeat schedule")

	heartBeatInterval, _ := configManager.GetConfigurationValue(configuration.HeartbeatInterval.String())
	if interval > 0 {
		interVal := fmt.Sprintf("%d", interval)
		heartBeatInterval = &interVal
	}

	_, err := cp.scheduler.Every(fmt.Sprintf("%ss", *heartBeatInterval)).Tag("heartbeat").Do(cp.sendHeartBeat)
	if err != nil {
		cp.logger.WithError(err).Errorf("Error scheduling heartbeat")
	}
}

// sendHeartBeat Send a setHeartbeat to the central system.
func (cp *ChargePoint) sendHeartBeat() error {
	cp.logger.Info("Sending a heartbeat")

	return cp.sendRequest(core.NewHeartbeatRequest(), func(confirmation ocpp.Response, protoError error) {
		if protoError != nil {
			return
		}

		cp.logger.Info("Sent heartbeat")
	})
}
