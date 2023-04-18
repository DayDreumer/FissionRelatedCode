package trace

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"time"

	"fission.tracing/tag"
)

type TraceID [16]byte

var noneTID TraceID

func (t TraceID) Exist() bool {
	return !bytes.Equal(t[:], noneTID[:])
}

func (t TraceID) String() string {
	tid := make([]byte, 16)
	tid = t[:]
	return hex.EncodeToString(tid)
}

type SpanID [8]byte

var noneSID SpanID

func (s SpanID) Exist() bool {
	return !bytes.Equal(s[:], noneSID[:])
}

func (s SpanID) String() string {
	sid := make([]byte, 8)
	sid = s[:]
	return hex.EncodeToString(sid)
}

type SpanContextInfo struct {
	TraceID        TraceID
	SpanID         SpanID
	ParentSpanID   SpanID // optional
	RemotelyCalled bool
}

func (sci SpanContextInfo) ConvertToSpanContext() SpanContext {
	return SpanContext{
		traceID:        sci.TraceID,
		spanID:         sci.SpanID,
		parentSpanID:   sci.ParentSpanID,
		remotelyCalled: sci.RemotelyCalled,
	}
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

func (sc SpanContext) TraceID() TraceID {
	return sc.traceID
}

func (sc SpanContext) SpanID() SpanID {
	return sc.spanID
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
	traceTag *tag.TagDict

	// spanContext is used to show relationships about childs Or followers.
	spanContext SpanContext

	// parent's SpanContext of current span.
	parentSpanContext SpanContext

	// number of current span's childs.
	childSpanCount int

	spanHandler *SpanHandler
	// tracer *Tracer
}

func (cs *CommonSpan) initSpanContext(parentSpanID SpanID) {
	if parentSpanID != (SpanID{}) {
		cs.spanContext.parentSpanID = parentSpanID
	}

}

func (cs *CommonSpan) StartTime() time.Time {
	return cs.startTime
}

func (cs *CommonSpan) EndTime() time.Time {
	return cs.endTime
}

func (cs *CommonSpan) SpanContext() SpanContext {
	return cs.spanContext
}

func (cs *CommonSpan) ParentSpanContext() SpanContext {
	return cs.parentSpanContext
}

func (cs *CommonSpan) AddTag(key tag.Key, value tag.Value) {
	cs.traceTag.Insert(key, value)
}

// func (cs *CommonSpan) ShowTag()

func (cs CommonSpan) backToTracer() {
	cs.spanHandler.Enqueue(cs)
}

func (cs *CommonSpan) End() {
	cs.endTime = GetEndTime(cs.StartTime())
	cs.backToTracer()
}
