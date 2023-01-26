package chargepoint

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/xBlaz3kx/ChargePi-go/internal/api/service"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	s "github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/users"
	uDb "github.com/xBlaz3kx/ChargePi-go/internal/users/database"
	"github.com/xBlaz3kx/ChargePi-go/pkg/logging"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/api/http"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/auth"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
	v16 "github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/v16"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/ocpp"
)

func CreateChargePoint(
	ctx context.Context,
	protocolVersion ocpp.ProtocolVersion,
	logger *log.Logger,
	manager evse.Manager,
	tagManager auth.TagManager,
	hardware settings.Hardware,
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

func SetupUserApi(db *badger.DB, api settings.Api, handler chargePoint.ChargePoint, tagManager auth.TagManager, manager evse.Manager, ocppVariableManager ocppConfigManager.Manager) {
	// User database layer
	userDb := uDb.NewUserDb(db)

	// User service layer
	userService := users.NewUserService(userDb)

	// Expose the API endpoints
	server := service.NewServer(api, handler, tagManager, manager, ocppVariableManager, userService)
	go server.Run()

	// Launch UI at http://localhost:4269/
	// The UI should be integrated for portability.
	ui := http.NewUi()
	go ui.Serve("0.0.0.0:4269")
}

// Run is an entrypoint with all the configuration needed. This is a blocking function.
func Run(debug bool, config *settings.Settings) {
	var (
		ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt)

		handler     chargePoint.ChargePoint
		logger      = log.StandardLogger()
		evseManager = evse.GetManager()

		// Settings
		hardware           = config.ChargePoint.Hardware
		connectionSettings = config.ChargePoint.ConnectionSettings
		chargePointInfo    = config.ChargePoint.Info
		protocolVersion    = connectionSettings.ProtocolVersion
		serverUrl          = util.CreateConnectionUrl(connectionSettings)
	)
	defer cancel()

	// Load/initialize a database for EVSE, tags, users and settings
	db := database.Get()

	tagManager := auth.NewTagManager(db)

	// Setup OCPP configuration from the database.
	s.SetupOcppConfiguration(db, configuration.ProtocolVersion(protocolVersion), core.ProfileName, reservation.ProfileName, remotetrigger.ProfileName, localauth.ProfileName)

	// Create the logger
	logging.Setup(logger, config.ChargePoint.Logging, debug)

	// Initialize all the EVSEs, if there are any
	err := evseManager.InitAll(ctx)
	if err != nil {
		logger.WithError(err).Fatal("Cannot add EVSEs")
	}

	// Create a context for the OCPP connection, so it can be dynamically reconnected.
	parentCtxForOcpp, parentCancel := context.WithCancel(ctx)
	defer parentCancel()

	// Create a Charge point and connect
	handler = CreateChargePoint(parentCtxForOcpp, protocolVersion, logger, evseManager, tagManager, hardware)
	handler.SetSettings(chargePointInfo)
	handler.SetConnectionSettings(connectionSettings)

	go handler.ListenForConnectorStatusChange(ctx, evseManager.GetNotificationChannel())
	go handler.Connect(parentCtxForOcpp, serverUrl)

	// Setup User API
	ocppVariableManager := ocppConfigManager.GetManager()
	SetupUserApi(db, config.Api, handler, tagManager, evseManager, ocppVariableManager)

	<-ctx.Done()
	handler.CleanUp(core.ReasonLocal)
	time.Sleep(time.Millisecond * 500)
	logger.Info("Exiting..")
}
