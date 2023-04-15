package auth

import (
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/pkg/util"
)

type tagManagerTestSuite struct {
	suite.Suite
}

func (s *tagManagerTestSuite) SetupTest() {}

func (s *tagManagerTestSuite) TestAddTag() {
	authListMock := NewLocalAuthListMock(s.T())
	authCacheMock := NewCacheMock(s.T())

	tagId := util.GenerateRandomTag()
	authCacheMock.OnAddTag(tagId, &types.IdTagInfo{}).Return(nil)

	tagManager := &TagManagerImpl{
		authList:             authListMock,
		cache:                authCacheMock,
		authCacheEnabled:     true,
		localAuthListEnabled: false,
	}

	err := tagManager.AddTag(tagId, &types.IdTagInfo{})
	s.Assert().NoError(err)
}

func (s *tagManagerTestSuite) TestGetTag() {
	authListMock := NewLocalAuthListMock(s.T())
	authCacheMock := NewCacheMock(s.T())

	tagId := util.GenerateRandomTag()
	authCacheMock.OnGetTag(tagId).Return(nil, nil)

	tagManager := &TagManagerImpl{
		authList:             authListMock,
		cache:                authCacheMock,
		authCacheEnabled:     true,
		localAuthListEnabled: false,
	}

	tagInfo, err := tagManager.GetTag(tagId)
	s.Assert().NoError(err)
	s.Assert().NotNil(tagInfo)
}

func (s *tagManagerTestSuite) TestGetTags() {
	authListMock := NewLocalAuthListMock(s.T())
	authCacheMock := NewCacheMock(s.T())

	// authCacheMock().Return([]localauth.AuthorizationData{})

	tagManager := &TagManagerImpl{
		authList:             authListMock,
		cache:                authCacheMock,
		authCacheEnabled:     true,
		localAuthListEnabled: false,
	}

	tags := tagManager.GetTags()
	s.Assert().NotEmpty(tags)
	s.Assert().Len(tags, 1)
}

func (s *tagManagerTestSuite) TestRemoveTag() {
	authListMock := NewLocalAuthListMock(s.T())
	authCacheMock := NewCacheMock(s.T())

	// authCacheMock.OnRemoveTag("").Return(nil)

	tagManager := &TagManagerImpl{
		authList:             authListMock,
		cache:                authCacheMock,
		authCacheEnabled:     true,
		localAuthListEnabled: true,
	}

	err := tagManager.RemoveTag("")
	s.Assert().NoError(err)
}

func (s *tagManagerTestSuite) TestClearCache() {
	authListMock := NewLocalAuthListMock(s.T())
	authCacheMock := NewCacheMock(s.T())

	tagManager := &TagManagerImpl{
		authList:             authListMock,
		cache:                authCacheMock,
		authCacheEnabled:     true,
		localAuthListEnabled: false,
	}

	tagManager.ClearCache()

	tags := tagManager.GetTags()
	s.Assert().Empty(tags)
}

func (s *tagManagerTestSuite) TestUpdateLocalAuthList() {
	authListMock := NewLocalAuthListMock(s.T())
	authCacheMock := NewCacheMock(s.T())

	tagManager := &TagManagerImpl{
		authList:             authListMock,
		cache:                authCacheMock,
		authCacheEnabled:     true,
		localAuthListEnabled: false,
	}

	tagManager.AddTag("", nil)
}

func (s *tagManagerTestSuite) TestSetMaxTags() {
	authListMock := NewLocalAuthListMock(s.T())
	authCacheMock := NewCacheMock(s.T())

	tagManager := &TagManagerImpl{
		authList:             authListMock,
		cache:                authCacheMock,
		authCacheEnabled:     true,
		localAuthListEnabled: false,
	}

	tagManager.SetMaxTags(1)
}

func TestTagManager(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(tagManagerTestSuite))
}
