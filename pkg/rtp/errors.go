package rtp

import "fmt"

// ErrIncompleteHeader describes error when header length is less than required
type ErrIncompleteHeader struct {
	Expected int
	Actual   int
}

func (e ErrIncompleteHeader) Error() string {
	return fmt.Sprintf("incoming buffer too short: %d < %d", e.Actual, e.Expected)
}

func newErrIncompleteHeader(expected int, actual int) error {
	return ErrIncompleteHeader{
		Expected: expected,
		Actual:   actual,
	}
}

// ErrVersionMismatch describes error when packet has an unknown version
type ErrVersionMismatch struct {
	Version uint8
}

func (e ErrVersionMismatch) Error() string {
	return fmt.Sprintf("RTP version mismatch: %d != %d", e.Version, Version)
}

// ErrNotEnoughBufferSpace describes error when buffer length is not enough for building packet
type ErrNotEnoughBufferSpace struct {
	Expected int
	Actual   int
}

func (e ErrNotEnoughBufferSpace) Error() string {
	return fmt.Sprintf("not enough buffer space: %d < %d", e.Actual, e.Expected)
}

func newErrNotEnoughBufferSpace(expected int, actual int) error {
	return ErrNotEnoughBufferSpace{
		Expected: expected,
		Actual:   actual,
	}
}

// ErrInvalidPayloadType happens when user set PT > MaxPayloadType
type ErrInvalidPayloadType struct {
	PayloadType uint8
}

func (e ErrInvalidPayloadType) Error() string {
	return fmt.Sprintf("invalid payload type: %d > %d", e.PayloadType, MaxPayloadType)
}

func newErrInvalidPayloadType(pt uint8) error {
	return ErrInvalidPayloadType{
		PayloadType: pt,
	}
}

// ErrCSRCLimitExceeded happens when len(CSRC) > MaxCSRC
type ErrCSRCLimitExceeded struct {
	Count int
}

func (e ErrCSRCLimitExceeded) Error() string {
	return fmt.Sprintf("CSRC limit exceeded: %d / %d", e.Count, MaxCSRC)
}

func newErrCSRCLimitExceeded(count int) error {
	return ErrCSRCLimitExceeded{
		count,
	}
}

// ErrPayloadIsMissing happens if incoming data packet hasn't payload
type ErrPayloadIsMissing struct {
}

func (e ErrPayloadIsMissing) Error() string {
	return "payload is missing"
}
