package auth

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type AuthCacheTestSuite struct {
	suite.Suite
	tag        *types.IdTagInfo
	blockedTag *types.IdTagInfo
	expiredTag *types.IdTagInfo
	authCache  *Cache
}

func (s *AuthCacheTestSuite) SetupTest() {
	s.authCache = NewAuthCache("./auth.json")
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
	s.authCache.SetMaxCachedTags(1)
	s.authCache.AddTag(s.tag.ParentIdTag, s.tag)

	s.Require().True(s.authCache.IsTagAuthorized(s.tag.ParentIdTag))

	overLimitTag := types.IdTagInfo{
		ParentIdTag: "testTag123",
		ExpiryDate:  types.NewDateTime(time.Now().Add(10 * time.Minute)),
		Status:      types.AuthorizationStatusAccepted,
	}

	// Test cached tag limit
	s.authCache.AddTag(overLimitTag.ParentIdTag, &overLimitTag)
	s.Require().False(s.authCache.IsTagAuthorized(overLimitTag.ParentIdTag))
}

func (s *AuthCacheTestSuite) TestIsTagAuthorized() {
	s.authCache.SetMaxCachedTags(5)

	s.authCache.AddTag(s.tag.ParentIdTag, s.tag)
	s.authCache.AddTag(s.blockedTag.ParentIdTag, s.blockedTag)
	s.authCache.AddTag(s.expiredTag.ParentIdTag, s.expiredTag)

	s.Require().True(s.authCache.IsTagAuthorized(s.tag.ParentIdTag))
	s.Require().False(s.authCache.IsTagAuthorized(s.blockedTag.ParentIdTag))
	s.Require().False(s.authCache.IsTagAuthorized(s.expiredTag.ParentIdTag))
}

func (s *AuthCacheTestSuite) TestRemoveCachedTags() {
	s.authCache.SetMaxCachedTags(5)

	s.authCache.AddTag(s.tag.ParentIdTag, s.tag)
	s.authCache.AddTag(s.blockedTag.ParentIdTag, s.blockedTag)
	s.authCache.AddTag(s.expiredTag.ParentIdTag, s.expiredTag)

	s.authCache.RemoveCachedTags()

	s.Require().Equal(2, s.authCache.cache.ItemCount())
}

func (s *AuthCacheTestSuite) TestRemoveTag() {
	s.authCache.SetMaxCachedTags(15)

	s.authCache.AddTag(s.tag.ParentIdTag, s.tag)
	s.authCache.AddTag(s.blockedTag.ParentIdTag, s.blockedTag)
	s.authCache.AddTag(s.expiredTag.ParentIdTag, s.expiredTag)

	s.authCache.RemoveTag(s.blockedTag.ParentIdTag)

	_, isFound := s.authCache.cache.Get(fmt.Sprintf("AuthTag%s", s.blockedTag.ParentIdTag))
	s.Require().False(isFound)

	_, isFound = s.authCache.cache.Get(fmt.Sprintf("AuthTag%s", s.tag.ParentIdTag))
	s.Require().True(isFound)

	s.authCache.RemoveTag("AuthTag1234")
	_, isFound = s.authCache.cache.Get(fmt.Sprintf("AuthTag%s", s.tag.ParentIdTag))
	s.Require().True(isFound)
}

func (s *AuthCacheTestSuite) TestSetMaxCachedTags() {
	s.authCache.SetMaxCachedTags(1)
	numCachedTags, isFound := s.authCache.cache.Get(MaxTagsKey)
	s.Require().True(isFound)
	s.Require().Equal(1, numCachedTags)

	s.authCache.SetMaxCachedTags(2)
	numCachedTags, isFound = s.authCache.cache.Get(MaxTagsKey)
	s.Require().True(isFound)
	s.Require().Equal(2, numCachedTags)

	s.authCache.SetMaxCachedTags(-1)
	numCachedTags, isFound = s.authCache.cache.Get(MaxTagsKey)
	s.Require().True(isFound)
	s.Require().Equal(2, numCachedTags)

	s.authCache.SetMaxCachedTags(0)
	numCachedTags, isFound = s.authCache.cache.Get(MaxTagsKey)
	s.Require().True(isFound)
	s.Require().Equal(2, numCachedTags)
}

func TestAuthCache(t *testing.T) {
	suite.Run(t, new(AuthCacheTestSuite))
}
