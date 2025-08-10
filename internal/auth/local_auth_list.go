package auth

import (
	"errors"

	"github.com/ChargePi/ChargePi-go/internal/auth/list"
	"github.com/ChargePi/ChargePi-go/pkg/util"
	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	log "github.com/sirupsen/logrus"
)

var (
	ErrTagLimitReached = errors.New("tat limit reached")
	ErrorTagNotFound   = errors.New("tag not found")
	ErrInvalidTagId    = errors.New("tag ID invalid")
	ErrTagNil          = errors.New("tag cannot be nil")
)

type (
	LocalAuthList interface {
		AddTag(tagId string, tagInfo *types.IdTagInfo) error
		UpdateTag(tagId string, tagInfo *types.IdTagInfo) error
		RemoveTag(tagId string) error
		RemoveAll()
		GetTag(tagId string) (*types.IdTagInfo, error)
		GetTags() ([]localauth.AuthorizationData, error)
		SetMaxTags(number int)
		GetVersion() int
		SetVersion(version int)
	}

	List struct {
		repository LocalAuthListRepository
		maxTags    int
		logger     log.FieldLogger
	}
)

type LocalAuthListRepository interface {
	AddTagToAuthList(tagId string, tagInfo *types.IdTagInfo) error
	RemoveAuthListTag(tagId string) error
	GetLocalAuthListTag(tagId string) (*types.IdTagInfo, error)
	GetLocalAuthListTags() ([]localauth.AuthorizationData, error)
	GetAuthListTagsForVersion(version int) ([]localauth.AuthorizationData, error)
	RemoveAuthListAllTagsForVersion(version int) error
	AddAuthList(list.LocalAuthListVersion) error
}

func newLocalAuthList(repository LocalAuthListRepository, maxTags int) *List {
	return &List{
		repository: repository,
		maxTags:    maxTags,
		logger:     log.StandardLogger().WithField("component", "local-auth-list"),
	}
}

// AddTag Add a tag to the local auth list.
func (l *List) AddTag(tagId string, tagInfo *types.IdTagInfo) error {
	logInfo := l.logger.WithField("tagId", tagId)
	logInfo.Debug("Adding a tag to local auth list")

	if stringUtils.IsEmpty(tagId) {
		return ErrInvalidTagId
	}

	if util.IsNilInterfaceOrPointer(tagInfo) {
		return ErrTagNil
	}

	err := ocppj.Validate.Struct(tagInfo)
	if err != nil {
		return err
	}

	tags, err := l.repository.GetLocalAuthListTags()
	if err != nil {
		return err
	}

	if len(tags) >= l.maxTags {
		return ErrTagLimitReached
	}

	return l.repository.AddTagToAuthList(tagId, tagInfo)
}

// RemoveTag Remove a tag with the ID from the Local Auth List.
func (l *List) RemoveTag(tagId string) error {
	logInfo := l.logger.WithField("tagId", tagId)
	logInfo.Debug("Removing a tag from local auth list")

	if stringUtils.IsEmpty(tagId) {
		return ErrInvalidTagId
	}

	return l.repository.RemoveAuthListTag(tagId)
}

// RemoveAll Remove all tags.
func (l *List) RemoveAll() {
	l.logger.Debugf("Removing local auth list")

	// Remove all cached keys from database
	err := l.repository.RemoveAuthListAllTagsForVersion(l.GetVersion())
	if err != nil {
		l.logger.WithError(err).Error("Error removing local auth list")
	}
}

// GetTag Get a tag from local auth list.
func (l *List) GetTag(tagId string) (*types.IdTagInfo, error) {
	logInfo := l.logger.WithField("tagId", tagId)
	logInfo.Info("Fetching the tag")

	if stringUtils.IsEmpty(tagId) {
		return nil, ErrInvalidTagId
	}

	return l.repository.GetLocalAuthListTag(tagId)
}

// GetTags Get all tags stored in the Local Auth store.
func (l *List) GetTags() ([]localauth.AuthorizationData, error) {
	l.logger.Infof("Fetching tags")

	tags, err := l.repository.GetLocalAuthListTags()
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// UpdateTag Update a tag in the Local Auth store.
func (l *List) UpdateTag(tagId string, tagInfo *types.IdTagInfo) error {
	logInfo := l.logger.WithField("tagId", tagId)
	logInfo.Info("Updating expectedTag")

	if stringUtils.IsEmpty(tagId) {
		return ErrInvalidTagId
	}

	if util.IsNilInterfaceOrPointer(tagInfo) {
		return ErrTagNil
	}

	err := ocppj.Validate.Struct(tagInfo)
	if err != nil {
		return err
	}

	// todo
	return nil
}

// GetVersion Get the current version of the Local Auth list.
func (l *List) GetVersion() int {
	l.logger.Info("Fetching list version")
	version := -1

	return version
}

// SetVersion Set the current version of the Local Auth list.
func (l *List) SetVersion(version int) {
	logInfo := l.logger.WithField("version", version)
	logInfo.Info("Updating list version")

	// todo
}

// SetMaxTags Set the maximum number of tags that can be stored in the Local Auth list.
func (l *List) SetMaxTags(number int) {
	if number >= 0 {
		l.logger.Debugf("Set max tags to %d", number)
		l.maxTags = number
	}
}
