package smartCharging

import (
	"sort"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

type ScheduleInterval struct {
	StartTime time.Time
	Duration  time.Duration
	Limit     float64
}

func CreateCompositeSchedule(connectorSchedules []*types.ChargingProfile) []ScheduleInterval {
	// Map to store the minimum limit for each time interval
	minLimitMap := make(map[time.Time]float64)

	// Loop through all the charging profiles for the connector
	for _, profile := range connectorSchedules {
		if profile == nil {
			continue
		}

		if profile.ValidFrom.After(time.Now()) || profile.ValidTo.Before(time.Now()) {
			continue
		}

		startTime := time.Now()

		// Loop through all the schedule intervals in the profile
		for _, period := range profile.ChargingSchedule.ChargingSchedulePeriod {
			endTime := startTime.Add(time.Duration(period.StartPeriod))

			// Check if the current period limit is less than the stored minimum limit
			if minLimit, exists := minLimitMap[startTime]; exists {
				if period.Limit < minLimit {
					minLimitMap[startTime] = period.Limit
				}
			} else {
				minLimitMap[startTime] = period.Limit
			}

			// Repeat the same process for the end time of the period
			if minLimit, exists := minLimitMap[endTime]; exists {
				if period.Limit < minLimit {
					minLimitMap[endTime] = period.Limit
				}
			} else {
				minLimitMap[endTime] = period.Limit
			}
		}
	}

	return toScheduleInterval(minLimitMap)
}

// Convert the map to a slice of schedule intervals
func toScheduleInterval(minLimitMap map[time.Time]float64) []ScheduleInterval {
	var compositeSchedule []ScheduleInterval
	var startTimes []time.Time

	for startTime := range minLimitMap {
		startTimes = append(startTimes, startTime)
	}
	sort.Slice(startTimes, func(i, j int) bool { return startTimes[i].Before(startTimes[j]) })

	for i := 0; i < len(startTimes)-1; i++ {
		startTime := startTimes[i]
		endTime := startTimes[i+1]
		limit := minLimitMap[startTime]
		duration := endTime.Sub(startTime)
		compositeSchedule = append(compositeSchedule, ScheduleInterval{
			StartTime: startTime,
			Duration:  duration,
			Limit:     limit,
		})
	}

	return compositeSchedule
}
