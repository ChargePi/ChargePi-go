package integration_tests

import (
	"context"
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint"
	v16 "github.com/xBlaz3kx/ChargePi-go/chargepoint/v16"
	"github.com/xBlaz3kx/ChargePi-go/components/connector"
	"github.com/xBlaz3kx/ChargePi-go/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/components/hardware/reader"
	"github.com/xBlaz3kx/ChargePi-go/components/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/data/settings"
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
}

func (s *chargePointTestSuite) SetupTest() {
	s.csMock = new(centralSystemMock)
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

func (s *chargePointTestSuite) setupChargePoint(ctx context.Context, lcd display.LCD, reader reader.Reader, connector ...connector.Connector) chargepoint.ChargePoint {
	chargePoint := v16.NewChargePoint(reader, lcd, nil, scheduler.GetScheduler(), nil)
	chargePoint.Init(ctx, &chargePointSettings)
	chargePoint.Connect(ctx, fmt.Sprintf("ws://%s:%d%s", centralSystemHost, centralSystemPort, centralSystemEndpoint))
	//chargePoint.AddConnectors()
	return chargePoint
}

func (s *chargePointTestSuite) TestStartTransaction() {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	)

	s.csMock.On("OnAuthorize", nil).
		Return(core.NewAuthorizationConfirmation(types.NewIdTagInfo(types.AuthorizationStatusAccepted)), nil)

	s.csMock.On("OnBootNotification").
		Return(core.NewBootNotificationConfirmation(types.NewDateTime(time.Now()), 500, core.RegistrationStatusAccepted), nil)

	s.csMock.On("OnHeartbeat").
		Return(core.NewHeartbeatConfirmation(types.NewDateTime(time.Now())), nil)

	s.csMock.On("OnMeterValues").
		Return(core.NewMeterValuesConfirmation(), nil)

	s.csMock.On("OnStatusNotification").
		Return(core.NewStatusNotificationConfirmation(), nil)

	s.csMock.On("OnStartTransaction").
		Return(core.NewStartTransactionConfirmation(types.NewIdTagInfo(types.AuthorizationStatusAccepted), 1), nil)

	s.csMock.On("OnStopTransaction").
		Return(core.NewStopTransactionConfirmation(), nil)

	// Start the central system
	go s.setupCentralSystem(s.csMock)

	time.Sleep(time.Second * 5)

	// Create and connect the Charge Point
	chargePoint := s.setupChargePoint(ctx, nil, nil)

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
	)

	// Mock central system
	s.csMock.On("OnAuthorize", nil).
		Return(core.NewAuthorizationConfirmation(types.NewIdTagInfo(types.AuthorizationStatusAccepted)), nil)

	s.csMock.On("OnBootNotification").
		Return(core.NewBootNotificationConfirmation(types.NewDateTime(time.Now()), 500, core.RegistrationStatusAccepted), nil)

	s.csMock.On("OnHeartbeat").
		Return(core.NewHeartbeatConfirmation(types.NewDateTime(time.Now())), nil)

	s.csMock.On("OnMeterValues").
		Return(core.NewMeterValuesConfirmation(), nil)

	s.csMock.On("OnStatusNotification").
		Return(core.NewStatusNotificationConfirmation(), nil)

	s.csMock.On("OnStartTransaction").
		Return(core.NewStartTransactionConfirmation(types.NewIdTagInfo(types.AuthorizationStatusAccepted), 1), nil)

	s.csMock.On("OnStopTransaction").
		Return(core.NewStopTransactionConfirmation(), nil)

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
	_ = s.setupChargePoint(ctx, nil, s.tagReader)

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
