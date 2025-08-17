package smartCharging

import (
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/samber/lo"
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
	profile := lo.MaxBy(profiles, func(profile types.ChargingProfile, profile2 types.ChargingProfile) bool {
		return float64(profile.StackLevel) > float64(profile2.StackLevel)
	})

	return &profile
}

func getProfilesWithPurpose(purpose types.ChargingProfilePurposeType, profiles []types.ChargingProfile) []types.ChargingProfile {
	return lo.Filter(profiles, func(profile types.ChargingProfile, _ int) bool {
		return profile.ChargingProfilePurpose == purpose
	})
}
