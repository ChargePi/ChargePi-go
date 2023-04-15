package v16

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/xBlaz3kx/ChargePi-go/internal/auth"
)

func (cp *ChargePoint) OnGetLocalListVersion(request *localauth.GetLocalListVersionRequest) (confirmation *localauth.GetLocalListVersionConfirmation, err error) {
	cp.logger.Infof("Received request %s", request.GetFeatureName())
	version := cp.tagManager.GetAuthListVersion()
	res := localauth.NewGetLocalListVersionConfirmation(version)
	return res, nil
}

func (cp *ChargePoint) OnSendLocalList(request *localauth.SendLocalListRequest) (confirmation *localauth.SendLocalListConfirmation, err error) {
	cp.logger.Infof("Received request %s", request.GetFeatureName())

	res := localauth.UpdateStatusFailed

	updateErr := cp.tagManager.UpdateLocalAuthList(request.ListVersion, request.UpdateType, request.LocalAuthorizationList)
	switch updateErr {
	case nil:
		res = localauth.UpdateStatusAccepted
	case auth.ErrLocalAuthListNotEnabled:
		res = localauth.UpdateStatusNotSupported
	}

	return localauth.NewSendLocalListConfirmation(res), nil
}
