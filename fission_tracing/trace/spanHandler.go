package trace

import "sync"

type SpanHandler struct {
	queue   chan CommonSpan
	spanSeq []CommonSpan
	mu      sync.Mutex
	sw      sync.WaitGroup
	stopCh  chan struct{}
}

func NewSpanHandler() SpanHandler {
	sp := SpanHandler{
		queue:   make(chan CommonSpan, 10),
		spanSeq: make([]CommonSpan, 0, 10),
		stopCh:  make(chan struct{}),
	}
	go sp.HandlerJob()
	return sp
}

func (sh *SpanHandler) OnEnd() {
	if sh.stopCh == nil {
		panic("error: SpanHandler is nil.")
	}
	close(sh.stopCh)
	sh.sw.Wait()
	sh.stopCh = nil
}

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
		}
	}
}

func (sh *SpanHandler) Enqueue(cs CommonSpan) bool {
	select {
	case <-sh.stopCh:
		return false
	default:
	}
	sh.queue <- cs
	return true
}
