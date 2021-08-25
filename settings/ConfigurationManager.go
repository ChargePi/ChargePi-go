package settings

import (
	"errors"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	goCache "github.com/patrickmn/go-cache"
	"github.com/xBlaz3kx/ChargePi-go/cache"
	"log"
	"sync"
)

type OCPPConfig struct {
	Version int
	Keys    []core.ConfigurationKey
}

// init Read the OCPP configuration from the configuration.json file.
func init() {
	once := sync.Once{}
	once.Do(func() {
		var ocppConfig OCPPConfig
		DecodeFile("configs/configuration.json", &ocppConfig)
		err := cache.Cache.Add("OCPPConfiguration", &ocppConfig, goCache.NoExpiration)
		if err != nil {
			panic(err)
		}
		log.Println("Added OCPP configuration to cache")
	})
}

// GetConfiguration Get the global configuration
func GetConfiguration() (*OCPPConfig, error) {
	configuration, isFound := cache.Cache.Get("OCPPConfiguration")
	if isFound {
		return configuration.(*OCPPConfig), nil
	}
	return nil, errors.New("configuration not found in cache")
}

// UpdateKey Update the configuration variable in the global configuration if it is not readonly.
func UpdateKey(key string, value string) (err error) {
	configuration, err := GetConfiguration()
	if err != nil {
		log.Println(err)
		return err
	}
	return configuration.UpdateKey(key, value)
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
		log.Println(err)
		return err
	}
	err = WriteToFile("configs/configuration.json", &configuration)
	if err != nil {
		log.Println(err)
	}
	return err
}

// UpdateKey Update the configuration variable in the configuration if it is not readonly.
func (config *OCPPConfig) UpdateKey(key string, value string) error {
	for i, configKey := range config.Keys {
		if configKey.Key == key {
			if !configKey.Readonly {
				config.Keys[i].Value = value
				return nil
			}
			return errors.New("attribute is read-only")
		}
	}
	return errors.New("key not found")
}

//GetConfigurationValue Get the value of specified configuration variable in String format.
func (config *OCPPConfig) GetConfigurationValue(key string) (string, error) {
	for _, configKey := range config.Keys {
		if configKey.Key == key {
			return configKey.Value, nil
		}
	}
	return "", errors.New("key not found")
}

// GetConfig Get the configuration
func (config *OCPPConfig) GetConfig() []core.ConfigurationKey {
	return config.Keys
}
