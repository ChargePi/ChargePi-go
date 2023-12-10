package database

import (
	"fmt"

	"github.com/xBlaz3kx/ChargePi-go/pkg/models/ocpp"
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

func GetOcppConfigurationKey(version ocpp.ProtocolVersion) []byte {
	return []byte(fmt.Sprintf("ocpp-configuration-%s", version))
}

func GetSettingsKey() []byte {
	return []byte(fmt.Sprintf("settings"))
}
