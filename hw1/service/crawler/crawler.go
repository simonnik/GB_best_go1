//go:generate mockery --name=Requester
package crawler

import (
	"context"
	"github.com/t0pep0/GB_best_go1/hw1/domain"
	"sync"
	"sync/atomic"
)

type Requester interface {
	Get(ctx context.Context, url string) (domain.Page, error)
}

type Struct struct {
	r        Requester
	Res      chan domain.CrawlResult
	Visited  map[string]struct{}
	Mu       sync.RWMutex
	MaxDepth uint64
}

func NewCrawler(r Requester, maxDepth uint64) *Struct { // (lint: revive)
	return &Struct{
		r:        r,
		Res:      make(chan domain.CrawlResult),
		Visited:  make(map[string]struct{}),
		Mu:       sync.RWMutex{},
		MaxDepth: maxDepth,
	}
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
			c.Res <- domain.CrawlResult{Err: err} // Записываем ошибку в канал, (lint: gocritic)
			return
		}
		c.Mu.Lock()
		c.Visited[url] = struct{}{} // Помечаем страницу просмотренной, (lint: gocritic)
		c.Mu.Unlock()
		c.Res <- domain.CrawlResult{ // Отправляем результаты в канал, (lint: gocritic)
			Title: page.GetTitle(),
			URL:   url,
		}
		for _, link := range page.GetLinks() {
			go c.Scan(ctx, link, depth+1) // На все полученные ссылки запускаем новую рутину сборки, (lint: gocritic)
		}
	}
}

func (c *Struct) ChanResult() <-chan domain.CrawlResult {
	return c.Res
}

func (c *Struct) IncreaseMaxDepth(depth uint64) {
	atomic.AddUint64(&c.MaxDepth, depth)
}
