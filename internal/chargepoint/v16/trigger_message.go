package v16

import (
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/ChargePi/ocpp-manager/ocpp_v16"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"

	"github.com/ChargePi/ChargePi-go/internal/evse"
	"github.com/ChargePi/ChargePi-go/internal/pkg/notifications"
)

func (cp *ChargePoint) OnTriggerMessage(request *remotetrigger.TriggerMessageRequest) (confirmation *remotetrigger.TriggerMessageConfirmation, err error) {
	logInfo := cp.logger.With(zap.String("feature", remotetrigger.TriggerMessageFeatureName))
	logInfo.Info("Received a request")

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
			logInfo.Info("Sending a status update for all connectors")

			// Send the status of all connectors after the response
			defer func() {
				time.Sleep(time.Second)
				for _, c := range cp.evseManager.GetEVSEs() {
					cpStatus, errCode := c.GetStatus()
					go cp.notifyStatus(c.GetEvseId(), cpStatus, errCode)
				}
			}()

			status = remotetrigger.TriggerMessageStatusAccepted

		default:

			logInfo.Info("Sending a status update for a connector")

			// Send a StatusNotification for a certain connector
			c, findErr := cp.evseManager.GetEVSE(*request.ConnectorId)
			if findErr == nil {
				defer func(c evse.EVSE) {
					time.Sleep(time.Second)
					cpStatus, errCode := c.GetStatus()
					cp.notifyStatus(c.GetEvseId(), cpStatus, errCode)
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
	// Get EVSE
	requestedEvse, err := cp.evseManager.GetEVSE(evseId)
	if err != nil {
		return err
	}

	// Get measurands from OCPP configuration
	value, err := cp.settingsManager.GetConfigurationValue(ocpp_v16.MeterValuesSampledData)
	if err != nil {
		return err
	}

	var measurands []types.Measurand
	for _, measurand := range strings.Split(*value, ",") {
		measurands = append(measurands, types.Measurand(measurand))
	}

	// Sample the power meter
	meter, err := requestedEvse.SamplePowerMeter(measurands)
	if err != nil {
		return err
	}

	// Send the meter values
	if cp.meterValuesChannel != nil {
		meterValues := types.MeterValue{
			Timestamp:    types.NewDateTime(time.Now()),
			SampledValue: meter,
		}

		cp.meterValuesChannel <- notifications.MeterValueNotification{MeterValues: []types.MeterValue{meterValues}, EvseId: evseId}
	}

	return nil
}
