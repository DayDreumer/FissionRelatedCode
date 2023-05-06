package test

import (
	"bytes"
	"encoding/json"
	"testing"

	trace "fission.tracing"
	tag "fission.tracing/tag"
)

const (
	timeFormart = "2006-01-02 15:04:05.000"
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

func TestSpanSerialization(t *testing.T) {
	testParentID := trace.SpanID{}
	testSpanContext := trace.SpanContext{}
	testSpanContext.InitWithParentSpanID(testParentID)
	testSpan := trace.NewSpanWithSpanContext("test", testSpanContext)
	keyArray := []tag.Key{
		"isEnd",
		"childs",
	}
	ValueArray := []tag.Value{
		tag.GetBoolValue(true),
		tag.GetIntValue(3),
	}
	for i := 0; i < 2; i++ {
		testSpan.AddTag(keyArray[i], ValueArray[i])
	}
	spanInfo, err := json.Marshal(testSpan)
	if err != nil {
		t.Fatalf("Span Serialization failed.")
	}
	// t.Logf("Span info is: %s", string(spanInfo))

	recvSpan := trace.CommonSpan{}
	err = json.Unmarshal(spanInfo, &recvSpan)
	if err != nil {
		t.Fatalf("Span Deserialization failed.")
	}
	// temp, err := json.Marshal(recvSpan)
	// t.Logf("recvSpan is: %v", string(temp))
	if recvSpan.Operatorname != testSpan.Operatorname {
		t.Fatalf("recvSpan's Operatoename is wrong.")
	}
	if rtid, tid := recvSpan.SpanContext().TraceID(), testSpan.SpanContext().TraceID(); !bytes.Equal(rtid[:], tid[:]) {
		t.Fatalf("receiver's TraceID is wrong, not equals to testSpan's TraceID.")
	}
	if rsid, sid := recvSpan.SpanContext().SpanID(), testSpan.SpanContext().SpanID(); !bytes.Equal(rsid[:], sid[:]) {
		t.Fatalf("receiver's SpanID is wrong, not equals to testSpan's SpanID.")
	}
	if recvSpan.StartTime().Format(timeFormart) != testSpan.StartTime().Format(timeFormart) {
		t.Fatalf("receiver's StartTime is wrong, not equals to testSpan's startTime.recv is %v, test is %v", recvSpan.StartTime(), testSpan.StartTime())
	}
	if recvSpan.EndTime().Format(timeFormart) != testSpan.EndTime().Format(timeFormart) {
		t.Fatalf("receiver's EndTime is wrong, not equals to testSpan's endTime.")
	}
	if v, ok := recvSpan.GetTag(keyArray[0]); ok {
		if v != ValueArray[0] {
			t.Logf("v's info is %s", v.String())
			t.Fatalf("recvSpan's tag[isEnd] is wrong.")
		}
	} else {
		t.Fatalf("Can't find tag[isEnd] in recvSpan Tag.")
	}
	if v, ok := recvSpan.GetTag(keyArray[1]); ok {
		if v != ValueArray[1] {
			t.Fatalf("recvSpan's tag[childs] is wrong.")
		}
	} else {
		t.Fatalf("Can't find tag[childs] in recvSpan Tag.")
	}
}
