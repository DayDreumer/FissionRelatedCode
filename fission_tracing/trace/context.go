package trace

import "context"

type contextKeyType int

const spanKey contextKeyType = iota

func InheritParentContext(ctx context.Context, currentSpan CommonSpan) context.Context {
	return context.WithValue(ctx, spanKey, currentSpan)
}

// func Inherit

func GetLastSpanFromContext(ctx context.Context) CommonSpan {
	if ctx == nil {
		return CommonSpan{}
	}
	if parentSpan, ok := ctx.Value(spanKey).(CommonSpan); ok {
		return parentSpan
	}
	return CommonSpan{
		Operatorname: "none",
	}
}
