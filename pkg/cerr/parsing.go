package cerr

import "fmt"

type InBufferTooShort struct {
	Actual   int
	Expected int
}

func (e InBufferTooShort) Error() string {
	return fmt.Sprintf("incoming buffer too short: %d < %d", e.Actual, e.Expected)
}

func NewInBufferTooShort(actual int, expected int) error {
	return InBufferTooShort{
		Actual:   actual,
		Expected: expected,
	}
}

type ProtocolVersionMismatch struct {
	Protocol string
	Actual   int
	Expected int
}

func (e ProtocolVersionMismatch) Error() string {
	return fmt.Sprintf("%s protocol version mismatch: %d != %d", e.Protocol, e.Actual, e.Expected)
}

func NewProtocolVersionMismatch(protocol string, actual int, expected int) error {
	return ProtocolVersionMismatch{
		Protocol: protocol,
		Actual:   actual,
		Expected: expected,
	}
}
