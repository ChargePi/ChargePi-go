package v16

import (
	"errors"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
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
	s.cp.settings = &settings.Settings{ChargePoint: settings.ChargePoint{
		Hardware: settings.Hardware{
			Display: settings.Display{
				IsEnabled: true,
			},
		},
	}}

	s.cp.sendToLCD(exampleMessage, exampleMessage1)

	// Disable LCD
	s.cp.settings.ChargePoint.Hardware.Display.IsEnabled = false
	s.cp.sendToLCD(exampleMessage, exampleMessage1)

	s.lcdMock.AssertNumberOfCalls(s.T(), "DisplayMessage", 1)
}

func (s *hardwareTestSuite) TestDisplayLedStatus() {
	s.indicatorMock.On("DisplayColor", 1, uint32(indicator.Blue)).Return(nil)
	s.indicatorMock.On("DisplayColor", 1, uint32(indicator.Red)).Return(nil)
	s.indicatorMock.On("DisplayColor", 1, uint32(indicator.Yellow)).Return(nil)
	s.indicatorMock.On("DisplayColor", 1, uint32(indicator.Green)).Return(nil)
	s.indicatorMock.On("DisplayColor", 1, uint32(indicator.Orange)).Return(nil)
	s.indicatorMock.On("DisplayColor", 1, uint32(indicator.Off)).Return(errors.New("invalid color")).Once()

	s.cp.indicator = s.indicatorMock
	s.cp.settings = &settings.Settings{ChargePoint: settings.ChargePoint{
		Hardware: settings.Hardware{
			LedIndicator: settings.LedIndicator{
				Enabled: true,
			},
		},
	}}

	// Ok statuses
	s.cp.displayLEDStatus(1, core.ChargePointStatusCharging)
	s.cp.displayLEDStatus(1, core.ChargePointStatusFinishing)
	s.cp.displayLEDStatus(1, core.ChargePointStatusAvailable)
	s.cp.displayLEDStatus(1, core.ChargePointStatusFaulted)
	s.cp.displayLEDStatus(1, core.ChargePointStatusUnavailable)
	s.cp.displayLEDStatus(1, core.ChargePointStatusReserved)
	// Invalid status
	s.cp.displayLEDStatus(1, "")

	time.Sleep(time.Second)

	s.indicatorMock.AssertNumberOfCalls(s.T(), "DisplayColor", 6)
}

func (s *hardwareTestSuite) TestIndicateCard() {
	s.indicatorMock.On("Blink", 1, 3, uint32(indicator.White)).Return(nil)
	s.indicatorMock.On("Blink", 1, 3, uint32(123)).Return(errors.New("invalid color"))

	s.cp.indicator = s.indicatorMock
	s.cp.settings = &settings.Settings{ChargePoint: settings.ChargePoint{
		Hardware: settings.Hardware{
			LedIndicator: settings.LedIndicator{
				Enabled: true,
			},
		},
	}}

	// Ok indication
	s.cp.indicateCard(1, indicator.White)

	// Invalid color
	s.cp.indicateCard(1, uint32(123))

	time.Sleep(time.Second)

	s.indicatorMock.AssertNumberOfCalls(s.T(), "Blink", 2)
}

func TestHardware(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	suite.Run(t, new(hardwareTestSuite))
}
