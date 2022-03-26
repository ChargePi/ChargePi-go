package cache

import (
	goCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var Cache *goCache.Cache

func init() {
	once := sync.Once{}
	once.Do(func() {
		log.Info("Creating cache..")
		GetCache()
	})
}

func GetCache() *goCache.Cache {
	if Cache == nil {
		Cache = goCache.New(5*time.Minute, 5*time.Minute)
	}
	return Cache
}
