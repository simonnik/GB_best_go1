package httpcl

import (
	"net/http"
	"time"
)

func NewHTTPClient(timeout time.Duration) *HTTPClient { // (lint: revive)
	cl := &http.Client{
		Timeout: timeout,
	}
	return &HTTPClient{
		cl: *cl,
	}
}

type HTTPClient struct { // (lint: revive)
	cl http.Client
}

func (h HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return h.cl.Do(req)
}
