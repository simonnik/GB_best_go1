package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	crw "github.com/t0pep0/GB_best_go1/hw1/crawler"
	"github.com/t0pep0/GB_best_go1/hw1/crawler/mocks"
	"testing"
	"time"
)

func TestNewRequester(t *testing.T) {
	sec := 5 * time.Second
	r := crw.NewRequester(sec)
	assert.NotNil(t, r)
	assert.Equal(t, sec, r.Timeout)
}

func TestNewCrawler(t *testing.T) {
	maxDepth := uint64(3)
	r := crw.NewRequester(5 * time.Second)
	cr := crw.NewCrawler(r, maxDepth)
	assert.NotNil(t, cr)
	assert.Equal(t, maxDepth, cr.MaxDepth)
}

func TestCrawlerScan(t *testing.T) {
	r := &mocks.Requester{}
	page := &mocks.Page{}
	cr := crw.NewCrawler(r, 3)
	ctx := context.Background()
	fetcher := []string{
		"https://github.com",
		"/features/codespaces",
		"#start-of-content",
		"/enterprise",
		"/explore",
		"/topics",
	}
	visitedUrlsCount := len(fetcher)
	page.On("GetTitle").Return("Crawler will be here")
	page.On("GetLinks").Return(fetcher)
	go cr.Scan(ctx, fetcher[0], 1)
	for url := range fetcher {
		r.On("Get", ctx, url).Return(page, nil)
	}

	var maxResult, maxErrors = 10, 5
	doFor := true
	for doFor {
		select {
		case res := <-cr.ChanResult():
			if res.Err != nil {
				maxErrors--
				if maxErrors <= 0 {
					doFor = false
				}
				break
			} else if len(res.Title) > 0 {
				doFor = false
			} else {
				maxResult--
				if maxResult <= 0 {
					doFor = false
				}
			}
		}
	}

	assert.Equal(t, visitedUrlsCount, len(cr.Visited))
	//assert.NotNil(t, cl)
	//assert.Equal(t, timeout, cl.cl.Timeout)
}
