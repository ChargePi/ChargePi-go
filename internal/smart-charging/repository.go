package smartCharging

import (
	"context"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

type ChargingScheduleRepository interface {
	AddProfile(ctx context.Context, profile *types.ChargingProfile) error
	GetProfile(ctx context.Context, profileId int) (*types.ChargingProfile, error)
	RemoveProfile(ctx context.Context, profileId int) error
	GetProfiles(ctx context.Context) ([]types.ChargingProfile, error)
}
