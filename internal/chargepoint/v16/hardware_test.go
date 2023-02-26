package v16

import (
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
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
	s.cp.logger = log.StandardLogger()
}

func (s *hardwareTestSuite) TestSendToLCD() {
	lcdMock := display.NewDisplayMock(s.T())
	s.cp.SetDisplay(lcdMock)
	lcdMock.OnCleanup().Return()

	s.cp.sendToLCD(exampleMessage, exampleMessage1)

	// Disable LCD
	s.cp.sendToLCD(exampleMessage, exampleMessage1)
}

func (s *hardwareTestSuite) TestDisplayLedStatus() {
	// Ok statuses
	s.cp.displayStatusChangeOnIndicator(1, core.ChargePointStatusCharging)
	s.cp.displayStatusChangeOnIndicator(1, core.ChargePointStatusFinishing)
	s.cp.displayStatusChangeOnIndicator(1, core.ChargePointStatusAvailable)
	s.cp.displayStatusChangeOnIndicator(1, core.ChargePointStatusFaulted)
	s.cp.displayStatusChangeOnIndicator(1, core.ChargePointStatusUnavailable)
	s.cp.displayStatusChangeOnIndicator(1, core.ChargePointStatusReserved)
	// Invalid status
	s.cp.displayStatusChangeOnIndicator(1, "")

	time.Sleep(time.Second)
}

func (s *hardwareTestSuite) TestIndicateCard() {
	// Ok indication
	s.cp.indicateCard(1, indicator.White)

	time.Sleep(time.Second)
}

func TestHardware(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	suite.Run(t, new(hardwareTestSuite))
}
