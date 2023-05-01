package trace

type ResultType interface {
	GetStatus() int
	GetResult() int
}

type UniResponse struct {
	TraceInfo string
	Result    ResultType
}
