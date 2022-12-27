package v16

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	data "github.com/xBlaz3kx/ChargePi-go/pkg/models/ocpp"
	configManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

// bootNotification Notify the central system that the charging point is online. Set the setHeartbeat interval and call restoreState.
// If the central system does not accept the charge point, exit the client.
func (cp *ChargePoint) bootNotification() {
	var (
		ocppInfo = cp.settings.ChargePoint.Info.OCPPInfo
		request  = core.BootNotificationRequest{
			ChargeBoxSerialNumber:   ocppInfo.ChargeBoxSerialNumber,
			ChargePointModel:        ocppInfo.Model,
			ChargePointSerialNumber: ocppInfo.ChargePointSerialNumber,
			ChargePointVendor:       ocppInfo.Vendor,
			FirmwareVersion:         "0.1.0",
			Iccid:                   ocppInfo.Iccid,
			Imsi:                    ocppInfo.Imsi,
		}
	)

	callback := func(confirmation ocpp.Response, protoError error) {
		bootConf := confirmation.(*core.BootNotificationConfirmation)

		switch bootConf.Status {
		case core.RegistrationStatusAccepted:
			cp.logger.Info("Accepted from the central system")
			cp.setHeartbeat(bootConf.Interval)
			cp.restoreState()
			_ = cp.sendChargePointInfo()
		case core.RegistrationStatusPending:
			cp.logger.Info("Registration status pending")
			//todo reschedule boot notification
		default:
			cp.logger.Fatal("Denied by the central system.")
		}
	}

	cp.logger.Info("Sending a boot notification")
	err := util.SendRequest(cp.chargePoint, request, callback)
	util.HandleRequestErr(err, "Error sending BootNotification")
}

func (cp *ChargePoint) setHeartbeat(interval int) {
	cp.logger.Debug("Setting a heartbeat schedule")

	heartBeatInterval, _ := configManager.GetConfigurationValue(configuration.HeartbeatInterval.String())
	if interval > 0 {
		interVal := fmt.Sprintf("%d", interval)
		heartBeatInterval = &interVal
	}

	_, err := scheduler.GetScheduler().Every(fmt.Sprintf("%ss", *heartBeatInterval)).Tag("heartbeat").Do(cp.sendHeartBeat)
	if err != nil {
		cp.logger.WithError(err).Errorf("Error scheduling heartbeat")
	}
}

// sendHeartBeat Send a setHeartbeat to the central system.
func (cp *ChargePoint) sendHeartBeat() error {
	return util.SendRequest(cp.chargePoint,
		core.NewHeartbeatRequest(),
		func(confirmation ocpp.Response, protoError error) {
			if protoError != nil {
				return
			}

			cp.logger.Info("Sent heartbeat")
		})
}

// sendHeartBeat Send a setHeartbeat to the central system.
func (cp *ChargePoint) sendChargePointInfo() error {
	cpInfo := cp.settings.ChargePoint.Info
	dataTransfer := core.NewDataTransferRequest(cpInfo.OCPPInfo.Vendor)
	dataTransfer.Data = data.NewChargePointInfo(cpInfo.Type, cpInfo.MaxPower)
	//dataTransfer.MessageId

	return util.SendRequest(cp.chargePoint,
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
