package list

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
)

type LocalAuthListVersion struct {
	// Version of the local authorization list.
	Version int `json:"version" yaml:"version" mapstructure:"version" validate:"min=1"`

	// List of authorization data entries.
	Tags []localauth.AuthorizationData `json:"tags" yaml:"tags" mapstructure:"tags" validate:"min=1"`
}
