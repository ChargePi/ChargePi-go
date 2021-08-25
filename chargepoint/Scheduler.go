package chargepoint

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
		scheduler = gocron.NewScheduler(time.UTC)
		//Set to execute jobs on first interval and not immediately
		scheduler.WaitForScheduleAll()
		// Start non-blocking
		scheduler.StartAsync()
	})
}

func GetScheduler() *gocron.Scheduler {
	return scheduler
}
