package evse

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
)

func (m *managerImpl) FindEVSEWithReservationId(reservationId int) (EVSE, error) {
	evseId, isFound := m.reservations[reservationId]
	if !isFound || evseId == nil {
		return nil, ErrReservationNotFound
	}

	evse, err := m.FindEVSE(*evseId)
	if err != nil {
		return nil, err
	}

	return evse, nil
}

func (m *managerImpl) Reserve(evseId int, connectorId *int, reservationId int, tagId string) error {
	logInfo := log.WithFields(log.Fields{
		"evseId":        evseId,
		"tagId":         tagId,
		"reservationId": reservationId,
	})
	logInfo.Debugf("Reserving evse")

	evse, err := m.FindEVSE(evseId)
	if err != nil {
		return err
	}

	evse.SetStatus(core.ChargePointStatusReserved, core.NoError)
	return nil
}

func (m *managerImpl) RemoveReservation(reservationId int) error {
	_, isFound := m.reservations[reservationId]
	if !isFound {
		return ErrReservationNotFound
	}

	evse, err := m.FindEVSEWithReservationId(reservationId)
	if err != nil {
		return err
	}

	logInfo := log.WithField("reservationId", reservationId)
	logInfo.Debugf("Removing reservation")

	m.reservations[reservationId] = nil
	evse.SetStatus(core.ChargePointStatusAvailable, core.NoError)
	return nil
}
