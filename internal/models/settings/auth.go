package settings

import "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"

type (
	Tag struct {
		TagId   string          `fig:"tagId" json:"tagId,omitempty" yaml:"tagId"`
		TagInfo types.IdTagInfo `fig:"tagInfo" json:"tagInfo" yaml:"tagInfo"`
	}

	AuthorizationFile struct {
		Version int   `fig:"Version" validation:"required" json:"version" yaml:"version"`
		Tags    []Tag `fig:"tags" json:"tags,omitempty" yaml:"tags"`
	}
)
