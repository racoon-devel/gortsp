package rtp

import "encoding/binary"

const (
	versionMask  = 0xC0
	versionShift = 6
	paddingMask  = 0x20
	paddingShift = 5
	extMask      = 0x10
	extShift     = 4
	ccMask       = 0xF
	markerShift  = 7
	markerMask   = 0x80
	ptMask       = 0x7F
)

// RawPacket is a buffer with entire RTP packet
// All RawPacket methods are unsafe. It can be used only if you understand what
// you do
type RawPacket []byte

/*
	RFC3550:

	0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |V=2|P|X|  CC   |M|     PT      |       sequence number         |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                           timestamp                           |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |           synchronization source (SSRC) identifier            |
   +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
   |            contributing source (CSRC) identifiers             |
   |                             ....                              |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*/

// ValidateHeader returns RTP packet header size and error if it's malformed
func (p RawPacket) ValidateHeader() (size int, err error) {
	expected := HeaderLength
	if len(p) < expected {
		err = newErrIncompleteHeader(expected, len(p))
		return
	}

	version := p.Version()
	if version != Version {
		err = ErrVersionMismatch{version}
		return
	}

	expected += int(p.CC()) * 4

	if p.X() {
		expected += HeaderExtensionLength
	}

	if len(p) < expected {
		err = newErrIncompleteHeader(expected, len(p))
		return
	}

	if p.X() {
		expected += int(p.ExtensionLength()) * 4
	}

	if len(p) < expected {
		err = newErrIncompleteHeader(expected, len(p))
		return
	}

	size = expected
	return
}

func (p RawPacket) Version() uint8 {
	return p[0] >> versionShift
}

func (p RawPacket) SetVersion(version uint8) {
	p[0] &^= versionMask
	p[0] |= version << versionShift & versionMask
}

func (p RawPacket) P() bool {
	return (p[0] & paddingMask >> paddingShift) > 0
}

func (p RawPacket) SetP(bit bool) {
	if bit {
		p[0] |= paddingMask
	} else {
		p[0] &^= paddingMask
	}
}

func (p RawPacket) X() bool {
	return (p[0] & extMask >> extShift) > 0
}

func (p RawPacket) SetX(bit bool) {
	if bit {
		p[0] |= extMask
	} else {
		p[0] &^= extMask
	}
}

func (p RawPacket) CC() uint8 {
	return p[0] & ccMask
}

func (p RawPacket) SetCC(cc uint8) {
	p[0] &^= ccMask
	p[0] |= cc & ccMask
}

func (p RawPacket) M() bool {
	return (p[1] >> markerShift) > 0
}

func (p RawPacket) SetM(bit bool) {
	if bit {
		p[1] |= markerMask
	} else {
		p[1] &^= markerMask
	}
}

func (p RawPacket) PT() uint8 {
	return p[1] & ptMask
}

func (p RawPacket) SetPT(pt uint8) {
	p[1] &^= ptMask
	p[1] |= pt & ptMask
}

func (p RawPacket) Seq() uint16 {
	return binary.BigEndian.Uint16(p[2:4])
}

func (p RawPacket) SetSeq(seq uint16) {
	binary.BigEndian.PutUint16(p[2:4], seq)
}

func (p RawPacket) Timestamp() uint32 {
	return binary.BigEndian.Uint32(p[4:8])
}

func (p RawPacket) SetTimestamp(ts uint32) {
	binary.BigEndian.PutUint32(p[4:8], ts)
}

func (p RawPacket) SSRC() uint32 {
	return binary.BigEndian.Uint32(p[8:12])
}

func (p RawPacket) SetSSRC(ssrc uint32) {
	binary.BigEndian.PutUint32(p[8:12], ssrc)
}

func (p RawPacket) CSRC(index int) uint32 {
	offset := HeaderLength + index*4
	return binary.BigEndian.Uint32(p[offset : offset+4])
}

func (p RawPacket) SetCSRC(index int, csrc uint32) {
	offset := HeaderLength + index*4
	binary.BigEndian.PutUint32(p[offset:offset+4], csrc)
}

func (p RawPacket) Padding() uint8 {
	return p[len(p)-1]
}

func (p RawPacket) SetPadding(padding uint8) {
	p[len(p)-1] = padding
}

func (p RawPacket) ExtensionProfile() uint16 {
	offset := HeaderLength + p.CC()*4
	return binary.BigEndian.Uint16(p[offset : offset+2])
}

func (p RawPacket) SetExtensionProfile(profile uint16) {
	offset := HeaderLength + p.CC()*4
	binary.BigEndian.PutUint16(p[offset:offset+2], profile)
}

func (p RawPacket) ExtensionLength() uint16 {
	offset := HeaderLength + p.CC()*4 + 2
	return binary.BigEndian.Uint16(p[offset : offset+2])
}

func (p RawPacket) SetExtensionLength(length uint16) {
	offset := HeaderLength + p.CC()*4 + 2
	binary.BigEndian.PutUint16(p[offset:offset+2], length)
}

func (p RawPacket) ExtensionHeader() []byte {
	offset := int(HeaderLength + p.CC()*4 + HeaderExtensionLength)
	return p[offset : offset+int(p.ExtensionLength()*4)]
}

func (p RawPacket) SetExtensionHeader(hdr []byte) {
	offset := int(HeaderLength + p.CC()*4 + HeaderExtensionLength)
	copy(p[offset:], hdr)
}
