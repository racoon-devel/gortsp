package rtp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRawPacket_Version(t *testing.T) {
	type testCase struct {
		b byte
		v uint8
	}

	testCases := []testCase{
		{
			b: 0x3f,
			v: 0,
		},
		{
			b: 0x7f,
			v: 1,
		},
		{
			b: 0x7f,
			v: 1,
		},
		{
			b: 0xff,
			v: 3,
		},
		{
			b: 0xc0,
			v: 3,
		},
	}

	for i, c := range testCases {
		p := RawPacket{c.b}
		assert.Equal(t, c.v, p.Version(), "testCase : %d", i+1)
	}
}

func TestRawPacket_SetVersion(t *testing.T) {
	type testCase struct {
		b byte
		v uint8
		e byte
	}

	testCases := []testCase{
		{
			b: 0xff,
			v: 3,
			e: 0xff,
		},
		{
			b: 0xff,
			v: 2,
			e: 0xbf,
		},
		{
			b: 0xff,
			v: 1,
			e: 0x7f,
		},
		{
			b: 0xff,
			v: 0,
			e: 0x3f,
		},
		{
			v: 3,
			e: 0xc0,
		},
		{
			v: 2,
			e: 0x80,
		},
		{
			v: 1,
			e: 0x40,
		},
		{},
	}

	for i, c := range testCases {
		p := RawPacket{c.b}
		p.SetVersion(c.v)
		assert.Equal(t, c.v, p.Version(), "testCase : %d", i+1)
		assert.Equal(t, c.e, p[0], "testCase : %d", i+1)
	}
}

func TestRawPacket_P(t *testing.T) {
	type testCase struct {
		b byte
		p bool
	}

	testCases := []testCase{
		{
			b: 0xff,
			p: true,
		},
		{
			b: 0xdf,
			p: false,
		},
		{
			b: 0x00,
			p: false,
		},
		{
			b: 0x20,
			p: true,
		},
	}

	for i, c := range testCases {
		p := RawPacket{c.b}
		assert.Equal(t, c.p, p.P(), "testCase : %d", i+1)
	}
}

func TestRawPacket_SetP(t *testing.T) {
	type testCase struct {
		b byte
		p bool
		e byte
	}

	testCases := []testCase{
		{
			b: 0xff,
			p: true,
			e: 0xff,
		},
		{
			b: 0xff,
			p: false,
			e: 0xdf,
		},
		{
			p: false,
		},
		{
			b: 0x20,
			p: true,
			e: 0x20,
		},
		{
			p: true,
			e: 0x20,
		},
	}

	for i, c := range testCases {
		p := RawPacket{c.b}
		p.SetP(c.p)
		assert.Equal(t, c.p, p.P(), "testCase : %d", i+1)
		assert.Equal(t, c.e, p[0], "testCase : %d", i+1)
	}
}

func TestRawPacket_X(t *testing.T) {
	type testCase struct {
		b byte
		x bool
	}

	testCases := []testCase{
		{
			b: 0xff,
			x: true,
		},
		{
			b: 0xef,
			x: false,
		},
		{
			b: 0x00,
			x: false,
		},
		{
			b: 0x10,
			x: true,
		},
	}

	for i, c := range testCases {
		p := RawPacket{c.b}
		assert.Equal(t, c.x, p.X(), "testCase : %d", i+1)
	}
}

func TestRawPacket_SetX(t *testing.T) {
	type testCase struct {
		b byte
		x bool
		e byte
	}

	testCases := []testCase{
		{
			b: 0xff,
			x: true,
			e: 0xff,
		},
		{
			b: 0xff,
			x: false,
			e: 0xef,
		},
		{
			x: false,
		},
		{
			b: 0x10,
			x: true,
			e: 0x10,
		},
		{
			x: true,
			e: 0x10,
		},
	}

	for i, c := range testCases {
		p := RawPacket{c.b}
		p.SetX(c.x)
		assert.Equal(t, c.x, p.X(), "testCase : %d", i+1)
		assert.Equal(t, c.e, p[0], "testCase : %d", i+1)
	}
}

func TestRawPacket_CC(t *testing.T) {
	type testCase struct {
		b  byte
		cc uint8
	}

	testCases := []testCase{
		{
			b:  0xff,
			cc: 0xf,
		},
		{
			b:  0x0f,
			cc: 0xf,
		},
		{
			b:  0x00,
			cc: 0x00,
		},
		{
			b:  0xf0,
			cc: 0x0,
		},
		{
			b:  0xfa,
			cc: 0xa,
		},
	}

	for i, c := range testCases {
		p := RawPacket{c.b}
		assert.Equal(t, c.cc, p.CC(), "testCase : %d", i+1)
	}
}

func TestRawPacket_SetCC(t *testing.T) {
	p := RawPacket{0xff}
	for i := 0; i < MaxCSRC; i++ {
		p.SetCC(uint8(i))
		assert.EqualValues(t, i, p.CC(), "testCase : %d", i+1)
		assert.EqualValues(t, 3, p.Version(), "testCase : %d", i+1)
		assert.True(t, p.P(), "testCase : %d", i+1)
		assert.True(t, p.X(), "testCase : %d", i+1)
	}

	p = RawPacket{0x00}
	for i := 0; i < MaxCSRC; i++ {
		p.SetCC(uint8(i))
		assert.EqualValues(t, i, p.CC(), "testCase : %d", i+1)
		assert.EqualValues(t, 0, p.Version(), "testCase : %d", i+1)
		assert.False(t, p.P(), "testCase : %d", i+1)
		assert.False(t, p.X(), "testCase : %d", i+1)
	}
}

func TestRawPacket_M(t *testing.T) {
	type testCase struct {
		b byte
		m bool
	}

	testCases := []testCase{
		{
			b: 0xff,
			m: true,
		},
		{
			b: 0x7f,
			m: false,
		},
		{
			b: 0x00,
			m: false,
		},
		{
			b: 0x80,
			m: true,
		},
	}

	for i, c := range testCases {
		p := RawPacket{0x00, c.b}
		assert.Equal(t, c.m, p.M(), "testCase : %d", i+1)
	}
}

func TestRawPacket_SetM(t *testing.T) {
	type testCase struct {
		b byte
		m bool
		e byte
	}

	testCases := []testCase{
		{
			b: 0xff,
			m: true,
			e: 0xff,
		},
		{
			b: 0xff,
			m: false,
			e: 0x7f,
		},
		{
			m: false,
		},
		{
			b: 0x80,
			m: true,
			e: 0x80,
		},
		{
			m: true,
			e: 0x80,
		},
	}

	for i, c := range testCases {
		p := RawPacket{0xff, c.b}
		p.SetM(c.m)
		assert.Equal(t, c.m, p.M(), "testCase : %d", i+1)
		assert.Equal(t, c.e, p[1], "testCase : %d", i+1)
		assert.EqualValues(t, 0xff, p[0], "testCase : %d", i+1)
	}
}

func TestRawPacket_PT(t *testing.T) {
	type testCase struct {
		b  byte
		pt uint8
	}

	testCases := []testCase{
		{
			b:  0xff,
			pt: 0x7f,
		},
		{
			b:  0x7f,
			pt: 0x7f,
		},
		{
			b:  0x00,
			pt: 0x00,
		},
		{
			b:  0x80,
			pt: 0x00,
		},
		{
			b:  0xd5,
			pt: 85,
		},
	}

	for i, c := range testCases {
		p := RawPacket{0x00, c.b}
		assert.Equal(t, c.pt, p.PT(), "testCase : %d", i+1)
	}
}

func TestRawPacket_SetPT(t *testing.T) {
	p := RawPacket{0xff, 0xff}
	p.SetPT(0xff)
	assert.EqualValues(t, 0xff, p[1])

	p.SetPT(0x00)
	assert.EqualValues(t, 0x80, p[1])

	p.SetM(false)

	p.SetPT(0xfc)
	assert.EqualValues(t, 0x7c, p[1])
}

func TestRawPacket_Seq(t *testing.T) {
	type testCase struct {
		b   []byte
		seq uint16
	}

	testCases := []testCase{
		{
			b:   []byte{0x07, 0xc8},
			seq: 1992,
		},
	}

	for i, c := range testCases {
		p := RawPacket{0xff, 0xff}
		p = append(p, c.b...)
		assert.Equal(t, c.seq, p.Seq(), "testCase : %d", i+1)
	}
}

func TestRawPacket_SetSeq(t *testing.T) {
	p := make(RawPacket, HeaderLength)
	p.SetSeq(0xf1f2)
	assert.EqualValues(t, 0xf1f2, p.Seq())
	assert.Zero(t, p.Timestamp())
	assert.Zero(t, p[1])
}

func TestRawPacket_Timestamp(t *testing.T) {
	type testCase struct {
		b  []byte
		ts uint32
	}

	testCases := []testCase{
		{
			b:  []byte{0x2c, 0x7f, 0x54, 0x5f},
			ts: 746542175,
		},
	}

	for i, c := range testCases {
		p := RawPacket{0xff, 0xff, 0xff, 0xff}
		p = append(p, c.b...)
		assert.Equal(t, c.ts, p.Timestamp(), "testCase : %d", i+1)
	}
}

func TestRawPacket_SetTimestamp(t *testing.T) {
	p := make(RawPacket, HeaderLength)
	p.SetTimestamp(0xf1f2f3f4)
	assert.EqualValues(t, 0xf1f2f3f4, p.Timestamp())
	assert.Zero(t, p.Seq())
	assert.Zero(t, p.SSRC())
}

func TestRawPacket_SSRC(t *testing.T) {
	type testCase struct {
		b    []byte
		ssrc uint32
	}

	testCases := []testCase{
		{
			b:    []byte{0xec, 0x17, 0x7d, 0xc8},
			ssrc: 0xec177dc8,
		},
	}

	for i, c := range testCases {
		p := RawPacket{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
		p = append(p, c.b...)
		assert.Equal(t, c.ssrc, p.SSRC(), "testCase : %d", i+1)
	}
}

func TestRawPacket_SetSSRC(t *testing.T) {
	p := make(RawPacket, HeaderLength)
	p.SetSSRC(0xf1f2f3f4)
	assert.EqualValues(t, 0xf1f2f3f4, p.SSRC())
	assert.Zero(t, p.Timestamp())
}

func TestRawPacket_CSRC(t *testing.T) {
	const cc = 4
	hdr := make(RawPacket, HeaderLength, HeaderLength+cc*4)
	csrcs := []byte{
		0x01, 0x02, 0x03, 0x04,
		0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c,
		0x0d, 0x0e, 0x0f, 0x10,
	}
	hdr = append(hdr, csrcs...)

	csrc := hdr.CSRC(0)
	assert.EqualValues(t, 0x1020304, csrc)

	csrc = hdr.CSRC(1)
	assert.EqualValues(t, 0x5060708, csrc)

	csrc = hdr.CSRC(2)
	assert.EqualValues(t, 0x90a0b0c, csrc)

	csrc = hdr.CSRC(3)
	assert.EqualValues(t, 0xd0e0f10, csrc)
}

func TestRawPacket_SetCSRC(t *testing.T) {
	const cc = 4
	hdr := make(RawPacket, HeaderLength+cc*4)
	csrcs := []uint32{
		0x1020304, 0x5060708, 0x90a0b0c, 0xd0e0f10,
	}

	for i, v := range csrcs {
		hdr.SetCSRC(i, v)
	}

	for i, v := range csrcs {
		assert.Equal(t, v, hdr.CSRC(i))
	}

	assert.Zero(t, hdr.SSRC())
}

func TestRawPacket_Padding(t *testing.T) {
	p := RawPacket{0x01, 0x02, 0x03, 0xff}
	assert.EqualValues(t, 0xff, p.Padding())
}

func TestRawPacket_SetPadding(t *testing.T) {
	p := RawPacket{0x01, 0x02, 0x03, 0xff}
	p.SetPadding(0xae)
	assert.EqualValues(t, 0xae, p.Padding())
	assert.Equal(t, RawPacket{0x01, 0x02, 0x03, 0xae}, p)
}

func TestRawPacket_ExtensionProfile(t *testing.T) {
	p := make(RawPacket, HeaderLength)
	p = append(p, 0xab, 0xcd)
	assert.EqualValues(t, 0xabcd, p.ExtensionProfile())

	const cc = 5
	p = make(RawPacket, HeaderLength+cc*4)
	p = append(p, 0xab, 0xcd)
	p[0] = cc
	assert.EqualValues(t, 0xabcd, p.ExtensionProfile())
}

func TestRawPacket_SetExtensionProfile(t *testing.T) {
	const cc = 5
	p := make(RawPacket, HeaderLength+cc*4)
	p = append(p, 0xab, 0xcd)
	p[0] = cc
	p.SetExtensionProfile(0xa1b2)
	assert.EqualValues(t, 0xa1b2, p.ExtensionProfile())
}

func TestRawPacket_ExtensionLength(t *testing.T) {
	p := make(RawPacket, HeaderLength+2)
	p = append(p, 0x1b, 0x2b)
	assert.EqualValues(t, 0x1b2b, p.ExtensionLength())

	const cc = 5
	p = make(RawPacket, HeaderLength+cc*4+2)
	p = append(p, 0x1b, 0x2b)
	p[0] = cc
	assert.EqualValues(t, 0x1b2b, p.ExtensionLength())
}

func TestRawPacket_SetExtensionLength(t *testing.T) {
	const cc = 5
	p := make(RawPacket, HeaderLength+cc*4+4)
	p[0] = cc
	p.SetExtensionLength(0x8ae2)
	assert.EqualValues(t, 0x8ae2, p.ExtensionLength())
}

func TestRawPacket_ExtensionHeader(t *testing.T) {
	ext := []byte{
		0x01, 0x02, 0x03, 0x04,
		0x05, 0x06, 0x07, 0x08,
	}
	p := make(RawPacket, HeaderLength+2)
	p = append(p, 0x00, 0x2)
	p = append(p, ext...)
	assert.Equal(t, ext, p.ExtensionHeader())

	const cc = 5
	p = make(RawPacket, HeaderLength+cc*4+2)
	p[0] = cc
	p = append(p, 0x00, 0x2)
	p = append(p, ext...)
	assert.Equal(t, ext, p.ExtensionHeader())
}

func TestRawPacket_SetExtensionHeader(t *testing.T) {
	ext := []byte{
		0x01, 0x02, 0x03, 0x04,
		0x05, 0x06, 0x07, 0x08,
	}

	const cc = 5
	p := make(RawPacket, HeaderLength+cc*4+HeaderExtensionLength+len(ext))
	p.SetCC(5)
	p.SetExtensionLength(2)
	p.SetExtensionHeader(ext)

	assert.Equal(t, ext, p.ExtensionHeader())
	assert.EqualValues(t, 2, p.ExtensionLength())
}

func TestRawPacket_Parse(t *testing.T) {
	type testCase struct {
		src     RawPacket
		p       bool
		x       bool
		m       bool
		seq     uint
		pt      uint
		ts      uint
		ssrc    uint
		csrc    []uint32
		padding uint

		extProfile uint
		ext        []byte
	}

	var testCases = []testCase{
		{
			// packet without flags
			src:  RawPacket{0x80, 0x60, 0x23, 0xcc, 0x64, 0x3c, 0xa6, 0x79, 0x6b, 0x8b, 0x45, 0x67, 0x7c, 0x81, 0xe0, 0x01},
			pt:   96,
			seq:  9164,
			ts:   1681696377,
			ssrc: 0x6b8b4567,
		},
		{
			// packet with marker bit
			src:  RawPacket{0x80, 0xe0, 0x23, 0xd7, 0x64, 0x3c, 0xb4, 0x89, 0x6b, 0x8b, 0x45, 0x67},
			m:    true,
			pt:   96,
			seq:  9175,
			ts:   1681699977,
			ssrc: 0x6b8b4567,
		},
		{
			// packet with padding
			src:     RawPacket{0xA0, 0x60, 0x24, 0x2f, 0x64, 0x3d, 0x79, 0x69, 0x6b, 0x8b, 0x45, 0x67, 0x00, 0xff, 0xff, 0xff, 0x04},
			p:       true,
			pt:      96,
			seq:     9263,
			ts:      1681750377,
			ssrc:    0x6b8b4567,
			padding: 4,
		},
		{
			// packet with padding and marker bit
			src:     RawPacket{0xA0, 0xe0, 0x24, 0x2f, 0x64, 0x3d, 0x79, 0x69, 0x6b, 0x8b, 0x45, 0x67, 0x00, 0xff, 0xff, 0xff, 0x04},
			p:       true,
			m:       true,
			pt:      96,
			seq:     9263,
			ts:      1681750377,
			ssrc:    0x6b8b4567,
			padding: 4,
		},
		{
			// packet with extension header
			src:        RawPacket{0x90, 0x60, 0x07, 0xb7, 0x2c, 0x7f, 0x54, 0x5f, 0xec, 0x17, 0x7d, 0xc8, 0xab, 0xac, 0x00, 0x03, 0xe5, 0xd3, 0x03, 0x75, 0x50, 0x1f, 0x38, 0x00, 0x00, 0x05, 0x00, 0x00},
			x:          true,
			pt:         96,
			seq:        1975,
			ts:         746542175,
			ssrc:       0xec177dc8,
			extProfile: 0xabac,
			ext:        []byte{0xe5, 0xd3, 0x03, 0x75, 0x50, 0x1f, 0x38, 0x00, 0x00, 0x05, 0x00, 0x00},
		},
	}

	for i, c := range testCases {
		assert.EqualValues(t, Version, c.src.Version(), "testCase : %d", i+1)
		assert.EqualValues(t, c.p, c.src.P(), "testCase : %d", i+1)
		assert.EqualValues(t, c.x, c.src.X(), "testCase : %d", i+1)
		assert.EqualValues(t, len(c.csrc), c.src.CC(), "testCase : %d", i+1)
		assert.EqualValues(t, c.m, c.src.M(), "testCase : %d", i+1)
		assert.EqualValues(t, c.pt, c.src.PT(), "testCase : %d", i+1)
		assert.EqualValues(t, c.seq, c.src.Seq(), "testCase : %d", i+1)
		assert.EqualValues(t, c.ts, c.src.Timestamp(), "testCase : %d", i+1)
		assert.EqualValues(t, c.ssrc, c.src.SSRC(), "testCase : %d", i+1)
		if c.p {
			assert.EqualValues(t, c.padding, c.src.Padding(), "testCase : %d", i+1)
		}
		if c.x {
			assert.EqualValues(t, c.extProfile, c.src.ExtensionProfile(), "testCase : %d", i+1)
			assert.EqualValues(t, len(c.ext), c.src.ExtensionLength()*4, "testCase : %d", i+1)
			assert.EqualValues(t, c.ext, c.src.ExtensionHeader(), "testCase : %d", i+1)
		}
	}
}
