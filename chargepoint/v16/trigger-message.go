package v16

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/reactivex/rxgo/v2"
	log "github.com/sirupsen/logrus"
	connector2 "github.com/xBlaz3kx/ChargePi-go/components/connector"
	"github.com/xBlaz3kx/ChargePi-go/components/scheduler"
)

func (handler *ChargePointHandler) OnTriggerMessage(request *remotetrigger.TriggerMessageRequest) (confirmation *remotetrigger.TriggerMessageConfirmation, err error) {
	log.Infof("received %s for %v", request.GetFeatureName(), request.RequestedMessage)
	status := remotetrigger.TriggerMessageStatusRejected

	switch request.RequestedMessage {
	case core.BootNotificationFeatureName:
		_, err = scheduler.GetScheduler().Every(5).Seconds().LimitRunsTo(1).Do(handler.bootNotification)
		if err != nil {
			break
		}

		status = remotetrigger.TriggerMessageStatusAccepted
		break
	case firmware.DiagnosticsStatusNotificationFeatureName, firmware.FirmwareStatusNotificationFeatureName:
		status = remotetrigger.TriggerMessageStatusNotImplemented
		break
	case core.HeartbeatFeatureName:
		_, err = scheduler.GetScheduler().Every(5).Seconds().LimitRunsTo(1).Do(handler.sendHeartBeat)
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
				for _, connector := range handler.connectorManager.GetConnectors() {
					if handler.connectorChannel != nil {
						handler.connectorChannel <- rxgo.Of(connector)
					}
				}
			}()

			status = remotetrigger.TriggerMessageStatusAccepted
			break
		}

		connectorID := *request.ConnectorId
		connector := handler.connectorManager.FindConnector(1, connectorID)
		if connector != nil {
			defer func(connector connector2.Connector) {
				handler.notifyConnectorStatus(connector)
			}(connector)

			status = remotetrigger.TriggerMessageStatusAccepted
		}
		break
	default:
		return remotetrigger.NewTriggerMessageConfirmation(remotetrigger.TriggerMessageStatusNotImplemented), nil
	}

	return remotetrigger.NewTriggerMessageConfirmation(status), nil
}
