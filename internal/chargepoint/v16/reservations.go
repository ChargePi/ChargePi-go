package v16

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/xBlaz3kx/ChargePi-go/pkg/util"
)

func (cp *ChargePoint) OnReserveNow(request *reservation.ReserveNowRequest) (confirmation *reservation.ReserveNowConfirmation, err error) {
	cp.logger.Infof("Received %s for %v", request.GetFeatureName(), request.ConnectorId)
	connector := cp.connectorManager.FindConnector(1, request.ConnectorId)

	if util.IsNilInterfaceOrPointer(connector) {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusUnavailable), nil
	} else if !connector.IsAvailable() {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusOccupied), nil
	}

	err = connector.ReserveConnector(request.ReservationId)
	if err != nil {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusRejected), nil
	}

	_, schedulerErr := cp.scheduler.At(request.ExpiryDate.Format("HH:MM")).Do(connector.RemoveReservation)
	if schedulerErr != nil {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusRejected), nil
	}

	return reservation.NewReserveNowConfirmation(reservation.ReservationStatusAccepted), nil
}

func (cp *ChargePoint) OnCancelReservation(request *reservation.CancelReservationRequest) (confirmation *reservation.CancelReservationConfirmation, err error) {
	cp.logger.Infof("Received %s for %v", request.GetFeatureName(), request.ReservationId)
	var (
		connector = cp.connectorManager.FindConnectorWithReservationId(request.ReservationId)
		status    = reservation.CancelReservationStatusAccepted
	)

	if util.IsNilInterfaceOrPointer(connector) {
		return reservation.NewCancelReservationConfirmation(reservation.CancelReservationStatusRejected), nil
	}

	err = connector.RemoveReservation()
	if err != nil {
		status = reservation.CancelReservationStatusRejected
	}

	return reservation.NewCancelReservationConfirmation(status), nil
}
