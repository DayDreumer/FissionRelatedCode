package trace

import (
	"math/rand"
	"sync"

	"github.com/DayDreumer/FissionRelatedCode/fission_tracing/trace"
)

type randomGenerator struct {
	mu  sync.Mutex
	ran rand.Rand
}

func (rg *randomGenerator) generateTraceID() trace.TraceID {
	traceID := trace.TraceID{}
	_, _ = rg.ran.Read(traceID[:])
	return traceID
}

func (rg *randomGenerator) generateSpanID() trace.SpanID {
	spanID := trace.SpanID{}
	_, _ = rg.ran.Read(spanID[:])
	return spanID
}

func (rg *randomGenerator) generateID() (trace.TraceID, trace.SpanID) {
	return rg.generateTraceID(), rg.generateSpanID()
}
