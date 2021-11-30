package config

import (
	"flag"
)

// Config - структура для конфигурации, (lint: gocritic)
type Config struct {
	MaxDepth   uint64
	MaxResults int
	MaxErrors  int
	URL        string // (lint:revive)
	Timeout    int
}

func NewConfig() *Config {
	var (
		url        = flag.String("url", "https://telegram.org", "URL of target source")
		maxDepth   = flag.Uint64("maxDepth", 3, "Max depth for links")
		maxResults = flag.Int("maxResults", 10, "Max results of links")
		maxErrors  = flag.Int("maxErrors", 5, "Max errors of results")
		timeout    = flag.Int("timeout", 10, "Timeout in seconds")
	)
	flag.Parse()

	return &Config{
		MaxDepth:   *maxDepth,
		MaxResults: *maxResults,
		MaxErrors:  *maxErrors,
		URL:        *url,
		Timeout:    *timeout,
	}
}
