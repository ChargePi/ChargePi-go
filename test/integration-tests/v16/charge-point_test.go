package v16

import (
	"context"
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ws"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/auth"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/reader"
	v16 "github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/v16"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/scheduler"
	setting "github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
	"github.com/xBlaz3kx/ChargePi-go/test"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
	"net/http"
	"strings"
	"testing"
	"time"
)

const (
	centralSystemEndpoint = ""
	centralSystemPort     = 7777
	centralSystemHost     = "localhost"

	tagId                     = "exampleTag"
	chargePointId             = "exampleChargePoint"
	protocolVersion           = settings.OCPP16
	ocppConfigurationFilePath = "../../../configs/configuration.json"
)

var chargePointSettings = settings.Settings{ChargePoint: settings.ChargePoint{
	Info: settings.Info{
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
		Display:      settings.Display{IsEnabled: false},
		TagReader:    settings.TagReader{IsEnabled: true},
		LedIndicator: settings.LedIndicator{IndicateCardRead: false},
	},
},
}

type chargePointTestSuite struct {
	suite.Suite
	centralSystem ocpp16.CentralSystem
	csMock        *centralSystemV16Mock
	tagReader     *test.ReaderMock
	display       *test.DisplayMock
	manager       *test.ManagerMock
}

func (s *chargePointTestSuite) SetupTest() {
	log.SetLevel(log.DebugLevel)

	s.csMock = new(centralSystemV16Mock)
	s.tagReader = new(test.ReaderMock)
	s.display = new(test.DisplayMock)
	s.manager = new(test.ManagerMock)

	// Setup OCPP configuration manager
	setting.SetupOcppConfigurationManager(
		ocppConfigurationFilePath,
		configuration.OCPP16,
		core.ProfileName,
		reservation.ProfileName)

	// Start the central system
	go s.setupCentralSystem(s.csMock)
}

func (s *chargePointTestSuite) setupCentralSystem(cs *centralSystemV16Mock) {
	s.setupCentralSystemCoreExpectations()

	log.Debug("Setting up central system")
	wsServer := ws.NewServer()
	wsServer.SetCheckOriginHandler(func(r *http.Request) bool {
		return true
	})

	s.centralSystem = ocpp16.NewCentralSystem(nil, wsServer)
	s.centralSystem.SetCoreHandler(cs)
	s.centralSystem.Start(centralSystemPort, centralSystemEndpoint+"/{ws}")
}

func (s *chargePointTestSuite) setupCentralSystemCoreExpectations() {
	// Mock central system
	s.csMock.On("OnAuthorize", "/"+chargePointId, mock.Anything).Run(func(args mock.Arguments) {
		payload := args.Get(1)
		s.Require().NotNil(payload)
		s.Require().IsType(&core.AuthorizeRequest{}, payload)
		if payload != nil {
			authRequest := payload.(*core.AuthorizeRequest)
			s.Assert().Equal(strings.ToLower(tagId), strings.ToLower(authRequest.IdTag))
		}
	}).Return(core.NewAuthorizationConfirmation(types.NewIdTagInfo(types.AuthorizationStatusAccepted)), nil)

	s.csMock.On("OnBootNotification", "/"+chargePointId, mock.Anything).Run(func(args mock.Arguments) {
		payload := args.Get(1)
		s.Require().NotNil(payload)
		s.Require().IsType(&core.BootNotificationRequest{}, payload)
		if payload != nil {
			bootNotification := payload.(*core.BootNotificationRequest)
			s.Assert().Equal("exampleVendor", bootNotification.ChargePointVendor)
			s.Assert().Equal("exampleModel", bootNotification.ChargePointModel)
		}
	}).Return(core.NewBootNotificationConfirmation(types.NewDateTime(time.Now()), 500, core.RegistrationStatusAccepted), nil)

	s.csMock.On("OnHeartbeat", "/"+chargePointId, mock.Anything).Return(core.NewHeartbeatConfirmation(types.NewDateTime(time.Now())), nil)

	s.csMock.On("OnMeterValues", "/"+chargePointId, mock.Anything).Run(func(args mock.Arguments) {

	}).Return(core.NewMeterValuesConfirmation(), nil)

	s.csMock.On("OnStatusNotification", "/"+chargePointId, mock.Anything).Run(func(args mock.Arguments) {
		payload := args.Get(1)
		s.Require().NotNil(payload)
		s.Require().IsType(&core.StatusNotificationRequest{}, payload)
		if payload != nil {
			statusNotificationRequest := payload.(*core.StatusNotificationRequest)
			//s.Assert().Subset(core.ConfigurationStatusAccepted, statusNotificationRequest.Status)
			s.Assert().Equal(1, statusNotificationRequest.ConnectorId)
			s.Assert().Equal(core.NoError, statusNotificationRequest.ErrorCode)
		}
	}).Return(core.NewStatusNotificationConfirmation(), nil)

	s.csMock.On("OnStartTransaction", "/"+chargePointId, mock.Anything).Run(func(args mock.Arguments) {
		payload := args.Get(1)
		s.Require().NotNil(payload)
		s.Require().IsType(&core.StartTransactionRequest{}, payload)
		if payload != nil {
			startTransactionRequest := payload.(*core.StartTransactionRequest)
			s.Assert().Equal(1, startTransactionRequest.ConnectorId)
			s.Assert().Equal(strings.ToLower(tagId), strings.ToLower(startTransactionRequest.IdTag))
		}
	}).Return(core.NewStartTransactionConfirmation(types.NewIdTagInfo(types.AuthorizationStatusAccepted), 1), nil)

	s.csMock.On("OnStopTransaction", "/"+chargePointId, mock.Anything).Run(func(args mock.Arguments) {
		payload := args.Get(1)
		s.Require().NotNil(payload)
		s.Require().IsType(&core.StopTransactionRequest{}, payload)
		if payload != nil {
			stopTransactionRequest := payload.(*core.StopTransactionRequest)
			s.Assert().Equal(1, stopTransactionRequest.TransactionId)
			//s.Assert().Equal(tagId, stopTransactionRequest.IdTag)
		}
	}).Return(core.NewStopTransactionConfirmation(), nil)
}

func (s *chargePointTestSuite) setupChargePoint(ctx context.Context, lcd display.Display, reader reader.Reader, connectorManager *test.ManagerMock) chargePoint.ChargePoint {
	cp := v16.NewChargePoint(
		connectorManager,
		scheduler.GetScheduler(),
		auth.NewAuthCache("./auth.json"),
		chargePoint.WithDisplay(lcd),
		chargePoint.WithReader(ctx, reader),
		chargePoint.WithLogger(log.StandardLogger()),
	)
	cp.SetSettings(&chargePointSettings)
	cp.Connect(ctx, fmt.Sprintf("ws://%s:%d%s", centralSystemHost, centralSystemPort, centralSystemEndpoint))
	return cp
}

func Test(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	suite.Run(t, new(chargePointTestSuite))
}
