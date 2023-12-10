package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	cfg "github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/ocpp"
)

var (
	evseFolderPath            *string
	ocppConfigurationFilePath *string
	ocppVersionFlag           *string
	authFilePath              *string
	importSettingsFilePath    *string
)

// importCmd represents the import command
func importCommand() *cobra.Command {
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "Import configurations to ChargePi.",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			importer := cfg.GetImporter()

			evseFlag := cmd.Flags().Lookup(settings.EvseFlag).Changed
			ocppFlag := cmd.Flags().Lookup(settings.OcppConfigPathFlag).Changed
			authFlag := cmd.Flags().Lookup(settings.AuthFileFlag).Changed
			settingsFlag := cmd.Flags().Lookup(settings.SettingsFlag).Changed

			if evseFlag {
				// If a directory is specified, (try to) import all the files in that directory.
				err := importer.ImportEVSESettingsFromPath(*evseFolderPath)
				if err != nil {
					return fmt.Errorf("could not import EVSE settings: %v", err)
				}
			}

			// If the flag was set, import OCPP configuration to the ChargePi
			if ocppFlag {
				err := importer.ImportOcppConfigurationFromPath(ocpp.ProtocolVersion(*ocppVersionFlag), *ocppConfigurationFilePath)
				if err != nil {
					return fmt.Errorf("could not import OCPP configuration: %v", err)
				}
			}

			// If the flag was set, import tags to the database.
			if authFlag {
				err := importer.ImportLocalAuthListFromPath(*authFilePath)
				if err != nil {
					return fmt.Errorf("could not import tags: %v", err)
				}
			}

			if settingsFlag {
				err := importer.ImportChargePointSettingsFromPath(*importSettingsFilePath)
				if err != nil {
					return fmt.Errorf("could not import settings: %v", err)
				}
			}

			return nil
		},
	}

	evseFolderPath = importCmd.Flags().String(settings.EvseFlag, "", "evse folder path")
	ocppConfigurationFilePath = importCmd.Flags().String(settings.OcppConfigPathFlag, "", "OCPP config file path")
	ocppVersionFlag = importCmd.Flags().StringP(settings.OcppVersion, "v", "1.6", "OCPP config file path")
	authFilePath = importCmd.Flags().String(settings.AuthFileFlag, "", "authorization file path")
	importSettingsFilePath = importCmd.Flags().String(settings.SettingsFlag, "", "settings file path")

	return importCmd
}
