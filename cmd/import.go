package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	cfg "github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
)

var (
	evseFolderPath            *string
	ocppConfigurationFilePath *string
	authFilePath              *string

	// importCmd represents the import command
	importCmd = &cobra.Command{
		Use:   "import",
		Short: "Import component configuration to ChargePi.",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			manager := cfg.GetImporter()
			conf := viper.New()

			evseFlag := cmd.Flags().Lookup(settings.EvseFlag).Changed
			ocppFlag := cmd.Flags().Lookup(settings.OcppConfigPathFlag).Changed
			authFlag := cmd.Flags().Lookup(settings.AuthFileFlag).Changed

			if evseFlag {
				// If a directory is specified, (try to) import all the files in that directory.
				err := cfg.ImportEVSEs(manager, conf, *evseFolderPath)
				if err != nil {
					log.WithError(err).Errorf("Could not import EVSE settings")
				}
			}

			// If the flag was set, import OCPP configuration to the ChargePi
			if ocppFlag {
				err := cfg.ImportOcppConfiguration(manager, conf, *ocppConfigurationFilePath)
				if err != nil {
					log.WithError(err).Errorf("Could not import OCPP settings")
				}
			}

			// If the flag was set, import tags to the database.
			if authFlag {
				err := cfg.ImportLocalAuthList(manager, conf, *authFilePath)
				if err != nil {
					log.WithError(err).Errorf("Could not import Local Auth List")
				}
			}

		},
	}
)

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVar(evseFolderPath, settings.EvseFlag, "", "evse folder path")
	importCmd.Flags().StringVar(ocppConfigurationFilePath, settings.OcppConfigPathFlag, "", "OCPP config file path")
	importCmd.Flags().StringVar(authFilePath, settings.AuthFileFlag, "", "authorization file path")
}
