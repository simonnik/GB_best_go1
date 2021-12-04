//go:generate mockery --name=Crawler
package domain

import (
	"context"
)

// Crawler - интерфейс (контракт) краулера, (lint: gocritic)
type Crawler interface {
	Scan(ctx context.Context, url string, depth int)
	ChanResult() <-chan CrawlResult
	IncreaseMaxDepth(depth uint64)
}

type CrawlResult struct {
	Err   error
	Title string
	URL   string // (lint: revive)
}
