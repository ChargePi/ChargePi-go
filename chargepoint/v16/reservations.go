package v16

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/xBlaz3kx/ChargePi-go/components/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"log"
)

func (handler *ChargePointHandler) OnReserveNow(request *reservation.ReserveNowRequest) (confirmation *reservation.ReserveNowConfirmation, err error) {
	log.Printf("Received %s for %v", request.GetFeatureName(), request.ConnectorId)
	connector := handler.connectorManager.FindConnector(1, request.ConnectorId)

	if data.IsNilInterfaceOrPointer(connector) {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusUnavailable), nil
	} else if !connector.IsAvailable() {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusOccupied), nil
	}

	err = connector.ReserveConnector(request.ReservationId)
	if err != nil {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusRejected), nil
	}

	_, schedulerErr := scheduler.GetScheduler().At(request.ExpiryDate.Format("HH:MM")).Do(connector.RemoveReservation)
	if schedulerErr != nil {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusRejected), nil
	}

	return reservation.NewReserveNowConfirmation(reservation.ReservationStatusAccepted), nil
}

func (handler *ChargePointHandler) OnCancelReservation(request *reservation.CancelReservationRequest) (confirmation *reservation.CancelReservationConfirmation, err error) {
	log.Printf("Received %s for %v", request.GetFeatureName(), request.ReservationId)
	var (
		connector = handler.connectorManager.FindConnectorWithReservationId(request.ReservationId)
		status    = reservation.CancelReservationStatusAccepted
	)

	if data.IsNilInterfaceOrPointer(connector) {
		return reservation.NewCancelReservationConfirmation(reservation.CancelReservationStatusRejected), nil
	}

	err = connector.RemoveReservation()
	if err != nil {
		status = reservation.CancelReservationStatusRejected
	}

	return reservation.NewCancelReservationConfirmation(status), nil
}
