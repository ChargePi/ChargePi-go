package v16

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/reactivex/rxgo/v2"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/connector"
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
		break
	case firmware.DiagnosticsStatusNotificationFeatureName, firmware.FirmwareStatusNotificationFeatureName:
		status = remotetrigger.TriggerMessageStatusNotImplemented
		break
	case core.HeartbeatFeatureName:
		_, err = cp.scheduler.Every(5).Seconds().LimitRunsTo(1).Do(cp.sendHeartBeat)
		if err != nil {
			break
		}

		status = remotetrigger.TriggerMessageStatusAccepted
		break
	case core.MeterValuesFeatureName:
		status = remotetrigger.TriggerMessageStatusNotImplemented
		break
	case core.StatusNotificationFeatureName:
		if request.ConnectorId == nil {
			// Send the status of all connectors after the response
			defer func() {
				for _, c := range cp.connectorManager.GetConnectors() {
					if cp.connectorChannel != nil {
						cp.connectorChannel <- rxgo.Of(c)
					}
				}
			}()

			status = remotetrigger.TriggerMessageStatusAccepted
			break
		}

		connectorID := *request.ConnectorId
		c := cp.connectorManager.FindConnector(1, connectorID)
		if c != nil {
			defer func(c connector.Connector) {
				cp.notifyConnectorStatus(c)
			}(c)

			status = remotetrigger.TriggerMessageStatusAccepted
		}
		break
	default:
		return remotetrigger.NewTriggerMessageConfirmation(remotetrigger.TriggerMessageStatusNotImplemented), nil
	}

	return remotetrigger.NewTriggerMessageConfirmation(status), nil
}
