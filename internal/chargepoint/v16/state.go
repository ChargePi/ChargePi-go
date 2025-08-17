package v16

import (
	"sync"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)

type State struct {
	mu           sync.Mutex
	availability core.AvailabilityType
	isConnected  bool
}

func (s *State) SetAvailability(availability core.AvailabilityType) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.availability = availability
}

func (s *State) GetAvailability() core.AvailabilityType {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.availability
}

func (s *State) SetConnected(isConnected bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isConnected = isConnected
}

func (s *State) IsConnected() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isConnected
}
