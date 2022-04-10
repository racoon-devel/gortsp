package rtp

import (
	"math"
)

// ExtensionHeader represents RTP extension header
type ExtensionHeader struct {
	Profile uint16
	Content []byte
}

// Header represents RTP packet header fields (RFC3550)
type Header struct {
	Marker         bool
	Padding        bool
	PayloadType    uint8
	SequenceNumber uint16
	Timestamp      uint32
	SSRC           uint32
	CSRC           []uint32
	Extension      *ExtensionHeader
}

// Size returns entire size of RTP header according to fields
func (h Header) Size() int {
	size := HeaderLength + len(h.CSRC)*4
	if h.Extension != nil {
		size += HeaderExtensionLength + int(math.Ceil(float64(len(h.Extension.Content))/4.))*4
	}
	return size
}

// Parse validates RTP header and fills fields
func (h *Header) Parse(data []byte) (n int, err error) {
	p := RawPacket(data)

	expected := HeaderLength
	if len(p) < expected {
		err = newErrIncompleteHeader(len(p), expected)
		return
	}

	version := p.Version()
	if version != Version {
		err = ErrVersionMismatch{version}
		return
	}

	csrcCount := int(p.CC())
	expected += csrcCount * 4

	if p.X() {
		expected += HeaderExtensionLength
	}

	if len(p) < expected {
		err = newErrIncompleteHeader(len(p), expected)
		return
	}

	if p.X() {
		expected += int(p.ExtensionLength()) * 4
	}

	if len(p) < expected {
		err = newErrIncompleteHeader(len(p), expected)
		return
	}

	h.Padding = p.P()
	h.Marker = p.M()
	h.PayloadType = p.PT()
	h.SequenceNumber = p.Seq()
	h.Timestamp = p.Timestamp()
	h.SSRC = p.SSRC()
	if csrcCount != 0 {
		h.CSRC = make([]uint32, csrcCount)
		for i := 0; i < csrcCount; i++ {
			h.CSRC[i] = p.CSRC(i)
		}
	} else {
		h.CSRC = nil
	}

	if p.X() {
		h.Extension = &ExtensionHeader{
			Profile: p.ExtensionProfile(),
			Content: p.ExtensionHeader(),
		}
	} else {
		h.Extension = nil
	}

	n = expected
	return
}

// Compose builds an RTP packet header
func (h Header) Compose() ([]byte, error) {
	buf := make([]byte, h.Size())
	_, err := h.ComposeTo(buf)
	return buf, err
}

// ComposeTo builds an RTP packet header to specified buffer
func (h Header) ComposeTo(buf []byte) (n int, err error) {
	if len(buf) < HeaderLength {
		err = newErrNotEnoughBufferSpace(len(buf), HeaderLength)
		return
	}

	if len(h.CSRC) > MaxCSRC {
		err = newErrCSRCLimitExceeded(len(h.CSRC))
		return
	}

	if h.PayloadType > MaxPayloadType {
		err = newErrInvalidPayloadType(h.PayloadType)
		return
	}

	size := h.Size()
	if len(buf) < size {
		err = newErrNotEnoughBufferSpace(len(buf), size)
		return
	}

	p := RawPacket(buf)

	p.SetVersion(Version)
	p.SetP(h.Padding)
	p.SetX(h.Extension != nil)
	p.SetCC(uint8(len(h.CSRC)))
	p.SetM(h.Marker)
	p.SetPT(h.PayloadType)
	p.SetSeq(h.SequenceNumber)
	p.SetTimestamp(h.Timestamp)
	p.SetSSRC(h.SSRC)
	for i, v := range h.CSRC {
		p.SetCSRC(i, v)
	}

	if h.Extension != nil {
		p.SetExtensionProfile(h.Extension.Profile)
		p.SetExtensionLength(uint16(math.Ceil(float64(len(h.Extension.Content)) / 4.)))
		p.SetExtensionHeader(h.Extension.Content)
	}

	n = size
	return
}
