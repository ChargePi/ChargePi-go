package v16

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
)

func (cp *ChargePoint) OnReserveNow(request *reservation.ReserveNowRequest) (confirmation *reservation.ReserveNowConfirmation, err error) {
	cp.logger.Infof("Received %s for %v", request.GetFeatureName(), request.ConnectorId)

	connector, err := cp.evseManager.FindEVSE(request.ConnectorId)
	if err != nil {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusUnavailable), nil
	}

	err = connector.Reserve(request.ReservationId, request.IdTag)
	switch err {
	case nil:
		timeFormat := fmt.Sprintf("%d:%d", request.ExpiryDate.Hour(), request.ExpiryDate.Minute())
		_, schedulerErr := cp.scheduler.Every(1).Day().At(timeFormat).LimitRunsTo(1).Do(connector.RemoveReservation)
		if schedulerErr != nil {
			return reservation.NewReserveNowConfirmation(reservation.ReservationStatusRejected), nil
		}

		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusAccepted), nil
	case evse.ErrConnectorStatusInvalid:
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusOccupied), nil
	default:
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusRejected), nil
	}
}

func (cp *ChargePoint) OnCancelReservation(request *reservation.CancelReservationRequest) (confirmation *reservation.CancelReservationConfirmation, err error) {
	cp.logger.Infof("Received %s for %v", request.GetFeatureName(), request.ReservationId)
	var (
		connector, rErr = cp.evseManager.FindEVSEWithReservationId(request.ReservationId)
		status          = reservation.CancelReservationStatusAccepted
	)

	if rErr != nil {
		return reservation.NewCancelReservationConfirmation(reservation.CancelReservationStatusRejected), nil
	}

	err = connector.RemoveReservation()
	if err != nil {
		status = reservation.CancelReservationStatusRejected
	}

	return reservation.NewCancelReservationConfirmation(status), nil
}
