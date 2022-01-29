package integration_tests

import (
	"context"
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint"
	v16 "github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/v16"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/auth"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/connector"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/hardware/reader"
	setting "github.com/xBlaz3kx/ChargePi-go/internal/components/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/scheduler"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
	"net/http"
	"testing"
	"time"
)

const (
	centralSystemEndpoint = ""
	centralSystemPort     = 7777
	centralSystemHost     = "localhost"

	tagId           = "exampleTag"
	chargePointId   = "exampleChargePoint"
	protocolVersion = settings.OCPP16
)

var chargePointSettings = settings.Settings{ChargePoint: settings.ChargePoint{
	Info: struct {
		Id              string            `fig:"Id" validate:"required"`
		ProtocolVersion string            `fig:"ProtocolVersion" default:"1.6"`
		ServerUri       string            `fig:"ServerUri" validate:"required"`
		MaxChargingTime int               `fig:"MaxChargingTime" default:"180"`
		OCPPInfo        settings.OCPPInfo `fig:"ocpp"`
	}{
		Id:              chargePointId,
		ProtocolVersion: string(protocolVersion),
		ServerUri:       fmt.Sprintf("ws://%s:%d%s", centralSystemHost, centralSystemPort, centralSystemEndpoint),
		MaxChargingTime: 15,
		OCPPInfo: settings.OCPPInfo{
			Vendor: "exampleVendor",
			Model:  "exampleModel",
		},
	},
	Logging: settings.Logging{},
	TLS:     settings.TLS{IsEnabled: false},
	Hardware: settings.Hardware{
		Lcd:          settings.Lcd{},
		TagReader:    settings.TagReader{},
		LedIndicator: settings.LedIndicator{},
	},
},
}

type chargePointTestSuite struct {
	suite.Suite
	centralSystem ocpp16.CentralSystem
	csMock        *centralSystemMock
	tagReader     *readerMock
	display       *displayMock
	manager       *managerMock
}

func (s *chargePointTestSuite) SetupTest() {
	s.csMock = new(centralSystemMock)
	s.tagReader = new(readerMock)
	s.display = new(displayMock)
	s.manager = new(managerMock)
}

func (s *chargePointTestSuite) setupCentralSystem(cs *centralSystemMock) {
	wsServer := ws.NewServer()
	wsServer.SetCheckOriginHandler(func(r *http.Request) bool {
		return true
	})

	s.centralSystem = ocpp16.NewCentralSystem(nil, wsServer)
	s.centralSystem.SetCoreHandler(cs)
	s.centralSystem.Start(centralSystemPort, centralSystemEndpoint+"/{ws}")
}

func (s *chargePointTestSuite) setupChargePoint(ctx context.Context, lcd display.LCD, reader reader.Reader, connectorManager *managerMock) chargepoint.ChargePoint {
	chargePoint := v16.NewChargePoint(reader, lcd, connectorManager, scheduler.GetScheduler(), auth.NewAuthCache("./auth.json"))
	chargePoint.Init(ctx, &chargePointSettings)
	chargePoint.Connect(ctx, fmt.Sprintf("ws://%s:%d%s", centralSystemHost, centralSystemPort, centralSystemEndpoint))
	return chargePoint
}

func (s *chargePointTestSuite) TestStartTransaction() {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
		conn        = new(connectorMock)
	)

	// Setup OCPP configuration manager
	setting.SetupOcppConfigurationManager(
		"../../configs/configuration.json",
		configuration.OCPP16,
		core.ProfileName,
		reservation.ProfileName)

	// Mock central system
	s.csMock.On("OnAuthorize", mock.Anything, nil).
		Return(core.NewAuthorizationConfirmation(types.NewIdTagInfo(types.AuthorizationStatusAccepted)), nil)

	s.csMock.On("OnBootNotification", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {

	}).Return(core.NewBootNotificationConfirmation(types.NewDateTime(time.Now()), 500, core.RegistrationStatusAccepted), nil)

	s.csMock.On("OnHeartbeat", mock.Anything, mock.Anything).
		Return(core.NewHeartbeatConfirmation(types.NewDateTime(time.Now())), nil)

	s.csMock.On("OnMeterValues", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {

	}).Return(core.NewMeterValuesConfirmation(), nil)

	s.csMock.On("OnStatusNotification", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {

	}).Return(core.NewStatusNotificationConfirmation(), nil)

	s.csMock.On("OnStartTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {

	}).Return(core.NewStartTransactionConfirmation(types.NewIdTagInfo(types.AuthorizationStatusAccepted), 1), nil)

	s.csMock.On("OnStopTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {

	}).Return(core.NewStopTransactionConfirmation(), nil)

	// Setup expectations
	conn.On("StartCharging", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	conn.On("ResumeCharging", mock.Anything).Return(nil, 0)
	conn.On("StopCharging", core.ReasonLocal).Return(nil)
	conn.On("RemoveReservation").Return(nil)
	conn.On("GetReservationId").Return(0)
	conn.On("GetTransaction").Return(0).Once()
	conn.On("GetTransaction").Return(1)
	conn.On("GetConnectorId").Return(1)
	conn.On("GetEvseId").Return(1)
	conn.On("CalculateSessionAvgEnergyConsumption").Return(30.0)
	conn.On("GetStatus").Return(core.ChargePointStatusAvailable, core.NoError)
	conn.On("IsAvailable").Return(true)
	conn.On("IsPreparing").Return(false)
	conn.On("IsCharging").Return(false)
	conn.On("IsReserved").Return(false)
	conn.On("IsUnavailable").Return(false)
	conn.On("GetMaxChargingTime").Return(15)
	conn.On("SetNotificationChannel", mock.Anything).Return()

	s.manager.On("GetConnectors").Return([]connector.Connector{conn})
	s.manager.On("FindConnector", 1, 1).Return(conn)
	s.manager.On("FindAvailableConnector").Return(conn)
	s.manager.On("FindConnectorWithTagId", tagId).Return(conn)
	s.manager.On("FindConnectorWithTransactionId", "1").Return(conn)
	s.manager.On("StartChargingConnector").Return()
	s.manager.On("StopChargingConnector").Return()
	s.manager.On("StopAllConnectors").Return()
	s.manager.On("AddConnector", conn).Return(nil)
	s.manager.On("AddConnectorFromSettings", mock.Anything).Return(nil)
	s.manager.On("AddConnectorsFromConfiguration", mock.Anything).Return(nil)
	s.manager.On("RestoreConnectorStatus", mock.Anything).Return(nil)
	s.manager.On("SetNotificationChannel").Return()

	// Start the central system
	go s.setupCentralSystem(s.csMock)

	time.Sleep(time.Second * 5)

	// Create and connect the Charge Point
	chargePoint := s.setupChargePoint(ctx, nil, nil, s.manager)

	// Start charging
	chargePoint.HandleChargingRequest(tagId)

	// Stop charging
	chargePoint.HandleChargingRequest(tagId)

Loop:
	for {
		select {
		case <-ctx.Done():
			cancel()
			break Loop
		}
	}
}

func (s *chargePointTestSuite) TestStartTransactionWithReader() {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
		conn        = new(connectorMock)
	)

	// Setup OCPP configuration manager
	setting.SetupOcppConfigurationManager(
		"../../configs/configuration.json",
		configuration.OCPP16,
		core.ProfileName,
		reservation.ProfileName)

	// Mock central system
	s.csMock.On("OnAuthorize", mock.Anything, nil).
		Return(core.NewAuthorizationConfirmation(types.NewIdTagInfo(types.AuthorizationStatusAccepted)), nil)

	s.csMock.On("OnBootNotification", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {

	}).Return(core.NewBootNotificationConfirmation(types.NewDateTime(time.Now()), 500, core.RegistrationStatusAccepted), nil)

	s.csMock.On("OnHeartbeat", mock.Anything, mock.Anything).
		Return(core.NewHeartbeatConfirmation(types.NewDateTime(time.Now())), nil)

	s.csMock.On("OnMeterValues", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {

	}).Return(core.NewMeterValuesConfirmation(), nil)

	s.csMock.On("OnStatusNotification", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {

	}).Return(core.NewStatusNotificationConfirmation(), nil)

	s.csMock.On("OnStartTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {

	}).Return(core.NewStartTransactionConfirmation(types.NewIdTagInfo(types.AuthorizationStatusAccepted), 1), nil)

	s.csMock.On("OnStopTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {

	}).Return(core.NewStopTransactionConfirmation(), nil)

	// Setup expectations
	conn.On("StartCharging", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	conn.On("ResumeCharging", mock.Anything).Return(nil, 0)
	conn.On("StopCharging", core.ReasonLocal).Return(nil)
	conn.On("RemoveReservation").Return(nil)
	conn.On("GetReservationId").Return(0)
	conn.On("GetTransaction").Return(1).Once()
	conn.On("GetConnectorId").Return(1)
	conn.On("GetEvseId").Return(1)
	conn.On("CalculateSessionAvgEnergyConsumption").Return(30.0)
	conn.On("GetStatus").Return(core.ChargePointStatusAvailable, core.NoError)
	conn.On("IsAvailable").Return(true)
	conn.On("IsPreparing").Return(false)
	conn.On("IsCharging").Return(false)
	conn.On("IsReserved").Return(false)
	conn.On("IsUnavailable").Return(false)
	conn.On("GetMaxChargingTime").Return(15)
	conn.On("SetNotificationChannel", mock.Anything).Return()

	s.manager.On("GetConnectors").Return([]connector.Connector{conn})
	s.manager.On("FindConnector", 1, 1).Return(conn)
	s.manager.On("FindAvailableConnector").Return(conn)
	s.manager.On("FindConnectorWithTagId", tagId).Return(conn)
	s.manager.On("FindConnectorWithTransactionId", "1").Return(conn)
	s.manager.On("StartChargingConnector").Return()
	s.manager.On("StopChargingConnector").Return()
	s.manager.On("StopAllConnectors").Return()
	s.manager.On("AddConnector", conn).Return(nil)
	s.manager.On("AddConnectorFromSettings", mock.Anything).Return(nil)
	s.manager.On("AddConnectorsFromConfiguration", mock.Anything).Return(nil)
	s.manager.On("RestoreConnectorStatus", mock.Anything).Return(nil)
	s.manager.On("SetNotificationChannel").Return()

	// Mock tagReader
	readerChannel := make(chan string)
	s.tagReader.On("ListenForTags").Return()
	s.tagReader.On("Cleanup").Return()
	s.tagReader.On("Reset").Return()
	s.tagReader.On("GetTagChannel").Return(readerChannel)

	// Start the central system
	go s.setupCentralSystem(s.csMock)

	time.Sleep(time.Second * 5)

	// Create and connect the Charge Point
	_ = s.setupChargePoint(ctx, s.display, s.tagReader, s.manager)

	// Simulate reading a card
	go func() {
		// Card read once - start charging
		readerChannel <- tagId

		time.Sleep(time.Second * 10)
		s.tagReader.AssertCalled(s.T(), "ListenForTags")
		s.tagReader.AssertCalled(s.T(), "GetTagChannel")

		// Card read second time - stop charging
		readerChannel <- tagId
		time.Sleep(time.Second * 10)

		cancel()
	}()

Loop:
	for {
		select {
		case <-ctx.Done():
			cancel()
			break Loop
		}
	}
}

func Test(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	suite.Run(t, new(chargePointTestSuite))
}
