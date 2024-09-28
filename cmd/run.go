package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ChargePi/ChargePi-go/internal/api/grpc"
	"github.com/ChargePi/ChargePi-go/internal/api/http"
	"github.com/ChargePi/ChargePi-go/internal/auth"
	"github.com/ChargePi/ChargePi-go/internal/chargepoint"
	"github.com/ChargePi/ChargePi-go/internal/diagnostics"
	"github.com/ChargePi/ChargePi-go/internal/evse/manager"
	"github.com/ChargePi/ChargePi-go/internal/pkg/badger"
	"github.com/ChargePi/ChargePi-go/internal/pkg/configuration"
	configManager "github.com/ChargePi/ChargePi-go/internal/pkg/configuration/manager"
	"github.com/ChargePi/ChargePi-go/internal/pkg/notifications"
	"github.com/ChargePi/ChargePi-go/internal/sessions"
	"github.com/ChargePi/ChargePi-go/internal/users"
	"github.com/ChargePi/ChargePi-go/pkg/observability"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	settingsFilePath string
)

// runCommand is the command for the ChargePi core.
func runCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:     "run",
		Short:   "Run the ChargePi core",
		Long:    ``,
		Version: chargepoint.FirmwareVersion,
		PreRun: func(cmd *cobra.Command, args []string) {
			configuration.InitSettings(settingsFilePath)
		},
		Run: func(cmd *cobra.Command, args []string) {
			debug := viper.GetBool(configuration.Debug)
			runtimeSettings := configuration.GetRuntimeSettings()

			// Run the charge point
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)
			defer cancel()

			// Create a logger
			logger := log.StandardLogger()
			observability.SetupLogging(logger, runtimeSettings.Logging, debug)

			// Create a database for EVSE settings and state, tags, users, sessions and settings
			db, err := badger.NewBadgerDb(*databasePath)
			if err != nil {
				logger.WithError(err).Fatal("Cannot create database")
			}

			// Get the persistent settings
			persistentSettings, err := db.GetSettings()
			if err != nil {
				logger.WithError(err).Fatal("Cannot read persistent settings")
			}

			var (
				handler            chargepoint.ChargePoint
				hardware           = persistentSettings.Hardware
				connectionSettings = persistentSettings.ConnectionSettings
				chargePointInfo    = persistentSettings.Info
				protocolVersion    = connectionSettings.ProtocolVersion
			)

			cfg, err := ocpp_v16.DefaultConfigurationFromProfiles(supportedOcppV16Profiles...)
			if err != nil {
				logger.WithError(err).Fatal("Cannot create OCPP configuration")
			}

			ocppVariableManager, err := ocpp_v16.NewV16ConfigurationManager(*cfg, supportedOcppV16Profiles...)
			if err != nil {
				logger.WithError(err).Fatal("Cannot create OCPP variable manager")
			}

			settingsManager, err := configManager.NewManager(db, db, ocppVariableManager)
			if err != nil {
				logger.WithError(err).Fatal("Cannot create settings manager")
			}

			evseManager, err := manager.NewManager(db, db, make(chan notifications.StatusNotification, 100))
			if err != nil {
				logger.WithError(err).Fatal("Cannot create EVSE manager")
			}

			diagnosticsManager, err := diagnostics.NewService()
			if err != nil {
				logger.WithError(err).Fatal("Cannot create diagnostics service")
			}

			tagManager := auth.NewManager(db, db)
			sessionManager, err := sessions.NewSessionService(db)
			if err != nil {
				logger.WithError(err).Fatal("Cannot create session service")
			}

			// User service
			userService := users.NewUserService(db)

			// Initialize all the EVSEs.
			err = evseManager.InitAll(ctx)
			if err != nil {
				logger.WithError(err).Fatal("Cannot init EVSEs")
			}

			// Setup GRPC API if enabled
			if runtimeSettings.GRPC.Enabled {
				server, err := grpc.NewServer(runtimeSettings.GRPC, handler, tagManager, evseManager, settingsManager, userService)
				if err != nil {
					logger.WithError(err).Fatal("Cannot create the API server")
				}
				defer server.Stop()

				go func() {
					err = server.Run()
					if err != nil {
						logger.WithError(err).Fatal("Cannot start the API server")
					}
				}()
			}

			// Start the HTTP server, for health checks and UI
			httpServer := http.NewServer(runtimeSettings.HTTP)
			httpServer.Serve(handler)

			// Create a context for the OCPP connection, so it can be dynamically reconnected.
			parentCtxForOcpp, parentCancel := context.WithCancel(ctx)
			defer parentCancel()

			// Create a charge point handler
			handler, err = NewChargePoint(
				parentCtxForOcpp,
				protocolVersion,
				logger,
				evseManager,
				tagManager,
				settingsManager,
				sessionManager,
				diagnosticsManager,
				hardware,
			)
			if err != nil {
				logger.WithError(err).Fatal("Unable to create charge point")
			}

			err = handler.SetSettings(chargePointInfo)
			if err != nil {
				logger.WithError(err).Fatal("Unable to set the charge point settings")
			}

			err = handler.SetConnectionSettings(connectionSettings)
			if err != nil {
				logger.WithError(err).Fatal("Unable to set the connection settings")
			}

			// Listen for connector status changes
			go handler.ListenForConnectorStatusChange(ctx)

			serverUrl, err := chargepoint.CreateConnectionUrl(connectionSettings)
			if err != nil {
				logger.WithError(err).Fatal("Cannot create connection URL")
			}

			// Connect the charge point to the backend
			err = handler.Connect(parentCtxForOcpp, serverUrl)
			if err != nil {
				logger.WithError(err).Fatal("Cannot connect to the central system")
			}

			// Wait for the termination signal
			<-ctx.Done()
			logger.Info("Shutting down ChargePi...")

			// Gracefully stop the charge point and the server
			_ = handler.Cleanup(core.ReasonLocal)
			_ = evseManager.Shutdown()
			db.Close()
			httpServer.Stop()
		},
	}

	runCmd.Flags().StringVar(&settingsFilePath, configuration.SettingsFlag, "", "mainSettings file path")
	runCmd.Flags().String(configuration.ApiAddressFlag, "localhost:4269", "listen address")
	databasePath = runCmd.Flags().String("database", configuration.DatabasePath, "database path")
	_ = viper.BindPFlag(configuration.ApiAddress, runCmd.Flags().Lookup(configuration.ApiAddressFlag))

	return runCmd
}
