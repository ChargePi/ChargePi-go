package cmd

import (
	"fmt"
	"strings"

	"github.com/ChargePi/ocpp-manager/ocpp_v16"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
)

var supportedOcppV16ConfigurationProfiles = []string{
	core.ProfileName,
	localauth.ProfileName,
}

func newDefaultOcppV16Configuration() (*ocpp_v16.Config, error) {
	cfg, err := ocpp_v16.DefaultConfigurationFromProfiles(supportedOcppV16ConfigurationProfiles...)
	if err != nil {
		return nil, err
	}

	upsertOcppV16ConfigKey(cfg, ocpp_v16.SupportedFeatureProfiles, strings.Join(supportedOcppV16Profiles, ", "), true)
	upsertOcppV16ConfigKey(cfg, ocpp_v16.AuthorizationCacheEnabled, "true", false)
	upsertOcppV16ConfigKey(cfg, ocpp_v16.ReserveConnectorZeroSupported, "false", true)
	upsertOcppV16ConfigKey(cfg, ocpp_v16.LightIntensity, "100", false)

	return cfg, nil
}

func upsertOcppV16ConfigKey(cfg *ocpp_v16.Config, key ocpp_v16.Key, value string, readonly bool) {
	for i := range cfg.Keys {
		if cfg.Keys[i].Key == key.String() {
			cfg.Keys[i].Readonly = readonly
			cfg.Keys[i].Value = &value
			return
		}
	}

	cfg.Keys = append(cfg.Keys, core.ConfigurationKey{
		Key:      key.String(),
		Readonly: readonly,
		Value:    &value,
	})
}

func getOcppV16ConfigValue(cfg *ocpp_v16.Config, key ocpp_v16.Key) (*string, error) {
	for _, item := range cfg.Keys {
		if item.Key == key.String() {
			return item.Value, nil
		}
	}

	return nil, fmt.Errorf("missing OCPP 1.6 configuration key %s", key)
}
