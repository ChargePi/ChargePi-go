package chargepoint

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/api"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/util"
	v16 "github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/v16"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/auth"
	connectorManager "github.com/xBlaz3kx/ChargePi-go/internal/components/connector-manager"
	s "github.com/xBlaz3kx/ChargePi-go/internal/components/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/grpc"
	"github.com/xBlaz3kx/ChargePi-go/pkg/logging"
	"github.com/xBlaz3kx/ChargePi-go/pkg/scheduler"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
	"os"
	"os/signal"
	"syscall"
)

func CreateChargePoint(
	ctx context.Context,
	protocolVersion settings.ProtocolVersion,
	logger *log.Logger,
	manager connectorManager.Manager,
	sch *gocron.Scheduler,
	authCache *auth.Cache,
	hardware settings.Hardware,
) chargePoint.ChargePoint {
	switch protocolVersion {
	case settings.OCPP16:
		// Create the client
		return v16.NewChargePoint(
			manager,
			sch,
			authCache,
			v16.WithDisplayFromSettings(ctx, hardware.Display),
			v16.WithReaderFromSettings(ctx, hardware.TagReader),
			v16.WithLogger(logger),
		)
	case settings.OCPP201:
		logger.Fatal("Version 2.0.1 is not supported yet.")
		return nil
	default:
		logger.WithField("protocolVersion", protocolVersion).Fatal("Protocol version not supported")
		return nil
	}
}

func Run(isDebug bool, config *settings.Settings, connectors []*settings.Connector, configurationFilePath, authFilePath string) {
	var (
		// ChargePoint components
		handler   chargePoint.ChargePoint
		authCache = auth.NewAuthCache(authFilePath)
		logger    = log.StandardLogger()
		manager   = connectorManager.GetManager()
		sch       = scheduler.GetScheduler()
		// Settings
		chargePointInfo = config.ChargePoint.Info
		hardware        = config.ChargePoint.Hardware
		serverUrl       = util.CreateConnectionUrl(config.ChargePoint)
		protocolVersion = settings.ProtocolVersion(chargePointInfo.ProtocolVersion)
		// Execution
		ctx, cancel = context.WithCancel(context.Background())
		quitChannel = make(chan os.Signal, 5)
	)

	defer cancel()
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	// Create the logger
	logging.Setup(logger, config.ChargePoint.Logging, isDebug)

	// Load tags
	go authCache.LoadAuthFile()

	// Setup OCPP configuration manager
	s.SetupOcppConfigurationManager(
		configurationFilePath,
		configuration.ProtocolVersion(chargePointInfo.ProtocolVersion),
		core.ProfileName,
		reservation.ProfileName)

	// Initialize the client
	handler = CreateChargePoint(ctx, protocolVersion, logger, manager, sch, authCache, hardware)
	handler.Init(config)
	handler.AddConnectors(connectors)

	// Finally, connect to the central system
	handler.Connect(ctx, serverUrl)

	if config.Api.Enabled {
		var (
			apiReceiveChannel = make(chan api.Message, 5)
			apiSendChannel    = make(chan api.Message, 5)
		)

		// Expose the API endpoints
		go func() {
			address := fmt.Sprintf("%s:%d", config.Api.Address, config.Api.Port)
			grpc.CreateAndRunGrpcServer(address, apiSendChannel, apiReceiveChannel)
		}()
	}

Loop:
	for {
		select {
		// Capture the terminate signal
		case <-quitChannel:
			break Loop
		case <-ctx.Done():
			break Loop
		}
	}

	handler.CleanUp(core.ReasonLocal)
}
