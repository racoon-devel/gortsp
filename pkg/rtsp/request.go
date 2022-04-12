package rtsp

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	urlpkg "net/url"
	"strconv"
)

// Request represents client RTSP request
type Request struct {
	// Method specifies the RTSP method (OPTIONS, DESCRIBE, ANNOUNCE, etc.).
	Method Method

	// URL specifies either the URI being requested (for server
	// requests) or the URL to access (for client requests).
	URL *urlpkg.URL

	// Header is a map of request headers
	Header http.Header

	// Body contains data, Content-Length header will be added automatically
	Body []byte
}

// NewRequest makes a new RTSP request with specified Method and URL
func NewRequest(method Method, url string) (*Request, error) {
	r := Request{Method: method}
	if !r.Method.IsValid() {
		return nil, ErrInvalidMethod{Method: r.Method}
	}

	var err error
	r.URL, err = urlpkg.Parse(url)
	if err != nil {
		return nil, err
	}

	if err = validateURL(r.URL); err != nil {
		return nil, err
	}

	return &r, nil
}

// Write writes an RTSP request to any io.Writer. If Body defined Content-Length header will be added automatically
func (r Request) Write(w io.Writer) error {
	if !r.Method.IsValid() {
		return ErrInvalidMethod{Method: r.Method}
	}
	if err := validateURL(r.URL); err != nil {
		return err
	}

	if len(r.Body) != 0 {
		r.Header.Add("Content-Length", strconv.Itoa(len(r.Body)))
	}

	bw := bufio.NewWriter(w)
	_, err := bw.WriteString(fmt.Sprintf("%s %s RTSP/1.0\r\n", r.Method, r.URL))
	if err != nil {
		return err
	}

	if err = r.Header.Write(bw); err != nil {
		return err
	}

	if _, err = bw.WriteString("\r\n"); err != nil {
		return err
	}

	if len(r.Body) != 0 {
		if _, err = bw.Write(r.Body); err != nil {
			return err
		}
	}

	return bw.Flush()
}

func (r Request) Parse(buf []byte) error {
	return nil
}
