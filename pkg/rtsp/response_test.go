package rtsp

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestResponse_Read(t *testing.T) {
	type testCase struct {
		r   Response
		raw string
		err bool
	}

	testCases := []testCase{
		{
			raw: "RTSP/1.0 200 OK\r\nCseq: 1\r\nPublic: DESCRIBE, GET_PARAMETER, SET_PARAMETER, SETUP, TEARDOWN, PLAY\r\n\r\n",
			r: Response{
				Status:     "OK",
				StatusCode: 200,
				ProtoMajor: 1,
				ProtoMinor: 0,
				Header: http.Header{
					"Cseq":   {"1"},
					"Public": {"DESCRIBE, GET_PARAMETER, SET_PARAMETER, SETUP, TEARDOWN, PLAY"},
				},
			},
		},
		// some cameras do not send char '\r'
		{
			raw: "RTSP/1.0 200 OK\nCseq: 1\nPublic: DESCRIBE, GET_PARAMETER, SET_PARAMETER, SETUP, TEARDOWN, PLAY\n\n",
			r: Response{
				Status:     "OK",
				StatusCode: 200,
				ProtoMajor: 1,
				ProtoMinor: 0,
				Header: http.Header{
					"Cseq":   {"1"},
					"Public": {"DESCRIBE, GET_PARAMETER, SET_PARAMETER, SETUP, TEARDOWN, PLAY"},
				},
			},
		},
		{
			raw: "RTSP/1.0 200 OKCseq: 1\r\nPublic: DESCRIBE, GET_PARAMETER, SET_PARAMETER, SETUP, TEARDOWN, PLAY\r\n\r\n",
			err: true,
		},
		{
			raw: "RTSP/1.0 200 OK\r\nCseq: 1\r\nPublic: DESCRIBE, GET_PARAMETER, SET_PARAMETER, SETUP, TEARDOWN, PLAY\r\n",
			err: true,
		},
		{
			raw: "RTSP/1.0 200 OK\r\nCseq 1\r\nPublic: DESCRIBE, GET_PARAMETER, SET_PARAMETER, SETUP, TEARDOWN, PLAY\r\n",
			err: true,
		},
		{
			raw: "RTSP/1.0 200 OK\r\nCseq: 2\r\nContent-Type: application/sdp\nContent-Length: 11\r\n\r\nI_LIKE_TITS",
			r: Response{
				Status:     "OK",
				StatusCode: Ok,
				ProtoMajor: 1,
				ProtoMinor: 0,
				Header: http.Header{
					"Cseq":           {"2"},
					"Content-Type":   {"application/sdp"},
					"Content-Length": {"11"},
				},
				Body: []byte("I_LIKE_TITS"),
			},
		},
	}

	for i, c := range testCases {
		var resp Response
		reader := bytes.NewReader([]byte(c.raw))
		br := bufio.NewReader(reader)
		err := resp.Read(br)
		if !c.err {
			assert.NoError(t, err, "testCase : %d", i+1)
			assert.Equal(t, c.r, resp, "testCase : %d", i+1)
		} else {
			assert.Error(t, err, "testCase : %d", i+1)
		}
	}
}

func TestResponse_Write(t *testing.T) {
	type testCase struct {
		r   Response
		raw string
		err bool
	}

	testCases := []testCase{
		{
			raw: "RTSP/1.0 200 OK\r\nCseq: 1\r\nPublic: DESCRIBE, GET_PARAMETER, SET_PARAMETER, SETUP, TEARDOWN, PLAY\r\n\r\n",
			r: Response{
				Status:     "OK",
				StatusCode: 200,
				ProtoMajor: 1,
				ProtoMinor: 0,
				Header: http.Header{
					"Cseq":   {"1"},
					"Public": {"DESCRIBE, GET_PARAMETER, SET_PARAMETER, SETUP, TEARDOWN, PLAY"},
				},
			},
		},
		{
			raw: "RTSP/1.0 200 OK\r\nContent-Length: 11\r\nContent-Type: application/sdp\r\nCseq: 2\r\n\r\nI_LIKE_TITS",
			r: Response{
				Status:     "OK",
				StatusCode: Ok,
				ProtoMajor: 1,
				ProtoMinor: 0,
				Header: http.Header{
					"Cseq":         {"2"},
					"Content-Type": {"application/sdp"},
				},
				Body: []byte("I_LIKE_TITS"),
			},
		},
	}

	for i, c := range testCases {
		buf := &bytes.Buffer{}
		err := c.r.Write(buf)
		if !c.err {
			assert.NoError(t, err, "testCase : %d", i+1)
			assert.Equal(t, c.raw, buf.String(), "testCase : %d", i+1)
		} else {
			assert.Error(t, err, "testCase : %d", i+1)
		}
	}
}
