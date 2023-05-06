package trace

import "fission.tracing/tag"

type UniResponse struct {
	TraceInfo string
	Status    int
	Value     tag.Value
}
