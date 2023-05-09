package trace

import "fission.tracing/tag"

// UniResponse 统一Response结构体
type UniResponse struct {
	TraceInfo string
	Status    int
	Value     tag.Value
}
