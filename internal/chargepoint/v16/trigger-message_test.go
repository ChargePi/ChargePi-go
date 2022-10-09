package v16

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/test"
	"testing"
	"time"
)

type triggerMessageTestSuite struct {
	suite.Suite
	cp *ChargePoint
}

func (s *triggerMessageTestSuite) SetupTest() {
	s.cp = &ChargePoint{
		logger:    log.StandardLogger(),
		scheduler: scheduler.GetScheduler(),
	}
	s.cp.scheduler.Clear()
}

func (s *triggerMessageTestSuite) TestTriggerMessage() {
	var (
		connectorMock = new(test.EvseMock)
		managerMock   = new(test.ManagerMock)
		connectorChan = make(chan chargePoint.StatusNotification, 2)
	)

	// Set manager expectations
	managerMock.On("FindEVSE", 1, connectorId).Return(connectorMock).Once()
	managerMock.On("GetEVSEs").Return([]evse.EVSE{connectorMock}).Once()

	s.cp.connectorManager = managerMock
	s.cp.connectorChannel = connectorChan

	numMessages := 0
	go func() {
		for {
			select {
			case <-connectorChan:
				numMessages++
			}
		}
	}()

	response, err := s.cp.OnTriggerMessage(remotetrigger.NewTriggerMessageRequest("something"))
	s.Assert().NoError(err)
	s.Assert().NotNil(response)
	s.Assert().EqualValues(remotetrigger.TriggerMessageStatusNotImplemented, response.Status)

	response, err = s.cp.OnTriggerMessage(remotetrigger.NewTriggerMessageRequest(core.MeterValuesFeatureName))
	s.Assert().NoError(err)
	s.Assert().NotNil(response)
	s.Assert().EqualValues(remotetrigger.TriggerMessageStatusNotImplemented, response.Status)

	response, err = s.cp.OnTriggerMessage(remotetrigger.NewTriggerMessageRequest(core.HeartbeatFeatureName))
	s.Assert().NoError(err)
	s.Assert().NotNil(response)
	s.Assert().EqualValues(remotetrigger.TriggerMessageStatusAccepted, response.Status)
	s.Assert().Len(s.cp.scheduler.Jobs(), 1)

	s.cp.scheduler.Clear()

	time.Sleep(time.Second)

	response, err = s.cp.OnTriggerMessage(remotetrigger.NewTriggerMessageRequest(core.BootNotificationFeatureName))
	s.Assert().NoError(err)
	s.Assert().NotNil(response)
	s.Assert().EqualValues(remotetrigger.TriggerMessageStatusAccepted, response.Status)
	s.Assert().Len(s.cp.scheduler.Jobs(), 1)

	s.cp.scheduler.Clear()

	// Get status of all connectors
	response, err = s.cp.OnTriggerMessage(remotetrigger.NewTriggerMessageRequest(core.StatusNotificationFeatureName))
	s.Assert().NoError(err)
	s.Assert().NotNil(response)
	s.Assert().EqualValues(remotetrigger.TriggerMessageStatusAccepted, response.Status)

	time.Sleep(time.Second * 3)
	s.Assert().EqualValues(1, numMessages)

	// Get status of a single connector
	/*request := remotetrigger.NewTriggerMessageRequest(core.StatusNotificationFeatureName)
	response, err = s.cp.OnTriggerMessage(request)
	s.Assert().NoError(err)
	s.Assert().NotNil(response)
	s.Assert().EqualValues(remotetrigger.TriggerMessageStatusAccepted, response.Status)*/
}

func TestTriggerMessage(t *testing.T) {
	suite.Run(t, new(triggerMessageTestSuite))
}
