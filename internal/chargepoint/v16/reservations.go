package v16

import (
	"errors"
	"fmt"
	"github.com/ChargePi/ChargePi-go/internal/evse/manager"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
)

func (cp *ChargePoint) OnReserveNow(request *reservation.ReserveNowRequest) (confirmation *reservation.ReserveNowConfirmation, err error) {
	cp.logger.Infof("Received %s for %v", request.GetFeatureName(), request.ConnectorId)

	err = cp.evseManager.Reserve(request.ConnectorId, nil, request.ReservationId, request.IdTag)
	switch {
	case err == nil:
		timeFormat := fmt.Sprintf("%d:%d", request.ExpiryDate.Hour(), request.ExpiryDate.Minute())
		_, schedulerErr := cp.scheduler.Every(1).Day().At(timeFormat).LimitRunsTo(1).Do(cp.evseManager.RemoveReservation, request.ReservationId)
		if schedulerErr != nil {
			return reservation.NewReserveNowConfirmation(reservation.ReservationStatusRejected), nil
		}

		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusAccepted), nil
	case errors.Is(err, manager.ErrConnectorStatusInvalid):
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusOccupied), nil
	default:
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusRejected), nil
	}
}

func (cp *ChargePoint) OnCancelReservation(request *reservation.CancelReservationRequest) (confirmation *reservation.CancelReservationConfirmation, err error) {
	cp.logger.Infof("Received %s for %v", request.GetFeatureName(), request.ReservationId)
	status := reservation.CancelReservationStatusAccepted

	err = cp.evseManager.RemoveReservation(request.ReservationId)
	switch err {
	case nil:
	default:
		status = reservation.CancelReservationStatusRejected
	}

	return reservation.NewCancelReservationConfirmation(status), nil
}
