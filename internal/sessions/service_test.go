package sessions

import (
	"testing"

	"go.uber.org/zap/zaptest"

	mock_sessions "github.com/ChargePi/ChargePi-go/gen/mocks/sessions"

	"github.com/stretchr/testify/suite"
)

type sessionServiceTestSuite struct {
	suite.Suite
	service               *Impl
	sessionRepositoryMock *mock_sessions.MockSessionRepository
}

func (s *sessionServiceTestSuite) SetupTest() {
	s.sessionRepositoryMock = mock_sessions.NewMockSessionRepository(s.T())

	service, err := NewSessionService(zaptest.NewLogger(s.T()), s.sessionRepositoryMock)
	s.Require().NoError(err)
	s.service = service
}

func (s *sessionServiceTestSuite) TestStartSession() {
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

func (s *sessionServiceTestSuite) TestStopSession() {
	tests := []struct {
		name string
	}{
		{
			name: "Session exists",
		},
		{
			name: "Session does not exist",
		},
		{
			name: "Database error",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *sessionServiceTestSuite) TestUpdateMeterValues() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *sessionServiceTestSuite) TestGetSession() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *sessionServiceTestSuite) TestGetSessionWithTransactionId() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *sessionServiceTestSuite) TestGetSessionWithTagId() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *sessionServiceTestSuite) TestAddTransactionIdToSession() {
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

func TestSession(t *testing.T) {
	suite.Run(t, new(sessionServiceTestSuite))
}
