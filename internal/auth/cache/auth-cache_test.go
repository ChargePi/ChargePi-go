package cache

import (
	"github.com/ChargePi/ChargePi-go/internal/auth"
	"testing"

	"github.com/ChargePi/ChargePi-go/internal/pkg/database"
	"github.com/ChargePi/ChargePi-go/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type authCacheTestSuite struct {
	suite.Suite
	authCache *BadgerCache
}

func (s *authCacheTestSuite) SetupTest() {
	db := database.Get()
	s.authCache = NewAuthCache(db)
	s.authCache.RemoveCachedTags()
}

func (s *authCacheTestSuite) TestAddTag() {
	s.authCache.SetMaxCachedTags(1)

	tagId := util.GenerateRandomTag()
	s.authCache.AddTag(tagId, auth.okTag)

	// Test cached tag limit
	tagId = util.GenerateRandomTag()
	s.authCache.AddTag(tagId, auth.expiredTag)
}

func (s *authCacheTestSuite) TestRemoveCachedTags() {
	tagId1 := util.GenerateRandomTag()
	s.authCache.AddTag(tagId1, auth.okTag)

	tagId2 := util.GenerateRandomTag()
	s.authCache.AddTag(tagId2, auth.expiredTag)

	tagId3 := util.GenerateRandomTag()
	s.authCache.AddTag(tagId3, auth.blockedTag)

	s.authCache.RemoveCachedTags()

	_, err := s.authCache.GetTag(tagId1)
	s.Assert().Error(err)

	_, err = s.authCache.GetTag(tagId2)
	s.Assert().Error(err)

	_, err = s.authCache.GetTag(tagId3)
	s.Assert().Error(err)
}

func (s *authCacheTestSuite) TestGetTag() {
	tagId := util.GenerateRandomTag()
	s.authCache.AddTag(tagId, auth.okTag)
	tag, err := s.authCache.GetTag(tagId)
	s.Assert().NoError(err)
	s.Assert().EqualValues(*auth.okTag, *tag)

	tagId = util.GenerateRandomTag()
	_, err = s.authCache.GetTag(tagId)
	s.Assert().Error(err)
}

func TestAuthCache(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(authCacheTestSuite))
}
