package ocpp

import (
	"errors"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrReadOnly    = errors.New("attribute is read-only")
	mandatoryKeys  = []Key{
		AuthorizeRemoteTxRequests,
		ClockAlignedDataInterval,
		ConnectionTimeOut,
		ConnectorPhaseRotation,
		GetConfigurationMaxKeys,
		HeartbeatInterval,
		LocalAuthorizeOffline,
		LocalPreAuthorize,
		MeterValuesAlignedData,
		MeterValuesSampledData,
		MeterValueSampleInterval,
		NumberOfConnectors,
		ResetRetries,
		StopTransactionOnEVSideDisconnect,
		StopTransactionOnInvalidId,
		StopTxnAlignedData,
		StopTxnSampledData,
		SupportedFeatureProfiles,
		TransactionMessageAttempts,
		TransactionMessageRetryInterval,
		UnlockConnectorOnEVSideDisconnect,
		WebSocketPingInterval,
		LocalAuthListEnabled,
		LocalAuthListMaxLength,
		SendLocalListMaxLength,
		ChargeProfileMaxStackLevel,
		ChargingScheduleAllowedChargingRateUnit,
		ChargingScheduleMaxPeriods,
		MaxChargingProfilesInstalled,
	}
)

const (
	keyNotFound                             = -1
	AuthorizeRemoteTxRequests               = Key("AuthorizeRemoteTxRequests")
	ClockAlignedDataInterval                = Key("ClockAlignedDataInterval")
	ConnectionTimeOut                       = Key("ConnectionTimeOut")
	ConnectorPhaseRotation                  = Key("ConnectorPhaseRotation")
	GetConfigurationMaxKeys                 = Key("GetConfigurationMaxKeys")
	HeartbeatInterval                       = Key("HeartbeatInterval")
	LocalAuthorizeOffline                   = Key("LocalAuthorizeOffline")
	LocalPreAuthorize                       = Key("LocalPreAuthorize")
	MeterValuesAlignedData                  = Key("MeterValuesAlignedData")
	MeterValuesSampledData                  = Key("MeterValuesSampledData")
	MeterValueSampleInterval                = Key("MeterValueSampleInterval")
	NumberOfConnectors                      = Key("NumberOfConnectors")
	ResetRetries                            = Key("ResetRetries")
	StopTransactionOnEVSideDisconnect       = Key("StopTransactionOnEVSideDisconnect")
	StopTransactionOnInvalidId              = Key("StopTransactionOnInvalidId")
	StopTxnAlignedData                      = Key("StopTxnAlignedData")
	StopTxnSampledData                      = Key("StopTxnSampledData")
	SupportedFeatureProfiles                = Key("SupportedFeatureProfiles")
	TransactionMessageAttempts              = Key("TransactionMessageAttempts")
	TransactionMessageRetryInterval         = Key("TransactionMessageRetryInterval")
	UnlockConnectorOnEVSideDisconnect       = Key("UnlockConnectorOnEVSideDisconnect")
	WebSocketPingInterval                   = Key("WebSocketPingInterval")
	LocalAuthListEnabled                    = Key("LocalAuthListEnabled")
	LocalAuthListMaxLength                  = Key("LocalAuthListMaxLength")
	SendLocalListMaxLength                  = Key("SendLocalListMaxLength")
	ChargeProfileMaxStackLevel              = Key("ChargeProfileMaxStackLevel")
	ChargingScheduleAllowedChargingRateUnit = Key("ChargingScheduleAllowedChargingRateUnit")
	ChargingScheduleMaxPeriods              = Key("ChargingScheduleMaxPeriods")
	MaxChargingProfilesInstalled            = Key("MaxChargingProfilesInstalled")
)

type (
	Key string

	Config struct {
		Version int                     `fig:"version" default:"1"`
		Keys    []core.ConfigurationKey `fig:"keys"`
	}
)

func (k Key) String() string {
	return string(k)
}

// UpdateKey Update the configuration variable in the configuration if it is not readonly.
func (config *Config) UpdateKey(key string, value string) error {
	log.Debugf("Updating key %s to %s", key, value)

	for i, configKey := range config.Keys {
		if configKey.Key == key {
			if !configKey.Readonly {
				config.Keys[i].Value = value
				return nil
			}

			return ErrReadOnly
		}
	}

	return ErrKeyNotFound
}

//GetConfigurationValue Get the value of specified configuration variable in String format.
func (config *Config) GetConfigurationValue(key string) (string, error) {
	log.Debugf("Getting key %s", key)

	for _, configKey := range config.Keys {
		if configKey.Key == key {
			return configKey.Value, nil
		}
	}

	return "", ErrKeyNotFound
}

// GetConfig Get the configuration
func (config *Config) GetConfig() []core.ConfigurationKey {
	return config.Keys
}
