package evse

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
)

func (m *managerImpl) GetEVSEWithReservationId(reservationId int) (EVSE, error) {
	logInfo := m.logger.WithField("reservationId", reservationId)
	logInfo.Debugf("Finding evse with reservation id")

	evseId, isFound := m.reservations[reservationId]
	if !isFound || evseId == nil {
		return nil, ErrReservationNotFound
	}

	evse, err := m.GetEVSE(*evseId)
	if err != nil {
		return nil, err
	}

	return evse, nil
}

func (m *managerImpl) Reserve(evseId int, connectorId *int, reservationId int, tagId string) error {
	logInfo := m.logger.WithFields(log.Fields{
		"evseId":        evseId,
		"tagId":         tagId,
		"reservationId": reservationId,
	})
	logInfo.Debugf("Reserving evse")

	evse, err := m.GetEVSE(evseId)
	if err != nil {
		return err
	}

	evse.SetStatus(core.ChargePointStatusReserved, core.NoError)
	return nil
}

func (m *managerImpl) RemoveReservation(reservationId int) error {
	logInfo := m.logger.WithField("reservationId", reservationId)
	logInfo.Debugf("Removing reservation")

	_, isFound := m.reservations[reservationId]
	if !isFound {
		return ErrReservationNotFound
	}

	evse, err := m.GetEVSEWithReservationId(reservationId)
	if err != nil {
		return err
	}

	m.reservations[reservationId] = nil
	evse.SetStatus(core.ChargePointStatusAvailable, core.NoError)
	return nil
}
