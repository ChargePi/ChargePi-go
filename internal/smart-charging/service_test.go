package smartCharging

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type smartChargingManagerTestSuite struct {
	suite.Suite
}

func (s *smartChargingManagerTestSuite) SetupTest() {
}

func (s *smartChargingManagerTestSuite) TestRestoreState() {
}

func (s *smartChargingManagerTestSuite) TestNotifyConnectorStatus() {
}

func TestSmartChargingManager(t *testing.T) {
	suite.Run(t, new(smartChargingManagerTestSuite))
}
