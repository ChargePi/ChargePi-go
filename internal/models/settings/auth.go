package settings

import "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"

type (
	AuthorizationFile struct {
		Version       int               `fig:"Version" validation:"required" json:"version,omitempty" yaml:"version"`
		MaxCachedTags int               `fig:"MaxCachedTags" validation:"required" json:"MaxCachedTags,omitempty" yaml:"MaxCachedTags"`
		Tags          []types.IdTagInfo `fig:"Tags" json:"tags,omitempty" yaml:"tags"`
	}
)
