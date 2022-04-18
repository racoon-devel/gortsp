package rtsp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

var (
	statusLineRegEx = regexp.MustCompile("^RTSP\\/(\\d).(\\d)[\\s]+(\\d\\d\\d)[\\s]+([\\w\\s]+)$")
)

// Response represents RTSP response
type Response struct {
	Status     string     // e.g. "OK"
	StatusCode StatusCode // e.g. 200

	ProtoMajor int // e.g. 1
	ProtoMinor int // e.g. 0

	// Header maps header keys to values. If the response had multiple
	// headers with the same key, they may be concatenated, with comma
	// delimiters.
	Header http.Header

	// Body contains response content if presented
	Body []byte
}

// Read reads and parses RTSP response header.
// The method does read response body if it's presented.
func (r *Response) Read(rd *bufio.Reader) error {
	statusLine, err := readLine(rd)
	if err != nil {
		return err
	}

	matches := statusLineRegEx.FindStringSubmatch(statusLine)
	if matches == nil {
		return errors.New("cannot parse status line")
	}

	r.ProtoMajor, _ = strconv.Atoi(matches[1])
	r.ProtoMinor, _ = strconv.Atoi(matches[2])
	intCode, _ := strconv.Atoi(matches[3])
	r.StatusCode = StatusCode(intCode)
	r.Status = matches[4]

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

func (r Response) Write(w io.Writer) error {
	bw := bufio.NewWriter(w)
	if _, err := bw.WriteString(fmt.Sprintf("RTSP/1.0 %d %s\r\n", r.StatusCode, r.Status)); err != nil {
		return err
	}

	if len(r.Body) > 0 {
		r.Header.Add("Content-Length", strconv.Itoa(len(r.Body)))
	}

	if err := r.Header.Write(bw); err != nil {
		return err
	}

	if _, err := bw.WriteString("\r\n"); err != nil {
		return err
	}

	if len(r.Body) > 0 {
		if _, err := bw.Write(r.Body); err != nil {
			return fmt.Errorf("write body failed: %w", err)
		}
	}

	return bw.Flush()
}

func (r Response) Seq() (uint64, error) {
	seqString := r.Header.Get("Cseq")
	seq, err := strconv.ParseUint(seqString, 10, 64)
	if err != nil {
		return 0, err
	}

	return seq, nil
}
