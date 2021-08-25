package cache

import (
	goCache "github.com/patrickmn/go-cache"
	"log"
	"sync"
	"time"
)

var Cache *goCache.Cache

//var AuthCache *goCache.Cache

func init() {
	once := sync.Once{}
	once.Do(func() {
		log.Println("Creating cache..")
		if Cache == nil {
			Cache = goCache.New(5*time.Minute, 10*time.Minute)
		}
		/*if AuthCache == nil {
			AuthCache = goCache.New(time.Minute*10, time.Minute*10)
		}*/
	})
}
