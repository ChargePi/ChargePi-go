package v16

import (
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
)

func (cp *ChargePoint) OnTriggerMessage(request *remotetrigger.TriggerMessageRequest) (confirmation *remotetrigger.TriggerMessageConfirmation, err error) {
	logInfo := cp.logger.WithFields(log.Fields{"feature": request.GetFeatureName(), "request": request.RequestedMessage})
	logInfo.Infof("Received a request")

	status := remotetrigger.TriggerMessageStatusRejected

	switch request.RequestedMessage {
	case core.BootNotificationFeatureName:

		// Send a BootNotification after the response
		defer cp.bootNotification()
		status = remotetrigger.TriggerMessageStatusAccepted

	case firmware.DiagnosticsStatusNotificationFeatureName, firmware.FirmwareStatusNotificationFeatureName:
		status = remotetrigger.TriggerMessageStatusNotImplemented

	case core.HeartbeatFeatureName:

		// Send a Heartbeat after the response
		defer cp.sendHeartBeat()
		status = remotetrigger.TriggerMessageStatusAccepted

	case core.MeterValuesFeatureName:

		switch request.ConnectorId {
		case nil:

			// Send the status of all connectors after the response
			defer func() {
				for _, evse := range cp.evseManager.GetEVSEs() {
					_ = cp.getMeasurements(evse.GetEvseId())
				}
			}()

			status = remotetrigger.TriggerMessageStatusAccepted
		default:

			// Send a MeterValues for a certain connector
			defer cp.getMeasurements(*request.ConnectorId)

			status = remotetrigger.TriggerMessageStatusAccepted
		}

	case core.StatusNotificationFeatureName:

		switch request.ConnectorId {
		case nil:
			logInfo.Infof("Sending a status update for all connectors")

			// Send the status of all connectors after the response
			defer func() {
				time.Sleep(time.Second)
				for _, c := range cp.evseManager.GetEVSEs() {
					cpStatus, errCode := c.GetStatus()
					go cp.notifyConnectorStatus(c.GetEvseId(), cpStatus, errCode)
				}
			}()

			status = remotetrigger.TriggerMessageStatusAccepted

		default:

			logInfo.Infof("Sending a status update for a connector")

			// Send a StatusNotification for a certain connector
			c, findErr := cp.evseManager.GetEVSE(*request.ConnectorId)
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

func (cp *ChargePoint) getMeasurements(evseId int) error {
	requestedEvse, err := cp.evseManager.GetEVSE(evseId)
	if err != nil {
		return err
	}

	// Todo sample
	_ = requestedEvse.SamplePowerMeter(util.GetTypesToSample())
	return nil
}
