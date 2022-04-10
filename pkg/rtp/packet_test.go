package rtp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPacket_Size(t *testing.T) {
	type testCase struct {
		p    Packet
		size int
	}

	testCases := []testCase{
		{
			p: Packet{
				Header:       Header{},
				Payload:      make([]byte, 10),
				PaddingBytes: 0,
			},
			size: HeaderLength + 10,
		},
		{
			p: Packet{
				Header: Header{
					Padding: true,
				},
				Payload:      make([]byte, 10),
				PaddingBytes: 0,
			},
			size: HeaderLength + 10 + 1,
		},
		{
			p: Packet{
				Header: Header{
					Padding: true,
				},
				Payload:      make([]byte, 10),
				PaddingBytes: 1,
			},
			size: HeaderLength + 10 + 2,
		},
	}

	for i, c := range testCases {
		assert.Equal(t, c.size, c.p.Size(), "testCase : %d", i+1)
	}
}

func TestPacket_Compose(t *testing.T) {
	p := Packet{
		Header: Header{
			Padding:        true,
			PayloadType:    96,
			SequenceNumber: 9164,
			Timestamp:      1681696377,
			SSRC:           0x6b8b4567,
		},
		Payload:      []byte{0x01, 0x02, 0x03},
		PaddingBytes: 2,
	}
	buf, err := p.Compose()
	assert.NoError(t, err)
	assert.Equal(t, p.Size(), len(buf))
}

func TestPacket_ComposeTo(t *testing.T) {
	type testCase struct {
		p   Packet
		r   []byte
		buf []byte
		err error
	}

	testCases := []testCase{
		{
			p: Packet{
				Header: Header{
					PayloadType:    96,
					SequenceNumber: 9164,
					Timestamp:      1681696377,
					SSRC:           0x6b8b4567,
				},
				Payload:      []byte{0x01, 0x02, 0x03},
				PaddingBytes: 3, // MUST be ignored
			},
			r:   []byte{0x80, 0x60, 0x23, 0xcc, 0x64, 0x3c, 0xa6, 0x79, 0x6b, 0x8b, 0x45, 0x67, 0x01, 0x02, 0x03},
			buf: make([]byte, 64),
		},
		// zero padding
		{
			p: Packet{
				Header: Header{
					Padding:        true,
					PayloadType:    96,
					SequenceNumber: 9164,
					Timestamp:      1681696377,
					SSRC:           0x6b8b4567,
				},
				Payload: []byte{0x01, 0x02, 0x03},
			},
			r:   []byte{0xa0, 0x60, 0x23, 0xcc, 0x64, 0x3c, 0xa6, 0x79, 0x6b, 0x8b, 0x45, 0x67, 0x01, 0x02, 0x03, 0x01},
			buf: make([]byte, 64),
		},
		// with padding
		{
			p: Packet{
				Header: Header{
					Padding:        true,
					PayloadType:    96,
					SequenceNumber: 9164,
					Timestamp:      1681696377,
					SSRC:           0x6b8b4567,
				},
				Payload:      []byte{0x01, 0x02, 0x03},
				PaddingBytes: 2,
			},
			r:   []byte{0xa0, 0x60, 0x23, 0xcc, 0x64, 0x3c, 0xa6, 0x79, 0x6b, 0x8b, 0x45, 0x67, 0x01, 0x02, 0x03, 0x00, 0x00, 0x03},
			buf: make([]byte, 64),
		},
		// not enough space for payload
		{
			p: Packet{
				Payload: []byte{0x01, 0x02, 0x03},
			},
			buf: make([]byte, HeaderLength+1),
			err: ErrNotEnoughBufferSpace{
				Actual:   HeaderLength + 1,
				Expected: HeaderLength + 3,
			},
		},
		// not enough space for padding
		{
			p: Packet{
				Header: Header{
					Padding: true,
				},
				Payload:      []byte{0x01, 0x02, 0x03},
				PaddingBytes: 2,
			},
			buf: make([]byte, HeaderLength+3),
			err: ErrNotEnoughBufferSpace{
				Actual:   HeaderLength + 3,
				Expected: HeaderLength + 3 + 2 + 1,
			},
		},
	}

	for i, c := range testCases {
		n, err := c.p.ComposeTo(c.buf)
		if c.err == nil {
			assert.NoError(t, err, "testCase : %d", i+1)
			assert.Equal(t, c.p.Size(), n, "testCase : %d", i+1)
			assert.Equal(t, c.r, c.buf[:n], "testCase : %d", i+1)
		} else {
			assert.ErrorIs(t, err, c.err, "testCase : %d", i+1)
		}
	}
}

func TestPacket_Parse(t *testing.T) {
	type testCase struct {
		r   []byte
		p   Packet
		err error
	}

	testCases := []testCase{
		{
			r: []byte{0x80, 0x60, 0x23, 0xcc, 0x64, 0x3c, 0xa6, 0x79, 0x6b, 0x8b, 0x45, 0x67, 0x01, 0x02, 0x03},
			p: Packet{
				Header: Header{
					PayloadType:    96,
					SequenceNumber: 9164,
					Timestamp:      1681696377,
					SSRC:           0x6b8b4567,
				},
				Payload: []byte{0x01, 0x02, 0x03},
			},
		},
		// with padding
		{
			r: []byte{0xa0, 0x60, 0x23, 0xcc, 0x64, 0x3c, 0xa6, 0x79, 0x6b, 0x8b, 0x45, 0x67, 0x01, 0x02, 0x03, 0x00, 0x00, 0x03},
			p: Packet{
				Header: Header{
					Padding:        true,
					PayloadType:    96,
					SequenceNumber: 9164,
					Timestamp:      1681696377,
					SSRC:           0x6b8b4567,
				},
				Payload:      []byte{0x01, 0x02, 0x03},
				PaddingBytes: 2,
			},
		},
		// invalid header
		{
			r: []byte{0xa0, 0x60, 0x23},
			err: ErrIncompleteHeader{
				Actual:   3,
				Expected: HeaderLength,
			},
		},
		// no payload
		{
			r:   []byte{0x80, 0x60, 0x23, 0xcc, 0x64, 0x3c, 0xa6, 0x79, 0x6b, 0x8b, 0x45, 0x67},
			err: ErrPayloadIsMissing{},
		},
		// no payload
		{
			r:   []byte{0xa0, 0x60, 0x23, 0xcc, 0x64, 0x3c, 0xa6, 0x79, 0x6b, 0x8b, 0x45, 0x67, 0x01},
			err: ErrPayloadIsMissing{},
		},
		// no payload
		{
			r:   []byte{0xa0, 0x60, 0x23, 0xcc, 0x64, 0x3c, 0xa6, 0x79, 0x6b, 0x8b, 0x45, 0x67, 0x00, 0x02},
			err: ErrPayloadIsMissing{},
		},
		// invalid padding
		{
			r:   []byte{0xa0, 0x60, 0x23, 0xcc, 0x64, 0x3c, 0xa6, 0x79, 0x6b, 0x8b, 0x45, 0x67, 0x01, 0x02, 0x03},
			err: ErrPayloadIsMissing{},
		},
		// zero padding
		{
			r: []byte{0xa0, 0x60, 0x23, 0xcc, 0x64, 0x3c, 0xa6, 0x79, 0x6b, 0x8b, 0x45, 0x67, 0x01, 0x02, 0x00},
			p: Packet{
				Header: Header{
					Padding:        true,
					PayloadType:    96,
					SequenceNumber: 9164,
					Timestamp:      1681696377,
					SSRC:           0x6b8b4567,
				},
				Payload: []byte{0x01, 0x02, 0x00},
			},
		},
	}

	var p Packet
	for i, c := range testCases {
		err := p.Parse(c.r)
		if c.err == nil {
			assert.NoError(t, err, "testCase : %d", i+1)
			assert.Equal(t, &c.p, &p, "testCase : %d", i+1)
		} else {
			assert.ErrorIs(t, err, c.err, "testCase : %d", i+1)
		}
	}
}
