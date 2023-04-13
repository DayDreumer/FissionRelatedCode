package trace

import (
	"context"
	"time"

	"fission.tracing/tag"
)

type Tracer struct {
	Name string
}

// start a new span
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
	}
	return newSpan
}

func (tr *Tracer) ExtractSpanContext() {

}

func (tr *Tracer) InjectSpanContext() {

}
