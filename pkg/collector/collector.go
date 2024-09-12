package collector

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
)

func NewCollector(logger *slog.Logger) *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"),
		colly.MaxDepth(2),
	)

	c.WithTransport(&http.Transport{
		IdleConnTimeout:       10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		logger.Info(fmt.Sprint("Visiting: ", r.URL.String()))
	})

	c.OnScraped(func(r *colly.Response) {
		logger.Info(fmt.Sprint("Finished: ", r.Request.URL.String()))
	})

	return c
}
