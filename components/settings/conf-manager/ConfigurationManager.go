package conf_manager

import (
	"errors"
	"github.com/kkyr/fig"
	types2 "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	goCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/components/cache"
	s "github.com/xBlaz3kx/ChargePi-go/components/settings"
	"github.com/xBlaz3kx/ChargePi-go/data/ocpp"
	"path/filepath"
	"strings"
)

var (
	ErrConfigNotFound = errors.New("configuration not found in cache")
)

// InitConfiguration load OCPP configuration from the persistence file.
func InitConfiguration() {
	var (
		ocppConfig            ocpp.Config
		configurationFilePath = ""
		err                   error
	)

	configurationPath, isFound := cache.Cache.Get("configurationFilePath")
	if isFound {
		configurationFilePath = configurationPath.(string)
	}

	err = fig.Load(&ocppConfig,
		fig.File(filepath.Base(configurationFilePath)),
		fig.Dirs(filepath.Dir(configurationFilePath)),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = cache.Cache.Add("OCPPConfiguration", &ocppConfig, goCache.NoExpiration)
	if err != nil {
		log.Println(err)
		return
	}

	log.Info("Added OCPP configuration to cache")
}

// GetConfiguration Get the global configuration
func GetConfiguration() (*ocpp.Config, error) {
	configuration, isFound := cache.Cache.Get("OCPPConfiguration")
	if isFound {
		return configuration.(*ocpp.Config), nil
	}

	return nil, ErrConfigNotFound
}

// UpdateKey Update the configuration variable in the global configuration if it is not readonly.
func UpdateKey(key string, value string) error {
	configuration, err := GetConfiguration()
	if err != nil {
		return err
	}

	err = configuration.UpdateKey(key, value)
	if err != nil {
		return err
	}

	return UpdateConfigurationFile()
}

// GetConfigurationValue Get the value of specified configuration variable from the global configuration in String format.
func GetConfigurationValue(key string) (string, error) {
	configuration, err := GetConfiguration()
	if err != nil {
		return "", err
	}

	return configuration.GetConfigurationValue(key)
}

// UpdateConfigurationFile Write/Rewrite the existing global configuration to the file.
func UpdateConfigurationFile() error {
	configuration, err := GetConfiguration()
	if err != nil {
		return err
	}

	value, isFound := cache.Cache.Get("configurationFilePath")
	if !isFound {
		return ErrConfigNotFound
	}

	return s.WriteToFile(value.(string), &configuration)
}

func GetTypesToSample() []types2.Measurand {
	var (
		measurands []types2.Measurand
		// Get the types to sample
		measurandsString, err = GetConfigurationValue("MeterValuesSampledData")
	)

	if err != nil {
		measurandsString = string(types2.MeasurandPowerActiveExport)
	}

	for _, measurand := range strings.Split(measurandsString, ",") {
		measurands = append(measurands, types2.Measurand(measurand))
	}

	return measurands
}
