package trace

import (
	"sync"
	"time"

	"github.com/DayDreumer/FissionRelatedCode/fission_tracing/trace"
)

type contextKeyType int

type TraceID [16]byte

type SpanID [8]byte

const spanKey contextKeyType = iota

type SpanContext struct {
	traceID        TraceID
	spanID         SpanID
	parentSpanID   SpanID // optional
	remotelyCalled bool
}

func (sc *SpanContext) init(parentSpanID SpanID) {
	if parentSpanID != (SpanID{}) {
		sc.parentSpanID = parentSpanID
	}
	ig := trace.RandomGenerator{}

}

type CommonSpan struct {
	name string

	// startTime is the time at which this span was started.
	startTime time.Time

	// endTime is the time at which this span was ended.
	endTime time.Time

	// status is the status of this span.
	// status Status

	mu sync.Mutex

	childSpanCount int

	spanContext SpanContext

	tracer *tracer
}

func (cs *CommonSpan) initSpanContext(parentSpanID SpanID) {
	if parentSpanID != (SpanID{}) {
		cs.spanContext.parentSpanID = parentSpanID
	}

}
