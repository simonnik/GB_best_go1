package crawler

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

//Crawler - интерфейс (контракт) краулера
type Crawler interface {
	Scan(ctx context.Context, url string, depth int)
	ChanResult() <-chan CrawlResult
	IncreaseMaxDepth(depth uint64)
}

type crawler struct {
	r        Requester
	Res      chan CrawlResult
	Visited  map[string]struct{}
	Mu       sync.RWMutex
	MaxDepth uint64
}

func NewCrawler(r Requester, maxDepth uint64) *crawler {
	return &crawler{
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
	Url   string
}

type Requester interface {
	Get(ctx context.Context, url string) (Page, error)
}

type requester struct {
	cl BaseHttpClient
}

func NewRequester(cl BaseHttpClient) requester {
	return requester{
		cl: cl,
	}
}

func (r requester) Get(ctx context.Context, url string) (Page, error) {
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
	return nil, nil
}

func NewHttpClient(timeout time.Duration) *HttpClient {
	cl := &http.Client{
		Timeout: timeout,
	}
	return &HttpClient{
		cl: *cl,
	}
}

type BaseHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HttpClient struct {
	cl http.Client
}

func (h HttpClient) Do(req *http.Request) (*http.Response, error) {
	return h.cl.Do(req)
}

func (c *crawler) Scan(ctx context.Context, url string, depth int) {
	if int(c.MaxDepth) == depth { //Проверяем то, что есть запас по глубине
		return
	}
	c.Mu.RLock()
	_, ok := c.Visited[url] //Проверяем, что мы ещё не смотрели эту страницу
	c.Mu.RUnlock()
	if ok {
		return
	}
	select {
	case <-ctx.Done(): //Если контекст завершен - прекращаем выполнение
		return
	default:
		page, err := c.r.Get(ctx, url) //Запрашиваем страницу через Requester
		if err != nil {
			c.Res <- CrawlResult{Err: err} //Записываем ошибку в канал
			return
		}
		c.Mu.Lock()
		c.Visited[url] = struct{}{} //Помечаем страницу просмотренной
		c.Mu.Unlock()
		c.Res <- CrawlResult{ //Отправляем результаты в канал
			Title: page.GetTitle(),
			Url:   url,
		}
		for _, link := range page.GetLinks() {
			go c.Scan(ctx, link, depth+1) //На все полученные ссылки запускаем новую рутину сборки
		}
	}
}

func (c *crawler) ChanResult() <-chan CrawlResult {
	return c.Res
}

func (c *crawler) IncreaseMaxDepth(depth uint64) {
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
