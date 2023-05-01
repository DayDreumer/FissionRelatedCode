package trace

type UniResponse struct {
	TraceInfo string
	status    int
	value     int
}

func (ur UniResponse) GetStatus() int {
	return ur.status
}

func (ur UniResponse) GetValue() int {
	return ur.value
}
