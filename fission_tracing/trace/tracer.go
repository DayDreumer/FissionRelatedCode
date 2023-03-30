package trace

import "context"

type tracer struct {
}

func (tr *tracer) Start(name string, ctx context.Context) (context.Context, string) {

}
