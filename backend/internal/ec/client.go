package ec

import (
	"crypto/tls"
	"net/http"
	"time"
)

type httpDoer interface {
	Do(r *http.Request) (*http.Response, error)
}

// NewHTTPClient returns new HTTP client with timeout
// passed as an argument and without SSL verification.
func NewHTTPClient(timeoutInSeconds int) *http.Client {
	return &http.Client{
		Timeout: ReturnDurationFromSeconds(timeoutInSeconds),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

// ReturnDurationFromSeconds by passing amount of seconds as an argument.
func ReturnDurationFromSeconds(value int) time.Duration {
	return time.Duration(value) * time.Second
}
