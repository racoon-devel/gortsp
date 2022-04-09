package rtp

import (
	"fmt"
	"github.com/racoon-devel/gortsp/pkg/cerr"
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
		size += extensionHeaderLength + len(h.Extension.Content)
	}
	return size
}

// Parse validates RTP header and fills fields
func (h *Header) Parse(data []byte) error {
	p := RawPacket(data)

	expected := HeaderLength
	if len(p) < expected {
		return cerr.NewInBufferTooShort(len(p), expected)
	}

	version := p.Version()
	if version != Version {
		return cerr.NewProtocolVersionMismatch("RTP", int(version), Version)
	}

	csrcCount := int(p.CC())
	expected += csrcCount * 4

	if p.X() {
		expected += extensionHeaderLength
	}

	if len(p) < expected {
		return cerr.NewInBufferTooShort(len(p), expected)
	}

	if p.X() {
		expected += int(p.ExtensionLength()) * 4
	}

	if len(p) < expected {
		return cerr.NewInBufferTooShort(len(p), expected)
	}

	h.Padding = p.P()
	h.Marker = p.M()
	h.PayloadType = p.PT()
	h.SequenceNumber = p.Seq()
	h.Timestamp = p.Timestamp()
	h.SSRC = p.SSRC()
	h.CSRC = make([]uint32, csrcCount)
	for i := 0; i < csrcCount; i++ {
		h.CSRC[i] = p.CSRC(i)
	}

	if p.X() {
		h.Extension = &ExtensionHeader{
			Profile: p.ExtensionProfile(),
			Content: p.ExtensionHeader(),
		}
	} else {
		h.Extension = nil
	}

	return nil
}

// Compose builds an RTP packet header
func (h Header) Compose() ([]byte, error) {
	buf := make([]byte, h.Size())
	return buf, h.ComposeTo(buf)
}

// ComposeTo builds an RTP packet header to specified buffer
func (h Header) ComposeTo(buf []byte) error {
	size := h.Size()
	if len(buf) < size {
		return cerr.NewOutBufferTooShort(len(buf), size)
	}

	if len(h.CSRC) > MaxCSRC {
		return fmt.Errorf("max CSRC capacity reached: %d > %d", h.CSRC, MaxCSRC)
	}

	if h.PayloadType > MaxPayloadType {
		return fmt.Errorf("invalid payload type: %d", h.PayloadType)
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
		p.SetExtensionLength(uint16(len(h.Extension.Content) / 4))
		p.SetExtensionHeader(h.Extension.Content)
	}

	return nil
}
