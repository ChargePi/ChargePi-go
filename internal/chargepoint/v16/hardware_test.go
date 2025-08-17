package v16

import (
	"testing"
	"time"

	"go.uber.org/zap/zaptest"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/suite"

	"github.com/ChargePi/ChargePi-go/pkg/hardware/indicator"
)

const (
	exampleMessage  = "exampleMessage"
	exampleMessage1 = "exampleMessage2"
)

type hardwareTestSuite struct {
	suite.Suite
	cp *ChargePoint
}

func (s *hardwareTestSuite) SetupTest() {
	s.cp = new(ChargePoint)
	s.cp.logger = zaptest.NewLogger(s.T())
}

func (s *hardwareTestSuite) TestSendToLCD() {
}

func (s *hardwareTestSuite) TestDisplayLedStatus() {
	// Ok statuses
	s.cp.indicateStatusChange(1, core.ChargePointStatusCharging)
	s.cp.indicateStatusChange(1, core.ChargePointStatusFinishing)
	s.cp.indicateStatusChange(1, core.ChargePointStatusAvailable)
	s.cp.indicateStatusChange(1, core.ChargePointStatusFaulted)
	s.cp.indicateStatusChange(1, core.ChargePointStatusUnavailable)
	s.cp.indicateStatusChange(1, core.ChargePointStatusReserved)
	// Invalid status
	s.cp.indicateStatusChange(1, "")

	time.Sleep(time.Second)
}

func (s *hardwareTestSuite) TestIndicateCard() {
	// Ok indication
	s.cp.indicateCardRead(1, indicator.White)

	time.Sleep(time.Second)
}

func TestHardware(t *testing.T) {
	suite.Run(t, new(hardwareTestSuite))
}

func Test_ColorMapping(t *testing.T) {

}
