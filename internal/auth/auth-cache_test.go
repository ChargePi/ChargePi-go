package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/ChargePi/ChargePi-go/internal/auth/mocks"
	"github.com/ChargePi/ChargePi-go/pkg/util"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

var (
	okTag = &types.IdTagInfo{
		ParentIdTag: "",
		ExpiryDate:  types.Now(),
		Status:      types.AuthorizationStatusAccepted,
	}

	blockedTag = &types.IdTagInfo{
		ExpiryDate: types.NewDateTime(time.Now().Add(40 * time.Minute)),
		Status:     types.AuthorizationStatusBlocked,
	}

	expiredTag = &types.IdTagInfo{
		ExpiryDate: types.NewDateTime(time.Date(1999, 1, 1, 1, 1, 1, 0, time.Local)),
		Status:     types.AuthorizationStatusAccepted,
	}
)

type authCacheTestSuite struct {
	suite.Suite
	authCache     *cache
	tagRepository *mocks.MockTagRepository
}

func (s *authCacheTestSuite) SetupTest() {
	s.tagRepository = mocks.NewMockTagRepository(s.T())
	s.authCache = newAuthCache(s.tagRepository)
}

func (s *authCacheTestSuite) TearDownTest() {}

func (s *authCacheTestSuite) TestAddTag() {
	tests := []struct {
		name     string
		tagLimit int
		tagId    string
		tag      *types.IdTagInfo
		wantErr  bool
	}{
		{
			name:     "Tag can be added",
			tagLimit: 1,
			tagId:    "1",
			tag:      okTag,
			wantErr:  false,
		},
		{
			name:     "Max limit reached",
			tag:      okTag,
			tagId:    "2",
			tagLimit: 3,
			wantErr:  true,
		},
		{
			name:     "Nil tag",
			tag:      nil,
			tagLimit: 1,
			tagId:    "3",
			wantErr:  true,
		},
		{
			name:     "Database error",
			tag:      okTag,
			tagLimit: 1,
			tagId:    "4",
			wantErr:  true,
		},
		{
			name: "Validation failed",
			tag: &types.IdTagInfo{
				Status: "Unknown status",
			},
			tagId:   "3",
			wantErr: true,
		},
		{
			name: "Invalid status",
			tag: &types.IdTagInfo{
				ExpiryDate: nil,
				Status:     "Unknown type",
			},
			tagId:   "5",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s.authCache.SetMaxCachedTags(tt.tagLimit)

			switch tt.name {
			case "Nil tag":
			case "Tag can be added":
				s.tagRepository.EXPECT().GetTags().Return([]*types.IdTagInfo{}, nil).Once()
				s.tagRepository.EXPECT().AddTag(tt.tagId, tt.tag).Return(nil).Once()
			case "Validation failed":
			case "Invalid status":
			case "Database error":
				s.tagRepository.EXPECT().GetTags().Return([]*types.IdTagInfo{}, nil).Once()
				s.tagRepository.EXPECT().AddTag(tt.tagId, tt.tag).Return(errors.New("database error")).Once()
			case "Max limit reached":
				// Append the max number of tags to the cache
				tags := []*types.IdTagInfo{}
				for i := 0; i < tt.tagLimit; i++ {
					tags = append(tags, okTag)
				}
				s.tagRepository.EXPECT().GetTags().Return(tags, nil).Once()
			}

			err := s.authCache.AddTag(tt.tagId, tt.tag)
			if tt.wantErr {
				s.Assert().Error(err)
			} else {
				s.Assert().NoError(err)
			}
		})
	}
}

func (s *authCacheTestSuite) TestRemoveCachedTags() {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "No tags to remove",
			wantErr: false,
		},
		{
			name:    "Remove tags",
			wantErr: false,
		},
		{
			name:    "Database error",
			wantErr: true,
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t *testing.T) {

			if test.name == "Database error" {
				s.tagRepository.EXPECT().RemoveTags().Return(errors.New("database error")).Once()
			} else {
				s.tagRepository.EXPECT().RemoveTags().Return(nil).Once()
			}

			err := s.authCache.RemoveCachedTags()
			if test.wantErr {
				s.Assert().Error(err)
			} else {
				s.Assert().NoError(err)
			}
		})
	}
}

func (s *authCacheTestSuite) TestGetTag() {
	tests := []struct {
		name        string
		expectedTag *types.IdTagInfo
		tagId       string
		wantErr     bool
	}{
		{
			name:  "Tag found",
			tagId: util.GenerateRandomTag(),
			expectedTag: &types.IdTagInfo{
				ExpiryDate: types.NewDateTime(time.Now().Add(10 * time.Minute)),
				Status:     types.AuthorizationStatusAccepted,
			},
			wantErr: false,
		},
		{
			name:        "Tag not found",
			tagId:       "randomTag123",
			expectedTag: nil,
			wantErr:     true,
		},
		{
			name:        "Database error",
			tagId:       util.GenerateRandomTag(),
			expectedTag: nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			if tt.expectedTag != nil {
				s.tagRepository.EXPECT().GetTag(tt.tagId).Return(tt.expectedTag, nil)
			} else if tt.name == "Database error" {
				s.tagRepository.EXPECT().GetTag(tt.tagId).Return(nil, errors.New("error"))
			} else {
				s.tagRepository.EXPECT().GetTag(tt.tagId).Return(nil, errors.New("not found"))
			}

			tag, err := s.authCache.GetTag(tt.tagId)
			if tt.wantErr {
				s.Assert().Error(err)
			} else {
				s.Assert().NoError(err)
				s.Assert().EqualValues(tt.expectedTag, tag)
			}
		})
	}
}

func TestAuthCache(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(authCacheTestSuite))
}
