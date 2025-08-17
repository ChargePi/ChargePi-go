package v16

import (
	"context"
	"errors"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"

	"github.com/ChargePi/ChargePi-go/internal/auth"
)

func (cp *ChargePoint) OnGetLocalListVersion(request *localauth.GetLocalListVersionRequest) (confirmation *localauth.GetLocalListVersionConfirmation, err error) {
	cp.logger.Sugar().Infof("Received request %s", request.GetFeatureName())
	version := cp.tagAuthService.GetAuthListVersion()

	res := localauth.NewGetLocalListVersionConfirmation(version)
	return res, nil
}

func (cp *ChargePoint) OnSendLocalList(request *localauth.SendLocalListRequest) (confirmation *localauth.SendLocalListConfirmation, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cp.logger.Sugar().Infof("Received request %s", request.GetFeatureName())

	res := localauth.UpdateStatusFailed

	updateErr := cp.tagAuthService.UpdateLocalAuthList(ctx, request.ListVersion, request.UpdateType, request.LocalAuthorizationList)
	switch {
	case updateErr == nil:
		res = localauth.UpdateStatusAccepted
	case errors.Is(updateErr, auth.ErrLocalAuthListDisabled):
		res = localauth.UpdateStatusNotSupported
	}

	return localauth.NewSendLocalListConfirmation(res), nil
}
