package v16

import (
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
)

func (cp *ChargePoint) OnTriggerMessage(request *remotetrigger.TriggerMessageRequest) (confirmation *remotetrigger.TriggerMessageConfirmation, err error) {
	logInfo := cp.logger.WithFields(log.Fields{"feature": request.GetFeatureName(), "request": request.RequestedMessage})
	logInfo.Infof("Received a request")
	status := remotetrigger.TriggerMessageStatusRejected

	switch request.RequestedMessage {
	case core.BootNotificationFeatureName:

		// Schedule to send a BootNotification per request
		_, err = cp.scheduler.Every(5).Seconds().LimitRunsTo(1).Do(cp.bootNotification)
		if err != nil {
			logInfo.WithError(err).Error("Cannot schedule a boot notification")
			break
		}

		status = remotetrigger.TriggerMessageStatusAccepted

	case firmware.DiagnosticsStatusNotificationFeatureName, firmware.FirmwareStatusNotificationFeatureName:
		status = remotetrigger.TriggerMessageStatusNotImplemented

	case core.HeartbeatFeatureName:
		// Schedule to send a Heartbeat per request
		_, err = cp.scheduler.Every(5).Seconds().LimitRunsTo(1).Do(cp.sendHeartBeat)
		if err != nil {
			logInfo.WithError(err).Error("Cannot schedule a heartbeat")
			break
		}

		status = remotetrigger.TriggerMessageStatusAccepted

	case core.MeterValuesFeatureName:
		status = remotetrigger.TriggerMessageStatusNotImplemented
	case core.StatusNotificationFeatureName:

		switch request.ConnectorId {
		// Send the status of all connectors after the response
		case nil:
			defer func() {
				time.Sleep(time.Second)
				for _, c := range cp.evseManager.GetEVSEs() {
					cpStatus, errCode := c.GetStatus()
					go cp.notifyConnectorStatus(c.GetEvseId(), cpStatus, errCode)
				}
			}()

			status = remotetrigger.TriggerMessageStatusAccepted
		default:
			// Send a StatusNotification for a certain connector
			evseId := *request.ConnectorId
			c, findErr := cp.evseManager.FindEVSE(evseId)
			if findErr == nil {
				defer func(c evse.EVSE) {
					time.Sleep(time.Second)
					cpStatus, errCode := c.GetStatus()
					cp.notifyConnectorStatus(c.GetEvseId(), cpStatus, errCode)
				}(c)

				status = remotetrigger.TriggerMessageStatusAccepted
			}
		}

	default:
		return remotetrigger.NewTriggerMessageConfirmation(remotetrigger.TriggerMessageStatusNotImplemented), nil
	}

	return remotetrigger.NewTriggerMessageConfirmation(status), nil
}
