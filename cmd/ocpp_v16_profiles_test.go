package cmd

import (
	"strings"
	"testing"

	"github.com/ChargePi/ocpp-manager/ocpp_v16"
)

func TestNewDefaultOcppV16Configuration(t *testing.T) {
	cfg, err := newDefaultOcppV16Configuration()
	if err != nil {
		t.Fatalf("newDefaultOcppV16Configuration returned error: %v", err)
	}

	profiles, err := getOcppV16ConfigValue(cfg, ocpp_v16.SupportedFeatureProfiles)
	if err != nil {
		t.Fatal(err)
	}
	if profiles == nil {
		t.Fatal("SupportedFeatureProfiles value is nil")
	}

	expectedProfiles := strings.Join(supportedOcppV16Profiles, ", ")
	if *profiles != expectedProfiles {
		t.Fatalf("SupportedFeatureProfiles = %q, want %q", *profiles, expectedProfiles)
	}

	reserveConnectorZeroSupported, err := getOcppV16ConfigValue(cfg, ocpp_v16.ReserveConnectorZeroSupported)
	if err != nil {
		t.Fatal(err)
	}
	if reserveConnectorZeroSupported == nil || *reserveConnectorZeroSupported != "false" {
		t.Fatalf("ReserveConnectorZeroSupported = %v, want false", reserveConnectorZeroSupported)
	}

	for _, key := range []ocpp_v16.Key{ocpp_v16.AuthorizationCacheEnabled, ocpp_v16.LightIntensity} {
		if _, err := getOcppV16ConfigValue(cfg, key); err != nil {
			t.Fatal(err)
		}
	}

	if _, err := ocpp_v16.NewV16ConfigurationManager(*cfg, supportedOcppV16ConfigurationProfiles...); err != nil {
		t.Fatalf("NewV16ConfigurationManager returned error: %v", err)
	}
}
