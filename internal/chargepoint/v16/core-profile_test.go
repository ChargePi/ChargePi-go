package v16

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	s "github.com/xBlaz3kx/ChargePi-go/internal/components/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/test"
	v16 "github.com/xBlaz3kx/ocppManager-go/v16"
	"testing"
)

type coreTestSuite struct {
	suite.Suite
	cp *ChargePoint
}

func (s *coreTestSuite) SetupTest() {
	s.cp = &ChargePoint{
		chargePoint: nil,
		scheduler:   scheduler.GetScheduler(),
		logger:      log.StandardLogger(),
	}
}

func (s *coreTestSuite) TestChangeAvailability() {
	availability, err := s.cp.OnChangeAvailability(core.NewChangeAvailabilityRequest(0, core.AvailabilityTypeOperative))
	s.Assert().NoError(err)
	s.Assert().EqualValues(core.AvailabilityStatusAccepted, availability.Status)

	availability, err = s.cp.OnChangeAvailability(core.NewChangeAvailabilityRequest(1, core.AvailabilityTypeOperative))
	s.Assert().NoError(err)
	s.Assert().EqualValues(core.AvailabilityStatusRejected, availability.Status)
}

func (s *coreTestSuite) TestOnChangeConfiguration() {
	// Ok case
	resp, err := s.cp.OnChangeConfiguration(core.NewChangeConfigurationRequest(v16.AuthorizationCacheEnabled.String(), "false"))
	s.Assert().NoError(err)
	s.Assert().EqualValues(core.ConfigurationStatusAccepted, resp.Status)

	// Invalid key
	resp, err = s.cp.OnChangeConfiguration(core.NewChangeConfigurationRequest("invalidKey", ""))
	s.Assert().NoError(err)
	s.Assert().EqualValues(core.ConfigurationStatusRejected, resp.Status)

	// Readonly key
	resp, err = s.cp.OnChangeConfiguration(core.NewChangeConfigurationRequest(v16.SupportedFeatureProfiles.String(), ""))
	s.Assert().NoError(err)
	s.Assert().EqualValues(core.ConfigurationStatusRejected, resp.Status)
}

func (s *coreTestSuite) TestOnClearCache() {}

func (s *coreTestSuite) TestOnDataTransfer() {
	resp, err := s.cp.OnDataTransfer(core.NewDataTransferRequest(""))
	s.Assert().NoError(err)
	s.Assert().EqualValues(core.DataTransferStatusRejected, resp.Status)
}

func (s *coreTestSuite) TestGetConfiguration() {
	// Get all configuration vars
	resp, err := s.cp.OnGetConfiguration(core.NewGetConfigurationRequest([]string{}))
	s.Assert().NoError(err)
	s.Assert().Len(resp.UnknownKey, 0)

	// Get only specific keys
	resp, err = s.cp.OnGetConfiguration(core.NewGetConfigurationRequest([]string{v16.SupportedFeatureProfiles.String()}))
	s.Assert().NoError(err)
	s.Assert().Len(resp.ConfigurationKey, 1)
	s.Assert().Len(resp.UnknownKey, 0)

	// Keys don't exist
	resp, err = s.cp.OnGetConfiguration(core.NewGetConfigurationRequest([]string{"nonExistingKey"}))
	s.Assert().NoError(err)
	s.Assert().Len(resp.UnknownKey, 1)
	s.Assert().Len(resp.ConfigurationKey, 0)
}

func (s *coreTestSuite) TestOnReset() {}

func (s *coreTestSuite) TestOnUnlockConnector() {}

func (s *coreTestSuite) TestOnRemoteStopTransaction() {
	var (
		connectorManager = new(test.ManagerMock)
		connector        = new(test.ConnectorMock)
		transactionId    = 1
		transactionIdStr = "1"
	)

	connector.On("IsCharging").Return(true).Once()
	connectorManager.On("FindConnectorWithTransactionId", transactionIdStr).Return(connector).Once()

	s.cp.connectorManager = connectorManager

	req := core.NewRemoteStopTransactionRequest(transactionId)
	response, err := s.cp.OnRemoteStopTransaction(req)
	s.Assert().NoError(err)
	s.Assert().EqualValues(types.RemoteStartStopStatusAccepted, response.Status)
	s.Assert().EqualValues(1, s.cp.scheduler.Len())

	s.cp.scheduler.Clear()

	// No transaction
	connectorManager.On("FindConnectorWithTransactionId", transactionIdStr).Return(nil).Once()
	req = core.NewRemoteStopTransactionRequest(transactionId)
	response, err = s.cp.OnRemoteStopTransaction(req)
	s.Assert().NoError(err)
	s.Assert().EqualValues(types.RemoteStartStopStatusRejected, response.Status)
	s.Assert().EqualValues(0, s.cp.scheduler.Len())

	s.cp.scheduler.Clear()

	// Connector not charging
	connector.On("IsCharging").Return(false).Once()
	connectorManager.On("FindConnectorWithTransactionId", transactionIdStr).Return(connector).Once()
	req = core.NewRemoteStopTransactionRequest(transactionId)
	response, err = s.cp.OnRemoteStopTransaction(req)
	s.Assert().NoError(err)
	s.Assert().EqualValues(types.RemoteStartStopStatusRejected, response.Status)
	s.Assert().EqualValues(0, s.cp.scheduler.Len())
}

func (s *coreTestSuite) TestOnRemoteStartTransaction() {
	var (
		connectorManager       = new(test.ManagerMock)
		connector              = new(test.ConnectorMock)
		connectorId            = 1
		nonExistingConnectorId = 14
	)

	connector.On("IsAvailable").Return(true).Twice()
	connectorManager.On("FindAvailableConnector").Return(connector).Once()

	s.cp.connectorManager = connectorManager

	req := core.NewRemoteStartTransactionRequest(tagId)
	transaction, err := s.cp.OnRemoteStartTransaction(req)
	s.Assert().NoError(err)
	s.Assert().EqualValues(types.RemoteStartStopStatusAccepted, transaction.Status)
	s.Assert().EqualValues(1, s.cp.scheduler.Len())

	s.cp.scheduler.Clear()

	// Start charging a specific connector
	connectorManager.On("FindConnector", 1, connectorId).Return(connector).Once()
	req = core.NewRemoteStartTransactionRequest(tagId)
	req.ConnectorId = &connectorId
	transaction, err = s.cp.OnRemoteStartTransaction(req)
	s.Assert().NoError(err)
	s.Assert().EqualValues(types.RemoteStartStopStatusAccepted, transaction.Status)
	s.Assert().EqualValues(1, s.cp.scheduler.Len())

	s.cp.scheduler.Clear()

	// No such connector exists
	connectorManager.On("FindConnector", 1, nonExistingConnectorId).Return(nil).Once()
	req = core.NewRemoteStartTransactionRequest(tagId)
	req.ConnectorId = &nonExistingConnectorId
	transaction, err = s.cp.OnRemoteStartTransaction(req)
	s.Assert().NoError(err)
	s.Assert().EqualValues(types.RemoteStartStopStatusRejected, transaction.Status)
	s.Assert().EqualValues(0, s.cp.scheduler.Len())

	s.cp.scheduler.Clear()

	// Connector not available
	connector.On("IsAvailable").Return(false).Once()
	connectorManager.On("FindConnector", 1, connectorId).Return(nil).Once()
	req = core.NewRemoteStartTransactionRequest(tagId)
	req.ConnectorId = &connectorId
	transaction, err = s.cp.OnRemoteStartTransaction(req)
	s.Assert().NoError(err)
	s.Assert().EqualValues(types.RemoteStartStopStatusRejected, transaction.Status)
	s.Assert().EqualValues(0, s.cp.scheduler.Len())
}

func TestCore(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// Setup OCPP configuration manager
	s.SetupOcppConfigurationManager(
		"../../../configs/configuration.json",
		"1.6",
		core.ProfileName,
		reservation.ProfileName)

	suite.Run(t, new(coreTestSuite))
}
