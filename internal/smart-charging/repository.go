package smartCharging

import "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"

type ChargingScheduleRepository interface {
	AddProfile(profile *types.ChargingProfile) error
	GetProfile(profileId int) (*types.ChargingProfile, error)
	RemoveProfile(profileId int) error
	GetProfiles() ([]types.ChargingProfile, error)
}
