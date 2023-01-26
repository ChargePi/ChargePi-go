package cmd

import (
	"github.com/spf13/cobra"
	connectorManager "github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export component settings of the ChargePi.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = connectorManager.GetManager()
		_ = ocppConfigManager.GetManager()
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.
	exportCmd.PersistentFlags().StringVar(evseFolderPath, settings.EvseFlag, "./configs/evses", "evse folder path")
	exportCmd.PersistentFlags().StringVar(configurationFilePath, settings.OcppConfigPathFlag, "./configs/settings.yaml", "OCPP config file path")
	exportCmd.PersistentFlags().StringVar(authFilePath, settings.AuthFileFlag, "./configs/authorization.yaml", "authorization file path")
}
