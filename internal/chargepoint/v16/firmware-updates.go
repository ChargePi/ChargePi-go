package v16

import "github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"

func (cp *ChargePoint) OnGetDiagnostics(request *firmware.GetDiagnosticsRequest) (confirmation *firmware.GetDiagnosticsConfirmation, err error) {

	// cp.diagnosticsManager.UploadLogs(request.Location, request.Retries, request.RetryInterval)
	return firmware.NewGetDiagnosticsConfirmation(), nil
}

func (cp *ChargePoint) OnUpdateFirmware(request *firmware.UpdateFirmwareRequest) (confirmation *firmware.UpdateFirmwareConfirmation, err error) {
	return firmware.NewUpdateFirmwareConfirmation(), nil
}
