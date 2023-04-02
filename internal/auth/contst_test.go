package auth

import (
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

const ()

var (
	okTag = &types.IdTagInfo{
		ExpiryDate: types.NewDateTime(time.Now().Add(10 * time.Minute)),
		Status:     types.AuthorizationStatusAccepted,
	}

	blockedTag = &types.IdTagInfo{
		ExpiryDate: types.NewDateTime(time.Now().Add(40 * time.Minute)),
		Status:     types.AuthorizationStatusBlocked,
	}

	expiredTag = &types.IdTagInfo{
		ExpiryDate: types.NewDateTime(time.Date(1999, 1, 1, 1, 1, 1, 0, time.Local)),
		Status:     types.AuthorizationStatusAccepted,
	}
)
