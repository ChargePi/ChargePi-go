package configuration

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/ChargePi/ChargePi-go/internal/api/grpc"
	"github.com/ChargePi/ChargePi-go/internal/api/http"
)

// Non-persistent settings, used to configure the runtime of the charge point
type RuntimeSettings struct {
	// GRPC settings for API
	GRPC grpc.Configuration `json:"grpc" yaml:"grpc" mapstructure:"grpc"`

	// HTTP settings for health checks and UI
	HTTP http.Configuration `json:"http" yaml:"http" mapstructure:"http"`
}

// GetRuntimeSettings gets runtime settings, such as API and UI settings.
func GetRuntimeSettings() *RuntimeSettings {
	logger := zap.L()
	logger.Info("Fetching runtime settings..")

	var conf RuntimeSettings

	err := viper.Unmarshal(&conf)
	if err != nil {
		logger.With(zap.Error(err)).Fatal("Cannot unmarshal settings")
	}

	validationErr := validator.New().Struct(conf)
	if validationErr != nil {
		logger.With(zap.Error(err)).Fatal("Invalid settings")
	}

	return &conf
}
