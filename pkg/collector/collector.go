package collector

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
)

func NewCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (compatible; RealEstateScraper/1.0)"),
		colly.MaxDepth(2),
	)

	c.WithTransport(&http.Transport{
		IdleConnTimeout:       10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL.String())
	})

	return c
}
