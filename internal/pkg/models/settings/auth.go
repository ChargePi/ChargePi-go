package settings

import "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"

type (
	Tag struct {
		TagId   string          `json:"tagId" yaml:"tagId"`
		TagInfo types.IdTagInfo `json:"tagInfo" yaml:"tagInfo"`
	}

	AuthorizationFile struct {
		Tags []Tag `json:"tags,omitempty" yaml:"tags"`
	}

	LocalAuthListFile struct {
		Version int   `json:"version" yaml:"version" validation:"required"`
		Tags    []Tag `json:"tags,omitempty" yaml:"tags"`
	}
)
