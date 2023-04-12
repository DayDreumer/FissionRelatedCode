package trace

import "context"

type tracer struct {
	name string
}

func (tr *tracer) Start(name string, ctx context.Context) (context.Context, string) {
	if ctx == nil {
		ctx = context.Background()
	}
	tr.name = name

}
