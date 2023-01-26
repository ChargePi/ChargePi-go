package auth

import (
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/suite"
)

type AuthCacheTestSuite struct {
	suite.Suite
	tag        *types.IdTagInfo
	blockedTag *types.IdTagInfo
	expiredTag *types.IdTagInfo
	authCache  *CacheImpl
}

func (s *AuthCacheTestSuite) SetupTest() {
	s.authCache = NewAuthCache(nil)
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

	overLimitTag := types.IdTagInfo{
		ParentIdTag: "testTag123",
		ExpiryDate:  types.NewDateTime(time.Now().Add(10 * time.Minute)),
		Status:      types.AuthorizationStatusAccepted,
	}

	// Test cached tag limit
	s.authCache.AddTag(overLimitTag.ParentIdTag, &overLimitTag)
}

func (s *AuthCacheTestSuite) TestRemoveCachedTags() {
	s.authCache.SetMaxCachedTags(5)

	s.authCache.AddTag(s.tag.ParentIdTag, s.tag)
	s.authCache.AddTag(s.blockedTag.ParentIdTag, s.blockedTag)
	s.authCache.AddTag(s.expiredTag.ParentIdTag, s.expiredTag)

	s.authCache.RemoveCachedTags()
}

func (s *AuthCacheTestSuite) TestRemoveTag() {
	s.authCache.SetMaxCachedTags(15)

	s.authCache.AddTag(s.tag.ParentIdTag, s.tag)
	s.authCache.AddTag(s.blockedTag.ParentIdTag, s.blockedTag)
	s.authCache.AddTag(s.expiredTag.ParentIdTag, s.expiredTag)
}

func TestAuthCache(t *testing.T) {
	suite.Run(t, new(AuthCacheTestSuite))
}
