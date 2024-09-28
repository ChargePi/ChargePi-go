package evse

import (
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/suite"
)

type stateTestSuite struct {
	suite.Suite
	state *state
}

func (s *stateTestSuite) SetupTest() {
	s.state = newState()
}

func (s *stateTestSuite) TearDownSuite() {}

func (s *stateTestSuite) TestGetAvailability() {
	tests := []struct {
		name                 string
		expectedAvailability core.AvailabilityType
	}{
		{
			name:                 "Availability should be inoperative by default",
			expectedAvailability: core.AvailabilityTypeInoperative,
		},
		{
			name:                 "Availability should be operative",
			expectedAvailability: core.AvailabilityTypeOperative,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			if tt.name == "Availability should be operative" {
				err := s.state.SetAvailability(core.AvailabilityTypeOperative)
				s.Require().NoError(err)
			}

			availability := s.state.GetAvailability()
			s.Equal(tt.expectedAvailability, availability)
		})
	}
}

func (s *stateTestSuite) TestSetAvailability() {
	tests := []struct {
		name         string
		availability core.AvailabilityType
		err          error
	}{
		{
			name:         "Set to inoperative",
			availability: core.AvailabilityTypeInoperative,
			err:          nil,
		},
		{
			name:         "Set to operative",
			availability: core.AvailabilityTypeOperative,
			err:          nil,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.state.SetAvailability(tt.availability)
			if tt.err != nil {
				s.Equal(tt.err, err)
				return
			}

			s.NoError(err)
			s.Equal(tt.availability, s.state.availability)
		})
	}
}

func (s *stateTestSuite) TestGetStatus() {
	tests := []struct {
		name              string
		expectedStatus    core.ChargePointStatus
		expectedErrorCode core.ChargePointErrorCode
	}{
		{
			name:              "Status and error code should be the default ones",
			expectedStatus:    core.ChargePointStatusUnavailable,
			expectedErrorCode: core.NoError,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			status, errorCode := s.state.GetStatus()
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedErrorCode, errorCode)
		})
	}
}

func (s *stateTestSuite) TestSetStatus() {
	tests := []struct {
		name      string
		status    core.ChargePointStatus
		errorCode core.ChargePointErrorCode
		error     error
	}{
		{
			name:      "Set status and error code",
			status:    core.ChargePointStatusAvailable,
			errorCode: core.NoError,
			error:     nil,
		},

		{
			name:      "Unable to set - invalid status",
			status:    "",
			errorCode: core.NoError,
			error:     nil,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.state.SetStatus(tt.status, tt.errorCode)
			if tt.error != nil {
				s.NoError(err)
				return
			}

			s.NoError(err)
			s.Equal(tt.status, s.state.status)
			s.Equal(tt.errorCode, s.state.errorCode)
		})
	}
}

func TestState(t *testing.T) {
	suite.Run(t, new(stateTestSuite))
}
