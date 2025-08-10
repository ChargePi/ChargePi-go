package v16

import (
	"go.uber.org/zap/zaptest"
	"testing"
	"time"

	"github.com/ChargePi/ChargePi-go/internal/pkg/scheduler"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/suite"
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
		logger:    zaptest.NewLogger(s.T()),
		scheduler: scheduler.NewScheduler(),
	}
}

func (s *reservationTestSuite) TestReservation() {
	var (
		expiryDate = types.NewDateTime(time.Now().Add(time.Minute))
	)

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
	response, err = s.cp.OnReserveNow(reservation.NewReserveNowRequest(connectorId, expiryDate, tagId, 2))
	s.Assert().NoError(err)
	s.Assert().NotNil(response)
	s.Assert().EqualValues(reservation.ReservationStatusRejected, response.Status)
}

func (s *reservationTestSuite) TestCancelReservation() {
	response, err := s.cp.OnCancelReservation(reservation.NewCancelReservationRequest(1))
	s.Assert().NoError(err)
	s.Assert().NotNil(response)
	s.Assert().EqualValues(reservation.CancelReservationStatusAccepted, response.Status)

	// Something went wrong with the reservation
	response, err = s.cp.OnCancelReservation(reservation.NewCancelReservationRequest(reservationId))
	s.Assert().NoError(err)
	s.Assert().NotNil(response)
	s.Assert().EqualValues(reservation.CancelReservationStatusRejected, response.Status)

	// No connector with the reservation
	response, err = s.cp.OnCancelReservation(reservation.NewCancelReservationRequest(2))
	s.Assert().NoError(err)
	s.Assert().NotNil(response)
	s.Assert().EqualValues(reservation.CancelReservationStatusRejected, response.Status)
}

func TestReservation(t *testing.T) {
	suite.Run(t, new(reservationTestSuite))
}
