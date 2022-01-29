package connector_manager

import (
	"errors"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/reactivex/rxgo/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/connector"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/hardware/power-meter"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/session"
	settings2 "github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"testing"
	"time"
)

type (
	connectorManagerTestSuite struct {
		suite.Suite
		connectorManager  Manager
		connector1        *connectorMock
		connector2        *connectorMock
		connector3        *connectorMock
		chSession         session.Session
		connectorSettings *settings2.Connector
	}

	connectorMock struct {
		mock.Mock
		connector.Connector
	}
)

func (m *connectorMock) StartCharging(transactionId string, tagId string) error {
	args := m.Called(transactionId, tagId)
	return args.Error(0)
}

func (m *connectorMock) ResumeCharging(session session.Session) (error, int) {
	args := m.Called(session)
	return args.Error(0), args.Int(1)
}

func (m *connectorMock) StopCharging(reason core.Reason) error {
	args := m.Called(reason)
	return args.Error(0)
}

func (m *connectorMock) SetNotificationChannel(notificationChannel chan<- rxgo.Item) {
	m.Called(notificationChannel)
}

func (m *connectorMock) ReserveConnector(reservationId int) error {
	args := m.Called(reservationId)
	return args.Error(0)
}

func (m *connectorMock) RemoveReservation() error {
	args := m.Called()
	return args.Error(0)
}

func (m *connectorMock) GetReservationId() int {
	args := m.Called()
	return args.Int(0)
}

func (m *connectorMock) GetTagId() string {
	args := m.Called()
	return args.String(0)
}

func (m *connectorMock) GetTransactionId() string {
	args := m.Called()
	return args.String(0)
}

func (m *connectorMock) GetConnectorId() int {
	args := m.Called()
	return args.Int(0)
}

func (m *connectorMock) GetEvseId() int {
	args := m.Called()
	return args.Int(0)
}

func (m *connectorMock) CalculateSessionAvgEnergyConsumption() float32 {
	args := m.Called()
	return args.Get(0).(float32)
}

func (m *connectorMock) SamplePowerMeter(measurands []types.Measurand) {
	m.Called(measurands)
}

func (m *connectorMock) SetStatus(status core.ChargePointStatus, errCode core.ChargePointErrorCode) {
	m.Called(status, errCode)
}

func (m *connectorMock) GetStatus() (core.ChargePointStatus, core.ChargePointErrorCode) {
	args := m.Called()
	return core.ChargePointStatus(args.String(0)), core.ChargePointErrorCode(args.String(1))
}

func (m *connectorMock) IsAvailable() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *connectorMock) IsPreparing() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *connectorMock) IsCharging() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *connectorMock) IsReserved() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *connectorMock) IsUnavailable() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *connectorMock) GetPowerMeter() power_meter.PowerMeter {
	args := m.Called()
	return args.Get(0).(power_meter.PowerMeter)
}

func (m *connectorMock) GetMaxChargingTime() int {
	args := m.Called()
	return args.Int(0)
}

func CreateNewConnectorMock(evseId, connectorId int, session session.Session) *connectorMock {
	connector1 := new(connectorMock)
	// Setup expectations
	connector1.On("StartCharging", "exampleTransactionId123", "exampleTag").Return(nil)
	connector1.On("ResumeCharging", session).Return(nil, 0)
	connector1.On("StopCharging", core.ReasonLocal).Return(nil)
	connector1.On("RemoveReservation", 123).Return(nil)
	connector1.On("GetReservationId", 123).Return(0)
	connector1.On("GetConnectorId").Return(connectorId)
	connector1.On("GetEvseId").Return(evseId)
	connector1.On("CalculateSessionAvgEnergyConsumption").Return(30.0)
	connector1.On("GetStatus").Return(core.ChargePointStatusAvailable, core.NoError)
	connector1.On("IsAvailable").Return(true)
	connector1.On("IsPreparing").Return(false)
	connector1.On("IsCharging").Return(false)
	connector1.On("IsReserved").Return(false)
	connector1.On("IsUnavailable").Return(false)
	connector1.On("GetMaxChargingTime").Return(15)
	connector1.On("SetNotificationChannel", mock.Anything).Return()
	return connector1
}

/*------------------- Manager Test Suite ------------------------------*/

func (suite *connectorManagerTestSuite) SetupTest() {
	suite.connectorManager = NewManager(nil)

	suite.chSession = session.Session{
		IsActive:      true,
		TransactionId: "transaction123",
		TagId:         "tagId",
		Started:       time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
		Consumption:   nil,
	}

	suite.connectorSettings = &settings2.Connector{
		EvseId:      1,
		ConnectorId: 1,
		Type:        "Schuko",
		Status:      "Available",
		Session: settings2.Session{
			IsActive:      false,
			TransactionId: "",
			TagId:         "",
			Started:       "",
			Consumption:   nil,
		},
		Relay: settings2.Relay{
			RelayPin:     14,
			InverseLogic: false,
		},
		PowerMeter: settings2.PowerMeter{
			Enabled:              false,
			Type:                 "",
			PowerMeterPin:        0,
			SpiBus:               0,
			Consumption:          0,
			ShuntOffset:          0,
			VoltageDividerOffset: 0,
		},
	}

	suite.connector1 = CreateNewConnectorMock(1, 1, suite.chSession)
	suite.connector2 = CreateNewConnectorMock(1, 2, suite.chSession)
	suite.connector3 = CreateNewConnectorMock(1, 3, suite.chSession)

	err := suite.connectorManager.AddConnector(suite.connector1)
	suite.Require().NoError(err)

}

func (suite *connectorManagerTestSuite) TestAddConnector() {
	// Duplicate connector
	err := suite.connectorManager.AddConnector(suite.connector1)
	suite.Require().Error(err)

	err = suite.connectorManager.AddConnector(suite.connector2)
	suite.Require().NoError(err)

	// Nil connector not allowed
	err = suite.connectorManager.AddConnector(nil)
	suite.Require().Error(err)

	suite.Require().Contains(suite.connectorManager.GetConnectors(), suite.connector1)
	suite.Require().Contains(suite.connectorManager.GetConnectors(), suite.connector2)
}

func (suite *connectorManagerTestSuite) TestStartChargingConnector() {
	tagId := "exampleTag"
	transactionId := "exampleTransactionId123"

	err := suite.connectorManager.StartChargingConnector(1, 1, tagId, transactionId)
	suite.Require().NoError(err)

	newConn := new(connectorMock)
	newConn.On("GetEvseId").Return(1)
	newConn.On("GetConnectorId").Return(4)
	newConn.On("SetNotificationChannel", mock.Anything).Return()
	err = suite.connectorManager.AddConnector(newConn)
	suite.Require().NoError(err)

	// Start charging returns an error
	newConn.On("StartCharging", transactionId, tagId).Return(errors.New("something went wrong"))
	err = suite.connectorManager.StartChargingConnector(1, 4, tagId, transactionId)
	suite.Require().Error(err)

	// No such connector
	err = suite.connectorManager.StartChargingConnector(1, 5, tagId, transactionId)
	suite.Require().Error(err)

	time.Sleep(time.Second)
	suite.connector1.AssertCalled(suite.T(), "StartCharging", transactionId, tagId)
	newConn.AssertCalled(suite.T(), "StartCharging", transactionId, tagId)
}

func (suite *connectorManagerTestSuite) TestAddConnectorFromSettings() {
	// Try to add duplicate connector
	err := suite.connectorManager.AddConnectorFromSettings(15, suite.connectorSettings)
	suite.Require().Error(err)

	// Try to add another connector
	err = suite.connectorManager.AddConnectorFromSettings(15, &settings2.Connector{
		EvseId:      1,
		ConnectorId: 2,
		Type:        "Schuko",
		Status:      "Available",
		Session:     settings2.Session{},
		Relay: settings2.Relay{
			RelayPin:     23,
			InverseLogic: false,
		},
		PowerMeter: settings2.PowerMeter{},
	})
	suite.Require().NoError(err)

	// Try to add another connector
	err = suite.connectorManager.AddConnectorFromSettings(15, nil)
	suite.Require().Error(err)
}

func (suite *connectorManagerTestSuite) TestGetConnectors() {
	connectors := suite.connectorManager.GetConnectors()
	suite.Require().Len(connectors, 1)

	// Try to add another connector
	err := suite.connectorManager.AddConnector(suite.connector3)
	suite.Require().NoError(err)

	connectors = suite.connectorManager.GetConnectors()
	suite.Require().Len(connectors, 2)

	suite.Require().Contains(suite.connectorManager.GetConnectors(), suite.connector1)
	suite.Require().Contains(suite.connectorManager.GetConnectors(), suite.connector3)
}

func (suite *connectorManagerTestSuite) TestFindConnectorWithId() {
	c := suite.connectorManager.FindConnector(1, 1)
	suite.Require().NotNil(c)
	suite.Require().Equal(1, c.GetConnectorId())
	suite.Require().Equal(1, c.GetEvseId())

	c = suite.connectorManager.FindConnector(2, 4)
	suite.Require().Nil(c)
}

func (suite *connectorManagerTestSuite) TestFindConnectorWithTagId() {
	tagId := "exampleTag"

	suite.connector1.On("GetTagId").Return(tagId)

	connectorWithTag := suite.connectorManager.FindConnectorWithTagId(tagId)
	suite.Require().NotNil(connectorWithTag)
	suite.Require().Equal(1, connectorWithTag.GetConnectorId())
	suite.Require().Equal(1, connectorWithTag.GetEvseId())

	connectorWithTag = suite.connectorManager.FindConnectorWithTagId("noConnectorWithTag")
	suite.Require().Nil(connectorWithTag)
}

func (suite *connectorManagerTestSuite) TestFindConnectorWithTransactionId() {
	transactionId := "exampleTransactionId123"

	suite.connector1.On("GetTransactionId").Return(transactionId)

	connectorWithTransaction := suite.connectorManager.FindConnectorWithTransactionId(transactionId)
	suite.Require().NotNil(connectorWithTransaction)
	suite.Require().Equal(1, connectorWithTransaction.GetConnectorId())
	suite.Require().Equal(1, connectorWithTransaction.GetEvseId())

	connectorWithTag := suite.connectorManager.FindConnectorWithTransactionId("noTransactionWithThisId")
	suite.Require().Nil(connectorWithTag)
}

func (suite *connectorManagerTestSuite) TestStopAllConnectors() {
	var (
		mConnector = suite.connector1
		newConn    = new(connectorMock)
	)

	newConn.On("GetConnectorId").Return(1)
	newConn.On("GetEvseId").Return(4)
	newConn.On("SetNotificationChannel", mock.Anything).Return()
	newConn.On("StopCharging", core.ReasonLocal).Return(errors.New("something happened"))

	// Add another connector
	err := suite.connectorManager.AddConnector(suite.connector3)
	suite.Require().NoError(err)

	// Stop all
	err = suite.connectorManager.StopAllConnectors(core.ReasonLocal)
	suite.Require().NoError(err)

	// Add another "faulty" connector
	err = suite.connectorManager.AddConnector(newConn)
	suite.Require().NoError(err)

	// Stop causes an error
	err = suite.connectorManager.StopAllConnectors(core.ReasonLocal)
	suite.Require().Error(err)

	time.Sleep(time.Millisecond * 100)
	mConnector.AssertCalled(suite.T(), "StopCharging", core.ReasonLocal)
	newConn.AssertCalled(suite.T(), "StopCharging", core.ReasonLocal)
}

func (suite *connectorManagerTestSuite) TestRestoreState() {

}

func TestConnectorManager(t *testing.T) {
	suite.Run(t, new(connectorManagerTestSuite))
}
