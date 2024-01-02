package auth

import (
	"errors"

	"github.com/dgraph-io/badger/v3"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
)

var (
	ErrLocalAuthListNotEnabled = errors.New("auth list not enabled")
	ErrCacheNotEnabled         = errors.New("auth cache not enabled")
	ErrTagNotFound             = errors.New("tag not found")
)

type (
	TagManager interface {
		AddTag(tagId string, tagInfo *types.IdTagInfo) error
		GetTag(tagId string) (*types.IdTagInfo, error)
		GetTags() []localauth.AuthorizationData
		RemoveTag(tagId string) error
		ClearCache() error
		SetMaxTags(number int)
		ToggleAuthCache(enabled bool)
		ToggleLocalAuthList(enabled bool)
		UpdateLocalAuthList(version int, updateType localauth.UpdateType, tags []localauth.AuthorizationData) error
		GetAuthListVersion() int
	}

	TagManagerImpl struct {
		authList             LocalAuthList
		cache                Cache
		authCacheEnabled     bool
		localAuthListEnabled bool
		logger               log.FieldLogger
	}
)

func NewTagManager(db *badger.DB) *TagManagerImpl {
	cache := NewAuthCache(db)
	authList := NewLocalAuthList(db, 10)

	return &TagManagerImpl{
		authCacheEnabled:     true,
		localAuthListEnabled: false,
		cache:                cache,
		authList:             authList,
		logger:               log.StandardLogger().WithField("component", "tag-manager"),
	}
}

// AddTag adds a tag to the auth cache, if enabled.
func (t *TagManagerImpl) AddTag(tagId string, tagInfo *types.IdTagInfo) error {
	t.logger.WithField("tagId", tagId).Debug("Adding a tag to system")

	if t.authCacheEnabled {
		t.cache.AddTag(tagId, tagInfo)
	}

	return nil
}

// ClearCache clears the auth cache, if enabled.
func (t *TagManagerImpl) ClearCache() error {
	t.logger.Debug("Clearing the tag cache")

	if t.authCacheEnabled {
		t.cache.RemoveCachedTags()
		return nil
	}

	return ErrCacheNotEnabled
}

// SetMaxTags sets the maximum number of tags that can be cached.
func (t *TagManagerImpl) SetMaxTags(number int) {
	t.logger.Debug("Setting the maximum number of stored tags")

	t.authList.SetMaxTags(number)
	t.cache.SetMaxCachedTags(number)
}

// GetTag returns a tag from either the Local Auth List or the auth cache. If both are disabled, an error is returned.
func (t *TagManagerImpl) GetTag(tagId string) (*types.IdTagInfo, error) {
	logInfo := t.logger.WithField("tagId", tagId)

	// Check the localAuthList first
	if t.localAuthListEnabled {
		logInfo.Infof("Getting the tag from localAuthList")
		tag, err := t.authList.GetTag(tagId)
		if err != nil {
			goto CheckCache
		}

		return tag, err
	}

CheckCache:
	// Check the cache
	if t.authCacheEnabled {
		logInfo.Infof("Getting the tag from authCache")
		return t.authList.GetTag(tagId)
	}

	return nil, ErrTagNotFound
}

// GetTags returns all tags (only from the Local Auth List). The cached tags are not returned.
func (t *TagManagerImpl) GetTags() []localauth.AuthorizationData {
	t.logger.Debug("Getting all tags from localAuthList")

	if !t.localAuthListEnabled {
		return []localauth.AuthorizationData{}
	}

	return t.authList.GetTags()
}

// GetAuthListVersion returns the current version of the local auth list.
func (t *TagManagerImpl) GetAuthListVersion() int {
	t.logger.Debug("Getting the local auth list version")

	if !t.localAuthListEnabled {
		return -1
	}

	return t.authList.GetVersion()
}

// RemoveTag removes a tag from the auth cache, if enabled.
func (t *TagManagerImpl) RemoveTag(tagId string) error {
	t.logger.WithField("tagId", tagId).Debug("Removing a tag from system")

	if !t.localAuthListEnabled {
		return ErrLocalAuthListNotEnabled
	}

	return t.authList.RemoveTag(tagId)
}

// UpdateLocalAuthList updates the local auth list with the given tags.
func (t *TagManagerImpl) UpdateLocalAuthList(version int, updateType localauth.UpdateType, tags []localauth.AuthorizationData) error {
	t.logger.WithField("version", version).
		WithField("updateType", updateType).
		Debug("Updating the local auth list")

	if !t.localAuthListEnabled {
		return ErrLocalAuthListNotEnabled
	}

	switch updateType {
	case localauth.UpdateTypeDifferential:

		for _, tag := range tags {
			err := t.authList.UpdateTag(tag.IdTag, tag.IdTagInfo)
			if err != nil {
				return err
			}
		}

	case localauth.UpdateTypeFull:
		t.authList.RemoveAll()
		for _, tag := range tags {
			err := t.authList.AddTag(tag.IdTag, tag.IdTagInfo)
			if err != nil {
				return err
			}
		}
	}

	t.authList.SetVersion(version)
	return nil
}

func (t *TagManagerImpl) ToggleAuthCache(enabled bool) {
	t.authCacheEnabled = enabled
}

func (t *TagManagerImpl) ToggleLocalAuthList(enabled bool) {
	t.localAuthListEnabled = enabled
}
