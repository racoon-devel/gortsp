package rtsp

import "net/http"

// Response represents RTSP response
type Response struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200

	ProtoMajor int // e.g. 1
	ProtoMinor int // e.g. 0

	// Header maps header keys to values. If the response had multiple
	// headers with the same key, they may be concatenated, with comma
	// delimiters.
	Header http.Header
}
