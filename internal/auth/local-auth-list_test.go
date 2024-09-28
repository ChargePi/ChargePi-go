package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/ChargePi/ChargePi-go/internal/auth/mocks"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

var (
	okAuthListTag = &types.IdTagInfo{
		ExpiryDate: types.Now(),
		Status:     types.AuthorizationStatusAccepted,
	}

	blockedAuthListTag = &types.IdTagInfo{
		ExpiryDate: types.NewDateTime(time.Now().Add(40 * time.Minute)),
		Status:     types.AuthorizationStatusBlocked,
	}

	expiredAuthListTag = &types.IdTagInfo{
		ExpiryDate: types.NewDateTime(time.Date(1999, 1, 1, 1, 1, 1, 0, time.Local)),
		Status:     types.AuthorizationStatusAccepted,
	}
)

type localAuthListTestSuite struct {
	suite.Suite
	authList       *List
	mockRepository *mocks.MockLocalAuthListRepository
}

func (s *localAuthListTestSuite) SetupTest() {
	s.mockRepository = mocks.NewMockLocalAuthListRepository(s.T())
	s.authList = newLocalAuthList(s.mockRepository, 10)
}

func (s *localAuthListTestSuite) TestAddTag() {
	tests := []struct {
		name  string
		tagId string
		tag   *types.IdTagInfo
		err   error
	}{
		{
			name:  "Added",
			tagId: "tagId",
			tag:   okAuthListTag,
			err:   nil,
		},
		{
			name:  "Max limit reached",
			tagId: "tagId",
			tag:   okAuthListTag,
			err:   ErrTagLimitReached,
		},
		{
			name:  "Invalid tag",
			tag:   nil,
			tagId: "tagId",
			err:   ErrTagNil,
		},
		{
			name:  "Empty Id",
			tag:   okAuthListTag,
			tagId: "",
			err:   ErrInvalidTagId,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			switch tt.name {
			case "Added":
				s.mockRepository.EXPECT().GetLocalAuthListTags().Return([]localauth.AuthorizationData{}, nil).Once()
				s.mockRepository.EXPECT().AddTagToAuthList(tt.tagId, tt.tag).Return(nil).Once()
			case "Max limit reached":
				ret := make([]localauth.AuthorizationData, 10)
				s.mockRepository.EXPECT().GetLocalAuthListTags().Return(ret, nil).Once()
			case "Invalid tag":
			case "Empty Id":
			}

			err := s.authList.AddTag(tt.tagId, tt.tag)
			if tt.err != nil {
				s.Assert().Equal(tt.err, err)
			} else {
				s.Assert().NoError(err)
			}
		})
	}
}

func (s *localAuthListTestSuite) TestUpdateTag() {
	tests := []struct {
		name  string
		tagId string
		tag   *types.IdTagInfo
		err   error
	}{
		{
			name:  "Updated",
			tagId: "tagId",
			tag:   okAuthListTag,
			err:   nil,
		},

		{
			name:  "Empty Id",
			tagId: "",
			tag:   okAuthListTag,
			err:   ErrInvalidTagId,
		},
		{
			name:  "Invalid tag",
			tagId: "tagId",
			tag:   okAuthListTag,
			err:   nil,
		},
		{
			name:  "Nil tag",
			tagId: "tagId",
			tag:   nil,
			err:   ErrTagNil,
		},
		{
			name:  "Database error",
			tagId: "tagId",
			tag:   okAuthListTag,
			err:   nil,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// todo finalize
			err := s.authList.UpdateTag(tt.tagId, tt.tag)
			if tt.err != nil {
				s.Assert().ErrorIs(tt.err, err)
			} else {
				s.Assert().NoError(err)
			}
		})
	}
}

func (s *localAuthListTestSuite) TestRemoveTag() {
	tests := []struct {
		name  string
		tagId string
		err   error
	}{
		{
			name:  "Removed",
			tagId: "exampleTag",
			err:   nil,
		},
		{
			name:  "Invalid ID",
			tagId: "",
			err:   ErrInvalidTagId,
		},
		{
			name:  "Database error",
			tagId: "exampleTag",
			err:   errors.New("database error"),
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Removed":
				s.mockRepository.EXPECT().RemoveAuthListTag(tt.tagId).Return(nil).Once()
			case "Invalid ID":
			case "Database error":
				s.mockRepository.EXPECT().RemoveAuthListTag(tt.tagId).Return(tt.err).Once()
			}

			err := s.authList.RemoveTag(tt.tagId)
			if tt.err != nil {
				s.Assert().ErrorIs(tt.err, err)
			} else {
				s.Assert().NoError(err)
			}
		})
	}
}

func (s *localAuthListTestSuite) TestRemoveAll() {
	tests := []struct {
		name         string
		expectedTags []localauth.AuthorizationData
	}{
		{},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// todo
		})
	}
}

func (s *localAuthListTestSuite) TestGetTag() {
	tests := []struct {
		name        string
		tagId       string
		expectedTag *types.IdTagInfo
		err         error
	}{
		{
			name:        "Tag found",
			tagId:       "tag1",
			expectedTag: okTag,
			err:         nil,
		},
		{
			name:        "Tag ID empty",
			tagId:       "",
			expectedTag: nil,
			err:         ErrInvalidTagId,
		},
		{
			name:        "Database error",
			tagId:       "tag1",
			expectedTag: nil,
			err:         errors.New("database error"),
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			switch tt.name {
			case "Tag found":
				s.mockRepository.EXPECT().GetLocalAuthListTag(tt.tagId).Return(tt.expectedTag, nil).Once()
			case "Tag ID empty":
			case "Database error":
				s.mockRepository.EXPECT().GetLocalAuthListTag(tt.tagId).Return(nil, tt.err).Once()
			}

			tag, err := s.authList.GetTag(tt.tagId)
			if tt.err != nil {
				s.Assert().ErrorIs(tt.err, err)
			} else {
				s.Assert().NoError(err)
				s.Assert().Equal(tt.expectedTag, tag)
			}
		})
	}
}

func (s *localAuthListTestSuite) TestGetTags() {
	tests := []struct {
		name         string
		expectedTags []localauth.AuthorizationData
		err          error
	}{
		{
			name: "Tags found",
			expectedTags: []localauth.AuthorizationData{
				{
					IdTag:     "1",
					IdTagInfo: okTag,
				},
				{
					IdTag:     "2",
					IdTagInfo: okTag,
				},
			},
		},
		{
			name: "No tags found",
		},
		{
			name: "Database error",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tagList, err := s.authList.GetTags()
			if tt.err != nil {
				s.Assert().ErrorIs(tt.err, err)
			} else {
				s.Assert().NoError(err)
				s.Assert().NotEmpty(tagList)
				s.Assert().ElementsMatch(tt.expectedTags, tagList)
			}
		})
	}
}

func (s *localAuthListTestSuite) TestSetMaxTags() {
	tests := []struct {
		name    string
		maxTags int
	}{
		{
			name:    "Max tags set",
			maxTags: 10,
		},
		{
			name:    "Negative value not permitted",
			maxTags: -1,
		},
		{
			name:    "Zero permitted",
			maxTags: 0,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s.authList.SetMaxTags(tt.maxTags)

			if tt.maxTags >= 0 {
				s.Assert().Equal(tt.maxTags, s.authList.maxTags)
			} else {
				s.Assert().NotEqual(tt.maxTags, s.authList.maxTags)
			}
		})
	}
}

func (s *localAuthListTestSuite) TestVersion() {
	tests := []struct {
		name            string
		version         int
		expectedVersion int
	}{
		{
			name:            "Version set",
			version:         10,
			expectedVersion: -1,
		},
		{
			name:            "Version not set",
			expectedVersion: -1,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Version set":
			case "Version not set":

			}
			version := s.authList.GetVersion()
			s.Assert().EqualValues(tt.expectedVersion, version)
		})
	}
}

func TestLocalAuth(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(localAuthListTestSuite))
}
