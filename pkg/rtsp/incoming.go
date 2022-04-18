package rtsp

import "github.com/racoon-devel/gortsp/pkg/rtp"

// IncomingRTP helps to receive RTP packet from stream
type IncomingRTP struct {
	Channel uint8
	Packet  rtp.RawPacket
}

// IncomingRTCP helps to receive RTCP packet from stream
type IncomingRTCP struct {
	Channel uint8
	Packet  []byte
}
