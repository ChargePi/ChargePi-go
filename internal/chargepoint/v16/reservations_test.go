package v16

import (
	"errors"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/test"
	"testing"
	"time"
)

const (
	reservationId = 1
	connectorId   = 1
	tagId         = "exampleTagId"
)

type reservationTestSuite struct {
	suite.Suite
	cp *ChargePoint
}

func (s *reservationTestSuite) SetupTest() {
	s.cp = &ChargePoint{
		logger:    log.StandardLogger(),
		scheduler: scheduler.GetScheduler(),
	}
}

func (s *reservationTestSuite) TestReservation() {
	var (
		connectorMock = new(test.EvseMock)
		managerMock   = new(test.ManagerMock)
		expiryDate    = types.NewDateTime(time.Now().Add(time.Minute))
	)

	// Set connector expectations
	connectorMock.On("ReserveEvse", reservationId, tagId).Return(nil).Once()
	connectorMock.On("IsAvailable").Return(true).Once()
	connectorMock.On("RemoveReservation").Return()

	// Set manager expectations
	managerMock.On("FindEVSE", connectorId).Return(connectorMock).Twice()
	// Connector not found
	managerMock.On("FindEVSE", 2).Return(nil).Once()
	s.cp.connectorManager = managerMock

	response, err := s.cp.OnReserveNow(reservation.NewReserveNowRequest(connectorId, expiryDate, tagId, reservationId))
	s.Assert().NoError(err)
	s.Assert().NotNil(response)
	s.Assert().EqualValues(reservation.ReservationStatusAccepted, response.Status)

	// No connector with connectorId
	response, err = s.cp.OnReserveNow(reservation.NewReserveNowRequest(2, expiryDate, tagId, reservationId))
	s.Assert().NoError(err)
	s.Assert().NotNil(response)
	s.Assert().EqualValues(reservation.ReservationStatusUnavailable, response.Status)

	// Unable to reserve for whatever reason
	connectorMock.On("IsAvailable").Return(true).Once()
	connectorMock.On("ReserveEvse", 2, tagId).Return(errors.New("unable to reserve the connector")).Once()
	response, err = s.cp.OnReserveNow(reservation.NewReserveNowRequest(connectorId, expiryDate, tagId, 2))
	s.Assert().NoError(err)
	s.Assert().NotNil(response)
	s.Assert().EqualValues(reservation.ReservationStatusRejected, response.Status)
}

func (s *reservationTestSuite) TestCancelReservation() {
	var (
		connectorMock = new(test.EvseMock)
		managerMock   = new(test.ManagerMock)
	)

	// Set connector expectations
	connectorMock.On("RemoveReservation").Return(nil).Once()

	// Set manager expectations
	managerMock.On("FindEVSEWithReservationId", reservationId).Return(connectorMock)
	s.cp.connectorManager = managerMock

	response, err := s.cp.OnCancelReservation(reservation.NewCancelReservationRequest(1))
	s.Assert().NoError(err)
	s.Assert().NotNil(response)
	s.Assert().EqualValues(reservation.CancelReservationStatusAccepted, response.Status)

	// Something went wrong with the reservation
	connectorMock.On("RemoveReservation").Return(errors.New("something went wrong")).Once()
	response, err = s.cp.OnCancelReservation(reservation.NewCancelReservationRequest(reservationId))
	s.Assert().NoError(err)
	s.Assert().NotNil(response)
	s.Assert().EqualValues(reservation.CancelReservationStatusRejected, response.Status)

	// No connector with the reservation
	managerMock.On("FindEVSEWithReservationId", 2).Return(nil)
	response, err = s.cp.OnCancelReservation(reservation.NewCancelReservationRequest(2))
	s.Assert().NoError(err)
	s.Assert().NotNil(response)
	s.Assert().EqualValues(reservation.CancelReservationStatusRejected, response.Status)
}

func TestReservation(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(reservationTestSuite))
}
