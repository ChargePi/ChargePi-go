package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/settings"
	"os"
)

const (
	debugFlag          = "debug"
	apiFlag            = "api"
	apiPortFlag        = "api-port"
	apiAddressFlag     = "api-address"
	settingsFlag       = "settings"
	connectorsFlag     = "connector-folder"
	authFileFlag       = "auth"
	ocppConfigPathFlag = "ocpp-config"
)

var (
	// Basic configuration settings
	configurationFilePath string
	connectorsFolderPath  string
	settingsFilePath      string
	authFilePath          string

	rootCmd = &cobra.Command{
		Use:   "chargepi",
		Short: "ChargePi is an open-source OCPP client.",
		Long:  ``,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	settings.InitSettings(settingsFilePath)

	var (
		isDebug      = viper.GetBool(settings.Debug)
		mainSettings = settings.GetSettings()
		connectors   = settings.GetConnectors(connectorsFolderPath)
	)

	chargepoint.Run(isDebug, mainSettings, connectors, configurationFilePath, authFilePath)
}

func setupFlags() {
	var (
		workingDirectory, _   = os.Getwd()
		connectorsFolderName  = fmt.Sprintf("%s/configs/connectors", workingDirectory)
		defaultConfigFileName = fmt.Sprintf("%s/configs/configuration.%s", workingDirectory, "json")
	)

	// Set flags
	rootCmd.PersistentFlags().StringVar(&settingsFilePath, settingsFlag, "", "config file path")
	rootCmd.PersistentFlags().StringVar(&connectorsFolderPath, connectorsFlag, connectorsFolderName, "connector folder path")
	rootCmd.PersistentFlags().StringVar(&configurationFilePath, ocppConfigPathFlag, defaultConfigFileName, "OCPP config file path")
	rootCmd.PersistentFlags().StringVar(&authFilePath, authFileFlag, "", "authorization file path")
	rootCmd.PersistentFlags().BoolP(debugFlag, "d", false, "debug mode")

	// Api flags
	rootCmd.PersistentFlags().BoolP(apiFlag, "a", false, "expose API")
	rootCmd.PersistentFlags().String(apiAddressFlag, "localhost", "address of the api")
	rootCmd.PersistentFlags().Int(apiPortFlag, 4269, "port for the API")

	// Bind flags to viper
	_ = viper.BindPFlag(settings.Debug, rootCmd.PersistentFlags().Lookup(debugFlag))
	_ = viper.BindPFlag(settings.ApiEnabled, rootCmd.PersistentFlags().Lookup(apiFlag))
	_ = viper.BindPFlag(settings.ApiAddress, rootCmd.PersistentFlags().Lookup(apiAddressFlag))
	_ = viper.BindPFlag(settings.ApiPort, rootCmd.PersistentFlags().Lookup(apiPortFlag))
}

func main() {
	setupFlags()
	err := rootCmd.Execute()
	if err != nil {
		log.WithError(err).Fatal("Unable to run")
	}
}
