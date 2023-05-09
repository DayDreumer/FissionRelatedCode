package trace

import "time"

// GetEndTime
//  @param startTime
//  @return time.Time
func GetEndTime(startTime time.Time) time.Time {
	interval := time.Since(startTime)
	return startTime.Add(interval)
}
