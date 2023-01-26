package v16

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
)

func (cp *ChargePoint) OnTriggerMessage(request *remotetrigger.TriggerMessageRequest) (confirmation *remotetrigger.TriggerMessageConfirmation, err error) {
	cp.logger.Infof("received %s for %v", request.GetFeatureName(), request.RequestedMessage)
	status := remotetrigger.TriggerMessageStatusRejected

	switch request.RequestedMessage {
	case core.BootNotificationFeatureName:
		_, err = cp.scheduler.Every(5).Seconds().LimitRunsTo(1).Do(cp.bootNotification)
		if err != nil {
			break
		}

		status = remotetrigger.TriggerMessageStatusAccepted
	case firmware.DiagnosticsStatusNotificationFeatureName, firmware.FirmwareStatusNotificationFeatureName:
		status = remotetrigger.TriggerMessageStatusNotImplemented
	case core.HeartbeatFeatureName:
		_, err = cp.scheduler.Every(5).Seconds().LimitRunsTo(1).Do(cp.sendHeartBeat)
		if err != nil {
			break
		}

		status = remotetrigger.TriggerMessageStatusAccepted
	case core.MeterValuesFeatureName:
		status = remotetrigger.TriggerMessageStatusNotImplemented
	case core.StatusNotificationFeatureName:
		if request.ConnectorId == nil {
			// Send the status of all connectors after the response
			defer func() {
				for _, c := range cp.evseManager.GetEVSEs() {
					if cp.evseManager.GetNotificationChannel() != nil {
						cpStatus, errCode := c.GetStatus()
						go cp.notifyConnectorStatus(c.GetEvseId(), cpStatus, errCode)
					}
				}
			}()

			status = remotetrigger.TriggerMessageStatusAccepted
			break
		}

		connectorID := *request.ConnectorId
		c, err := cp.evseManager.FindEVSE(connectorID)
		if err == nil {
			defer func(c evse.EVSE) {
				cpStatus, errCode := c.GetStatus()
				cp.notifyConnectorStatus(c.GetEvseId(), cpStatus, errCode)
			}(c)

			status = remotetrigger.TriggerMessageStatusAccepted
		}
	default:
		return remotetrigger.NewTriggerMessageConfirmation(remotetrigger.TriggerMessageStatusNotImplemented), nil
	}

	return remotetrigger.NewTriggerMessageConfirmation(status), nil
}
