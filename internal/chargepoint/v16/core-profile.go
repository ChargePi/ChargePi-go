package v16

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	"go.uber.org/zap"
)

func (cp *ChargePoint) OnDataTransfer(request *core.DataTransferRequest) (confirmation *core.DataTransferConfirmation, err error) {
	cp.logger.Sugar().Infof("Received request %s", request.GetFeatureName())
	response := core.NewDataTransferConfirmation(core.DataTransferStatusRejected)

	// Supporting direct display control over custom data transfer messages, based on the messages in OCPP 2.0.1.
	settings, err := cp.settingsManager.GetChargePointSettings()
	if err != nil {
		return response, nil
	}

	if request.VendorId != settings.Info.OCPPDetails.Vendor {
		return core.NewDataTransferConfirmation(core.DataTransferStatusUnknownVendorId), nil
	}

	// Validate message id
	switch request.MessageId {
	case display.ClearDisplayMessageFeatureName:
		_ = request.Data.(display.ClearDisplayRequest)

		err := cp.display.Clear()
		if err != nil {
			cp.logger.With(zap.Error(err)).Warn("Failed to clear display")
		}
	case display.NotifyDisplayMessagesFeatureName:
		_ = request.Data.(display.NotifyDisplayMessagesRequest)
	case display.GetDisplayMessagesFeatureName:
		_ = request.Data.(display.GetDisplayMessagesRequest)
	case display.SetDisplayMessageFeatureName:
		req := request.Data.(display.SetDisplayMessageRequest)

		displayErr := cp.DisplayMessage(req.Message)
		if displayErr != nil {
			cp.logger.With(zap.Error(err)).Warn("Failed to display requested message")
		}
	default:
		response.Status = core.DataTransferStatusUnknownMessageId
	}

	return response, nil
}
