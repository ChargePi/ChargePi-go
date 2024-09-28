package sessions

import (
	"testing"

	"github.com/ChargePi/ChargePi-go/internal/sessions/mocks"
	"github.com/stretchr/testify/suite"
)

type sessionServiceTestSuite struct {
	suite.Suite
	service               *Impl
	sessionRepositoryMock *mocks.MockSessionRepository
}

func (s *sessionServiceTestSuite) SetupTest() {
	s.sessionRepositoryMock = mocks.NewMockSessionRepository(s.T())

	service, err := NewSessionService(s.sessionRepositoryMock)
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
