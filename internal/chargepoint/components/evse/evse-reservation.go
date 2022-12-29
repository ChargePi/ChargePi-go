package evse

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
)

func (evse *Impl) Reserve(reservationId int, tagId string) error {
	logInfo := log.WithFields(log.Fields{
		"evseId": evse.evseId,
		"tagId":  tagId,
	})
	logInfo.Debugf("Reserving evse for id %d", reservationId)

	if reservationId <= 0 {
		return ErrInvalidReservationId
	}

	if !evse.IsAvailable() {
		return ErrInvalidStatus
	}

	evse.reservationId = &reservationId
	evse.SetStatus(core.ChargePointStatusReserved, core.NoError)
	return nil
}

func (evse *Impl) RemoveReservation() error {
	if !evse.IsReserved() {
		return ErrInvalidStatus
	}

	logInfo := log.WithFields(log.Fields{
		"evseId": evse.evseId,
	})
	logInfo.Debugf("Removing reservation")

	evse.reservationId = nil
	evse.SetStatus(core.ChargePointStatusAvailable, core.NoError)
	return nil
}

func (evse *Impl) GetReservationId() int {
	if util.IsNilInterfaceOrPointer(evse.reservationId) {
		return -1
	}

	return *evse.reservationId
}
