package v16

import (
	"errors"

	"github.com/ChargePi/ChargePi-go/pkg/util"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
)

func (cp *ChargePoint) OnGetDiagnostics(request *firmware.GetDiagnosticsRequest) (confirmation *firmware.GetDiagnosticsConfirmation, err error) {
	go func() {
		err = cp.diagnosticsService.UploadLogs(request.Location, &request.StartTime.Time, &request.StopTime.Time, 0, 1)
	}()

	return firmware.NewGetDiagnosticsConfirmation(), nil
}

func (cp *ChargePoint) OnUpdateFirmware(request *firmware.UpdateFirmwareRequest) (confirmation *firmware.UpdateFirmwareConfirmation, err error) {
	// Check if it is run on docker - don't update.
	if util.IsRunningInContainer() {
		return firmware.NewUpdateFirmwareConfirmation(), errors.New("cannot update firmware on Docker")
	}

	// Todo
	return firmware.NewUpdateFirmwareConfirmation(), nil
}
