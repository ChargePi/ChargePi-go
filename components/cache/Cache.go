package cache

import (
	goCache "github.com/patrickmn/go-cache"
	"log"
	"sync"
	"time"
)

var Cache *goCache.Cache

func init() {
	once := sync.Once{}
	once.Do(func() {
		log.Println("Creating cache..")
		GetCache()
	})
}

func GetCache() *goCache.Cache {
	if Cache == nil {
		Cache = goCache.New(5*time.Minute, 5*time.Minute)
	}
	return Cache
}
