package main

import (
	"net/http"
)

// UserAgentRoundTripper wraps an RoundTripper to attach User-Agent header.
type UserAgentRoundTripper struct {
	RoundTripper http.RoundTripper
	UserAgent    string
}

func (u *UserAgentRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", u.UserAgent)
	return u.RoundTripper.RoundTrip(req)
}
