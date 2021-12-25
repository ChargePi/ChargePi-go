package auth

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/suite"
	cache2 "github.com/xBlaz3kx/ChargePi-go/components/cache"
	"testing"
	"time"
)

type AuthCacheTestSuite struct {
	suite.Suite
	tag        *types.IdTagInfo
	blockedTag *types.IdTagInfo
	expiredTag *types.IdTagInfo
}

func (s *AuthCacheTestSuite) SetupTest() {
	cache2.Cache = cache.New(time.Minute*10, time.Minute*10)

	authCache = GetAuthCache()
	authCache.Set("AuthCacheVersion", 1, cache.NoExpiration)
	s.tag = &types.IdTagInfo{
		ParentIdTag: "123",
		ExpiryDate:  types.NewDateTime(time.Now().Add(10 * time.Minute)),
		Status:      types.AuthorizationStatusAccepted,
	}

	s.blockedTag = &types.IdTagInfo{
		ParentIdTag: "BlockedTag123",
		ExpiryDate:  types.NewDateTime(time.Now().Add(40 * time.Minute)),
		Status:      types.AuthorizationStatusBlocked,
	}

	s.expiredTag = &types.IdTagInfo{
		ParentIdTag: "ExpiredTag123",
		ExpiryDate:  types.NewDateTime(time.Date(1999, 1, 1, 1, 1, 1, 0, time.Local)),
		Status:      types.AuthorizationStatusAccepted,
	}

}

func (s *AuthCacheTestSuite) TestAddTag() {
	SetMaxCachedTags(1)
	AddTag(s.tag.ParentIdTag, s.tag)

	cachedTag, isFound := authCache.Get(fmt.Sprintf("AuthTag%s", s.tag.ParentIdTag))
	s.Require().True(isFound)
	s.Require().Equal(s.tag, cachedTag)

	overLimitTag := types.IdTagInfo{
		ParentIdTag: "1234",
		ExpiryDate:  types.NewDateTime(time.Now().Add(10 * time.Minute)),
		Status:      types.AuthorizationStatusAccepted,
	}

	// Test cached tag limit
	AddTag(overLimitTag.ParentIdTag, &overLimitTag)

	_, isFound = authCache.Get(fmt.Sprintf("AuthTag%s", overLimitTag.ParentIdTag))
	s.Require().False(isFound)
}

func (s *AuthCacheTestSuite) TestIsTagAuthorized() {
	SetMaxCachedTags(5)

	AddTag(s.tag.ParentIdTag, s.tag)
	AddTag(s.blockedTag.ParentIdTag, s.blockedTag)
	AddTag(s.expiredTag.ParentIdTag, s.expiredTag)

	s.Require().True(IsTagAuthorized(s.tag.ParentIdTag))
	s.Require().False(IsTagAuthorized(s.blockedTag.ParentIdTag))
	s.Require().False(IsTagAuthorized(s.expiredTag.ParentIdTag))
}

func (s *AuthCacheTestSuite) TestRemoveCachedTags() {
	SetMaxCachedTags(5)

	AddTag(s.tag.ParentIdTag, s.tag)
	AddTag(s.blockedTag.ParentIdTag, s.blockedTag)
	AddTag(s.expiredTag.ParentIdTag, s.expiredTag)

	RemoveCachedTags()

	s.Require().Equal(2, authCache.ItemCount())
}

func (s *AuthCacheTestSuite) TestRemoveTag() {
	SetMaxCachedTags(15)

	AddTag(s.tag.ParentIdTag, s.tag)
	AddTag(s.blockedTag.ParentIdTag, s.blockedTag)
	AddTag(s.expiredTag.ParentIdTag, s.expiredTag)

	RemoveTag(s.blockedTag.ParentIdTag)

	_, isFound := authCache.Get(fmt.Sprintf("AuthTag%s", s.blockedTag.ParentIdTag))
	s.Require().False(isFound)

	_, isFound = authCache.Get(fmt.Sprintf("AuthTag%s", s.tag.ParentIdTag))
	s.Require().True(isFound)

	RemoveTag("AuthTag1234")
	_, isFound = authCache.Get(fmt.Sprintf("AuthTag%s", s.tag.ParentIdTag))
	s.Require().True(isFound)
}

func (s *AuthCacheTestSuite) TestSetMaxCachedTags() {
	SetMaxCachedTags(1)
	numCachedTags, isFound := authCache.Get("AuthCacheMaxTags")
	s.Require().True(isFound)
	s.Require().Equal(1, numCachedTags)

	SetMaxCachedTags(2)
	numCachedTags, isFound = authCache.Get("AuthCacheMaxTags")
	s.Require().True(isFound)
	s.Require().Equal(2, numCachedTags)

	SetMaxCachedTags(-1)
	numCachedTags, isFound = authCache.Get("AuthCacheMaxTags")
	s.Require().True(isFound)
	s.Require().Equal(2, numCachedTags)

	SetMaxCachedTags(0)
	numCachedTags, isFound = authCache.Get("AuthCacheMaxTags")
	s.Require().True(isFound)
	s.Require().Equal(2, numCachedTags)
}

func TestAuthCache(t *testing.T) {
	suite.Run(t, new(AuthCacheTestSuite))
}
