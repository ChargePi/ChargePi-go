package v16

import (
	"errors"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/test"
	"testing"
	"time"
)

const (
	exampleMessage  = "exampleMessage"
	exampleMessage1 = "exampleMessage2"
)

type hardwareTestSuite struct {
	suite.Suite
	cp            *ChargePoint
	lcdMock       *test.DisplayMock
	indicatorMock *test.IndicatorMock
}

func (s *hardwareTestSuite) SetupTest() {
	s.lcdMock = new(test.DisplayMock)
	s.indicatorMock = new(test.IndicatorMock)
	s.cp = new(ChargePoint)
	s.cp.logger = log.StandardLogger()
}

func (s *hardwareTestSuite) TestSendToLCD() {
	s.cp.display = s.lcdMock

	s.cp.sendToLCD(exampleMessage, exampleMessage1)

	// Disable LCD
	s.cp.sendToLCD(exampleMessage, exampleMessage1)

	s.lcdMock.AssertNumberOfCalls(s.T(), "DisplayMessage", 1)
}

func (s *hardwareTestSuite) TestDisplayLedStatus() {
	s.indicatorMock.On("DisplayColor", 1, indicator.Blue).Return(nil)
	s.indicatorMock.On("DisplayColor", 1, indicator.Red).Return(nil)
	s.indicatorMock.On("DisplayColor", 1, indicator.Yellow).Return(nil)
	s.indicatorMock.On("DisplayColor", 1, indicator.Green).Return(nil)
	s.indicatorMock.On("DisplayColor", 1, indicator.Orange).Return(nil)
	s.indicatorMock.On("DisplayColor", 1, indicator.Off).Return(errors.New("invalid color")).Once()

	s.cp.indicator = s.indicatorMock

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

	s.indicatorMock.AssertNumberOfCalls(s.T(), "DisplayColor", 6)
}

func (s *hardwareTestSuite) TestIndicateCard() {
	s.indicatorMock.On("Blink", 1, 3, indicator.White).Return(nil)
	s.indicatorMock.On("Blink", 1, 3, uint32(123)).Return(errors.New("invalid color"))

	s.cp.indicator = s.indicatorMock

	// Ok indication
	s.cp.indicateCard(1, indicator.White)

	time.Sleep(time.Second)
	s.indicatorMock.AssertNumberOfCalls(s.T(), "Blink", 2)
}

func TestHardware(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	suite.Run(t, new(hardwareTestSuite))
}
