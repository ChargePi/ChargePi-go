package auth

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

var (
	ErrLocalAuthListDisabled = errors.New("auth list disabled")
	ErrCacheDisabled         = errors.New("auth cache disabled")
	ErrTagNotFound           = errors.New("tag not found")
)

type Service interface {
	CacheTag(ctx context.Context, tagId string, tagInfo *types.IdTagInfo) error
	GetTag(ctx context.Context, tagId string) (*types.IdTagInfo, error)
	GetTags(ctx context.Context) ([]localauth.AuthorizationData, error)
	RemoveTag(ctx context.Context, tagId string) error
	ClearCache(ctx context.Context) error
	SetMaxTags(ctx context.Context, number int)
	ToggleAuthCache(enabled bool)
	ToggleLocalAuthList(enabled bool)
	UpdateLocalAuthList(ctx context.Context, version int, updateType localauth.UpdateType, tags []localauth.AuthorizationData) error
	GetAuthListVersion() int
}

type ManagerV1 struct {
	authList             LocalAuthList
	cache                Cache
	authCacheEnabled     bool
	localAuthListEnabled bool
	logger               *zap.Logger
}

func NewManager(logger *zap.Logger, localAuthRepository LocalAuthListRepository, tagRepository TagRepository) *ManagerV1 {
	cache := newAuthCache(logger, tagRepository)
	authList := newLocalAuthList(logger, localAuthRepository, 10)

	return &ManagerV1{
		// Cache is enabled by default
		authCacheEnabled: true,
		// Local Auth List is disabled by default
		localAuthListEnabled: false,
		cache:                cache,
		authList:             authList,
		logger:               logger.Named("auth_manager"),
	}
}

// CacheTag caches a tag in the auth cache, if enabled.
func (t *ManagerV1) CacheTag(ctx context.Context, tagId string, tagInfo *types.IdTagInfo) error {
	t.logger.With(zap.String("tagId", tagId)).Debug("Adding a tag to system")

	if t.authCacheEnabled {
		err := t.cache.AddTag(ctx, tagId, tagInfo)
		if err != nil {
			return err
		}
	}

	return ErrCacheDisabled
}

// ClearCache clears the auth cache, if enabled.
func (t *ManagerV1) ClearCache(ctx context.Context) error {
	t.logger.Debug("Clearing the tag cache")

	if t.authCacheEnabled {
		return t.cache.RemoveCachedTags(ctx)
	}

	return ErrCacheDisabled
}

// SetMaxTags sets the maximum number of tags that can be cached.
func (t *ManagerV1) SetMaxTags(ctx context.Context, number int) {
	t.logger.Debug("Setting the maximum number of stored tags")

	if t.authCacheEnabled {
		t.cache.SetMaxCachedTags(ctx, number)
	}

	if t.localAuthListEnabled {
		t.authList.SetMaxTags(number)
	}
}

// GetTag returns a tag with id from either the Local Auth List or the auth cache.
// If both are disabled, an error is returned.
func (t *ManagerV1) GetTag(ctx context.Context, tagId string) (*types.IdTagInfo, error) {
	logger := t.logger.With(zap.String("tagId", tagId))

	// Check the localAuthList first
	if t.localAuthListEnabled {
		logger.Info("Getting the tag from local auth list")
		tag, err := t.authList.GetTag(ctx, tagId)
		if err != nil {
			goto CheckCache
		}

		return tag, err
	}

CheckCache:
	// Check the cache
	if t.authCacheEnabled {
		logger.Info("Getting the tag from cache")
		return t.cache.GetTag(ctx, tagId)
	}

	return nil, ErrTagNotFound
}

// GetTags returns all tags from the local auth list.
func (t *ManagerV1) GetTags(ctx context.Context) ([]localauth.AuthorizationData, error) {
	t.logger.Debug("Getting all tags from local auth list")

	if !t.localAuthListEnabled {
		return nil, ErrLocalAuthListDisabled
	}

	return t.authList.GetTags(ctx)
}

// GetAuthListVersion returns the current version of the local auth list.
func (t *ManagerV1) GetAuthListVersion() int {
	t.logger.Debug("Getting the local auth list version")

	if !t.localAuthListEnabled {
		return -1
	}

	return t.authList.GetVersion()
}

// RemoveTag removes a tag from the auth cache, if enabled.
func (t *ManagerV1) RemoveTag(ctx context.Context, tagId string) error {
	t.logger.With(zap.String("tagId", tagId)).Debug("Removing a tag from system")

	if !t.localAuthListEnabled {
		return ErrLocalAuthListDisabled
	}

	return t.authList.RemoveTag(ctx, tagId)
}

// UpdateLocalAuthList updates the local auth list with the given tags.
func (t *ManagerV1) UpdateLocalAuthList(ctx context.Context, version int, updateType localauth.UpdateType, tags []localauth.AuthorizationData) error {
	t.logger.With(zap.Int("version", version), zap.String("updateType", string(updateType))).
		Debug("Updating the local auth list")

	if !t.localAuthListEnabled {
		return ErrLocalAuthListDisabled
	}

	switch updateType {
	case localauth.UpdateTypeDifferential:

		for _, tag := range tags {
			err := t.authList.UpdateTag(ctx, tag.IdTag, tag.IdTagInfo)
			if err != nil {
				return err
			}
		}

	case localauth.UpdateTypeFull:
		t.authList.RemoveAll(ctx)
		for _, tag := range tags {
			err := t.authList.AddTag(ctx, tag.IdTag, tag.IdTagInfo)
			if err != nil {
				return err
			}
		}
	}

	t.authList.SetVersion(version)
	return nil
}

func (t *ManagerV1) ToggleAuthCache(enabled bool) {
	t.authCacheEnabled = enabled
}

func (t *ManagerV1) ToggleLocalAuthList(enabled bool) {
	t.localAuthListEnabled = enabled
}
