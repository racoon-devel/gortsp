package rtsp

const (
	incomingItemsCapacity = 100
)

type request struct {
	req  *Request
	resp chan interface{}
	seq  uint64
}
