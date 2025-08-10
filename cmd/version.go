package cmd

import (
	chargePoint "github.com/ChargePi/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
