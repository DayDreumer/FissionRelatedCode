package trace

import (
	"context"
	"encoding/json"
)

type contextKeyType int

const spanKey contextKeyType = iota

// inherit parent context and set spanKey's value to current span.
func InheritParentContext(ctx context.Context, currentSpan CommonSpan) context.Context {
	return context.WithValue(ctx, spanKey, currentSpan)
}

// func Inherit

func GetLastSpanFromContext(ctx context.Context) CommonSpan {
	if ctx == nil {
		return CommonSpan{
			Operatorname: "none",
		}
	}
	if parentSpan, ok := ctx.Value(spanKey).(CommonSpan); ok {
		return parentSpan
	}
	return CommonSpan{
		Operatorname: "none",
	}
}

func GetLastSpanContextFromContext(ctx context.Context) SpanContext {
	return GetLastSpanFromContext(ctx).spanContext
}

func GetInnerSpanForTest(ctx context.Context) CommonSpan {
	return ctx.Value(spanKey).(CommonSpan)
}

func InheritContextFromCaller(parentSpanInfo string) context.Context {
	var parentSpan CommonSpan
	err := json.Unmarshal([]byte(parentSpanInfo), &parentSpan)
	if err != nil {
		panic(err)
	}
	return InheritParentContext(context.Background(), parentSpan)
}
