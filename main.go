package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint"
	"os"
)

const (
	debugFlag = "debug"
	apiFlag   = "api"
)

var (
	// Basic configuration settings
	configurationFilePath string
	connectorsFolderPath  string
	settingsFilePath      string
	authFilePath          string
	isDebug               = false

	// Settings for exposing the API
	exposeApi  = false
	apiAddress string
	apiPort    int

	rootCmd = &cobra.Command{
		Use:   "chargepi",
		Short: "ChargePi is an open-source OCPP client.",
		Long:  ``,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	chargepoint.Run(isDebug, settingsFilePath, configurationFilePath, connectorsFolderPath, authFilePath, exposeApi, apiAddress, apiPort)
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
	// Api flags
	rootCmd.PersistentFlags().BoolVarP(&exposeApi, apiFlag, "a", false, "expose API")
	rootCmd.PersistentFlags().StringVar(&apiAddress, "api-address", "localhost", "address of the api")
	rootCmd.PersistentFlags().IntVar(&apiPort, "port", 4269, "port for the API")

	err := rootCmd.Execute()
	if err != nil {
		log.WithError(err).Fatal("Unable to run")
	}
}
