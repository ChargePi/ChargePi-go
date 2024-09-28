package display

import (
	"testing"

	"github.com/ChargePi/ChargePi-go/pkg/hardware/display/mocks"
	"github.com/stretchr/testify/suite"
)

type strategyTestSuite struct {
	suite.Suite
	displays []*mocks.MockDisplay
}

func (s *strategyTestSuite) SetupTest() {
	s.displays = []*mocks.MockDisplay{
		mocks.NewMockDisplay(s.T()),
		mocks.NewMockDisplay(s.T()),
		mocks.NewMockDisplay(s.T()),
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
