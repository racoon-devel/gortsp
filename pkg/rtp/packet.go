package rtp

import (
	"github.com/racoon-devel/gortsp/pkg/cerr"
)

// Packet represents RTP packet
type Packet struct {
	Header  Header
	Payload []byte
}

func (p Packet) Size() int {
	return p.Header.Size() + len(p.Payload)
}

// Compose builds full RTP packet
func (p Packet) Compose() ([]byte, error) {
	buf := make([]byte, p.Size())
	return buf, p.ComposeTo(buf)
}

// ComposeTo builds full RTP packet to a specified buffer
func (p Packet) ComposeTo(buf []byte) error {
	size := p.Size()
	if len(buf) < size {
		return cerr.NewOutBufferTooShort(len(buf), size)
	}

	if err := p.Header.ComposeTo(buf); err != nil {
		return err
	}

	copy(buf[size:], p.Payload)

	return nil
}

// Parse parses data buffer and fill Header's fields and Payload
func (p *Packet) Parse(data []byte) error {
	if err := p.Header.Parse(data); err != nil {
		return err
	}

	p.Payload = data[p.Header.Size():]

	return nil
}
