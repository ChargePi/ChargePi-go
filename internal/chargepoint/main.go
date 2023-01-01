package chargepoint

import (
	"context"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	settings2 "github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"os"
	"os/signal"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/api/grpc"
	"github.com/xBlaz3kx/ChargePi-go/internal/api/http"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/auth"
	connectorManager "github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
	v16 "github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/v16"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/scheduler"
	s "github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"github.com/xBlaz3kx/ChargePi-go/pkg/logging"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/ocpp"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

func CreateChargePoint(
	ctx context.Context,
	protocolVersion ocpp.ProtocolVersion,
	logger *log.Logger,
	manager connectorManager.Manager,
	sch *gocron.Scheduler,
	tagManager auth.TagManager,
	hardware settings2.Hardware,
) chargePoint.ChargePoint {

	// Create a status indicator if enabled
	statusIndicator := indicator.NewIndicator(len(manager.GetEVSEs()), hardware.LedIndicator)

	// Create additional components based on the configuration
	opts := []chargePoint.Options{
		chargePoint.WithDisplayFromSettings(hardware.Display),
		chargePoint.WithReaderFromSettings(ctx, hardware.TagReader),
		chargePoint.WithLogger(logger),
		chargePoint.WithIndicator(statusIndicator),
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
func Run(isDebug bool, config *settings2.Settings, connectors []*settings2.EVSE, configurationFilePath, localAuthListFilePath string) {
	var (
		// ChargePoint components
		handler    chargePoint.ChargePoint
		logger     = log.StandardLogger()
		tagManager = auth.NewTagManager(localAuthListFilePath)
		manager    = connectorManager.GetManager()
		sch        = scheduler.GetScheduler()

		// Settings
		hardware           = config.ChargePoint.Hardware
		connectionSettings = config.ChargePoint.ConnectionSettings
		chargePointInfo    = config.ChargePoint.Info
		protocolVersion    = connectionSettings.ProtocolVersion
		serverUrl          = util.CreateConnectionUrl(connectionSettings)

		ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt)
	)

	defer cancel()

	// Create the logger
	logging.Setup(logger, config.ChargePoint.Logging, isDebug)

	// Setup OCPP configuration variables manager
	s.SetupOcppConfigurationManager(
		configurationFilePath,
		configuration.ProtocolVersion(connectionSettings.ProtocolVersion),
		core.ProfileName,
		reservation.ProfileName,
		remotetrigger.ProfileName,
		localauth.ProfileName,
	)

	// Load the local auth list of tags
	go func() {
		err := tagManager.ReadLocalAuthList()
		if err != nil {
			logger.WithError(err).Error("Cannot read local auth list")
		}
	}()

	// Add EVSEs, they will run standalone
	err := manager.AddEVSEsFromSettings(ctx, chargePointInfo.MaxChargingTime, connectors)
	if err != nil {
		logger.WithError(err).Fatal("Cannot add EVSEs")
	}

	// Create a new context just for the OCPP connection, so it can be dynamically rebooted
	parentCtxForOcpp, parentCancel := context.WithCancel(ctx)
	defer parentCancel()

	// Create an OCPP client and connect
	handler = CreateChargePoint(parentCtxForOcpp, protocolVersion, logger, manager, sch, tagManager, hardware)
	handler.SetSettings(chargePointInfo)
	handler.SetConnectionSettings(connectionSettings)
	go handler.ListenForConnectorStatusChange(ctx, manager.GetNotificationChannel())
	go handler.Connect(parentCtxForOcpp, serverUrl)

	// Expose the API endpoints
	server := grpc.NewServer(config.Api, handler, tagManager, manager, ocppConfigManager.GetManager(), nil)
	go server.Run()

	// Expose the ui at http://localhost:4269/
	ui := http.NewUi()
	go ui.Serve("0.0.0.0:4269")

	<-ctx.Done()
	handler.CleanUp(core.ReasonLocal)
	time.Sleep(time.Millisecond * 500)
}
