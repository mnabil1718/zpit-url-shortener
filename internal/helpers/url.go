package helpers

import (
	"errors"
	urlib "net/url"
	"strings"
)

var (
	ErrInvalidURL = errors.New("Invalid URL format")
)

func ValidateURL(url string) error {
	url = strings.Trim(url, " ")
	u, err := urlib.Parse(url)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return ErrInvalidURL
	}

	// ensure host has a valid TLD
	host := u.Hostname()
	parts := strings.Split(host, ".")
	if len(parts) < 2 || parts[len(parts)-1] == "" {
		return ErrInvalidURL
	}

	return nil
}
