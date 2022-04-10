package rtp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHeader_Parse(t *testing.T) {
	type testCase struct {
		raw    []byte
		header Header
		err    error
	}

	var testCases = []testCase{
		// packet length less than 12 bytes
		{
			raw: []byte{0x80, 0x60, 0x23, 0xcc, 0x64, 0x3c, 0xa6, 0x79, 0x6b, 0x8b, 0x45},
			err: ErrIncompleteHeader{
				Actual:   HeaderLength - 1,
				Expected: HeaderLength,
			},
		},
		// invalid version
		{
			raw: []byte{0x7F, 0x60, 0x23, 0xcc, 0x64, 0x3c, 0xa6, 0x79, 0x6b, 0x8b, 0x45, 0x67},
			err: ErrVersionMismatch{1},
		},
		// simple packet
		{
			raw: []byte{0x80, 0x60, 0x23, 0xcc, 0x64, 0x3c, 0xa6, 0x79, 0x6b, 0x8b, 0x45, 0x67},
			header: Header{
				PayloadType:    96,
				SequenceNumber: 9164,
				Timestamp:      1681696377,
				SSRC:           0x6b8b4567,
			},
		},
		// simple packet + flags
		{
			raw: []byte{0xA0, 0xe0, 0x24, 0x2f, 0x64, 0x3d, 0x79, 0x69, 0x6b, 0x8b, 0x45, 0x67},
			header: Header{
				Marker:         true,
				Padding:        true,
				PayloadType:    96,
				SequenceNumber: 9263,
				Timestamp:      1681750377,
				SSRC:           0x6b8b4567,
			},
		},
		// CSRC is missing
		{
			raw: []byte{0xA1, 0xe0, 0x24, 0x2f, 0x64, 0x3d, 0x79, 0x69, 0x6b, 0x8b, 0x45, 0x67},
			err: ErrIncompleteHeader{
				Actual:   HeaderLength,
				Expected: HeaderLength + 4,
			},
		},
		// incomplete CSRC
		{
			raw: []byte{0xA2, 0xe0, 0x24, 0x2f, 0x64, 0x3d, 0x79, 0x69, 0x6b, 0x8b, 0x45, 0x67, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
			err: ErrIncompleteHeader{
				Actual:   HeaderLength + 2*4 - 1,
				Expected: HeaderLength + 2*4,
			},
		},
		// simple packet + flags + 2 csrc
		{
			raw: []byte{0xA2, 0xe0, 0x24, 0x2f, 0x64, 0x3d, 0x79, 0x69, 0x6b, 0x8b, 0x45, 0x67, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			header: Header{
				Marker:         true,
				Padding:        true,
				PayloadType:    96,
				SequenceNumber: 9263,
				Timestamp:      1681750377,
				SSRC:           0x6b8b4567,
				CSRC: []uint32{
					0x01020304,
					0x05060708,
				},
			},
		},
		// packet with X bit, but no extension header
		{
			raw: []byte{0x90, 0x60, 0x07, 0xb7, 0x2c, 0x7f, 0x54, 0x5f, 0xec, 0x17, 0x7d, 0xc8},
			err: ErrIncompleteHeader{
				Actual:   HeaderLength,
				Expected: HeaderLength + HeaderExtensionLength,
			},
		},
		// incomplete extension header
		{
			raw: []byte{0x90, 0x60, 0x07, 0xb7, 0x2c, 0x7f, 0x54, 0x5f, 0xec, 0x17, 0x7d, 0xc8, 0xab, 0xac, 0x00, 0x03, 0xe5, 0xd3, 0x03, 0x75, 0x50, 0x1f, 0x38, 0x00, 0x00, 0x05, 0x00},
			err: ErrIncompleteHeader{
				Actual:   HeaderLength + HeaderExtensionLength + 3*4 - 1,
				Expected: HeaderLength + HeaderExtensionLength + 3*4,
			},
		},
		// packet with extension header
		{
			raw: []byte{0x90, 0x60, 0x07, 0xb7, 0x2c, 0x7f, 0x54, 0x5f, 0xec, 0x17, 0x7d, 0xc8, 0xab, 0xac, 0x00, 0x03, 0xe5, 0xd3, 0x03, 0x75, 0x50, 0x1f, 0x38, 0x00, 0x00, 0x05, 0x00, 0x00},
			header: Header{
				PayloadType:    96,
				SequenceNumber: 1975,
				Timestamp:      746542175,
				SSRC:           0xec177dc8,
				Extension: &ExtensionHeader{
					Profile: 0xabac,
					Content: []byte{0xe5, 0xd3, 0x03, 0x75, 0x50, 0x1f, 0x38, 0x00, 0x00, 0x05, 0x00, 0x00},
				},
			},
		},
		// packet with extension header and 2 CSRC
		{
			raw: []byte{0x92, 0x60, 0x07, 0xb7, 0x2c, 0x7f, 0x54, 0x5f, 0xec, 0x17, 0x7d, 0xc8, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0xab, 0xac, 0x00, 0x03, 0xe5, 0xd3, 0x03, 0x75, 0x50, 0x1f, 0x38, 0x00, 0x00, 0x05, 0x00, 0x00},
			header: Header{
				PayloadType:    96,
				SequenceNumber: 1975,
				Timestamp:      746542175,
				SSRC:           0xec177dc8,
				CSRC: []uint32{
					0x01020304,
					0x05060708,
				},
				Extension: &ExtensionHeader{
					Profile: 0xabac,
					Content: []byte{0xe5, 0xd3, 0x03, 0x75, 0x50, 0x1f, 0x38, 0x00, 0x00, 0x05, 0x00, 0x00},
				},
			},
		},
	}

	var h Header
	for i, c := range testCases {
		n, err := h.Parse(c.raw)
		if c.err == nil {
			assert.NoError(t, err, "testCase : %d", i+1)
			assert.Equal(t, c.header.Size(), n, "testCase : %d", i+1)
			assert.Equal(t, &c.header, &h, "testCase : %d", i+1)
		} else {
			assert.ErrorIs(t, err, c.err, "testCase : %d", i+1)
		}
	}
}

func TestHeader_Size(t *testing.T) {
	type testCase struct {
		Header
		size int
	}

	testCases := []testCase{
		{
			size: HeaderLength,
		},
		{
			Header: Header{
				CSRC:      []uint32{0x8390, 0x6748, 0x7463},
				Extension: nil,
			},
			size: HeaderLength + 3*4,
		},
		{
			Header: Header{
				Extension: &ExtensionHeader{
					Content: []byte{1, 2, 3, 4, 5, 6, 7, 8},
				},
			},
			size: HeaderLength + HeaderExtensionLength + 8,
		},
		{
			Header: Header{
				CSRC: []uint32{0x8390, 0x6748, 0x7463},
				Extension: &ExtensionHeader{
					Content: []byte{1, 2, 3, 4, 5, 6, 7, 8},
				},
			},
			size: HeaderLength + 3*4 + HeaderExtensionLength + 8,
		},
		{
			Header: Header{
				CSRC: []uint32{0x8390, 0x6748, 0x7463},
				Extension: &ExtensionHeader{
					Content: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
				},
			},
			size: HeaderLength + 3*4 + HeaderExtensionLength + 12,
		},
	}

	for i, c := range testCases {
		assert.Equal(t, c.Size(), c.size, "testCase : %d", i+1)
	}
}

func TestHeader_Compose(t *testing.T) {

	h := Header{
		Marker:         true,
		PayloadType:    28,
		SequenceNumber: 54,
		Timestamp:      12345,
		SSRC:           0xabcdef,
		CSRC:           []uint32{0x01020304, 0x05060708},
		Extension: &ExtensionHeader{
			Profile: 0xabcd,
			Content: []byte{2, 3, 4},
		},
	}
	buf, err := h.Compose()
	assert.NoError(t, err)
	assert.Equal(t, h.Size(), len(buf))
}

func TestHeader_ComposeTo(t *testing.T) {
	type testCase struct {
		header  Header
		buf     []byte
		result  []byte
		err     error
		padding bool
	}

	testCases := []testCase{
		{
			buf: make([]byte, HeaderLength-1),
			err: ErrNotEnoughBufferSpace{
				Actual:   HeaderLength - 1,
				Expected: HeaderLength,
			},
		},
		{
			buf:    make([]byte, HeaderLength),
			result: []byte{0x80, 0x60, 0x23, 0xcc, 0x64, 0x3c, 0xa6, 0x79, 0x6b, 0x8b, 0x45, 0x67},
			header: Header{
				PayloadType:    96,
				SequenceNumber: 9164,
				Timestamp:      1681696377,
				SSRC:           0x6b8b4567,
			},
		},
		{
			buf:    make([]byte, HeaderLength),
			result: []byte{0xA0, 0xe0, 0x24, 0x2f, 0x64, 0x3d, 0x79, 0x69, 0x6b, 0x8b, 0x45, 0x67},
			header: Header{
				Marker:         true,
				Padding:        true,
				PayloadType:    96,
				SequenceNumber: 9263,
				Timestamp:      1681750377,
				SSRC:           0x6b8b4567,
			},
		},
		// not enough length for CSRC
		{
			buf:    make([]byte, HeaderLength),
			header: Header{CSRC: []uint32{0}},
			err: ErrNotEnoughBufferSpace{
				Actual:   HeaderLength,
				Expected: HeaderLength + 4,
			},
		},
		{
			buf:    make([]byte, 64),
			result: []byte{0xA2, 0xe0, 0x24, 0x2f, 0x64, 0x3d, 0x79, 0x69, 0x6b, 0x8b, 0x45, 0x67, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			header: Header{
				Marker:         true,
				Padding:        true,
				PayloadType:    96,
				SequenceNumber: 9263,
				Timestamp:      1681750377,
				SSRC:           0x6b8b4567,
				CSRC: []uint32{
					0x01020304,
					0x05060708,
				},
			},
		},
		// invalid PT
		{
			buf:    make([]byte, HeaderLength),
			header: Header{PayloadType: 128},
			err:    ErrInvalidPayloadType{128},
		},
		// CSRC limit reached
		{
			buf:    make([]byte, HeaderLength),
			header: Header{CSRC: make([]uint32, MaxCSRC+1)},
			err:    ErrCSRCLimitExceeded{MaxCSRC + 1},
		},
		// not enough length for extension header
		{
			buf: make([]byte, HeaderLength+HeaderExtensionLength+11),
			header: Header{
				PayloadType:    96,
				SequenceNumber: 1975,
				Timestamp:      746542175,
				SSRC:           0xec177dc8,
				Extension: &ExtensionHeader{
					Profile: 0xabac,
					Content: []byte{0xe5, 0xd3, 0x03, 0x75, 0x50, 0x1f, 0x38, 0x00, 0x00, 0x05, 0x00, 0x00},
				},
			},
			err: ErrNotEnoughBufferSpace{
				Actual:   HeaderLength + HeaderExtensionLength + 11,
				Expected: HeaderLength + HeaderExtensionLength + 12,
			},
		},
		{
			buf:    make([]byte, HeaderLength+HeaderExtensionLength+12),
			result: []byte{0x90, 0x60, 0x07, 0xb7, 0x2c, 0x7f, 0x54, 0x5f, 0xec, 0x17, 0x7d, 0xc8, 0xab, 0xac, 0x00, 0x03, 0xe5, 0xd3, 0x03, 0x75, 0x50, 0x1f, 0x38, 0x00, 0x00, 0x05, 0x00, 0x00},
			header: Header{
				PayloadType:    96,
				SequenceNumber: 1975,
				Timestamp:      746542175,
				SSRC:           0xec177dc8,
				Extension: &ExtensionHeader{
					Profile: 0xabac,
					Content: []byte{0xe5, 0xd3, 0x03, 0x75, 0x50, 0x1f, 0x38, 0x00, 0x00, 0x05, 0x00, 0x00},
				},
			},
		},
		// extension header padding
		{
			buf:    make([]byte, HeaderLength+HeaderExtensionLength+16),
			result: []byte{0x90, 0x60, 0x07, 0xb7, 0x2c, 0x7f, 0x54, 0x5f, 0xec, 0x17, 0x7d, 0xc8, 0xab, 0xac, 0x00, 0x04, 0xe5, 0xd3, 0x03, 0x75, 0x50, 0x1f, 0x38, 0x00, 0x00, 0x05, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00},
			header: Header{
				PayloadType:    96,
				SequenceNumber: 1975,
				Timestamp:      746542175,
				SSRC:           0xec177dc8,
				Extension: &ExtensionHeader{
					Profile: 0xabac,
					Content: []byte{0xe5, 0xd3, 0x03, 0x75, 0x50, 0x1f, 0x38, 0x00, 0x00, 0x05, 0x00, 0x00, 0xFF},
				},
			},
			padding: true,
		},
		{
			buf:    make([]byte, 128),
			result: []byte{0x92, 0x60, 0x07, 0xb7, 0x2c, 0x7f, 0x54, 0x5f, 0xec, 0x17, 0x7d, 0xc8, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0xab, 0xac, 0x00, 0x03, 0xe5, 0xd3, 0x03, 0x75, 0x50, 0x1f, 0x38, 0x00, 0x00, 0x05, 0x00, 0x00},
			header: Header{
				PayloadType:    96,
				SequenceNumber: 1975,
				Timestamp:      746542175,
				SSRC:           0xec177dc8,
				CSRC: []uint32{
					0x01020304,
					0x05060708,
				},
				Extension: &ExtensionHeader{
					Profile: 0xabac,
					Content: []byte{0xe5, 0xd3, 0x03, 0x75, 0x50, 0x1f, 0x38, 0x00, 0x00, 0x05, 0x00, 0x00},
				},
			},
		},
	}

	for i, c := range testCases {
		n, err := c.header.ComposeTo(c.buf)
		if c.err == nil {
			assert.Equal(t, c.header.Size(), n, "testCase : %d", i+1)
			assert.NoError(t, err, "testCase : %d", i+1)
			assert.Equal(t, c.result, c.buf[:c.header.Size()], "testCase : %d", i+1)
			// extra check
			h := Header{}
			_, err = h.Parse(c.result)
			assert.NoError(t, err, "testCase : %d", i+1)
			// if extension header padding filled by zeroes - Headers are not equal
			if !c.padding {
				assert.Equal(t, &c.header, &h, "testCase : %d", i+1)
			}
		} else {
			assert.ErrorIs(t, err, c.err, "testCase : %d", i+1)
		}
	}
}
