package rtsp

import (
	"fmt"
	"net/url"
)

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
