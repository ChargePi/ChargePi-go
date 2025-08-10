package chargepoint

import (
	"context"
	"time"

	"github.com/ChargePi/ChargePi-go/internal/auth"
	"github.com/ChargePi/ChargePi-go/internal/diagnostics"
	"github.com/ChargePi/ChargePi-go/internal/evse/manager"
	"github.com/ChargePi/ChargePi-go/internal/pkg/database"
	chargePoint "github.com/ChargePi/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/ChargePi/ChargePi-go/internal/pkg/models/settings"
	cfg "github.com/ChargePi/ChargePi-go/internal/pkg/settings"
	"github.com/ChargePi/ChargePi-go/internal/pkg/util"
	database2 "github.com/ChargePi/ChargePi-go/internal/sessions/pkg/database"
	"github.com/ChargePi/ChargePi-go/internal/sessions/service/session"
	"github.com/ChargePi/ChargePi-go/pkg/observability/logging"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"go.uber.org/zap"
)

func Run(ctx context.Context, debug bool, config *settings.Settings) {
	// Create a logger
	logger := logging.SetupZap(config.ChargePoint.Logging, debug)
	zap.ReplaceGlobals(logger)
	defer logging.Sync(logger)

	// Create named loggers for components
	evseLogger := logger.Named("evse")

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

	settingsManager := cfg.GetManager()

	evseManager := manager.GetManager()
	diagnosticsManager := diagnostics.NewManager(logger)
	tagManager := auth.NewTagManager(logger, db)
	sessionRepository := database2.NewSessionBadgerDb(logger, db)
	sessionManager := session.NewSessionManager(logger, sessionRepository)

	// Initialize all the EVSEs
	err := evseManager.InitAll(ctx)
	if err != nil {
		evseLogger.Fatal("Cannot add EVSEs", zap.Error(err))
	}

	// Create a context for the OCPP connection, so it can be dynamically reconnected.
	parentCtxForOcpp, parentCancel := context.WithCancel(ctx)
	defer parentCancel()

	// Set the settings

	handler = CreateChargePoint(parentCtxForOcpp, protocolVersion, logger, evseManager, tagManager, sessionManager, settingsManager, diagnosticsManager, hardware)
	err = handler.SetSettings(chargePointInfo)
	if err != nil {
		logger.Fatal("Unable to set the charge point settings", zap.Error(err))
	}

	err = handler.SetConnectionSettings(connectionSettings)
	if err != nil {
		logger.Fatal("Unable to set the connection settings", zap.Error(err))
	}

	// Listen for connector status changes
	go handler.ListenForConnectorStatusChange(ctx, evseManager.GetNotificationChannel())

	// Start the UI and API
	go SetupApi(db, config.Api, handler, tagManager, evseManager, settingsManager)
	go SetupUi(config.Ui)
	go setupHealthcheck()

	// Connect to the backend system
	go handler.Connect(parentCtxForOcpp, serverUrl)

	<-ctx.Done()
	handler.CleanUp(core.ReasonLocal)
	time.Sleep(time.Millisecond * 500)
	logger.Info("Shutting down ChargePi...")
}
