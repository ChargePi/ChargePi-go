package auth

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type tagManagerTestSuite struct {
	suite.Suite
}

func (s *tagManagerTestSuite) SetupTest() {}

func (s *tagManagerTestSuite) TestAddTag() {
	authListMock := NewLocalAuthListMock(s.T())
	authCacheMock := NewCacheMock(s.T())

	tagManager := &TagManagerImpl{
		authList:             authListMock,
		cache:                authCacheMock,
		authCacheEnabled:     false,
		localAuthListEnabled: false,
	}

	err := tagManager.AddTag("", nil)
	s.Assert().NoError(err)
}

func (s *tagManagerTestSuite) TestGetTag() {
	authListMock := NewLocalAuthListMock(s.T())
	authCacheMock := NewCacheMock(s.T())

	tagManager := &TagManagerImpl{
		authList:             authListMock,
		cache:                authCacheMock,
		authCacheEnabled:     false,
		localAuthListEnabled: false,
	}

	err := tagManager.AddTag("", nil)
	s.Assert().NoError(err)
}

func (s *tagManagerTestSuite) TestGetTags() {
	authListMock := NewLocalAuthListMock(s.T())
	authCacheMock := NewCacheMock(s.T())

	tagManager := &TagManagerImpl{
		authList:             authListMock,
		cache:                authCacheMock,
		authCacheEnabled:     false,
		localAuthListEnabled: false,
	}

	err := tagManager.AddTag("", nil)
	s.Assert().NoError(err)
}

func (s *tagManagerTestSuite) TestRemoveTag() {
	authListMock := NewLocalAuthListMock(s.T())
	authCacheMock := NewCacheMock(s.T())

	tagManager := &TagManagerImpl{
		authList:             authListMock,
		cache:                authCacheMock,
		authCacheEnabled:     false,
		localAuthListEnabled: false,
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
		authCacheEnabled:     false,
		localAuthListEnabled: false,
	}

	tagManager.ClearCache()
}

func (s *tagManagerTestSuite) TestUpdateLocalAuthList() {
	authListMock := NewLocalAuthListMock(s.T())
	authCacheMock := NewCacheMock(s.T())

	tagManager := &TagManagerImpl{
		authList:             authListMock,
		cache:                authCacheMock,
		authCacheEnabled:     false,
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
		authCacheEnabled:     false,
		localAuthListEnabled: false,
	}

	tagManager.SetMaxTags(1)
}

func TestTagManager(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(tagManagerTestSuite))
}
