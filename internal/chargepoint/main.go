package chargepoint

import (
	"context"
	"github.com/go-co-op/gocron"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/auth"
	connectorManager "github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	v16 "github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/v16"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/grpc"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/scheduler"
	s "github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"github.com/xBlaz3kx/ChargePi-go/pkg/logging"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/ocpp"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
	"os"
	"os/signal"
	"time"
)

func CreateChargePoint(
	ctx context.Context,
	protocolVersion ocpp.ProtocolVersion,
	logger *log.Logger,
	manager connectorManager.Manager,
	sch *gocron.Scheduler,
	tagManager auth.TagManager,
	hardware settings.Hardware,
) chargePoint.ChargePoint {

	opts := []chargePoint.Options{
		chargePoint.WithDisplayFromSettings(hardware.Display),
		chargePoint.WithReaderFromSettings(ctx, hardware.TagReader),
		chargePoint.WithLogger(logger),
	}

	switch protocolVersion {
	case ocpp.OCPP16:
		// Create the client
		return v16.NewChargePoint(
			manager,
			sch,
			tagManager,
			opts...,
		)
	case ocpp.OCPP201:
		logger.Fatal("Version 2.0.1 is not supported yet.")
		return nil
	default:
		logger.WithField("protocolVersion", protocolVersion).Fatal("Protocol version not supported")
		return nil
	}
}

// Run is an entrypoint with all the configuration needed. This is a blocking function.
func Run(isDebug bool, config *settings.Settings, connectors []*settings.EVSE, configurationFilePath, localAuthListFilePath string) {
	var (
		// ChargePoint components
		handler    chargePoint.ChargePoint
		tagManager = auth.NewTagManager(localAuthListFilePath)
		logger     = log.StandardLogger()
		manager    = connectorManager.GetManager()
		sch        = scheduler.GetScheduler()
		// Settings
		hardware           = config.ChargePoint.Hardware
		connectionSettings = config.ChargePoint.ConnectionSettings
		serverUrl          = util.CreateConnectionUrl(connectionSettings)
		protocolVersion    = connectionSettings.ProtocolVersion

		ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt)
	)

	defer cancel()

	// Create the logger
	logging.Setup(logger, config.ChargePoint.Logging, isDebug)

	// Setup OCPP configuration manager
	s.SetupOcppConfigurationManager(
		configurationFilePath,
		configuration.ProtocolVersion(connectionSettings.ProtocolVersion),
		core.ProfileName,
		reservation.ProfileName,
	)

	// Load tags
	go func() {
		err := tagManager.ReadLocalAuthList()
		if err != nil {

		}
	}()

	// Create a new context just for the OCPP connection, so it can be dynamically rebooted
	parentCtxForOcpp, parentCancel := context.WithCancel(ctx)
	defer parentCancel()

	// Initialize the client
	handler = CreateChargePoint(parentCtxForOcpp, protocolVersion, logger, manager, sch, tagManager, hardware)
	handler.SetSettings(config)
	handler.AddEVSEs(connectors)

	// Finally, connect to the central system
	handler.Connect(ctx, serverUrl)


	<-ctx.Done()
	handler.CleanUp(core.ReasonLocal)
	time.Sleep(time.Millisecond * 500)
}
