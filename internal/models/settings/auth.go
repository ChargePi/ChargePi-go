package settings

import "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"

type (
	AuthorizationFile struct {
		Version       int               `fig:"Version" validation:"required"`
		MaxCachedTags int               `fig:"MaxCachedTags" validation:"required"`
		Tags          []types.IdTagInfo `fig:"Tags"`
	}
)
