package chargepoint

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/xBlaz3kx/ChargePi-go/internal/sessions/service/session"
	"github.com/xBlaz3kx/ChargePi-go/pkg/logging"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/xBlaz3kx/ChargePi-go/internal/auth"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	cfg "github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
)

func Run(debug bool, config *settings.Settings) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Create a logger
	logger := log.StandardLogger()
	logging.Setup(logger, config.ChargePoint.Logging, debug)

	var (
		handler            chargePoint.ChargePoint
		hardware           = config.ChargePoint.Hardware
		connectionSettings = config.ChargePoint.ConnectionSettings
		chargePointInfo    = config.ChargePoint.Info
		protocolVersion    = connectionSettings.ProtocolVersion
		serverUrl          = util.CreateConnectionUrl(connectionSettings)
	)

	// Create a database for EVSE, tags, users and settings
	db := database.Get()

	// Setup OCPP configuration from the database
	settingsManager := cfg.GetManager()
	settingsManager.SetupOcppConfiguration(
		configuration.ProtocolVersion(protocolVersion),
		core.ProfileName,
		reservation.ProfileName,
		remotetrigger.ProfileName,
		localauth.ProfileName,
	)
	ocppVariableManager := ocppConfigManager.GetManager()
	evseManager := evse.GetManager()
	tagManager := auth.NewTagManager(db)
	sessionManager := session.NewSessionManager(nil)

	// Initialize all the EVSEs
	err := evseManager.InitAll(ctx)
	if err != nil {
		logger.WithError(err).Fatal("Cannot add EVSEs")
	}

	// Create a context for the OCPP connection, so it can be dynamically reconnected.
	parentCtxForOcpp, parentCancel := context.WithCancel(ctx)
	defer parentCancel()

	// Set the settings and connect to the backend system
	handler = CreateChargePoint(parentCtxForOcpp, protocolVersion, logger, evseManager, tagManager, sessionManager, hardware)
	handler.SetSettings(chargePointInfo)
	handler.SetConnectionSettings(connectionSettings)

	// Listen for connector status changes
	go handler.ListenForConnectorStatusChange(ctx, evseManager.GetNotificationChannel())

	// Start the UI and API
	go SetupApi(db, config.Api, handler, tagManager, evseManager, ocppVariableManager)
	go SetupUi(config.Ui)
	go setupApplicationHealthcheck()

	// Connect to the backend system
	go handler.Connect(parentCtxForOcpp, serverUrl)

	<-ctx.Done()
	handler.CleanUp(core.ReasonLocal)
	time.Sleep(time.Millisecond * 500)
	logger.Info("Shutting down ChargePi...")
}
