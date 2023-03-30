package trace

import (
	"math/rand"
	"sync"
)

type randomGenerator struct {
	mu  sync.Mutex
	ran rand.Rand
}

func (rg *randomGenerator) generateTraceID() ([]byte, error) {
	// var traceID TraceID

}
