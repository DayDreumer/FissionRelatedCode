package test

import (
	"bytes"
	"encoding/json"
	"testing"

	trace "fission.tracing"
)

func TestSpanContextSerialization(t *testing.T) {
	testParentID := trace.SpanID{}
	testSpanContext := trace.SpanContext{}
	testSpanContext.InitWithParentSpanID(testParentID)
	spaninfo, err := json.Marshal(testSpanContext)
	if err != nil {
		t.Fatalf("SpanContext Serialization failed.")
	} else {
		t.Logf("SpanContext info: %s", string(spaninfo))
	}
	recv_spanContext := trace.SpanContext{}
	json.Unmarshal(spaninfo, &recv_spanContext)
	if rtid, tid := recv_spanContext.TraceID(), testSpanContext.TraceID(); !bytes.Equal(rtid[:], tid[:]) {
		t.Fatalf("receiver's TraceID is wrong, not equals to testSpanContext's TraceID.")
	}
	if rsid, sid := recv_spanContext.SpanID(), testSpanContext.SpanID(); !bytes.Equal(rsid[:], sid[:]) {
		t.Fatalf("receiver's SpanID is wrong, not equals to testSpanContext's SpanID.")
	}
}
