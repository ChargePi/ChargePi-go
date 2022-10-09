package evse

import (
	"errors"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/session"
	settingsModel "github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/test"
	"testing"
	"time"
)

type (
	connectorManagerTestSuite struct {
		suite.Suite
		connectorManager  Manager
		connector1        *test.EvseMock
		connector2        *test.EvseMock
		connector3        *test.EvseMock
		chSession         session.Session
		connectorSettings *settingsModel.EVSE
	}
)

func CreateNewConnectorMock(evseId, connectorId int, session session.Session) *test.EvseMock {
	connector1 := new(test.EvseMock)
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
	connector1.On("SetMeterValuesChannel", mock.Anything).Return()
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

	suite.connectorSettings = &settingsModel.EVSE{
		EvseId: 1,
		Status: "Available",
		Session: settingsModel.Session{
			IsActive: false,
		},
		EVCC: settingsModel.EVCC{
			RelayPin:     14,
			InverseLogic: false,
		},
		PowerMeter: settingsModel.PowerMeter{
			Enabled: false,
		},
	}

	suite.connector1 = CreateNewConnectorMock(1, 1, suite.chSession)
	suite.connector2 = CreateNewConnectorMock(1, 2, suite.chSession)
	suite.connector3 = CreateNewConnectorMock(1, 3, suite.chSession)

	err := suite.connectorManager.AddEVSE(suite.connector1)
	suite.Require().NoError(err)

}

func (suite *connectorManagerTestSuite) TestAddConnector() {
	// Duplicate connector
	err := suite.connectorManager.AddEVSE(suite.connector1)
	suite.Require().Error(err)

	err = suite.connectorManager.AddEVSE(suite.connector2)
	suite.Require().NoError(err)

	// Nil connector not allowed
	err = suite.connectorManager.AddEVSE(nil)
	suite.Require().Error(err)

	suite.Require().Contains(suite.connectorManager.GetEVSEs(), suite.connector1)
	suite.Require().Contains(suite.connectorManager.GetEVSEs(), suite.connector2)
}

func (suite *connectorManagerTestSuite) TestStartChargingConnector() {
	tagId := "exampleTag"
	transactionId := "exampleTransactionId123"

	err := suite.connectorManager.StartCharging(1, tagId, transactionId)
	suite.Require().NoError(err)

	newConn := new(test.EvseMock)
	newConn.On("GetEvseId").Return(1)
	newConn.On("GetConnectorId").Return(4)
	newConn.On("SetNotificationChannel", mock.Anything).Return()
	newConn.On("SetMeterValuesChannel", mock.Anything).Return()
	err = suite.connectorManager.AddEVSE(newConn)
	suite.Require().NoError(err)

	// Start charging returns an error
	newConn.On("StartCharging", transactionId, tagId).Return(errors.New("something went wrong"))
	err = suite.connectorManager.StartCharging(4, tagId, transactionId)
	suite.Require().Error(err)

	// No such connector
	err = suite.connectorManager.StartCharging(5, tagId, transactionId)
	suite.Require().Error(err)

	time.Sleep(time.Second)
	suite.connector1.AssertCalled(suite.T(), "StartCharging", transactionId, tagId)
	newConn.AssertCalled(suite.T(), "StartCharging", transactionId, tagId)
}

func (suite *connectorManagerTestSuite) TestAddConnectorFromSettings() {
	if testing.Short() {
		return
	}

	// Try to add duplicate connector
	err := suite.connectorManager.AddEVSEFromSettings(15, suite.connectorSettings)
	suite.Require().Error(err)

	// Try to add another connector
	err = suite.connectorManager.AddEVSEFromSettings(15, &settingsModel.EVSE{
		EvseId:  2,
		Status:  "Available",
		Session: settingsModel.Session{},
		EVCC: settingsModel.EVCC{
			RelayPin:     23,
			InverseLogic: false,
		},
		PowerMeter: settingsModel.PowerMeter{},
	})
	suite.Require().NoError(err)

	// Try to add another connector
	err = suite.connectorManager.AddEVSEFromSettings(15, nil)
	suite.Require().Error(err)
}

func (suite *connectorManagerTestSuite) TestGetConnectors() {
	connectors := suite.connectorManager.GetEVSEs()
	suite.Require().Len(connectors, 1)

	// Try to add another connector
	err := suite.connectorManager.AddEVSE(suite.connector3)
	suite.Require().NoError(err)

	connectors = suite.connectorManager.GetEVSEs()
	suite.Require().Len(connectors, 2)

	suite.Require().Contains(suite.connectorManager.GetEVSEs(), suite.connector1)
	suite.Require().Contains(suite.connectorManager.GetEVSEs(), suite.connector3)
}

func (suite *connectorManagerTestSuite) TestFindConnectorWithId() {
	c := suite.connectorManager.FindEVSE(1)
	suite.Require().NotNil(c)
	suite.Require().Equal(1, c.GetEvseId())

	c = suite.connectorManager.FindEVSE(2)
	suite.Require().Nil(c)
}

func (suite *connectorManagerTestSuite) TestFindConnectorWithTagId() {
	tagId := "exampleTag"

	suite.connector1.On("GetTagId").Return(tagId)

	connectorWithTag := suite.connectorManager.FindEVSEWithTagId(tagId)
	suite.Require().NotNil(connectorWithTag)
	suite.Require().Equal(1, connectorWithTag.GetEvseId())

	connectorWithTag = suite.connectorManager.FindEVSEWithTagId("noConnectorWithTag")
	suite.Require().Nil(connectorWithTag)
}

func (suite *connectorManagerTestSuite) TestFindConnectorWithTransactionId() {
	transactionId := "exampleTransactionId123"

	suite.connector1.On("GetTransactionId").Return(transactionId)

	connectorWithTransaction := suite.connectorManager.FindEVSEWithTransactionId(transactionId)
	suite.Require().NotNil(connectorWithTransaction)
	suite.Require().Equal(1, connectorWithTransaction.GetEvseId())

	connectorWithTag := suite.connectorManager.FindEVSEWithTransactionId("noTransactionWithThisId")
	suite.Require().Nil(connectorWithTag)
}

func (suite *connectorManagerTestSuite) TestStopAllConnectors() {
	var (
		mConnector = suite.connector1
		newConn    = new(test.EvseMock)
	)

	newConn.On("GetConnectorId").Return(1)
	newConn.On("GetEvseId").Return(4)
	newConn.On("SetNotificationChannel", mock.Anything).Return()
	newConn.On("SetMeterValuesChannel", mock.Anything).Return()
	newConn.On("StopCharging", core.ReasonLocal).Return(errors.New("something happened"))

	// Add another connector
	err := suite.connectorManager.AddEVSE(suite.connector3)
	suite.Require().NoError(err)

	// Stop all
	err = suite.connectorManager.StopAllEVSEs(core.ReasonLocal)
	suite.Require().NoError(err)

	// Add another "faulty" connector
	err = suite.connectorManager.AddEVSE(newConn)
	suite.Require().NoError(err)

	// Stop causes an error
	err = suite.connectorManager.StopAllEVSEs(core.ReasonLocal)
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
