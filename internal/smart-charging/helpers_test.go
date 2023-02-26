package smartCharging

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type helpersTestSuite struct {
	suite.Suite
}

func (s *helpersTestSuite) SetupTest() {
}

func (s *helpersTestSuite) TestRestoreState() {
}

func (s *helpersTestSuite) TestNotifyConnectorStatus() {
}

func TestHelpers(t *testing.T) {
	suite.Run(t, new(helpersTestSuite))
}
