package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	chargePoint "github.com/ChargePi/ChargePi-go/internal/chargepoint"
)

// versionCmd represents the version command
func versionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Version of ChargePi",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			zap.L().Info("ChargePi version", zap.String("version", chargePoint.FirmwareVersion))
		},
	}
}
