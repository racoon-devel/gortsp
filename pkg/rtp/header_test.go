package rtp

import (
	"github.com/racoon-devel/gortsp/pkg/cerr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func makeBuffer(size int) []byte {
	return make([]byte, size)
}

func TestHeader_Parse(t *testing.T) {
	var p Packet

	for i := 0; i < HeaderLength; i++ {
		assert.ErrorAs(t, p.Parse(makeBuffer(i)), &cerr.InBufferTooShort{})
	}

	assert.ErrorAs(t, p.Parse(makeBuffer(HeaderLength)), &cerr.ProtocolVersionMismatch{})
}
