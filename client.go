package gortsp

import (
	"net/url"
	"time"
)

// A Client is a high-level RTSP client
type Client struct {
	// ReadTimeout specifies session read operation timeout
	ReadTimeout time.Duration

	// WriteTimeout specifies session write operation timeout
	WriteTimeout time.Duration

	url *url.URL
}

// Run starts receive stream from URL
func (c *Client) Run(streamURL string) error {
	var err error
	c.url, err = url.Parse(streamURL)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) setDefaults() {
	if c.ReadTimeout == 0 {
		c.ReadTimeout = readTimeout
	}

	if c.WriteTimeout == 0 {
		c.WriteTimeout = writeTimeout
	}
}
