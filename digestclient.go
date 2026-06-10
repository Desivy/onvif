package onvif

import (
	"context"
	"net/http"
	"time"

	"github.com/icholy/digest"
)

// DigestClient represents an HTTP client used for making requests authenticated
// with HTTP digest authentication (RFC 2617 / RFC 7616).
// It supports both MD5 and SHA-256 digest schemes automatically.
type DigestClient struct {
	client *http.Client
}

// NewDigestClient returns a DigestClient that wraps a given standard library
// http Client with the given username and password, using icholy/digest as the
// RFC-compliant digest transport.
func NewDigestClient(baseClient *http.Client, username, password string) *DigestClient {
	t := &digest.Transport{
		Username: username,
		Password: password,
	}
	if baseClient != nil && baseClient.Transport != nil {
		t.Transport = baseClient.Transport
	}
	timeout := time.Duration(0)
	if baseClient != nil {
		timeout = baseClient.Timeout
	}
	return &DigestClient{
		client: &http.Client{
			Transport: t,
			Timeout:   timeout,
		},
	}
}

// Do performs an HTTP request using digest authentication.
func (dc *DigestClient) Do(httpMethod, endpoint, soap string) (*http.Response, error) {
	return dc.DoWithContext(context.Background(), httpMethod, endpoint, soap)
}

// DoWithContext is Do with request cancellation: the HTTP call aborts as soon
// as ctx is cancelled instead of waiting for the client timeout.
func (dc *DigestClient) DoWithContext(ctx context.Context, httpMethod, endpoint, soap string) (*http.Response, error) {
	req, err := createHttpRequest(ctx, httpMethod, endpoint, soap)
	if err != nil {
		return nil, err
	}
	return dc.client.Do(req)
}
