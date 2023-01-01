package auth

import (
	"errors"
	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/kkyr/fig"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	settingsData "github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"path/filepath"
	"strings"
	"sync"
)

const (
	VersionKey = "LocalAuthListVersion"
)

var (
	ErrTagLimitReached = errors.New("tag limit reached")
)

type (
	LocalAuthList interface {
		AddTag(tagId string, tagInfo *types.IdTagInfo) error
		UpdateTag(tagId string, tagInfo *types.IdTagInfo) error
		RemoveTag(tagId string) error
		RemoveAll()
		GetTag(tagId string) (*types.IdTagInfo, error)
		GetTags() []localauth.AuthorizationData
		SetMaxTags(number int)
		GetVersion() int
		SetVersion(version int)
		LoadFromFile() error
		WriteToFile() error
	}

	LocalAuthListImpl struct {
		filePath string
		tags     sync.Map
		numTags  int
		maxTags  int
	}
)

func NewLocalAuthList(filePath string, maxTags int) *LocalAuthListImpl {
	return &LocalAuthListImpl{
		numTags:  0,
		filePath: filePath,
		tags:     sync.Map{},
		maxTags:  maxTags,
	}
}

// AddTag Add a tag to the global authorization cache.
func (l *LocalAuthListImpl) AddTag(tagId string, tagInfo *types.IdTagInfo) error {
	if l.numTags+1 >= l.maxTags {
		return ErrTagLimitReached
	}

	// Add a tag if it doesn't exist in the cache already
	l.tags.Store(tagId, *tagInfo)
	// Update the file
	return l.WriteToFile()
}

// RemoveTag Remove a tag from the global authorization cache.
func (l *LocalAuthListImpl) RemoveTag(tagId string) error {
	l.tags.Delete(tagId)
	return nil
}

// RemoveAll Remove all tags.
func (l *LocalAuthListImpl) RemoveAll() {
	version, isVersionFound := l.tags.Load(VersionKey)
	if !isVersionFound {
		version = 1
	}

	l.tags = sync.Map{}
	l.SetVersion(version.(int))
}

// GetTag Get a tag
func (l *LocalAuthListImpl) GetTag(tagId string) (*types.IdTagInfo, error) {
	log.Infof("Fetching the tag %s", tagId)

	tagObject, isFound := l.tags.Load(tagId)
	if isFound {
		tagInfo := tagObject.(types.IdTagInfo)
		return &tagInfo, nil
	}

	return nil, ErrTagNotFound
}

// GetTags Get all stored tags.
func (l *LocalAuthListImpl) GetTags() []localauth.AuthorizationData {
	log.Infof("Fetching tags")
	var tags []localauth.AuthorizationData

	l.tags.Range(func(key, value interface{}) bool {
		if !stringUtils.Contains(key.(string), VersionKey) {
			tagInfo := value.(types.IdTagInfo)
			tag := localauth.AuthorizationData{
				IdTag:     key.(string),
				IdTagInfo: &tagInfo,
			}

			tags = append(tags, tag)
		}

		return false
	})

	return tags
}

func (l *LocalAuthListImpl) UpdateTag(tagId string, tagInfo *types.IdTagInfo) error {
	//TODO implement me
	panic("implement me")
}

func (l *LocalAuthListImpl) GetVersion() int {
	version, isVersionFound := l.tags.Load(VersionKey)
	if isVersionFound {
		return version.(int)
	}

	return 0
}

func (l *LocalAuthListImpl) SetVersion(version int) {
	l.tags.Store(VersionKey, version)
}

func (l *LocalAuthListImpl) SetMaxTags(number int) {
	if number > 0 {
		l.maxTags = number
	}
}

// LoadFromFile loads tags from the cache file
func (l *LocalAuthListImpl) LoadFromFile() error {
	var (
		auth settingsData.LocalAuthListFile
		err  error
	)

	err = fig.Load(&auth,
		fig.File(filepath.Base(l.filePath)),
		fig.Dirs(filepath.Dir(l.filePath)))
	if err != nil {
		//todo temporary fix - tags with ExpiryDate won't unmarshall successfully
		log.WithError(err).Errorf("Unable to load authorization file")
		return err
	}

	l.tags.Store(VersionKey, auth.Version)

	loadTags(&l.tags, auth.Tags)
	l.numTags = len(auth.Tags)

	log.WithField("version", auth.Version).Infof("Read local authorization file")
	return nil
}

func (l *LocalAuthListImpl) WriteToFile() error {
	log.Debug("Writing local authorization list to a file")
	var (
		authTags                []settingsData.Tag
		version, isVersionFound = l.tags.Load(VersionKey)
	)

	if !isVersionFound {
		version = 1
	}

	l.tags.Range(func(key, value interface{}) bool {
		if !strings.Contains(key.(string), VersionKey) {
			tag := settingsData.Tag{
				TagId:   key.(string),
				TagInfo: value.(types.IdTagInfo),
			}
			authTags = append(authTags, tag)
		}

		return false
	})

	err := util.WriteToFile(l.filePath, settingsData.LocalAuthListFile{
		Version: version.(int),
		Tags:    authTags,
	})
	if err != nil {
		log.WithError(err).Errorf("Error updating local auth list file")
	}

	return err
}

// loadTags loads the tags into the cache
func loadTags(cache *sync.Map, tags []settingsData.Tag) {
	if tags != nil {
		for _, tag := range tags {
			log.Tracef("Adding tag: %v", tag)
			cache.Store(tag.TagId, tag.TagInfo)
		}
	}
}
