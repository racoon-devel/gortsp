package rtsp

import (
	"context"
	"github.com/racoon-devel/gortsp/internal/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSession_Do(t *testing.T) {
	mock := mocks.NewConnMock()
	defer mock.Close()

	s := NewSession(mock.Client(), context.Background())
	defer s.Close()

	req := Request{
		Method: Options,
		URL:    mustParse("rtsp://127.0.0.1:554/"),
		Header: map[string][]string{
			"User-Agent": {"go-rtsp"},
			"Host":       {"127.0.0.1"},
		},
	}

	expected := &Response{
		Status:     "OK",
		StatusCode: Ok,
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header: http.Header{
			"Cseq":           {"1"},
			"Content-Type":   {"application/sdp"},
			"Content-Length": {"11"},
		},
		Body: []byte("I_LIKE_TITS"),
	}

	mock.ExpectWrite(t, []byte("OPTIONS rtsp://127.0.0.1:554/ RTSP/1.0\r\nCseq: 1\r\nHost: 127.0.0.1\r\nUser-Agent: go-rtsp\r\n\r\n"))
	mock.ExpectRead(t, []byte("RTSP/1.0 200 OK\r\nContent-Length: 11\r\nContent-Type: application/sdp\r\nCseq: 1\r\n\r\nI_LIKE_TITS"))
	resp, err := s.Do(&req)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)

}
