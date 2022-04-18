package rtsp

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestNewRequest(t *testing.T) {
	// invalid method
	_, err := NewRequest("", "rtsp://127.0.0.1:554/")
	assert.ErrorIs(t, err, ErrMethodMustBeSet)

	// bad URL
	_, err = NewRequest(Describe, "8086")
	assert.ErrorAs(t, err, &ErrInvalidURL)
}

func mustParse(rawURL string) *url.URL {
	r, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return r
}

func TestRequest_Write(t *testing.T) {
	type testCase struct {
		r      Request
		result string
		err    error
	}

	testCases := []testCase{
		{
			r: Request{
				Method: Options,
				URL:    mustParse("rtsp://127.0.0.1:554/"),
			},
			result: "OPTIONS rtsp://127.0.0.1:554/ RTSP/1.0\r\n\r\n",
		},
		{
			r: Request{
				Method: Options,
				URL:    mustParse("rtsp://127.0.0.1:554/"),
				Header: map[string][]string{
					"User-Agent": {"go-rtsp"},
					"Host":       {"127.0.0.1"},
				},
			},
			result: "OPTIONS rtsp://127.0.0.1:554/ RTSP/1.0\r\nHost: 127.0.0.1\r\nUser-Agent: go-rtsp\r\n\r\n",
		},
		{
			r: Request{
				Method: Options,
				URL:    mustParse("rtsp://127.0.0.1:554/"),
				Header: map[string][]string{
					"User-Agent": {"go-rtsp"},
					"Host":       {"127.0.0.1"},
				},
				Body: []byte("I_LIKE_TITS"),
			},

			result: "OPTIONS rtsp://127.0.0.1:554/ RTSP/1.0\r\nContent-Length: 11\r\nHost: 127.0.0.1\r\nUser-Agent: go-rtsp\r\n\r\nI_LIKE_TITS",
		},
		// empty method
		{
			r: Request{
				URL: mustParse("rtsp://127.0.0.1:554/"),
			},
			err: ErrMethodMustBeSet,
		},
		{
			r: Request{
				Method: Options,
			},
			err: ErrInvalidURL,
		},
	}

	for i, c := range testCases {
		buf := &bytes.Buffer{}
		err := c.r.Write(buf)
		if c.err == nil {
			assert.NoError(t, err, "testCase : %d", i+1)
			assert.EqualValues(t, c.result, buf.String(), "testCase : %d", i+1)
		} else {
			assert.ErrorIs(t, err, c.err, "testCase : %d", i+1)
		}
	}
}

func TestRequest_Read(t *testing.T) {
	type testCase struct {
		r   Request
		raw string
		e   bool
	}

	testCases := []testCase{
		{
			r: Request{
				Method:     Options,
				URL:        mustParse("rtsp://127.0.0.1:554/"),
				ProtoMajor: 1,
			},
			raw: "OPTIONS rtsp://127.0.0.1:554/ RTSP/1.0\r\n\r\n",
		},
		{
			r: Request{
				Method: Options,
				URL:    mustParse("rtsp://127.0.0.1:554/"),
				Header: map[string][]string{
					"User-Agent": {"go-rtsp"},
					"Host":       {"127.0.0.1"},
				},
				ProtoMajor: 1,
			},
			raw: "OPTIONS rtsp://127.0.0.1:554/ RTSP/1.0\r\nHost: 127.0.0.1\r\nUser-Agent: go-rtsp\r\n\r\n",
		},
		{
			r: Request{
				Method: Options,
				URL:    mustParse("rtsp://127.0.0.1:554/"),
				Header: map[string][]string{
					"User-Agent":     {"go-rtsp"},
					"Host":           {"127.0.0.1"},
					"Content-Length": {"11"},
				},
				Body:       []byte("I_LIKE_TITS"),
				ProtoMajor: 1,
			},

			raw: "OPTIONS rtsp://127.0.0.1:554/ RTSP/1.0\r\nContent-Length: 11\r\nHost: 127.0.0.1\r\nUser-Agent: go-rtsp\r\n\r\nI_LIKE_TITS",
		},
	}

	var req Request
	for i, c := range testCases {
		reader := bytes.NewReader([]byte(c.raw))
		br := bufio.NewReader(reader)
		err := req.Read(br)
		if !c.e {
			assert.NoError(t, err, "testCase : %d", i+1)
			assert.EqualValues(t, c.r, req, "testCase : %d", i+1)
		} else {
			assert.Error(t, err, "testCase : %d", i+1)
		}
	}
}
