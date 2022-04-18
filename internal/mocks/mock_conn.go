package mocks

import (
	mockConn "github.com/jordwest/mock-conn"
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"testing"
)

type lambda = func() bool

type Conn struct {
	m  *mockConn.Conn
	ch chan lambda
}

func NewConnMock() *Conn {
	c := &Conn{
		m:  mockConn.NewConn(),
		ch: make(chan lambda, 100),
	}

	go func() {
		for l := range c.ch {
			if !l() {
				return
			}
		}
	}()

	return c
}

func (c Conn) Client() net.Conn {
	return c.m.Client
}

func (c *Conn) ExpectWrite(t *testing.T, data []byte) {
	c.ch <- func() bool {
		buf := make([]byte, len(data))
		_, err := io.ReadFull(c.m.Server, buf)
		assert.NoError(t, err)
		assert.Equal(t, data, buf)
		return true
	}
}

func (c *Conn) ExpectRead(t *testing.T, data []byte) {
	c.ch <- func() bool {
		_, err := c.m.Server.Write(data)
		assert.NoError(t, err)
		return true
	}
}

func (c *Conn) Close() {
	c.ch <- func() bool {
		return false
	}
	_ = c.m.Close()
}
