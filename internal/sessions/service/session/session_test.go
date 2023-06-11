package session

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type sessionServiceTestSuite struct {
	suite.Suite
}

func (s *sessionServiceTestSuite) SetupTest() {

}

func TestSession(t *testing.T) {
	suite.Run(t, new(sessionServiceTestSuite))
}
