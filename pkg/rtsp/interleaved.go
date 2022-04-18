package rtsp

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
)

// InterleavedHeader represents RTSP interleaved header
type InterleavedHeader struct {
	Channel uint8
	Length  uint16
}

func (h InterleavedHeader) Write(w io.Writer) error {
	buf := make([]byte, InterleavedHeaderSize)
	buf[0] = MagicSymbol
	buf[1] = h.Channel
	binary.BigEndian.PutUint16(buf[2:InterleavedHeaderSize], h.Length)
	if _, err := w.Write(buf); err != nil {
		return err
	}

	return nil
}

func (h *InterleavedHeader) Read(rd *bufio.Reader) error {
	buf := make([]byte, InterleavedHeaderSize)
	if _, err := io.ReadFull(rd, buf); err != nil {
		return err
	}

	if buf[0] != MagicSymbol {
		return errors.New("invalid signature of interleaved header")
	}

	h.Channel = buf[1]
	h.Length = binary.BigEndian.Uint16(buf[2:InterleavedHeaderSize])

	return nil
}
