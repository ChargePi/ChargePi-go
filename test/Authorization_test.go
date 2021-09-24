package test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/patrickmn/go-cache"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"testing"
	"time"
)

func TestAddTag(t *testing.T) {
	type args struct {
		tagId   string
		tagInfo *types.IdTagInfo
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "AddTag",
			args: args{
				tagId: "123",
				tagInfo: &types.IdTagInfo{
					ParentIdTag: "123",
					ExpiryDate:  types.NewDateTime(time.Now().Add(10 * time.Minute)),
					Status:      types.AuthorizationStatusAccepted,
				},
			},
		}, {
			name: "TagLimitReached",
			args: args{
				tagId: "123",
				tagInfo: &types.IdTagInfo{
					ParentIdTag: "123",
					ExpiryDate:  types.NewDateTime(time.Now().Add(10 * time.Minute)),
					Status:      types.AuthorizationStatusAccepted,
				},
			},
		},
	}
	data.AuthCache = cache.New(time.Minute*10, time.Minute*10)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "AddTag":
				data.SetMaxCachedTags(3)
				break
			case "TagLimitReached":
				data.SetMaxCachedTags(1)
				for i := 0; i < 5; i++ {
					data.AddTag(fmt.Sprintf("%d", i), &types.IdTagInfo{
						ExpiryDate:  nil,
						ParentIdTag: fmt.Sprintf("%d", i),
						Status:      types.AuthorizationStatusAccepted,
					})
				}
				if data.AuthCache.ItemCount()+2 != 5 {
					t.Errorf("Items added to cache regardless of set limit")
				}
				break
			}
			data.AddTag(tt.args.tagId, tt.args.tagInfo)
			cacheData, exists := data.AuthCache.Get(fmt.Sprintf("AuthTag%s", tt.args.tagId))
			if (!exists || cacheData != tt.args.tagInfo) && tt.name != "TagLimitReached" {
				t.Errorf("AddTag() did not insert the tag")
			}
		})
	}
}

func TestIsTagAuthorized(t *testing.T) {
	type args struct {
		tagId string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "TagInCacheAndAuthorized",
			args: args{
				tagId: "123",
			},
			want: true,
		}, {
			name: "TagInCacheAndNotAuthorized",
			args: args{
				tagId: "123",
			},
			want: false,
		}, {
			name: "TagNotInCache",
			args: args{
				tagId: "1234",
			},
			want: false,
		},
	}
	data.AuthCache = cache.New(time.Minute*10, time.Minute*10)
	data.SetMaxCachedTags(10)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "TagInCacheAndAuthorized":
				data.AddTag(tt.args.tagId, &types.IdTagInfo{
					ParentIdTag: "",
					ExpiryDate:  types.NewDateTime(time.Now().Add(40 * time.Minute)),
					Status:      types.AuthorizationStatusAccepted,
				})
				break
			case "TagInCacheAndNotAuthorized":
				data.AddTag(
					tt.args.tagId,
					&types.IdTagInfo{
						ParentIdTag: "",
						ExpiryDate:  types.NewDateTime(time.Now().Add(40 * time.Minute)),
						Status:      types.AuthorizationStatusBlocked,
					})
				break
			case "TagNotInCache":
				break
			}
			if got := data.IsTagAuthorized(tt.args.tagId); got != tt.want {
				t.Errorf("IsTagAuthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveCachedTags(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "JustTags"},
		{name: "OtherElements"},
		{name: "Empty"},
	}
	data.AuthCache = cache.New(time.Minute*10, time.Minute*10)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "JustTags":
				break
			case "OtherElements":
				break
			case "Empty":
				break
			}
		})
	}
}

func TestRemoveTag(t *testing.T) {
	type args struct {
		tagId string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TagInCache",
			args: args{tagId: "123"},
		}, {
			name: "TagNotInCache",
			args: args{tagId: "2134"},
		},
	}
	data.AuthCache = cache.New(time.Minute*10, time.Minute*10)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestSetMaxCachedTags(t *testing.T) {
	type args struct {
		number int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "NoTagCached",
			args: args{number: 0},
		}, {
			name: "MaxOneTagCached",
			args: args{number: 1},
		}, {
			name: "FewTagsCached",
			args: args{number: 10},
		},
	}
	data.AuthCache = cache.New(time.Minute*10, time.Minute*10)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data.RemoveCachedTags()
			data.SetMaxCachedTags(tt.args.number)
			maxTags, isFound := data.AuthCache.Get("AuthCacheMaxTags")
			if isFound {
				if maxTags != tt.args.number {
					t.Errorf("Variable not set to cache")
				}
			} else if tt.name != "NoTagCached" {
				t.Errorf("Did not find the variable in cache")
			}
		})
	}
}
