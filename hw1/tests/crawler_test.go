package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	crw "github.com/t0pep0/GB_best_go1/hw1/crawler"
	"github.com/t0pep0/GB_best_go1/hw1/crawler/mocks"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

var testPage = `<!DOCTYPE html>
		<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Test</title>
			</head>
			<body>
				<p><a href="https://github.com">github.com</a></p>
			</body>
		</html>`

func mockNewPage() (crw.Page, error) {
	return crw.NewPage(ioutil.NopCloser(strings.NewReader(testPage)))
}

func mockNewHttpClient() *mocks.BaseHttpClient {
	return &mocks.BaseHttpClient{}
}
func TestNewRequester(t *testing.T) {
	r := crw.NewRequester(mockNewHttpClient())
	assert.NotNil(t, r)
}

func TestNewCrawler(t *testing.T) {
	maxDepth := uint64(3)
	r := crw.NewRequester(mockNewHttpClient())
	cr := crw.NewCrawler(r, maxDepth)
	assert.NotNil(t, cr)
	assert.Equal(t, maxDepth, cr.MaxDepth)
}

func TestNewPage(t *testing.T) {
	page, _ := mockNewPage()
	assert.NotNil(t, page)
}

func TestCrawlerScan(t *testing.T) {
	r := &mocks.Requester{}
	page := &mocks.Page{}
	cr := crw.NewCrawler(r, 6)
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
	cr := crw.NewCrawler(r, 1)
	crRes := crw.CrawlResult{
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

func TestRequesterGet(t *testing.T) {
	cl := mockNewHttpClient()
	r := crw.NewRequester(cl)
	ctx := context.Background()
	url := "https://github.com"
	req, _ := http.NewRequest("GET", url, nil)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader(testPage)),
	}
	cl.On("Do", req).Return(resp, nil)
	page, err := r.Get(ctx, url)

	if err != nil {
		t.Fatalf("Failed get page: %v", err)
	}

	assert.NotNil(t, page)
}

func TestCrawlerIncreaseMaxDepth(t *testing.T) {
	r := &mocks.Requester{}
	cr := crw.NewCrawler(r, 2)
	chDepth := uint64(6)

	cr.IncreaseMaxDepth(4)

	assert.Equal(t, chDepth, cr.MaxDepth)
}

func TestPageGetTitle(t *testing.T) {
	page, _ := mockNewPage()

	assert.Equal(t, "Test", page.GetTitle())
	assert.Equal(t, []string{"https://github.com"}, page.GetLinks())
}

func TestPageGetLinks(t *testing.T) {
	page, _ := mockNewPage()

	assert.Equal(t, []string{"https://github.com"}, page.GetLinks())
}
