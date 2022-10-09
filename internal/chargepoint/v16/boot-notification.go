package v16

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	configManager "github.com/xBlaz3kx/ocppManager-go"
	v16 "github.com/xBlaz3kx/ocppManager-go/v16"
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

	heartBeatInterval, _ := configManager.GetConfigurationValue(v16.HeartbeatInterval.String())
	if interval > 0 {
		heartBeatInterval = fmt.Sprintf("%d", interval)
	}

	heartBeatInterval = fmt.Sprintf("%ss", heartBeatInterval)
	_, err := scheduler.GetScheduler().Every(heartBeatInterval).Tag("heartbeat").Do(cp.sendHeartBeat)
	if err != nil {
		cp.logger.WithError(err).Errorf("Error scheduling heartbeat")
	}
}

// sendHeartBeat Send a setHeartbeat to the central system.
func (cp *ChargePoint) sendHeartBeat() error {
	return util.SendRequest(cp.chargePoint,
		core.NewHeartbeatRequest(),
		func(confirmation ocpp.Response, protoError error) {
			cp.logger.Info("Sent heartbeat")
		})
}
