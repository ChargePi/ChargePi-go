package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
	"os"
)

func init() {
	cobra.OnInitialize(setupFlags)
}

func setupFlags() {
	var (
		workingDirectory, _   = os.Getwd()
		evseFolderName        = fmt.Sprintf("%s/configs/evses", workingDirectory)
		defaultConfigFileName = fmt.Sprintf("%s/configs/configuration.%s", workingDirectory, "json")
	)

	// Set flags
	rootCmd.PersistentFlags().StringVar(&settingsFilePath, settings.SettingsFlag, "", "config file path")
	rootCmd.PersistentFlags().StringVar(&evseFolderPath, settings.EvseFlag, evseFolderName, "evse folder path")
	rootCmd.PersistentFlags().StringVar(&configurationFilePath, settings.OcppConfigPathFlag, defaultConfigFileName, "OCPP config file path")
	rootCmd.PersistentFlags().StringVar(&authFilePath, settings.AuthFileFlag, "./configs/auth.json", "authorization file path")
	rootCmd.PersistentFlags().BoolP(settings.DebugFlag, "d", false, "debug mode")

	// Bind flags to viper
	_ = viper.BindPFlag(settings.Debug, rootCmd.PersistentFlags().Lookup(settings.DebugFlag))
	setupApiCfg()
}

func setupApiCfg() {
	rootCmd.PersistentFlags().String(settings.ApiAddressFlag, "localhost", "address of the api")
	rootCmd.PersistentFlags().Int(settings.ApiPortFlag, 4269, "port for the API")

	_ = viper.BindPFlag(settings.ApiAddress, rootCmd.PersistentFlags().Lookup(settings.ApiAddressFlag))
	_ = viper.BindPFlag(settings.ApiPort, rootCmd.PersistentFlags().Lookup(settings.ApiPortFlag))
}
