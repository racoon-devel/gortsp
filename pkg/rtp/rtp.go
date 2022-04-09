package rtp

// Parse parses buffer and returns RTP packet if it is possible
func Parse(buf []byte) (*Packet, error) {
	var p Packet
	if err := p.Parse(buf); err != nil {
		return nil, err
	}
	return &p, nil
}
