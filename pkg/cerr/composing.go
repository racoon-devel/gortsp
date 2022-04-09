package cerr

import "fmt"

type OutBufferTooShort struct {
	Actual   int
	Expected int
}

func (e OutBufferTooShort) Error() string {
	return fmt.Sprintf("outgoing buffer too short: %d < %d", e.Actual, e.Expected)
}

func NewOutBufferTooShort(actual int, expected int) error {
	return OutBufferTooShort{
		Actual:   actual,
		Expected: expected,
	}
}
