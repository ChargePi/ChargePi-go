package v16

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/components/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/components/settings/conf-manager"
)

// bootNotification Notify the central system that the charging point is online. Set the setHeartbeat interval and call restoreState.
// If the central system does not accept the charge point, exit the client.
func (handler *ChargePointHandler) bootNotification() {
	var (
		ocppInfo = handler.Settings.ChargePoint.Info.OCPPInfo
		request  = core.BootNotificationRequest{
			ChargeBoxSerialNumber:   ocppInfo.ChargeBoxSerialNumber,
			ChargePointModel:        ocppInfo.Model,
			ChargePointSerialNumber: ocppInfo.ChargePointSerialNumber,
			ChargePointVendor:       ocppInfo.Vendor,
			FirmwareVersion:         "1.0",
			Iccid:                   ocppInfo.Iccid,
			Imsi:                    ocppInfo.Imsi,
		}
	)

	callback := func(confirmation ocpp.Response, protoError error) {
		bootConf := confirmation.(*core.BootNotificationConfirmation)

		switch bootConf.Status {
		case core.RegistrationStatusAccepted:
			log.Info("Notified and accepted from the central system")
			handler.setHeartbeat(bootConf.Interval)
			handler.restoreState()
			break
		case core.RegistrationStatusPending:
			//todo reschedule boot notification
			break
		default:
			log.Fatal("Denied by the central system.")
		}
	}

	handler.SendRequest(request, callback)
}

func (handler *ChargePointHandler) setHeartbeat(interval int) {
	log.Infof("Setting a heartbeat schedule")

	heartBeatInterval, _ := conf_manager.GetConfigurationValue("HeartbeatInterval")
	if interval > 0 {
		heartBeatInterval = fmt.Sprintf("%d", interval)
	}

	heartBeatInterval = fmt.Sprintf("%ss", heartBeatInterval)
	_, err := scheduler.GetScheduler().Every(heartBeatInterval).Tag("heartbeat").Do(handler.sendHeartBeat)
	if err != nil {
		log.Errorf("%v", err)
	}
}

// sendHeartBeat Send a setHeartbeat to the central system.
func (handler *ChargePointHandler) sendHeartBeat() error {
	return handler.SendRequest(
		core.NewHeartbeatRequest(),
		func(confirmation ocpp.Response, protoError error) {
			log.Info("Sent heartbeat")
		})
}
