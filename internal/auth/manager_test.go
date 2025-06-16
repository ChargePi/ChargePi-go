package auth

import (
	"testing"

	mock_cache "github.com/ChargePi/ChargePi-go/gen/mocks/auth/cache"
	mock_list "github.com/ChargePi/ChargePi-go/gen/mocks/auth/list"
	"github.com/ChargePi/ChargePi-go/pkg/util"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type tagManagerTestSuite struct {
	suite.Suite
	authListMock  *mock_list.MockLocalAuthList
	authCacheMock *mock_cache.MockCache

	tagManager *ManagerV1
}

func (s *tagManagerTestSuite) SetupTest() {
	s.authListMock = mock_list.NewMockLocalAuthList(s.T())
	s.authCacheMock = mock_cache.NewMockCache(s.T())

	s.tagManager = &ManagerV1{
		authList: s.authListMock,
		cache:    s.authCacheMock,
	}
}

func (s *tagManagerTestSuite) TestAddTag() {
	tagId := util.GenerateRandomTag()
	s.authCacheMock.EXPECT().AddTag(tagId, &types.IdTagInfo{}).Return()

	s.tagManager.authCacheEnabled = true
	s.tagManager.localAuthListEnabled = false

	err := s.tagManager.AddTag(tagId, &types.IdTagInfo{})
	s.Assert().NoError(err)
}

func (s *tagManagerTestSuite) TestGetTag() {
	tagId := util.GenerateRandomTag()
	s.authCacheMock.EXPECT().GetTag(tagId).Return(nil, nil)

	s.tagManager.authCacheEnabled = true
	s.tagManager.localAuthListEnabled = false

	tagInfo, err := s.tagManager.GetTag(tagId)
	s.Assert().NoError(err)
	s.Assert().NotNil(tagInfo)
}

func (s *tagManagerTestSuite) TestGetTags() {

	s.tagManager.authCacheEnabled = true
	s.tagManager.localAuthListEnabled = false

	tags := s.tagManager.GetTags()
	s.Assert().NotEmpty(tags)
	s.Assert().Len(tags, 1)
}

func (s *tagManagerTestSuite) TestRemoveTag() {
	authListMock := mock_list.NewMockLocalAuthList(s.T())
	authCacheMock := mock_cache.NewMockCache(s.T())

	// authCacheMock.OnRemoveTag("").Return(nil)

	tagManager := &ManagerV1{
		authList:             authListMock,
		cache:                authCacheMock,
		authCacheEnabled:     true,
		localAuthListEnabled: true,
	}

	err := tagManager.RemoveTag("")
	s.Assert().NoError(err)
}

func (s *tagManagerTestSuite) TestClearCache() {
	err := s.tagManager.ClearCache()
	s.Assert().NoError(err)

	tags := s.tagManager.GetTags()
	s.Assert().Empty(tags)
}

func (s *tagManagerTestSuite) TestUpdateLocalAuthList() {

	_ = s.tagManager.AddTag("", nil)
}

func (s *tagManagerTestSuite) TestSetMaxTags() {
	s.tagManager.SetMaxTags(1)
}

func TestTagManager(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(tagManagerTestSuite))
}
