package trace_test

import (
	"context"
	"reflect"
	"testing"

	trace "fission.tracing"
)

func TestRootStart(t *testing.T) {
	initCtx := context.Background()
	tr := &trace.Tracer{
		Name: "test",
	}
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
}
