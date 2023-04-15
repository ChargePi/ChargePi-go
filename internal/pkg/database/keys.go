package database

import (
	"fmt"

	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

func GetEvseKey(evseId int) string {
	return fmt.Sprintf("evse-%d", evseId)
}

func GetLocalAuthTagPrefix(tagId string) []byte {
	return []byte(fmt.Sprintf("auth-tag-%s", tagId))
}

func GetLocalAuthVersion() []byte {
	return []byte("auth-version")
}

func GetSmartChargingProfile(profileId int) []byte {
	return []byte(fmt.Sprintf("profile-%d", profileId))
}

func GetOcppConfigurationKey(version configuration.ProtocolVersion) []byte {
	return []byte(fmt.Sprintf("ocpp-configuration-%s", version))
}

func GetSettingsKey() []byte {
	return []byte(fmt.Sprintf("settings"))
}
