package rtp

// Packet represents RTP packet
type Packet struct {
	Header  Header
	Payload []byte

	// PaddingBytes is a count of extra bytes excluding (!) padding octet
	// It has the meaning only if Header.Padding is true
	PaddingBytes uint8
}

// Size returns full serialized RTP packet entire size
func (p Packet) Size() int {
	if p.Header.Padding {
		return p.Header.Size() + len(p.Payload) + 1 + int(p.PaddingBytes)
	}
	return p.Header.Size() + len(p.Payload)
}

// Compose builds full RTP packet
func (p Packet) Compose() ([]byte, error) {
	buf := make([]byte, p.Size())
	_, err := p.ComposeTo(buf)
	return buf, err
}

// ComposeTo builds full RTP packet to a specified buffer
func (p Packet) ComposeTo(buf []byte) (int, error) {
	size := p.Size()
	if len(buf) < size {
		return 0, newErrNotEnoughBufferSpace(len(buf), size)
	}

	n, err := p.Header.ComposeTo(buf)
	if err != nil {
		return 0, err
	}

	copy(buf[n:], p.Payload)
	n += len(p.Payload)

	// write padding
	if p.Header.Padding {
		for i := 0; i < int(p.PaddingBytes); i++ {
			buf[n+i] = 0
		}
		n += int(p.PaddingBytes)
		buf[n] = p.PaddingBytes + 1
		n++
	}

	return n, nil
}

// Parse parses data buffer and fill Header's fields and Payload
func (p *Packet) Parse(data []byte) error {
	n, err := p.Header.Parse(data)
	if err != nil {
		return err
	}

	if len(data) <= n {
		return ErrPayloadIsMissing{}
	}

	if p.Header.Padding && RawPacket(data).Padding() != 0 {
		p.PaddingBytes = RawPacket(data).Padding() - 1
		padding := int(p.PaddingBytes) + 1
		if len(data) <= n+padding {
			return ErrPayloadIsMissing{}
		}
		p.Payload = make([]byte, len(data)-n-padding)
		copy(p.Payload, data[n:len(data)-padding])
	} else {
		p.Payload = make([]byte, len(data)-n)
		copy(p.Payload, data[n:])
		p.PaddingBytes = 0
	}

	return nil
}
