package trace

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"time"

	"fission.tracing/tag"
)

const (
	timeFormart = "2006-01-02 15:04:05.000"
)

type TraceID [16]byte

var noneTID TraceID

// Exist
//
//	@receiver t
//	@return bool
func (t TraceID) Exist() bool {
	return !bytes.Equal(t[:], noneTID[:])
}

// String
//
//	@receiver t
//	@return string
func (t TraceID) String() string {
	tid := make([]byte, 16)
	tid = t[:]
	return hex.EncodeToString(tid)
}

type SpanID [8]byte

var noneSID SpanID

// Exist
//
//	@receiver s
//	@return bool
func (s SpanID) Exist() bool {
	return !bytes.Equal(s[:], noneSID[:])
}

// String
//
//	@receiver s
//	@return string
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

// ConvertToSpanContext
//
//	@receiver sci
//	@return SpanContext
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

// InitWithParentSpanID
//
//	@receiver sc
//	@param parentSpanID
func (sc *SpanContext) InitWithParentSpanID(parentSpanID SpanID) {
	if parentSpanID != (SpanID{}) {
		sc.parentSpanID = parentSpanID
	}
	ig := NewRandomGenerator()
	sc.traceID = ig.generateTraceID()
	sc.spanID = ig.generateSpanID()

}

// TraceID
//
//	@receiver sc
//	@return TraceID
func (sc SpanContext) TraceID() TraceID {
	return sc.traceID
}

// SpanID
//
//	@receiver sc
//	@return SpanID
func (sc SpanContext) SpanID() SpanID {
	return sc.spanID
}

// MarshalJSON SpanContext自定义序列化
//
//	@receiver sc
//	@return []byte
//	@return error
func (sc SpanContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(SpanContextInfo{
		TraceID:        sc.traceID,
		SpanID:         sc.spanID,
		ParentSpanID:   sc.parentSpanID,
		RemotelyCalled: sc.remotelyCalled,
	})
}

// UnmarshalJSON SpanContext自定义反序列化
//
//	@receiver sc
//	@param data
//	@return error
func (sc *SpanContext) UnmarshalJSON(data []byte) error {
	sci := SpanContextInfo{}
	json.Unmarshal(data, &sci)
	sc.traceID = sci.TraceID
	sc.spanID = sci.SpanID
	sc.parentSpanID = sci.ParentSpanID
	sc.remotelyCalled = sci.RemotelyCalled
	return nil
}

type CommonSpanInfo struct {
	Operatorname            string
	StartTime               string
	EndTime                 string
	TraceTag                *tag.TagDict
	CommonSpanContext       SpanContext
	CommonSpanParentContext SpanContext
	ChildSpanCount          int
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
}

// NewSpan
//
//	@param name
//	@return CommonSpan
func NewSpan(name string) CommonSpan {
	return CommonSpan{
		Operatorname:   name,
		startTime:      time.Now(),
		traceTag:       tag.NewTagDict(),
		childSpanCount: 0,
		spanHandler:    NewSpanHandler(),
	}
}

// NewSpanWithSpanContext
//
//	@param name
//	@param sc
//	@return CommonSpan
func NewSpanWithSpanContext(name string, sc SpanContext) CommonSpan {
	cs := NewSpan(name)
	cs.spanContext = sc
	return cs
}

// NewSpanWithMultiContext
//
//	@param name
//	@param sc
//	@param psc
//	@return CommonSpan
func NewSpanWithMultiContext(name string, sc SpanContext, psc SpanContext) CommonSpan {
	cs := NewSpanWithSpanContext(name, sc)
	cs.parentSpanContext = psc
	return cs
}

// initSpanContext
//
//	@receiver cs
//	@param parentSpanID
func (cs *CommonSpan) initSpanContext(parentSpanID SpanID) {
	if parentSpanID != (SpanID{}) {
		cs.spanContext.parentSpanID = parentSpanID
	}

}

// StartTime
//
//	@receiver cs
//	@return time.Time
func (cs *CommonSpan) StartTime() time.Time {
	return cs.startTime
}

func (cs *CommonSpan) EndTime() time.Time {
	return cs.endTime
}

// SpanContext
//
//	@receiver cs
//	@return SpanContext
func (cs *CommonSpan) SpanContext() SpanContext {
	return cs.spanContext
}

// ParentSpanContext
//
//	@receiver cs
//	@return SpanContext
func (cs *CommonSpan) ParentSpanContext() SpanContext {
	return cs.parentSpanContext
}

// AddTag
//
//	@receiver cs
//	@param key
//	@param value
func (cs *CommonSpan) AddTag(key tag.Key, value tag.Value) {
	cs.traceTag.Insert(key, value)
}

// GetTag
//
//	@receiver cs
//	@param key
//	@return tag.Value
//	@return bool
func (cs *CommonSpan) GetTag(key tag.Key) (tag.Value, bool) {
	if v, ok := cs.traceTag.Search(key); ok {
		return v, true
	}
	return tag.GetNoneValue(), false
}

// backToTracer
//
//	@receiver cs
func (cs CommonSpan) backToTracer() {
	cs.spanHandler.Enqueue(cs)
}

// End
//
//	@receiver cs
func (cs *CommonSpan) End() {
	cs.endTime = GetEndTime(cs.StartTime())
	cs.backToTracer()
}

// MarshalJSON CommonSpan自定义序列化
//
//	@receiver cs
//	@return []byte
//	@return error
func (cs CommonSpan) MarshalJSON() ([]byte, error) {
	return json.Marshal(CommonSpanInfo{
		Operatorname:            cs.Operatorname,
		StartTime:               string(cs.startTime.Format(timeFormart)),
		EndTime:                 string(cs.endTime.Format(timeFormart)),
		TraceTag:                cs.traceTag,
		CommonSpanContext:       cs.spanContext,
		CommonSpanParentContext: cs.parentSpanContext,
		ChildSpanCount:          cs.childSpanCount,
	})
}

// UnmarshalJSON CommonSpan自定义反序列化
//
//	@receiver cs
//	@param data
//	@return error
func (cs *CommonSpan) UnmarshalJSON(data []byte) error {
	csi := CommonSpanInfo{}
	err := json.Unmarshal(data, &csi)
	if err != nil {
		return err
	}
	cs.Operatorname = csi.Operatorname
	cs.startTime, _ = time.ParseInLocation(timeFormart, csi.StartTime, time.Local)
	cs.endTime, _ = time.ParseInLocation(timeFormart, csi.EndTime, time.Local)
	cs.traceTag = csi.TraceTag
	cs.spanContext = csi.CommonSpanContext
	cs.parentSpanContext = csi.CommonSpanParentContext
	cs.childSpanCount = csi.ChildSpanCount
	return nil
}
