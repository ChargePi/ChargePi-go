package manager

import (
	"github.com/ChargePi/ChargePi-go/internal/evse"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"go.uber.org/zap"
)

func (m *managerImpl) GetEVSEWithReservationId(reservationId int) (evse.EVSE, error) {
	logInfo := m.logger.With(zap.Int("reservationId", reservationId))
	logInfo.Debug("Finding evse with reservation id")

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
	logInfo := m.logger.With(
		zap.Int("evseId", evseId),
		zap.Int("reservationId", reservationId),
		zap.String("tagId", tagId))
	logInfo.Debug("Reserving evse")

	evse, err := m.GetEVSE(evseId)
	if err != nil {
		return err
	}

	evse.SetStatus(core.ChargePointStatusReserved, core.NoError)
	return nil
}

func (m *managerImpl) RemoveReservation(reservationId int) error {
	logInfo := m.logger.With(zap.Int("reservationId", reservationId))
	logInfo.Debug("Removing reservation")

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
