package auth

import (
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/database"
)

type authCacheTestSuite struct {
	suite.Suite
	tag        *types.IdTagInfo
	blockedTag *types.IdTagInfo
	expiredTag *types.IdTagInfo
	authCache  Cache
}

func (s *authCacheTestSuite) SetupTest() {
	db := database.Get()
	s.authCache = NewAuthCache(db)
	s.authCache.RemoveCachedTags()

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

func (s *authCacheTestSuite) TestAddTag() {
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

func (s *authCacheTestSuite) TestRemoveCachedTags() {
	s.authCache.SetMaxCachedTags(5)

	s.authCache.AddTag(s.tag.ParentIdTag, s.tag)
	s.authCache.AddTag(s.blockedTag.ParentIdTag, s.blockedTag)
	s.authCache.AddTag(s.expiredTag.ParentIdTag, s.expiredTag)

	s.authCache.RemoveCachedTags()
}

func (s *authCacheTestSuite) TestGetTag() {
	_, err := s.authCache.GetTag("")
	s.Assert().NoError(err)
}

func TestAuthCache(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(authCacheTestSuite))
}
