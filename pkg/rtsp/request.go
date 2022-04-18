package rtsp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	urlpkg "net/url"
	"regexp"
	"strconv"
)

var (
	requestLineRegEx = regexp.MustCompile("^([A-Z-_]+)\\s([^\\s]+)\\sRTSP\\/(\\d).(\\d)$")
)

// Request represents client RTSP request
type Request struct {
	// Method specifies the RTSP method (OPTIONS, DESCRIBE, ANNOUNCE, etc.).
	Method Method

	// URL specifies the URL to access.
	URL *urlpkg.URL

	ProtoMajor int // eg 1
	ProtoMinor int // eg 0

	// Header is a map of request headers
	Header http.Header

	// Body contains data, Content-Length header will be added automatically
	Body []byte
}

// NewRequest makes a new RTSP request with specified Method and URL
func NewRequest(method Method, url string) (*Request, error) {
	r := Request{Method: method}
	if r.Method == "" {
		return nil, ErrMethodMustBeSet
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
	if r.Method == "" {
		return ErrMethodMustBeSet
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

// Read reads and parses RTSP request
func (r *Request) Read(rd *bufio.Reader) error {
	requestLine, err := readLine(rd)
	if err != nil {
		return err
	}

	matches := requestLineRegEx.FindStringSubmatch(requestLine)
	if matches == nil {
		return errors.New("cannot parse request line")
	}

	r.Method = Method(matches[1])
	r.URL, err = urlpkg.Parse(matches[2])
	if err != nil {
		return fmt.Errorf("parse URL failed: %w", err)
	}
	r.ProtoMajor, _ = strconv.Atoi(matches[3])
	r.ProtoMinor, _ = strconv.Atoi(matches[4])

	r.Header, err = readHeaders(rd)
	if err != nil {
		return err
	}

	r.Body, err = readBody(rd, r.Header)
	if err != nil {
		return fmt.Errorf("read body failed: %w", err)
	}

	return nil
}

func (r Request) Seq() (uint64, error) {
	seqString := r.Header.Get("Cseq")
	seq, err := strconv.ParseUint(seqString, 10, 64)
	if err != nil {
		return 0, err
	}

	return seq, nil
}
