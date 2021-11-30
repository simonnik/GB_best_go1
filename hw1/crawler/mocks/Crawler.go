// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	crawler "github.com/t0pep0/GB_best_go1/hw1/crawler"
)

// Crawler is an autogenerated mock type for the Crawler type
type Crawler struct {
	mock.Mock
}

// ChanResult provides a mock function with given fields:
func (_m *Crawler) ChanResult() <-chan crawler.CrawlResult {
	ret := _m.Called()

	var r0 <-chan crawler.CrawlResult
	if rf, ok := ret.Get(0).(func() <-chan crawler.CrawlResult); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan crawler.CrawlResult)
		}
	}

	return r0
}

// IncreaseMaxDepth provides a mock function with given fields: depth
func (_m *Crawler) IncreaseMaxDepth(depth uint64) {
	_m.Called(depth)
}

// Scan provides a mock function with given fields: ctx, url, depth
func (_m *Crawler) Scan(ctx context.Context, url string, depth int) {
	_m.Called(ctx, url, depth)
}
