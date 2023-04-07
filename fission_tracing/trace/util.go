package trace

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"sync"
)

func TraceidToString(traceID TraceID) string {
	recv := make([]byte, 16)
	recv = traceID[:]
	return hex.EncodeToString(recv)
}

func SpanidToString(spanID SpanID) string {
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

func (rg *RandomGenerator) generateTraceID() TraceID {
	traceID := TraceID{}
	_, _ = rg.ran.Read(traceID[:])
	return traceID
}

func (rg *RandomGenerator) generateSpanID() SpanID {
	spanID := SpanID{}
	_, _ = rg.ran.Read(spanID[:])
	return spanID
}

func (rg *RandomGenerator) generateID() (TraceID, SpanID) {
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
