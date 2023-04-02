package settings

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
)

type AuthList struct {
	Version int                           `json:"version" yaml:"version" mapstructure:"version"`
	Tags    []localauth.AuthorizationData `json:"tags" yaml:"tags" mapstructure:"tags"`
}
