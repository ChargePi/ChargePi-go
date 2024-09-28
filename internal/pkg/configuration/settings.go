package configuration

import (
	"github.com/ChargePi/ChargePi-go/internal/api/grpc"
	"github.com/ChargePi/ChargePi-go/internal/api/http"
	"github.com/ChargePi/ChargePi-go/pkg/observability"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Non-persistent settings, used to configure the runtime of the charge point
type RuntimeSettings struct {
	// GRPC settings for API
	GRPC grpc.Configuration `json:"grpc" yaml:"grpc" mapstructure:"grpc"`

	// HTTP settings for health checks and UI
	HTTP http.Configuration `json:"http" yaml:"http" mapstructure:"http"`

	// Logging settings
	Logging observability.Logging `json:"logging" yaml:"logging" mapstructure:"logging"`
}

// GetRuntimeSettings gets runtime settings, such as API and UI settings.
func GetRuntimeSettings() *RuntimeSettings {
	log.Info("Fetching runtime settings..")

	var conf RuntimeSettings

	err := viper.Unmarshal(&conf)
	if err != nil {
		log.WithError(err).Fatalf("Cannot unmarshal settings")
	}

	validationErr := validator.New().Struct(conf)
	if validationErr != nil {
		log.WithError(validationErr).Fatalf("Invalid settings")
	}

	return &conf
}
