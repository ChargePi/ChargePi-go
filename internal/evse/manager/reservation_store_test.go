package manager

import (
	"testing"

	"github.com/ChargePi/ChargePi-go/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type reservationTestSuite struct {
	suite.Suite
	store *reservations
}

func (s *reservationTestSuite) SetupTest() {
	s.store = newReservationsStore()
}

func (s *reservationTestSuite) TestGetEVSEWithReservationId() {
	tests := []struct {
		name                string
		reservationId       int
		expectedReservation reservation
		wantErr             bool
		error               error
	}{
		{
			name:          "Valid reservation",
			reservationId: 1,
			expectedReservation: reservation{
				EvseId:      1,
				ConnectorId: nil,
				TagId:       "1",
			},
			wantErr: false,
		},
		{
			name:          "Reservation not found",
			reservationId: 2,
			wantErr:       true,
			error:         ErrReservationNotFound,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			if tt.name == "Valid reservation" {
				s.store.reservations[tt.reservationId] = tt.expectedReservation
			}

			reservation, err := s.store.GetEVSEWithReservationId(tt.reservationId)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.error)
				assert.Nil(t, reservation)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, reservation)
			assert.Equal(t, tt.expectedReservation, *reservation)
		})
	}
}

func (s *reservationTestSuite) TestReserve() {
	tests := []struct {
		name          string
		reservationId int
		reservation   reservation
		error         error
	}{
		{
			name:          "Valid reservation",
			reservationId: 1,
			reservation: reservation{
				EvseId:      1,
				ConnectorId: nil,
				TagId:       util.GenerateRandomTag(),
			},
			error: nil,
		},
		{
			name:          "Reservation already present",
			reservationId: 1,
			reservation: reservation{
				EvseId:      2,
				ConnectorId: nil,
				TagId:       util.GenerateRandomTag(),
			},
			error: ErrReservationNotFound,
		},
		{
			name:          "EVSE already reserved",
			reservationId: 2,
			reservation: reservation{
				EvseId:      1,
				ConnectorId: nil,
				TagId:       util.GenerateRandomTag(),
			},
			error: ErrInvalidReservation,
		},
		{
			name:          "Validation failed- Invalid evse ID",
			reservationId: 2,
			reservation: reservation{
				EvseId:      -1,
				ConnectorId: nil,
				TagId:       util.GenerateRandomTag(),
			},
			error: ErrValidationFailed,
		},
		{
			name:          "Validation failed - Invalid tag ID",
			reservationId: 2,
			reservation: reservation{
				EvseId:      2,
				ConnectorId: nil,
				TagId:       "",
			},
			error: ErrValidationFailed,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.store.Reserve(tt.reservationId, tt.reservation)
			if tt.error != nil {
				assert.ErrorIs(t, err, tt.error)
				return
			}

			assert.NoError(t, err)
			assert.Contains(t, s.store.reservations, tt.reservationId)
			assert.Equal(t, s.store.reservations[tt.reservationId], tt.reservation)
		})
	}
}

func (s *reservationTestSuite) TestRemoveReservation() {
	tests := []struct {
		name          string
		reservationId int
		reservation   reservation
		error         error
	}{
		{
			name:          "Reservation present",
			reservationId: 1,
			reservation: reservation{
				EvseId:      1,
				ConnectorId: nil,
				TagId:       util.GenerateRandomTag(),
			},
			error: nil,
		},
		{
			name:          "Reservation not present",
			reservationId: 1,
			reservation: reservation{
				EvseId:      2,
				ConnectorId: nil,
				TagId:       util.GenerateRandomTag(),
			},
			error: ErrReservationNotFound,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			if tt.name == "Reservation present" {
				s.store.reservations[tt.reservationId] = tt.reservation
			}

			err := s.store.RemoveReservation(tt.reservationId)
			if tt.error != nil {
				assert.ErrorIs(t, err, tt.error)
				return
			}

			assert.NoError(t, err)
			assert.NotContains(t, s.store.reservations, tt.reservationId)
			assert.Empty(t, s.store.reservations[tt.reservationId])
		})
	}
}

func TestReservations(t *testing.T) {
	suite.Run(t, new(reservationTestSuite))
}
