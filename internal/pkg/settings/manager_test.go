package settings

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type managerTestSuite struct {
	suite.Suite
}

func (s *managerTestSuite) SetupTest() {
}

func TestManager(t *testing.T) {
	suite.Run(t, new(managerTestSuite))
}
