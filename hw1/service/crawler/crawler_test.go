package crawler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/t0pep0/GB_best_go1/hw1/domain"
	pageMock "github.com/t0pep0/GB_best_go1/hw1/domain/mocks"
	"github.com/t0pep0/GB_best_go1/hw1/service/crawler/mocks"
	"testing"
)

func TestNewCrawler(t *testing.T) {
	maxDepth := uint64(3)
	r := new(mocks.Requester)
	cr := NewCrawler(r, maxDepth)
	assert.NotNil(t, cr)
	assert.Equal(t, maxDepth, cr.MaxDepth)
}

func TestCrawlerScan(t *testing.T) {
	r := &mocks.Requester{}
	page := &pageMock.Page{}
	cr := NewCrawler(r, 6)
	ctx := context.Background()
	fetcher := map[string][]string{
		"https://github.com": {
			"/features/codespaces",
			"#start-of-content",
		},
	}
	visitedUrlsCount := len(fetcher)
	page.On("GetTitle").Return("Crawler will be here")
	page.On("GetLinks").Return(fetcher["https://github.com"])
	r.On("Get", ctx, "https://github.com").Return(page, nil)
	for _, url := range fetcher["https://github.com"] {
		r.On("Get", ctx, url).Return(page, nil)
	}
	go cr.Scan(ctx, "https://github.com", 1)

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
}

func TestCrawlerChanResult(t *testing.T) {
	r := &mocks.Requester{}
	cr := NewCrawler(r, 1)
	crRes := domain.CrawlResult{
		Err:   nil,
		Title: "Test title",
		URL:   "https://github.com",
	}
	go func() {
		cr.Res <- crRes
	}()

	res := <-cr.ChanResult()
	assert.Equal(t, crRes, res)
}

func TestCrawlerIncreaseMaxDepth(t *testing.T) {
	r := &mocks.Requester{}
	cr := NewCrawler(r, 2)
	chDepth := uint64(6)

	cr.IncreaseMaxDepth(4)

	assert.Equal(t, chDepth, cr.MaxDepth)
}
