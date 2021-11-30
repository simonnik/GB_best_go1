package crawler

import (
	"context"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PuerkitoBio/goquery" // (lint: goimports)
)

// Crawler - интерфейс (контракт) краулера, (lint: gocritic)
type Crawler interface {
	Scan(ctx context.Context, url string, depth int)
	ChanResult() <-chan CrawlResult
	IncreaseMaxDepth(depth uint64)
}

type Struct struct {
	r        Requester
	Res      chan CrawlResult
	Visited  map[string]struct{}
	Mu       sync.RWMutex
	MaxDepth uint64
}

func NewCrawler(r Requester, maxDepth uint64) *Struct { // (lint: revive)
	return &Struct{
		r:        r,
		Res:      make(chan CrawlResult),
		Visited:  make(map[string]struct{}),
		Mu:       sync.RWMutex{},
		MaxDepth: maxDepth,
	}
}

type CrawlResult struct {
	Err   error
	Title string
	URL   string // (lint: revive)
}

type Requester interface {
	Get(ctx context.Context, url string) (Page, error)
}

type RequesterStruct struct {
	cl BaseHTTPClient
}

func NewRequester(cl BaseHTTPClient) RequesterStruct { // (lint: revive)
	return RequesterStruct{
		cl: cl,
	}
}

func (r RequesterStruct) Get(ctx context.Context, url string) (Page, error) {
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
		page, err := NewPage(body.Body)
		if err != nil {
			return nil, err
		}
		return page, nil
	}
	// (lint: govet)
}

func NewHTTPClient(timeout time.Duration) *HTTPClient { // (lint: revive)
	cl := &http.Client{
		Timeout: timeout,
	}
	return &HTTPClient{
		cl: *cl,
	}
}

type BaseHTTPClient interface { // (lint: revive)
	Do(req *http.Request) (*http.Response, error)
}

type HTTPClient struct { // (lint: revive)
	cl http.Client
}

func (h HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return h.cl.Do(req)
}

func (c *Struct) Scan(ctx context.Context, url string, depth int) {
	if int(c.MaxDepth) == depth { // Проверяем то, что есть запас по глубине, (lint: gocritic)
		return
	}
	c.Mu.RLock()
	_, ok := c.Visited[url] // Проверяем, что мы ещё не смотрели эту страницу, (lint: gocritic)
	c.Mu.RUnlock()
	if ok {
		return
	}
	select {
	case <-ctx.Done(): // Если контекст завершен - прекращаем выполнение, (lint: gocritic)
		return
	default:
		page, err := c.r.Get(ctx, url) // Запрашиваем страницу через Requester, (lint: gocritic)
		if err != nil {
			c.Res <- CrawlResult{Err: err} // Записываем ошибку в канал, (lint: gocritic)
			return
		}
		c.Mu.Lock()
		c.Visited[url] = struct{}{} // Помечаем страницу просмотренной, (lint: gocritic)
		c.Mu.Unlock()
		c.Res <- CrawlResult{ // Отправляем результаты в канал, (lint: gocritic)
			Title: page.GetTitle(),
			URL:   url,
		}
		for _, link := range page.GetLinks() {
			go c.Scan(ctx, link, depth+1) // На все полученные ссылки запускаем новую рутину сборки, (lint: gocritic)
		}
	}
}

func (c *Struct) ChanResult() <-chan CrawlResult {
	return c.Res
}

func (c *Struct) IncreaseMaxDepth(depth uint64) {
	atomic.AddUint64(&c.MaxDepth, depth)
}

type Page interface {
	GetTitle() string
	GetLinks() []string
}

type page struct {
	doc *goquery.Document
}

func NewPage(raw io.Reader) (Page, error) {
	doc, err := goquery.NewDocumentFromReader(raw)
	if err != nil {
		return nil, err
	}
	return &page{doc: doc}, nil
}

func (p *page) GetTitle() string {
	return p.doc.Find("title").First().Text()
}

func (p *page) GetLinks() []string {
	var urls []string
	p.doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if ok {
			urls = append(urls, url)
		}
	})
	return urls
}
