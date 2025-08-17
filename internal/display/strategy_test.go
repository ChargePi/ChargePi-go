package display

import (
	"testing"

	mock_display "github.com/ChargePi/ChargePi-go/gen/mocks/pkg/hardware/display"

	"github.com/stretchr/testify/suite"
)

type strategyTestSuite struct {
	suite.Suite
	displays []*mock_display.MockDisplay
}

func (s *strategyTestSuite) SetupTest() {
	s.displays = []*mock_display.MockDisplay{
		mock_display.NewMockDisplay(s.T()),
		mock_display.NewMockDisplay(s.T()),
		mock_display.NewMockDisplay(s.T()),
	}
}

func (s *strategyTestSuite) TestNewTranslatorStrategy() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *strategyTestSuite) TestNewDisplayOnlyTranslatorStrategy() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestTranslatorStrategy(t *testing.T) {
	suite.Run(t, new(strategyTestSuite))
}

func Test_displayConcurrently(t *testing.T) {
	// todo
}
