package manager

import (
	"github.com/ChargePi/ChargePi-go/internal/chargepoint"
	"github.com/ChargePi/ChargePi-go/internal/evse/manager"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/go-playground/validator/v10"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
)

type (
	Manager interface {
		ocpp_v16.Manager
		GetChargePointSettings() (*chargepoint.Settings, error)
		SetChargePointSettings(settings chargepoint.Settings) error
	}

	V1 struct {
		evseSettingsRepository manager.EvseSettingsRepository
		settingsRepository     SettingsRepository
		ocpp16VariableManager  ocpp_v16.Manager
		logger                 log.FieldLogger
	}
)

func NewManager(
	evseSettingsRepository manager.EvseSettingsRepository,
	settingsRepository SettingsRepository,
	ocpp16VariableManager ocpp_v16.Manager,
) (*V1, error) {
	return &V1{
		evseSettingsRepository: evseSettingsRepository,
		settingsRepository:     settingsRepository,
		ocpp16VariableManager:  ocpp16VariableManager,
		logger:                 log.WithField("component", "settings-manager"),
	}, nil
}

func (i *V1) GetChargePointSettings() (*chargepoint.Settings, error) {
	i.logger.Debug("Getting charge point settings")

	settings, err := i.settingsRepository.GetSettings()
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (i *V1) SetChargePointSettings(settings chargepoint.Settings) error {
	i.logger.Debug("Setting charge point settings")

	// Validate the settings
	validationErr := validator.New().Struct(settings)
	if validationErr != nil {
		return validationErr
	}

	return i.settingsRepository.UpdateSettings(settings)
}

func (i *V1) SetMandatoryKeys(mandatoryKeys []ocpp_v16.Key) error {
	return i.ocpp16VariableManager.SetMandatoryKeys(mandatoryKeys)
}

func (i *V1) GetMandatoryKeys() []ocpp_v16.Key {
	return i.ocpp16VariableManager.GetMandatoryKeys()
}

func (i *V1) RegisterCustomKeyValidator(keyValidator ocpp_v16.KeyValidator) {
	i.ocpp16VariableManager.RegisterCustomKeyValidator(keyValidator)
}

func (i *V1) ValidateKey(key ocpp_v16.Key, value *string) error {
	return i.ocpp16VariableManager.ValidateKey(key, value)
}

func (i *V1) UpdateKey(key ocpp_v16.Key, value *string) error {
	// todo update the configuration in the database as well
	return i.ocpp16VariableManager.UpdateKey(key, value)
}

func (i *V1) OnUpdateKey(key ocpp_v16.Key, handler ocpp_v16.OnUpdateHandler) error {
	return i.ocpp16VariableManager.OnUpdateKey(key, handler)
}

func (i *V1) GetConfigurationValue(key ocpp_v16.Key) (*string, error) {
	return i.ocpp16VariableManager.GetConfigurationValue(key)
}

func (i *V1) SetConfiguration(configuration ocpp_v16.Config) error {
	// todo update the configuration in the database as well
	return i.ocpp16VariableManager.SetConfiguration(configuration)
}

func (i *V1) GetConfiguration() ([]core.ConfigurationKey, error) {
	return i.ocpp16VariableManager.GetConfiguration()
}
