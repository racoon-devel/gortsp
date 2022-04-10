package rtp

const (

	// Version is an RTP specification version
	Version = 2

	// HeaderLength is a minimal RTP packet header length
	HeaderLength = 12

	// MaxCSRC is a maximum CSRC count per packet
	MaxCSRC = 15

	// MaxPayloadType is a higher value of PayloadType
	MaxPayloadType = 127

	// HeaderExtensionLength is a header extension describer
	//    0                   1                   2                   3
	//    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//   |      defined by profile       |           length              |
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	HeaderExtensionLength = 4
)
