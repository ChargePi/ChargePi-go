package v16

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type chargePointTestSuite struct {
	suite.Suite
	cp *ChargePoint
}

func (s *chargePointTestSuite) SetupTest() {
}

func (s *chargePointTestSuite) TestRestoreState() {
}

func (s *chargePointTestSuite) TestNotifyConnectorStatus() {
}

func TestChargePoint(t *testing.T) {
	suite.Run(t, new(chargePointTestSuite))
}
