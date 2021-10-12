package test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/require"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"testing"
	"time"
)

func TestAddTag(t *testing.T) {
	require := require.New(t)
	data.AuthCache = cache.New(time.Minute*10, time.Minute*10)

	data.AuthCache.Set("AuthCacheVersion", 1, cache.NoExpiration)
	data.SetMaxCachedTags(1)

	tag := types.IdTagInfo{
		ParentIdTag: "123",
		ExpiryDate:  types.NewDateTime(time.Now().Add(10 * time.Minute)),
		Status:      types.AuthorizationStatusAccepted,
	}

	data.AddTag(tag.ParentIdTag, &tag)

	cachedTag, isFound := data.AuthCache.Get(fmt.Sprintf("AuthTag%s", tag.ParentIdTag))

	require.True(isFound)
	require.Equal(tag, cachedTag)

	overLimitTag := types.IdTagInfo{
		ParentIdTag: "1234",
		ExpiryDate:  types.NewDateTime(time.Now().Add(10 * time.Minute)),
		Status:      types.AuthorizationStatusAccepted,
	}

	// test cached tag limit
	data.AddTag(overLimitTag.ParentIdTag, &overLimitTag)

	_, isFound = data.AuthCache.Get(fmt.Sprintf("AuthTag%s", overLimitTag.ParentIdTag))

	require.False(isFound)
}

func TestIsTagAuthorized(t *testing.T) {
	require := require.New(t)
	data.AuthCache = cache.New(time.Minute*10, time.Minute*10)

	okTag := types.IdTagInfo{
		ParentIdTag: "OkTag123",
		ExpiryDate:  types.NewDateTime(time.Now().Add(40 * time.Minute)),
		Status:      types.AuthorizationStatusAccepted,
	}

	blockedTag := types.IdTagInfo{
		ParentIdTag: "BlockedTag123",
		ExpiryDate:  types.NewDateTime(time.Now().Add(40 * time.Minute)),
		Status:      types.AuthorizationStatusBlocked,
	}

	expiredTag := types.IdTagInfo{
		ParentIdTag: "ExpiredTag123",
		ExpiryDate:  types.NewDateTime(time.Date(1999, 1, 1, 1, 1, 1, 0, time.Local)),
		Status:      types.AuthorizationStatusAccepted,
	}

	data.SetMaxCachedTags(5)

	data.AddTag(okTag.ParentIdTag, &okTag)
	data.AddTag(blockedTag.ParentIdTag, &blockedTag)
	data.AddTag(expiredTag.ParentIdTag, &expiredTag)

	require.True(data.IsTagAuthorized(okTag.ParentIdTag))
	require.False(data.IsTagAuthorized(blockedTag.ParentIdTag))
	require.False(data.IsTagAuthorized(expiredTag.ParentIdTag))

}

func TestRemoveCachedTags(t *testing.T) {
	require := require.New(t)
	data.AuthCache = cache.New(time.Minute*10, time.Minute*10)

	okTag := types.IdTagInfo{
		ParentIdTag: "OkTag123",
		ExpiryDate:  types.NewDateTime(time.Now().Add(40 * time.Minute)),
		Status:      types.AuthorizationStatusAccepted,
	}

	blockedTag := types.IdTagInfo{
		ParentIdTag: "BlockedTag123",
		ExpiryDate:  types.NewDateTime(time.Now().Add(40 * time.Minute)),
		Status:      types.AuthorizationStatusBlocked,
	}

	expiredTag := types.IdTagInfo{
		ParentIdTag: "ExpiredTag123",
		ExpiryDate:  types.NewDateTime(time.Date(1999, 1, 1, 1, 1, 1, 0, time.Local)),
		Status:      types.AuthorizationStatusAccepted,
	}

	data.SetMaxCachedTags(5)

	data.AddTag(okTag.ParentIdTag, &okTag)
	data.AddTag(blockedTag.ParentIdTag, &blockedTag)
	data.AddTag(expiredTag.ParentIdTag, &expiredTag)

	data.RemoveCachedTags()

	require.Equal(2, data.AuthCache.ItemCount())
}

func TestRemoveTag(t *testing.T) {
	require := require.New(t)
	data.AuthCache = cache.New(time.Minute*10, time.Minute*10)

	okTag := types.IdTagInfo{
		ParentIdTag: "OkTag123",
		ExpiryDate:  types.NewDateTime(time.Now().Add(40 * time.Minute)),
		Status:      types.AuthorizationStatusAccepted,
	}

	blockedTag := types.IdTagInfo{
		ParentIdTag: "BlockedTag123",
		ExpiryDate:  types.NewDateTime(time.Now().Add(40 * time.Minute)),
		Status:      types.AuthorizationStatusBlocked,
	}

	expiredTag := types.IdTagInfo{
		ParentIdTag: "ExpiredTag123",
		ExpiryDate:  types.NewDateTime(time.Date(1999, 1, 1, 1, 1, 1, 0, time.Local)),
		Status:      types.AuthorizationStatusAccepted,
	}

	data.SetMaxCachedTags(15)

	data.AddTag(okTag.ParentIdTag, &okTag)
	data.AddTag(blockedTag.ParentIdTag, &blockedTag)
	data.AddTag(expiredTag.ParentIdTag, &expiredTag)

	data.RemoveTag(blockedTag.ParentIdTag)

	_, isFound := data.AuthCache.Get(fmt.Sprintf("AuthTag%s", blockedTag.ParentIdTag))
	require.False(isFound)

	_, isFound = data.AuthCache.Get(fmt.Sprintf("AuthTag%s", okTag.ParentIdTag))
	require.True(isFound)

	data.RemoveTag("AuthTag1234")
	_, isFound = data.AuthCache.Get(fmt.Sprintf("AuthTag%s", okTag.ParentIdTag))
	require.True(isFound)

}

func TestSetMaxCachedTags(t *testing.T) {
	require := require.New(t)
	data.AuthCache = cache.New(time.Minute*10, time.Minute*10)

	data.SetMaxCachedTags(1)
	numCachedTags, isFound := data.AuthCache.Get("AuthCacheMaxTags")
	require.True(isFound)
	require.Equal(1, numCachedTags)

	data.SetMaxCachedTags(2)
	numCachedTags, isFound = data.AuthCache.Get("AuthCacheMaxTags")
	require.True(isFound)
	require.Equal(2, numCachedTags)

	data.SetMaxCachedTags(-1)
	numCachedTags, isFound = data.AuthCache.Get("AuthCacheMaxTags")
	require.True(isFound)
	require.Equal(2, numCachedTags)

	data.SetMaxCachedTags(0)
	numCachedTags, isFound = data.AuthCache.Get("AuthCacheMaxTags")
	require.True(isFound)
	require.Equal(2, numCachedTags)

}
