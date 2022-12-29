package auth

import (
	"errors"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
	"strconv"
)

var (
	ErrLocalAuthListNotEnabled = errors.New("auth list not enabled")
	ErrTagNotFound             = errors.New("tag not found")
)

type (
	TagManager interface {
		AddTag(tagId string, tagInfo *types.IdTagInfo) error
		GetTag(tagId string) (*types.IdTagInfo, error)
		GetTags() []localauth.AuthorizationData
		ClearCache()
		SetMaxTags(number int)
		UpdateLocalAuthList(version int, updateType localauth.UpdateType, tags []localauth.AuthorizationData) error
		GetAuthListVersion() int
		ReadLocalAuthList() error
		WriteLocalAuthList() error
	}

	TagManagerImpl struct {
		authList             LocalAuthList
		cache                Cache
		authCacheEnabled     bool
		localAuthListEnabled bool
	}
)

func NewTagManager(filePath string) *TagManagerImpl {
	authCacheEnabled, cacheErr := ocppConfigManager.GetConfigurationValue(configuration.AuthorizationCacheEnabled.String())
	localListLength, err := ocppConfigManager.GetConfigurationValue(configuration.LocalAuthListMaxLength.String())

	if cacheErr != nil {
	}

	if err != nil {
		length := "0"
		localListLength = &length
	}

	maxTags, err := strconv.Atoi(*localListLength)
	if err != nil {
		maxTags = 0
	}

	cache := NewAuthCache()
	authList := NewLocalAuthList(filePath, maxTags)

	return &TagManagerImpl{
		authCacheEnabled: authCacheEnabled != nil && *authCacheEnabled == "true",
		cache:            cache,
		authList:         authList,
	}
}

func (t *TagManagerImpl) AddTag(tagId string, tagInfo *types.IdTagInfo) error {
	if t.authCacheEnabled {
		t.cache.AddTag(tagId, tagInfo)
	}

	if t.localAuthListEnabled {
		return t.authList.AddTag(tagId, tagInfo)
	}

	return nil
}

func (t *TagManagerImpl) ClearCache() {
	if t.authCacheEnabled {
		t.cache.RemoveCachedTags()
	}
}

func (t *TagManagerImpl) SetMaxTags(number int) {
	t.authList.SetMaxTags(number)
	t.cache.SetMaxCachedTags(number)
}

func (t *TagManagerImpl) GetTag(tagId string) (*types.IdTagInfo, error) {
	logInfo := log.WithField("tagId", tagId)

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

func (t *TagManagerImpl) GetTags() []localauth.AuthorizationData {
	if !t.localAuthListEnabled {
		return []localauth.AuthorizationData{}
	}

	return t.authList.GetTags()
}

func (t *TagManagerImpl) GetAuthListVersion() int {
	if !t.localAuthListEnabled {
		return -1
	}

	return t.authList.GetVersion()
}

func (t *TagManagerImpl) UpdateLocalAuthList(version int, updateType localauth.UpdateType, tags []localauth.AuthorizationData) error {
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

func (t *TagManagerImpl) ReadLocalAuthList() error {
	if !t.localAuthListEnabled {
		return ErrLocalAuthListNotEnabled
	}

	return t.authList.LoadFromFile()
}

func (t *TagManagerImpl) WriteLocalAuthList() error {
	if !t.localAuthListEnabled {
		return ErrLocalAuthListNotEnabled
	}

	return t.authList.WriteToFile()
}
