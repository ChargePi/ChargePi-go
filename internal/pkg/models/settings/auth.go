package settings

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
)

type AuthList struct {
	Version int
	Tags    []localauth.AuthorizationData
}
