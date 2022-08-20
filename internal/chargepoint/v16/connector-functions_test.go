package v16

import (
	"context"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/internal/models"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/test"
	ocppManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
	"testing"
	"time"
)

var (
	ocppConfig = configuration.Config{
		Version: 1,
		Keys: []core.ConfigurationKey{
			{
				Key:      "AllowOfflineTxForUnknownId",
				Readonly: false,
				Value:    "false",
			},
			{
				Key:      "AuthorizationCacheEnabled",
				Readonly: false,
				Value:    "false",
			},
			{
				Key:      "AuthorizeRemoteTxRequests",
				Readonly: false,
				Value:    "false",
			},
			{
				Key:      "ClockAlignedDataInterval",
				Readonly: false,
				Value:    "0",
			},
			{
				Key:      "ConnectionTimeOut",
				Readonly: false,
				Value:    "50",
			},
			{
				Key:      "GetConfigurationMaxKeys",
				Readonly: false,
				Value:    "30",
			},
			{
				Key:      "HeartbeatInterval",
				Readonly: false,
				Value:    "60",
			},
			{
				Key:      "LocalAuthorizeOffline",
				Readonly: false,
				Value:    "true",
			},
			{
				Key:      "LocalPreAuthorize",
				Readonly: false,
				Value:    "true",
			},
			{
				Key:      "MaxEnergyOnInvalidId",
				Readonly: false,
				Value:    "0",
			},
			{
				Key:      "MeterValuesSampledData",
				Readonly: false,
				Value:    "Power.Active.Import",
			},
			{
				Key:      "MeterValuesAlignedData",
				Readonly: false,
				Value:    "false",
			},
			{
				Key:      "NumberOfConnectors",
				Readonly: false,
				Value:    "6",
			},
			{
				Key:      "MeterValueSampleInterval",
				Readonly: false,
				Value:    "60",
			},
			{
				Key:      "ResetRetries",
				Readonly: false,
				Value:    "3",
			},
			{
				Key:      "ConnectorPhaseRotation",
				Readonly: false,
				Value:    "0.RST, 1.RST, 2.RTS",
			},
			{
				Key:      "StopTransactionOnEVSideDisconnect",
				Readonly: false,
				Value:    "true",
			},
			{
				Key:      "StopTransactionOnInvalidId",
				Readonly: false,
				Value:    "true",
			},
			{
				Key:      "StopTxnAlignedData",
				Readonly: false,
			},
			{
				Key:      "StopTxnSampledData",
				Readonly: false,
			},
			{
				Key:      "SupportedFeatureProfiles",
				Readonly: true,
				Value:    "Core, LocalAuthListManagement, Reservation, RemoteTrigger",
			},
			{
				Key:      "TransactionMessageAttempts",
				Readonly: false,
				Value:    "3",
			},
			{
				Key:      "TransactionMessageRetryInterval",
				Readonly: false,
				Value:    "60",
			},
			{
				Key:      "UnlockConnectorOnEVSideDisconnect",
				Readonly: false,
				Value:    "true",
			},
			{
				Key:      "ReserveConnectorZeroSupported",
				Readonly: false,
				Value:    "false",
			},
			{
				Key:      "SendLocalListMaxLength",
				Readonly: false,
				Value:    "20",
			},
			{
				Key:      "LocalAuthListEnabled",
				Readonly: false,
				Value:    "true",
			},
			{
				Key:      "LocalAuthListMaxLength",
				Readonly: false,
				Value:    "20",
			},
		},
	}
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
		channel     = make(chan models.Message)
		lcdMock     = new(test.DisplayMock)
	)

	lcdMock.On("GetLcdChannel").Return(channel)
	s.cp.display = lcdMock
	s.cp.settings = &settings.Settings{ChargePoint: settings.ChargePoint{
		Hardware: settings.Hardware{
			Display: settings.Display{
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
		connectorMock = new(test.EvseMock)
	)

	connectorMock.On("GetStatus").Return("Available", "NoError")
	connectorMock.On("GetEvseId").Return(1)

	chargePoint.On("SendRequestAsync", mock.Anything).Run(func(args mock.Arguments) {
		s.Assert().IsType(&core.StatusNotificationRequest{}, args.Get(0))
		notification := args.Get(0).(*core.StatusNotificationRequest)
		s.Assert().EqualValues(connectorId, notification.ConnectorId)
		s.Assert().EqualValues(core.ChargePointStatusAvailable, notification.Status)
	}).Return(core.NewStatusNotificationConfirmation(), nil, nil)
	s.cp.chargePoint = chargePoint

	s.cp.notifyConnectorStatus(1, core.ChargePointStatusAvailable, core.NoError)

	chargePoint.AssertNumberOfCalls(s.T(), "SendRequestAsync", 1)
}

func TestConnectorFunctions(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	err := ocppManager.GetManager().SetConfiguration(ocppConfig)
	assert.NoError(t, err)

	suite.Run(t, new(connectorFunctionsTestSuite))
}
