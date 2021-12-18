package scheduler

import (
	"github.com/go-co-op/gocron"
	"log"
	"sync"
	"time"
)

var scheduler *gocron.Scheduler

func init() {
	once := sync.Once{}
	once.Do(func() {
		log.Println("Initializing scheduler")
		GetScheduler()
	})
}

func GetScheduler() *gocron.Scheduler {
	if scheduler == nil {
		scheduler = gocron.NewScheduler(time.UTC)
		// Set to execute jobs on first interval and not immediately
		scheduler.WaitForScheduleAll()
		// Start non-blocking
		scheduler.StartAsync()
	}

	return scheduler
}
