package trace

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"time"

	"fission.tracing/tag"
)

var (

	// globalTracer	全局Tracer
	globalTracer = defaultTracer()
)

type Tracer struct {
	Name        string
	spanhandler *SpanHandler
}

// NewTracer
//
//	@param name
//	@return *Tracer
func NewTracer(name string) *Tracer {
	if globalTracer == nil {
		return &Tracer{
			Name:        name,
			spanhandler: NewSpanHandler(),
		}
	}
	gt := globalTracer.Load().(*Tracer)
	if gt.Name == "none" {
		gt.Name = name
		globalTracer.Store(gt)
	}
	return gt
}

// defaultTracer
//
//	@return *atomic.Value
func defaultTracer() *atomic.Value {
	v := &atomic.Value{}
	v.Store(&Tracer{
		Name:        "none",
		spanhandler: NewSpanHandler(),
	})
	return v
}

// GetSpanHandlerForTest
//
//	@receiver tr
//	@return *SpanHandler
func (tr *Tracer) GetSpanHandlerForTest() *SpanHandler {
	return tr.spanhandler
}

// Start start a new span
//
//	@receiver tr
//	@param name
//	@param ctx
//	@return context.Context
//	@return CommonSpan
func (tr *Tracer) Start(name string, ctx context.Context) (context.Context, CommonSpan) {
	if ctx == nil {
		ctx = context.Background()
	}

	// if lastSpan exists, current span should be its new child.
	if lastSpan := GetLastSpanFromContext(ctx); lastSpan.Operatorname != "none" {
		lastSpan.childSpanCount++
		ctx = InheritParentContext(ctx, lastSpan)
	}

	newSpan := tr.getNewSpan(ctx, name)
	return InheritParentContext(ctx, newSpan), newSpan
}

func (tr *Tracer) End() (int, string) {
	tr.spanhandler.OnEnd()
	list, ok := tr.ExtractSpanList()
	if !ok {
		panic("extract SpanList error.")
	}
	return len(tr.spanhandler.spanSeq), list
}

// getNewSpan
//
//	@receiver tr
//	@param ctx
//	@param name
//	@return CommonSpan
func (tr *Tracer) getNewSpan(ctx context.Context, name string) CommonSpan {
	//	get last SpanContext to inherit its traceID
	var parentSC SpanContext = GetLastSpanContextFromContext(ctx)
	var tid TraceID
	var sid SpanID
	ig := NewRandomGenerator()
	if parentSC.TraceID().Exist() {
		tid = parentSC.TraceID()
		sid = ig.generateSpanID()
	} else {
		tid, sid = ig.generateID()
	}

	// construct new SpanContextInfo
	sci := SpanContextInfo{
		TraceID:        tid,
		SpanID:         sid,
		RemotelyCalled: parentSC.remotelyCalled,
	}
	// making sure whether parent's SpanID Exists
	if parentSC.parentSpanID.Exist() {
		sci.ParentSpanID = parentSC.parentSpanID
	}
	sc := sci.ConvertToSpanContext()
	newSpan := CommonSpan{
		Operatorname:      name,
		startTime:         time.Now(),
		traceTag:          tag.NewTagDict(),
		spanContext:       sc,
		parentSpanContext: parentSC,
		childSpanCount:    0,
		spanHandler:       NewSpanHandler(),
	}
	newSpan.spanHandler = tr.spanhandler
	go newSpan.spanHandler.HandlerJob()
	return newSpan
}

// ExtractSpanList
//
//	@receiver tr
//	@return string
//	@return bool
func (tr *Tracer) ExtractSpanList() (string, bool) {
	// spanList := make([]CommonSpan, len(tr.spanhandler.spanSeq))
	output, err := json.Marshal(tr.spanhandler.spanSeq)
	if err != nil {
		return "", false
	}
	return string(output), true
}

func (tr *Tracer) InjectSpanList(spanList string) error {
	recvSpanList := make([]CommonSpan, 10)
	err := json.Unmarshal([]byte(spanList), &recvSpanList)
	if err != nil {
		return err
	}
	tr.spanhandler.spanSeq = append(tr.spanhandler.spanSeq, recvSpanList...)
	return nil
}
