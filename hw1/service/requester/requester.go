//go:generate mockery --name=BaseHTTPClient
package requester

import (
	"context"
	"github.com/t0pep0/GB_best_go1/hw1/domain"
	"github.com/t0pep0/GB_best_go1/hw1/service/page"
	"net/http"
)

type BaseHTTPClient interface { // (lint: revive)
	Do(req *http.Request) (*http.Response, error)
}
type RequesterStruct struct {
	cl BaseHTTPClient
}

func NewRequester(cl BaseHTTPClient) RequesterStruct { // (lint: revive)
	return RequesterStruct{
		cl: cl,
	}
}

func (r RequesterStruct) Get(ctx context.Context, url string) (domain.Page, error) {
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		body, err := r.cl.Do(req)
		if err != nil {
			return nil, err
		}
		defer body.Body.Close()
		p, err := page.NewPage(body.Body)
		if err != nil {
			return nil, err
		}
		return p, nil
	}
	// (lint: govet)
}
