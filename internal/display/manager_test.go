package display

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ChargePi/ChargePi-go/internal/display/i18n"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/display/mocks"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	"github.com/stretchr/testify/suite"
)

type displayManagerTestSuite struct {
	suite.Suite
	manager      *DisplayManager
	displayMocks map[string]*mocks.MockDisplay
}

func (s *displayManagerTestSuite) SetupSuite() {
	manager, err := NewDisplayManager()
	s.Require().NoError(err)
	s.manager = manager
	s.displayMocks = make(map[string]*mocks.MockDisplay)
}

func (s *displayManagerTestSuite) TearDownSuite() {}

func (s *displayManagerTestSuite) SetupTest() {
	s.manager.messageQueue = newMessageQueue()
	s.manager.scheduler.Clear()

	// Create 3 mock displays
	for i := 0; i < 3; i++ {
		mockDisplay := mocks.NewMockDisplay(s.T())
		s.manager.displays[fmt.Sprintf("%d", i)] = mockDisplay
		s.displayMocks[fmt.Sprintf("%d", i)] = mockDisplay
	}

	err := s.manager.SetStrategy(&displayOnlyStrategy{})
	s.Require().NoError(err)
}

func (s *displayManagerTestSuite) TestDisplayMessage() {
	tests := []struct {
		name    string
		message display.MessageInfo
		err     error
	}{
		{
			name: "Successfully displayed message",
		},
		{
			name: "Display message",
		},
		{
			name: "Display message",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *displayManagerTestSuite) TestStrategy() {
	translator, err := i18n.NewTranslator(i18n.Settings{})
	s.Require().NoError(err)

	tests := []struct {
		name     string
		strategy Strategy
		message  display.MessageInfo
		err      error
	}{
		{
			name:     "Display message using displayOnly strategy",
			strategy: &displayOnlyStrategy{},
		},
		{
			name: "Display message using translation strategy",
			strategy: &translateAndDisplayStrategy{
				translator: translator,
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.manager.SetStrategy(tt.strategy)
			s.Require().NoError(err)

			err = s.manager.DisplayMessage(tt.message)
			if tt.err != nil {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *displayManagerTestSuite) TestRemoveMessage() {
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

func (s *displayManagerTestSuite) TestGetCurrentMessage() {
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

func (s *displayManagerTestSuite) TestCleanup() {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Successful cleanup",
			err:  nil,
		},
		{
			name: "Cleanup failed - scheduler error",
			err:  nil,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err := s.manager.Cleanup(ctx)
			if tt.err != nil {
				s.ErrorIs(err, tt.err)
			} else {
				s.NoError(err)

				// Verify that the scheduler is not running and has no jobs
				s.Empty(s.manager.scheduler.Jobs())
				s.False(s.manager.scheduler.IsRunning())

				// Verify that all displays have been cleaned up
				for _, mock := range s.displayMocks {
					mock.AssertNumberOfCalls(t, "Cleanup", 1)
				}
			}
		})
	}
}

func TestDisplayManager(t *testing.T) {
	suite.Run(t, new(displayManagerTestSuite))
}
