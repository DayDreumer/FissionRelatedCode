package trace

import (
	"context"
	"encoding/json"
)

type contextKeyType int

const spanKey contextKeyType = iota

// InheritParentContext inherit parent context and set spanKey's value to current span.
//
//	@param ctx
//	@param currentSpan
//	@return context.Context
func InheritParentContext(ctx context.Context, currentSpan CommonSpan) context.Context {
	return context.WithValue(ctx, spanKey, currentSpan)
}

// func Inherit

// GetLastSpanFromContext
//
//	@param ctx
//	@return CommonSpan
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

// GetLastSpanContextFromContext
//
//	@param ctx
//	@return SpanContext
func GetLastSpanContextFromContext(ctx context.Context) SpanContext {
	return GetLastSpanFromContext(ctx).spanContext
}

// GetInnerSpanForTest
//
//	@param ctx
//	@return CommonSpan
func GetInnerSpanForTest(ctx context.Context) CommonSpan {
	return ctx.Value(spanKey).(CommonSpan)
}

// InheritContextFromCaller
//
//	@param parentSpanInfo
//	@return context.Context
func InheritContextFromCaller(parentSpanInfo string) context.Context {
	var parentSpan CommonSpan
	err := json.Unmarshal([]byte(parentSpanInfo), &parentSpan)
	if err != nil {
		panic(err)
	}
	return InheritParentContext(context.Background(), parentSpan)
}
