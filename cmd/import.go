package cmd

import (
	"fmt"

	"github.com/agrison/go-commons-lang/stringUtils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	connectorManager "github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

var (
	// importCmd represents the import command
	importCmd = &cobra.Command{
		Use:   "import",
		Short: "Import component configuration to ChargePi.",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			_ = connectorManager.GetManager()
			_ = ocppConfigManager.GetManager()
		},
	}

	evseFolderPath        *string
	configurationFilePath *string
	authFilePath          *string
)

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVar(evseFolderPath, settings.EvseFlag, "./configs/evses", "evse folder path")
	importCmd.Flags().StringVar(configurationFilePath, settings.OcppConfigPathFlag, "./configs/settings.yaml", "OCPP config file path")
	importCmd.Flags().StringVar(authFilePath, settings.AuthFileFlag, "./configs/authorization.yaml", "authorization file path")
}

func loadOcppConfigurationFromFile(filePath string, version configuration.ProtocolVersion, supportedProfiles ...string) {
	ocppConfigManager.SetFilePath(filePath)
	ocppConfigManager.SetVersion(version)
	ocppConfigManager.SetSupportedProfiles(supportedProfiles...)

	// Load the configuration
	err := ocppConfigManager.LoadConfiguration()
	if err != nil {
		log.WithError(err).Fatalf("Cannot load OCPP configuration")
	}
}

// loadEVSEFromPath loads a connector from file
func loadEVSEFromPath(name, path string) (*settings.EVSE, error) {
	// Read the evse settings from the file in the directory
	var (
		cfg  = viper.New()
		evse settings.EVSE
	)

	readConfiguration(cfg, name, "yaml", path)

	err := cfg.Unmarshal(&evse)
	if err != nil {
		log.WithError(err).Errorf("Cannot read evse file")
		return nil, err
	}

	log.Debugf("Loaded evse from %s", path)
	_ = fmt.Sprintf("evse%d", evse.EvseId)
	// todo store in the database

	return &evse, nil
}

func readConfiguration(viper *viper.Viper, fileName, extension, filePath string) {
	viper.SetConfigName(fileName)
	viper.SetConfigType(extension)
	viper.AddConfigPath(settings.CurrentFolder)
	viper.AddConfigPath(settings.EvseFolder)
	viper.AddConfigPath(settings.DockerFolder)

	if stringUtils.IsNotEmpty(filePath) {
		viper.SetConfigFile(filePath)
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Fatalf("Cannot parse config file")
	}

	log.Debugf("Using configuration file: %s", viper.ConfigFileUsed())
}
