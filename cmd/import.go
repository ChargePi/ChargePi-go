package cmd

import (
	"fmt"

	"github.com/ChargePi/ChargePi-go/internal/chargepoint"
	"github.com/ChargePi/ChargePi-go/internal/pkg/badger"
	cfg "github.com/ChargePi/ChargePi-go/internal/pkg/configuration"
	importer2 "github.com/ChargePi/ChargePi-go/internal/pkg/configuration/importer"
	"github.com/ChargePi/ChargePi-go/internal/pkg/configuration/manager"
	"github.com/ChargePi/ChargePi-go/pkg/ocpp"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/spf13/cobra"
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
		Use:     "import",
		Short:   "Import configurations to ChargePi.",
		Long:    ``,
		Version: chargepoint.FirmwareVersion,
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := badger.NewBadgerDb(*databasePath)
			if err != nil {
				return fmt.Errorf("could not create database: %v", err)
			}

			configurationManager, err := ocpp_v16.NewV16ConfigurationManager(ocpp_v16.NewEmptyConfiguration())
			if err != nil {
				return fmt.Errorf("could not create OCPP configuration manager: %v", err)
			}

			settingsManager, err := manager.NewManager(db, db, configurationManager)
			if err != nil {
				return fmt.Errorf("could not create settings manager: %v", err)
			}

			importer := importer2.NewImporter(settingsManager, db, db)

			evseFlag := cmd.Flags().Lookup(cfg.EvseFlag).Changed
			ocppFlag := cmd.Flags().Lookup(cfg.OcppConfigPathFlag).Changed
			authFlag := cmd.Flags().Lookup(cfg.AuthFileFlag).Changed
			settingsFlag := cmd.Flags().Lookup(cfg.SettingsFlag).Changed

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

	evseFolderPath = importCmd.Flags().String(cfg.EvseFlag, "", "evse folder path")
	ocppConfigurationFilePath = importCmd.Flags().String(cfg.OcppConfigPathFlag, "", "OCPP config file path")
	ocppVersionFlag = importCmd.Flags().StringP(cfg.OcppVersion, "v", "1.6", "OCPP config file path")
	authFilePath = importCmd.Flags().String(cfg.AuthFileFlag, "", "authorization file path")
	importSettingsFilePath = importCmd.Flags().String(cfg.SettingsFlag, "", "settings file path")

	return importCmd
}
