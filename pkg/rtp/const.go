package rtp

const (

	// Version is an RTP specification version
	Version = 2

	// HeaderLength is a minimal RTP packet header length
	HeaderLength = 12

	// MaxCSRC maximum CSRC count per packet
	MaxCSRC = 16

	// MaxPayloadType higher value of PayloadType
	MaxPayloadType = 127

	extensionHeaderLength = 4
)
