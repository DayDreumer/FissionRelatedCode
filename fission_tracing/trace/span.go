package trace

import (
	"encoding/json"
	"sync"
	"time"

	"fission.tracing/tag"
)

type TraceID [16]byte

type SpanID [8]byte

type SpanContextInfo struct {
	TraceID        TraceID
	SpanID         SpanID
	ParentSpanID   SpanID // optional
	RemotelyCalled bool
}

type SpanContext struct {
	traceID        TraceID
	spanID         SpanID
	parentSpanID   SpanID // optional
	remotelyCalled bool
}

var _ json.Marshaler = SpanContext{}

func (sc *SpanContext) initWithParentSpanID(parentSpanID SpanID) {
	if parentSpanID != (SpanID{}) {
		sc.parentSpanID = parentSpanID
	}
	ig := RandomGenerator{}
	sc.traceID = ig.generateTraceID()
	sc.spanID = ig.generateSpanID()

}

func (sc SpanContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(SpanContextInfo{
		TraceID:        sc.traceID,
		SpanID:         sc.spanID,
		ParentSpanID:   sc.parentSpanID,
		RemotelyCalled: sc.remotelyCalled,
	})
}

type CommonSpan struct {
	Operatorname string

	// startTime is the time at which this span was started.
	startTime time.Time

	// endTime is the time at which this span was ended.
	endTime time.Time

	// status is the status of this span.
	// status Status

	// traceTag is used to store span tag
	tarceTag *tag.TagDict

	// spanContext is used to show relationships about childs Or followers.
	spanContext SpanContext

	// parent's SpanContext of current span.
	parentSpanContext SpanContext

	// number of current span's childs.
	childSpanCount int

	mu sync.Mutex

	tracer *tracer
}

func (cs *CommonSpan) initSpanContext(parentSpanID SpanID) {
	if parentSpanID != (SpanID{}) {
		cs.spanContext.parentSpanID = parentSpanID
	}

}
