package manager

import (
	mock_manager "github.com/ChargePi/ChargePi-go/gen/mocks/evse/manager"
	mock_sessions "github.com/ChargePi/ChargePi-go/gen/mocks/sessions"
	"testing"

	"github.com/ChargePi/ChargePi-go/internal/pkg/notifications"
	"github.com/ChargePi/ChargePi-go/internal/sessions/mocks"
	"github.com/stretchr/testify/suite"
)

type managerTestSuite struct {
	suite.Suite
	manager             *Impl
	sessionRepository   *mock_sessions.MockSessionRepository
	evseRepository      *mock_manager.MockEvseSettingsRepository
	notificationChannel chan notifications.StatusNotification
}

func (s *managerTestSuite) SetupTest() {
	s.sessionRepository = mock_sessions.NewMockSessionRepository(s.T())
	s.evseRepository = mock_manager.NewMockEvseSettingsRepository(s.T())
	s.notificationChannel = make(chan notifications.StatusNotification)

	manager, err := NewManager(s.sessionRepository, s.evseRepository, s.notificationChannel)
	s.Require().NoError(err)

	s.manager = manager
}

func (s *managerTestSuite) TearDownTest() {
	close(s.notificationChannel)
}

func (s *managerTestSuite) TestInitAll() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestAddEVSE() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestAddEVSEFromSettings() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestGetCurrentConsumption() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestUnlockConnector() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestFindEVSE() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestFindAvailableEVSE() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestFindEVSEWithReservationId() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestFindEVSEWithTransactionId() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestFindEVSEWithTagId() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestStartCharging() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestStopCharging() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestStopAll() {
	tests := []struct {
		name string
	}{{}}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestRestoreEVSEs() {
	tests := []struct {
		name string
	}{
		{},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestReserve() {
	tests := []struct {
		name string
	}{
		{},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *managerTestSuite) TestRemoveReservation() {
	tests := []struct {
		name          string
		reservationId int
		error         error
	}{
		{
			name: "Remove reservation",
		},
		{
			name: "Reservation does not exist",
		},
		{
			name: "Reservation exists but no EVSE details",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.manager.RemoveReservation(tt.reservationId)
			if tt.error != nil {
				s.ErrorIs(err, tt.error)
			} else {
				s.NoError(err)
			}
		})
	}
}

func TestManager(t *testing.T) {
	suite.Run(t, new(managerTestSuite))
}
