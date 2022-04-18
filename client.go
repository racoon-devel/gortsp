package gortsp

import (
	"context"
	"fmt"
	"github.com/racoon-devel/gortsp/pkg/rtsp"
	"net"
	"net/http"
	urlpkg "net/url"
)

type Client struct {
	UserAgent string

	url *urlpkg.URL
	s   *rtsp.Session
}

func (c *Client) Run(url string) error {
	return c.RunWithContext(url, context.Background())
}

func (c *Client) RunWithContext(url string, ctx context.Context) error {
	u, err := urlpkg.Parse(url)
	if err != nil {
		return err
	}
	if u.Scheme != "rtsp" || u.Host == "" {
		return rtsp.ErrInvalidURL
	}

	conn, err := net.Dial("tcp", u.Host)
	if err != nil {
		return err
	}

	c.url = u
	c.s = rtsp.NewSession(conn, ctx)

	return nil
}

func (c *Client) Receive() error {
	_, err := c.do(rtsp.Options, nil, nil)
	if err != nil {
		return fmt.Errorf("do OPTIONS failed: %w", err)
	}
	// todo: processing options

	resp, err := c.do(rtsp.Describe, http.Header{"Accept": {"application/sdp"}}, nil)
	if err != nil {
		return fmt.Errorf("do DESCRIBE failed: %w", err)
	}
}

func (c *Client) do(method rtsp.Method, headers http.Header, body []byte) (*rtsp.Response, error) {
	req := rtsp.Request{
		Method: method,
		URL:    c.url,
		Header: headers,
		Body:   body,
	}

	req.Header.Add("User-Agent", c.UserAgent)

	return c.s.Do(&req)
}
