package v16

import (
	"context"
	"errors"
	"strconv"

	"github.com/ChargePi/ChargePi-go/internal/chargepoint"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/indicator"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/go-playground/validator/v10"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)

func (cp *ChargePoint) OnChangeConfiguration(request *core.ChangeConfigurationRequest) (confirmation *core.ChangeConfigurationConfirmation, err error) {
	cp.logger.Infof("Received request %s", request.GetFeatureName())
	var response = core.ConfigurationStatusRejected

	// Update the desired key. Validation is done by the settings manager and is configured in the setupCustomConfigurationValidation function.
	// After a key is updated, a handler function is called. All the cases are covered in the settings.go file.
	err = cp.settingsManager.UpdateKey(ocpp_v16.Key(request.Key), &request.Value)
	if err == nil {
		response = core.ConfigurationStatusAccepted
	}

	return core.NewChangeConfigurationConfirmation(response), nil
}

func (cp *ChargePoint) OnGetConfiguration(request *core.GetConfigurationRequest) (confirmation *core.GetConfigurationConfirmation, err error) {
	cp.logger.Infof("Received request %s", request.GetFeatureName())

	var (
		unknownKeys []string
		configArray = []core.ConfigurationKey{}
		response    = core.NewGetConfigurationConfirmation(configArray)
	)

	configuration, confErr := cp.settingsManager.GetConfiguration()
	if confErr != nil || configuration == nil {
		return response, nil
	}

	configArray = configuration

	// Get all configuration variables
	if request.Key == nil || len(request.Key) == 0 {
		response.ConfigurationKey = configArray
		response.UnknownKey = unknownKeys
		return response, nil
	}

	configArray2 := []core.ConfigurationKey{}

	// Get only the requested variables
	for _, key := range request.Key {

		// Note: redundant looping, should've just created an ocpp.ConfigurationKey function
		// Check if the key exists
		_, keyErr := cp.settingsManager.GetConfigurationValue(ocpp_v16.Key(key))
		if keyErr != nil {
			unknownKeys = append(unknownKeys, key)
			continue
		}

		// Key should exist, therefore find it in the config
		for _, configurationKey := range configArray {
			if key == configurationKey.Key {
				configArray2 = append(configArray2, configurationKey)
			}
		}
	}

	response.ConfigurationKey = configArray2
	response.UnknownKey = unknownKeys
	return response, nil
}

// setupLocalAuthListConfigurationValidation sets up the configuration validation for local auth list.
func (cp *ChargePoint) setupLocalAuthListConfigurationValidation() error {
	err := cp.settingsManager.OnUpdateKey(ocpp_v16.LocalAuthListMaxLength, func(value *string) error {
		if value == nil {
			return errors.New("value is nil")
		}

		val, err := strconv.Atoi(*value)
		if err != nil {
			return errors.New("invalid value for LocalAuthListMaxLength")
		}

		cp.tagAuthService.SetMaxTags(val)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// setupSmartChargingConfigurationValidation sets up the configuration validation for smart charging.
func (cp *ChargePoint) setupSmartChargingConfigurationValidation() error {
	err := cp.settingsManager.OnUpdateKey(ocpp_v16.MaxChargingProfilesInstalled, func(value *string) error {
		return nil
	})
	if err != nil {
		return err
	}

	err = cp.settingsManager.OnUpdateKey(ocpp_v16.ChargeProfileMaxStackLevel, func(value *string) error {
		return nil
	})
	if err != nil {
		return err
	}

	err = cp.settingsManager.OnUpdateKey(ocpp_v16.ChargingScheduleAllowedChargingRateUnit, func(value *string) error {
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// setupCoreConfigurationValidation sets up the configuration validation for custom variables.
func (cp *ChargePoint) setupCustomConfigurationValidation() error {
	err := cp.settingsManager.OnUpdateKey("", func(value *string) error {
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// setupCoreConfigurationValidation sets up validation for the core configuration.
func (cp *ChargePoint) setupCoreConfigurationValidation() error {
	err := cp.settingsManager.OnUpdateKey(ocpp_v16.AuthorizationCacheEnabled, func(value *string) error {
		if value == nil {
			return errors.New("value is nil")
		}

		cp.tagAuthService.ToggleAuthCache(*value == "true")
		return nil
	})
	if err != nil {
		return err
	}

	err = cp.settingsManager.OnUpdateKey(ocpp_v16.HeartbeatInterval, func(value *string) error {
		if value == nil {
			return errors.New("value is nil")
		}

		val, err := strconv.Atoi(*value)
		if err != nil {
			return errors.New("invalid value for LocalAuthListMaxLength")
		}

		cp.setHeartbeat(val)
		return nil
	})
	if err != nil {
		return err
	}

	err = cp.settingsManager.OnUpdateKey(ocpp_v16.LightIntensity, func(value *string) error {
		if value == nil {
			return errors.New("value is nil")
		}

		brightness, err := strconv.Atoi(*value)
		if err != nil {
			return errors.New("invalid value for LightIntensity")
		}

		return cp.indicator.SetBrightness(brightness)
	})
	if err != nil {
		return err
	}

	return nil
}

func (cp *ChargePoint) SetSettings(settings chargepoint.Info) error {
	err := validator.New().StructCtx(context.Background(), settings)
	if err != nil {
		return err
	}

	cp.logger.Debug("Setting charge point settings")
	cp.info = settings
	return nil
}

func (cp *ChargePoint) GetConnectionSettings() chargepoint.ConnectionSettings {
	return cp.connectionSettings
}

func (cp *ChargePoint) SetConnectionSettings(settings chargepoint.ConnectionSettings) error {
	err := validator.New().StructCtx(context.Background(), settings)
	if err != nil {
		return err
	}

	cp.logger.Debug("Setting backend connection settings")
	cp.connectionSettings = settings
	return nil
}

func (cp *ChargePoint) GetSettings() chargepoint.Info {
	return cp.info
}

func (cp *ChargePoint) GetIndicatorSettings() indicator.StatusMapping {
	return cp.indicatorMapping
}

func (cp *ChargePoint) SetIndicatorSettings(settings indicator.StatusMapping) error {
	err := validator.New().StructCtx(context.Background(), settings)
	if err != nil {
		return err
	}

	cp.logger.Debug("Setting indicator settings")
	cp.indicatorMapping = settings
	return nil
}
