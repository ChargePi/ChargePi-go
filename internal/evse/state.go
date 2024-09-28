package evse

import (
	"sync"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)

type state struct {
	mu sync.Mutex

	// availability represents the availability of the charge point.
	availability core.AvailabilityType

	// status represents the status of the charge point.
	status core.ChargePointStatus

	// errorCode represents the error code of the charge point.
	errorCode core.ChargePointErrorCode
}

func newState() *state {
	return &state{
		availability: core.AvailabilityTypeInoperative,
		status:       core.ChargePointStatusUnavailable,
		errorCode:    core.NoError,
	}
}

func (s *state) GetAvailability() core.AvailabilityType {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.availability
}

func (s *state) SetAvailability(availability core.AvailabilityType) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// todo validate?

	s.availability = availability
	return nil
}

func (s *state) GetStatus() (core.ChargePointStatus, core.ChargePointErrorCode) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.status, s.errorCode
}

func (s *state) SetStatus(status core.ChargePointStatus, errorCode core.ChargePointErrorCode) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// todo validate?

	s.status = status
	s.errorCode = errorCode
	return nil
}
