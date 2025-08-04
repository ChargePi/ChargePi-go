package cache

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"testing"
	"time"

	"github.com/ChargePi/ChargePi-go/internal/pkg/database"
	"github.com/ChargePi/ChargePi-go/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

var (
	okTag = &types.IdTagInfo{
		ExpiryDate: types.NewDateTime(time.Now().Add(10 * time.Minute)),
		Status:     types.AuthorizationStatusAccepted,
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
	s.authCache.AddTag(tagId, okTag)

	// Test cached tag limit
	tagId = util.GenerateRandomTag()
	s.authCache.AddTag(tagId, expiredTag)
}

func (s *authCacheTestSuite) TestRemoveCachedTags() {
	tagId1 := util.GenerateRandomTag()
	s.authCache.AddTag(tagId1, okTag)

	tagId2 := util.GenerateRandomTag()
	s.authCache.AddTag(tagId2, expiredTag)

	tagId3 := util.GenerateRandomTag()
	s.authCache.AddTag(tagId3, blockedTag)

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
	s.authCache.AddTag(tagId, okTag)
	tag, err := s.authCache.GetTag(tagId)
	s.Assert().NoError(err)
	s.Assert().EqualValues(*okTag, *tag)

	tagId = util.GenerateRandomTag()
	_, err = s.authCache.GetTag(tagId)
	s.Assert().Error(err)
}

func TestAuthCache(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(authCacheTestSuite))
}
