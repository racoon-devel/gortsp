package rtsp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var headerRegEx = regexp.MustCompile("^([\\w-]+):[\\s]*([^\\n]+)$")

func readLine(rd *bufio.Reader) (string, error) {
	line, err := rd.ReadString('\n')
	if err != nil {
		return "", err
	}

	line = strings.TrimRight(line, "\r\n")
	return line, nil
}

func readHeaders(rd *bufio.Reader) (http.Header, error) {
	var h http.Header
	for {
		line, err := readLine(rd)
		if err != nil {
			return nil, err
		}

		if line == "" {
			break
		}

		matches := headerRegEx.FindStringSubmatch(line)
		if matches == nil {
			return nil, errors.New("cannot parse header")
		}

		if h == nil {
			h = make(http.Header)
		}
		h.Add(matches[1], matches[2])
	}

	return h, nil
}

func readBody(rd *bufio.Reader, h http.Header) ([]byte, error) {
	if contentLength := h.Get("Content-Length"); contentLength != "" {
		length, err := strconv.Atoi(contentLength)
		if err != nil {
			return nil, fmt.Errorf("Content-Length header malformed: %w", err)
		}
		body := make([]byte, length)
		_, err = io.ReadFull(rd, body)
		return body, err
	} else {
		return nil, nil
	}
}

func validateURL(u *url.URL) error {
	if u == nil {
		return ErrInvalidURL
	}

	if u.Scheme != "rtsp" {
		return fmt.Errorf("%w: invalid schema: %s", ErrInvalidURL, u.Scheme)
	}

	if u.Host == "" {
		return fmt.Errorf("%w: empty host", ErrInvalidURL)
	}

	if u.Port() == "" {
		return fmt.Errorf("%w: empty port", ErrInvalidURL)
	}

	return nil
}
