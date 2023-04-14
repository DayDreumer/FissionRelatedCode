package trace

import "time"

func GetEndTime(startTime time.Time) time.Time {
	interval := time.Since(startTime)
	return startTime.Add(interval)
}
