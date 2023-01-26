package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
)

func NewScheduler() *gocron.Scheduler {
	scheduler := gocron.NewScheduler(time.UTC)
	// Set to execute jobs on first interval and not immediately
	scheduler.WaitForScheduleAll()
	// Start non-blocking
	scheduler.StartAsync()

	return scheduler
}
