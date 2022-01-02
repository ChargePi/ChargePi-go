package main

import (
	"context"
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/v16"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/auth"
	connectorManager "github.com/xBlaz3kx/ChargePi-go/internal/components/connector-manager"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/hardware/reader"
	s "github.com/xBlaz3kx/ChargePi-go/internal/components/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/cache"
	"github.com/xBlaz3kx/ChargePi-go/pkg/logging"
	"github.com/xBlaz3kx/ChargePi-go/pkg/scheduler"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	debugFlag = "debug"
)

var (
	configurationFilePath string
	connectorsFolderPath  string
	settingsFilePath      string
	authFilePath          string

	isDebug = false

	rootCmd = &cobra.Command{
		Use:   "chargepi",
		Short: "ChargePi is an open-source OCPP client.",
		Long:  ``,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	var (
		handler     chargepoint.ChargePoint
		config      *settings.Settings
		tagReader   reader.Reader
		lcd         display.LCD
		mem         = cache.GetCache()
		authCache   = auth.NewAuthCache(authFilePath)
		err         error
		quitChannel = make(chan os.Signal, 1)
		ctx, cancel = context.WithCancel(context.Background())
	)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	// Read settings file and cache it
	config = s.GetSettings(mem, settingsFilePath)
	go authCache.LoadAuthFile(authFilePath)

	var (
		chargePointInfo = config.ChargePoint.Info
		hardware        = config.ChargePoint.Hardware
		serverUrl       = fmt.Sprintf("ws://%s", chargePointInfo.ServerUri)
		protocolVersion = settings.ProtocolVersion(chargePointInfo.ProtocolVersion)
	)

	// Create the logger
	logging.Setup(config.ChargePoint.Logging, isDebug)

	// Create hardware components based on settings
	tagReader, err = reader.NewTagReader(hardware.TagReader)
	if err == nil {
		go tagReader.ListenForTags(ctx)
	}

	lcd, err = display.NewDisplay(hardware.Lcd)
	if err == nil {
		go lcd.ListenForMessages(ctx)
	}

	if config.ChargePoint.TLS.IsEnabled {
		// Replace insecure Websockets
		serverUrl = strings.Replace(serverUrl, "ws", "wss", 1)
	}

	// Setup OCPP configuration manager
	s.SetupOcppConfigurationManager(
		configurationFilePath,
		configuration.ProtocolVersion(config.ChargePoint.Info.ProtocolVersion),
		core.ProfileName,
		reservation.ProfileName)

	switch protocolVersion {
	case settings.OCPP16:
		// Create the client
		handler = v16.NewChargePoint(tagReader, lcd, connectorManager.GetManager(), scheduler.GetScheduler(), authCache)
		break
	case settings.OCPP201:
		log.Fatal("Version 2.0.1 is not supported yet.")
	default:
		log.WithField("protocolVersion", protocolVersion).Fatal("Protocol version not supported")
	}

	// Initialize and connect to the Central System
	handler.Init(ctx, config)

	// Add connectors
	connectors := s.GetConnectors(mem, connectorsFolderPath)
	handler.AddConnectors(connectors)

	// Finally, connect to the central system
	handler.Connect(ctx, serverUrl)

Loop:
	for {
		select {
		// Capture the terminate signal
		case <-quitChannel:
			handler.CleanUp(core.ReasonLocal)
			cancel()
			break Loop
		case <-ctx.Done():
			handler.CleanUp(core.ReasonPowerLoss)
			cancel()
			break Loop
		}
	}
}

func main() {
	var (
		workingDirectory, _     = os.Getwd()
		defaultConfigFileName   = fmt.Sprintf("%s/configs/configuration.%s", workingDirectory, "json")
		defaultSettingsFileName = fmt.Sprintf("%s/configs/settings.%s", workingDirectory, "json")
		connectorsFolderName    = fmt.Sprintf("%s/configs/connectors", workingDirectory)
		defaultAuthFileName     = fmt.Sprintf("%s/configs/auth.%s", workingDirectory, "json")
	)

	// Set flags
	rootCmd.PersistentFlags().StringVar(&settingsFilePath, "settings", defaultSettingsFileName, "config file path")
	rootCmd.PersistentFlags().StringVar(&connectorsFolderPath, "connector-folder", connectorsFolderName, "connector folder path")
	rootCmd.PersistentFlags().StringVar(&configurationFilePath, "ocpp-config", defaultConfigFileName, "OCPP config file path")
	rootCmd.PersistentFlags().StringVar(&authFilePath, "auth", defaultAuthFileName, "authorization file path")
	rootCmd.PersistentFlags().BoolVarP(&isDebug, debugFlag, "d", false, "debug mode")

	err := rootCmd.Execute()
	if err != nil {
		log.WithError(err).Fatal("Unable to run")
	}
}
