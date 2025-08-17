package auth

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ocppj"

	"github.com/ChargePi/ChargePi-go/internal/auth/list"
	"github.com/ChargePi/ChargePi-go/pkg/util"
)

var (
	ErrTagLimitReached = errors.New("tat limit reached")
	ErrorTagNotFound   = errors.New("tag not found")
	ErrInvalidTagId    = errors.New("tag ID invalid")
	ErrTagNil          = errors.New("tag cannot be nil")
)

type (
	LocalAuthList interface {
		AddTag(ctx context.Context, tagId string, tagInfo *types.IdTagInfo) error
		UpdateTag(ctx context.Context, tagId string, tagInfo *types.IdTagInfo) error
		RemoveTag(ctx context.Context, tagId string) error
		RemoveAll(ctx context.Context)
		GetTag(ctx context.Context, tagId string) (*types.IdTagInfo, error)
		GetTags(ctx context.Context) ([]localauth.AuthorizationData, error)
		SetMaxTags(number int)
		GetVersion() int
		SetVersion(version int)
	}

	List struct {
		repository LocalAuthListRepository
		maxTags    int
		logger     *zap.Logger
	}
)

type LocalAuthListRepository interface {
	AddTagToAuthList(ctx context.Context, tagId string, tagInfo *types.IdTagInfo) error
	RemoveAuthListTag(ctx context.Context, tagId string) error
	GetLocalAuthListTag(ctx context.Context, tagId string) (*types.IdTagInfo, error)
	GetLocalAuthListTags(ctx context.Context) ([]localauth.AuthorizationData, error)
	GetAuthListTagsForVersion(ctx context.Context, version int) ([]localauth.AuthorizationData, error)
	RemoveAuthListAllTagsForVersion(ctx context.Context, version int) error
	AddAuthList(ctx context.Context, localAuthListVersion list.LocalAuthListVersion) error
}

func newLocalAuthList(logger *zap.Logger, repository LocalAuthListRepository, maxTags int) *List {
	return &List{
		repository: repository,
		maxTags:    maxTags,
		logger:     logger.Named("local_auth_list"),
	}
}

// AddTag Add a tag to the local auth list.
func (l *List) AddTag(ctx context.Context, tagId string, tagInfo *types.IdTagInfo) error {
	logger := l.logger.With(zap.String("tagId", tagId), zap.Any("tagInfo", tagInfo))
	logger.Debug("Adding a tag to local auth list")

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

	tags, err := l.repository.GetLocalAuthListTags(ctx)
	if err != nil {
		return err
	}

	if len(tags) >= l.maxTags {
		return ErrTagLimitReached
	}

	return l.repository.AddTagToAuthList(ctx, tagId, tagInfo)
}

// RemoveTag Remove a tag with the ID from the Local Auth List.
func (l *List) RemoveTag(ctx context.Context, tagId string) error {
	logger := l.logger.With(zap.String("tagId", tagId))
	logger.Debug("Removing a tag from local auth list")

	if stringUtils.IsEmpty(tagId) {
		return ErrInvalidTagId
	}

	return l.repository.RemoveAuthListTag(ctx, tagId)
}

// RemoveAll Remove all tags.
func (l *List) RemoveAll(ctx context.Context) {
	l.logger.Debug("Removing local auth list")

	// Remove all cached keys from database
	err := l.repository.RemoveAuthListAllTagsForVersion(ctx, l.GetVersion())
	if err != nil {
		l.logger.With(zap.Error(err)).Error("Error removing local auth list")
	}
}

// GetTag Get a tag from local auth list.
func (l *List) GetTag(ctx context.Context, tagId string) (*types.IdTagInfo, error) {
	logger := l.logger.With(zap.String("tagId", tagId))
	logger.Info("Fetching the tag")

	if stringUtils.IsEmpty(tagId) {
		return nil, ErrInvalidTagId
	}

	return l.repository.GetLocalAuthListTag(ctx, tagId)
}

// GetTags Get all tags stored in the Local Auth store.
func (l *List) GetTags(ctx context.Context) ([]localauth.AuthorizationData, error) {
	l.logger.Info("Fetching tags")

	tags, err := l.repository.GetLocalAuthListTags(ctx)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// UpdateTag Update a tag in the Local Auth store.
func (l *List) UpdateTag(ctx context.Context, tagId string, tagInfo *types.IdTagInfo) error {
	logger := l.logger.With(zap.String("tagId", tagId))
	logger.Info("Updating tag")

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
	logInfo := l.logger.With(zap.Int("version", version))
	logInfo.Info("Updating list version")

	// todo
}

// SetMaxTags Set the maximum number of tags that can be stored in the Local Auth list.
func (l *List) SetMaxTags(number int) {
	if number >= 0 {
		l.logger.Sugar().Debugf("Set max tags to %d", number)
		l.maxTags = number
	}
}
