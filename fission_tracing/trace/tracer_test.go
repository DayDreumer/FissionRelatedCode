package trace_test

import (
	"bytes"
	"context"
	"reflect"
	"testing"
	"time"

	trace "fission.tracing"
)

func TestRootStart(t *testing.T) {
	initCtx := context.Background()
	tr := trace.NewTracer("test")
	nextCtx, span := tr.Start("Start", initCtx)
	if span.Operatorname != "Start" {
		t.Fatalf("span's Operatorname is wrong.")
	}

	if !span.SpanContext().TraceID().Exist() {
		t.Fatalf("SpanContext's TraceID doesn't exist.")
	} else {
		t.Logf("SpanContext's TraceID is %s", span.SpanContext().TraceID().String())
	}

	if !span.SpanContext().SpanID().Exist() {
		t.Fatalf("SpanContext's SpanID doesn't exist.")
	} else {
		t.Logf("SpanContext's SpanID is %s", span.SpanContext().SpanID().String())
	}

	if !reflect.DeepEqual(span.ParentSpanContext(), trace.SpanContext{}) {
		t.Fatalf("span's parentSpanContext is not empty.")
	} else {
		t.Logf("span's parentSpanContext is empty.")
	}

	if trace.GetInnerSpanForTest(nextCtx).Operatorname != span.Operatorname {
		t.Fatalf("nextCtx's Span name is wrong.")
	} else {
		t.Logf("nextCtx's span name is right, it's called %s", trace.GetInnerSpanForTest(nextCtx).Operatorname)
	}
	span.End()
	t.Logf("start time is %v, end time is  %v.", span.StartTime(), span.EndTime())
	sh := tr.GetSpanHandlerForTest()
	if sh.Len() != 1 {
		t.Fatalf("wrong number in SpanHandler, it should be 1.")
	}
}

func TestChildStart(t *testing.T) {
	initCtx := context.Background()
	tr := trace.NewTracer("test1")
	nextCtx, span := tr.Start("Start", initCtx)
	childSpan := ChildDo(nextCtx)
	time.Sleep(1 * time.Second)
	span.End()
	if childSpan.Operatorname != "ChildDo" {
		t.Fatalf("child span's Operatorname is wrong.")
	}
	childTid, rootTid := childSpan.SpanContext().TraceID(), span.SpanContext().TraceID()
	if !bytes.Equal(childTid[:], rootTid[:]) {
		t.Fatalf("child's TraceID is wrong.")
	}
	if !childSpan.StartTime().After(span.StartTime()) || !span.EndTime().After(childSpan.EndTime()) {
		t.Fatalf("wrong time in childSpan.")
	}
	n := tr.End()
	t.Logf("number of span is %d", n)
}

func ChildDo(ctx context.Context) trace.CommonSpan {
	tr := trace.NewTracer("test2")
	_, span := tr.Start("ChildDo", ctx)
	// t.Logf("[test]child's startTime is %v", span.StartTime())
	/*
		wrong use at 'defer span.End()'
		defer logic:
			final -> return x
			x = span
			span.End()
			return x
	*/
	span.End()
	return span
}
