package cmd

import (
	"context"
	"github.com/ChargePi/ChargePi-go/internal/pkg/configuration"
	"github.com/ChargePi/ChargePi-go/pkg/observability"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

var rootCmd = &cobra.Command{
	Use:   "chargepi",
	Short: "ChargePi is an open-source Charge point project.",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	cobra.OnInitialize(func() {
		// observability.SetupLogging(log.StandardLogger(), observability.Logging{}, viper.GetBool(configuration.Debug))
	})

	rootCmd.AddCommand(runCommand())
	rootCmd.AddCommand(versionCommand())
	rootCmd.AddCommand(exportCommand())
	rootCmd.AddCommand(importCommand())

	rootCmd.PersistentFlags().BoolP(configuration.DebugFlag, "d", false, "debug mode")
	_ = viper.BindPFlag(configuration.Debug, rootCmd.PersistentFlags().Lookup(configuration.DebugFlag))
}

func Execute() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)
	defer cancel()

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		zap.L().Fatal("Unable to run", zap.Error(err))
	}
}
