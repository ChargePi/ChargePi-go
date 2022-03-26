package v16

import (
	"context"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/hardware/display"
	s "github.com/xBlaz3kx/ChargePi-go/internal/components/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/test"
	"testing"
	"time"
)

type connectorFunctionsTestSuite struct {
	suite.Suite
	cp *ChargePoint
}

func (s *connectorFunctionsTestSuite) SetupTest() {
	s.cp = &ChargePoint{
		logger: log.StandardLogger(),
	}
}

func (s *connectorFunctionsTestSuite) TestAddConnectors() {
	//todo
}

func (s *connectorFunctionsTestSuite) TestRestoreState() {
	//todo
}

func (s *connectorFunctionsTestSuite) TestDisplayConnectorStatus() {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
		channel     = make(chan display.LCDMessage)
		lcdMock     = new(test.DisplayMock)
	)

	lcdMock.On("GetLcdChannel").Return(channel)
	s.cp.LCD = lcdMock
	s.cp.Settings = &settings.Settings{ChargePoint: settings.ChargePoint{
		Hardware: settings.Hardware{
			Lcd: settings.Lcd{
				IsEnabled: true,
				Language:  "en",
			},
		},
	}}

	go func() {
		time.Sleep(time.Millisecond * 100)
		s.cp.displayConnectorStatus(connectorId, core.ChargePointStatusAvailable)

		time.Sleep(time.Millisecond * 100)
		s.cp.displayConnectorStatus(connectorId, core.ChargePointStatusCharging)

		time.Sleep(time.Millisecond * 100)
		s.cp.displayConnectorStatus(connectorId, core.ChargePointStatusFinishing)
	}()

	numMessages := 0
Loop:
	for {
		select {
		case msg := <-channel:
			numMessages++
			log.Debugf("Received message from channel %v", msg)
			s.Condition(func() (success bool) {
				switch numMessages {
				case 1:
					return s.Contains(msg.Messages, "available.")
				case 2:
					return s.Contains(msg.Messages, "Started charging") &&
						s.Contains(msg.Messages, "at 1.")
				case 3:
					return s.Contains(msg.Messages, "Stopped charging")
				default:
					s.Fail("Invalid message number")
					return false
				}
			})

			if numMessages == 3 {
				cancel()
			}
			break
		case <-ctx.Done():
			break Loop
		}
	}

	cancel()
}

func (s *connectorFunctionsTestSuite) TestNotifyConnectorStatus() {
	var (
		chargePoint   = new(chargePointMock)
		connectorMock = new(test.ConnectorMock)
	)

	connectorMock.On("GetStatus").Return("Available", "NoError")
	connectorMock.On("GetConnectorId").Return(1)

	chargePoint.On("SendRequestAsync", mock.Anything).Run(func(args mock.Arguments) {
		s.Assert().IsType(&core.StatusNotificationRequest{}, args.Get(0))
		notification := args.Get(0).(*core.StatusNotificationRequest)
		s.Assert().EqualValues(connectorId, notification.ConnectorId)
		s.Assert().EqualValues(core.ChargePointStatusAvailable, notification.Status)
	}).Return(core.NewStatusNotificationConfirmation(), nil, nil)
	s.cp.chargePoint = chargePoint

	s.cp.notifyConnectorStatus(connectorMock)
	s.cp.notifyConnectorStatus(nil)

	chargePoint.AssertNumberOfCalls(s.T(), "SendRequestAsync", 1)
}

func TestConnectorFunctions(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// Setup OCPP configuration manager
	s.SetupOcppConfigurationManager(
		"../../../configs/configuration.json",
		"1.6",
		core.ProfileName,
		reservation.ProfileName)

	suite.Run(t, new(connectorFunctionsTestSuite))
}
