package trace

import "sync"

type SpanHandler struct {
	queue   chan CommonSpan
	spanSeq []CommonSpan
	mu      sync.Mutex
	sw      sync.WaitGroup
	stopCh  chan struct{}
}

// NewSpanHandler
//  @return *SpanHandler
func NewSpanHandler() *SpanHandler {
	sp := &SpanHandler{
		queue:   make(chan CommonSpan, 10),
		spanSeq: make([]CommonSpan, 0, 10),
		stopCh:  make(chan struct{}),
	}
	go sp.HandlerJob()
	return sp
}

// OnEnd
//  @receiver sh
func (sh *SpanHandler) OnEnd() {
	if sh.stopCh == nil {
		panic("error: SpanHandler is nil.")
	}
	close(sh.stopCh)
	sh.sw.Wait()
	sh.stopCh = nil
}

// HandlerJob
//  @receiver sh
func (sh *SpanHandler) HandlerJob() {
	sh.sw.Add(1)
	for {
		select {
		case <-sh.stopCh:
			sh.sw.Done()
			return
		case s := <-sh.queue:
			sh.mu.Lock()
			//	Do Task on Span
			sh.spanSeq = append(sh.spanSeq, s)
			sh.mu.Unlock()
		}
	}
}

// Enqueue
//  @receiver sh
//  @param cs
//  @return bool
func (sh *SpanHandler) Enqueue(cs CommonSpan) bool {
	select {
	case <-sh.stopCh:
		return false
	default:
	}
	sh.queue <- cs
	return true
}

// Len
//  @receiver sh
//  @return int
func (sh *SpanHandler) Len() int {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	return len(sh.spanSeq)
}
