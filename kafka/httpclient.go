package kafka

import (
	"net/http"
	"net/url"
	"time"
)

// HTTPClient defines the client for HTTP communication
type HTTPClient struct {
	url        string
	httpClient *http.Client
	endPoint   string
}

// NewClient returns a new instance of HTTPClient
func NewHTTPClient(url, endpoint string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		url:        url,
		httpClient: &http.Client{Timeout: timeout},
		endPoint:   endpoint,
	}
}

// GetUrl returns the URL of a HTTPClient
func (hc *HTTPClient) GetUrl() string {
	u := url.URL{
		Scheme: "http",
		Host:   hc.url,
		Path:   hc.endPoint,
	}
	return u.String()
}
