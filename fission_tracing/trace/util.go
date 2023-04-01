package trace

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"sync"

	"fission.tracing/trace"
)

func TraceidToString(traceID trace.TraceID) string {
	recv := make([]byte, 16)
	recv = traceID[:]
	return hex.EncodeToString(recv)
}

func SpanidToString(spanID trace.SpanID) string {
	recv := make([]byte, 8)
	recv = spanID[:]
	return hex.EncodeToString(recv)
}

type RandomGenerator struct {
	mu  sync.Mutex
	ran *rand.Rand
}

func NewRandomGenerator() RandomGenerator {
	rg := RandomGenerator{}
	var rgSeed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &rgSeed)
	rg.ran = rand.New(rand.NewSource(rgSeed))
	return rg
}

func (rg *RandomGenerator) generateTraceID() trace.TraceID {
	traceID := trace.TraceID{}
	_, _ = rg.ran.Read(traceID[:])
	return traceID
}

func (rg *RandomGenerator) generateSpanID() trace.SpanID {
	spanID := trace.SpanID{}
	_, _ = rg.ran.Read(spanID[:])
	return spanID
}

func (rg *RandomGenerator) generateID() (trace.TraceID, trace.SpanID) {
	return rg.generateTraceID(), rg.generateSpanID()
}

type Coder struct {
	Type string
	// Object json.RawMessage
	ParentSpanID SpanID
}

func (c *Coder) encode() string {
	res, err := json.Marshal(c)
	if err != nil {
		return "fail"
	}
	return string(res)
}

func (c *Coder) decode(val string) {
	json.Unmarshal([]byte(val), c)
}
