package smartCharging

import (
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

func getValidProfiles(profiles []types.ChargingProfile) []types.ChargingProfile {
	ret := []types.ChargingProfile{}

	for _, profile := range profiles {

		// Exception is TxProfile
		if profile.ChargingProfilePurpose == types.ChargingProfilePurposeTxProfile {
			ret = append(ret, profile)
			continue
		}

		// Valid if the validity dates are not set
		if profile.ValidFrom == nil && profile.ValidTo == nil {
			ret = append(ret, profile)
			continue
		}

		// Check if the date hasn't expired yet first
		if profile.ValidTo != nil && time.Now().Before(profile.ValidTo.Time) {
			ret = append(ret, profile)
			continue
		}

		// Check if it is even valid
		if profile.ValidFrom != nil && time.Now().After(profile.ValidFrom.Time) {
			ret = append(ret, profile)
			continue
		}
	}

	return ret
}

func getProfileWithHighestStack(profiles []types.ChargingProfile) *types.ChargingProfile {
	var ret *types.ChargingProfile

	switch len(profiles) {
	case 0:
		return nil
	case 1:
		return &profiles[0]
	}

	maxStackLevel := profiles[0].StackLevel

	for _, profile := range profiles[1:] {
		if profile.StackLevel > maxStackLevel {
			maxStackLevel = profile.StackLevel
			ret = &profile
		}
	}

	return ret
}

func getProfilesWithPurpose(purpose types.ChargingProfilePurposeType, profiles []types.ChargingProfile) []types.ChargingProfile {
	ret := []types.ChargingProfile{}

	for _, profile := range profiles {
		if profile.ChargingProfilePurpose == purpose {
			ret = append(ret, profile)
		}
	}

	return ret
}
