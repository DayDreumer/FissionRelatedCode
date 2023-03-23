package trace

import (
	"time"
)

type FormalSpan struct {
	name string

	// startTime is the time at which this span was started.
	startTime time.Time

	// endTime is the time at which this span was ended.
	endTime time.Time

	// status is the status of this span.
	// status Status

	childSpanCount int

	tracer *tracer
}