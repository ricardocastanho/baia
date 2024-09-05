package scraper

import (
	"baia/pkg/collector"
	"context"
	"fmt"

	"github.com/gocolly/colly/v2"
)

// PerfilScraper implements the RealEstateScraperInterface for the "Perfil" real estate website.
type PerfilScraper struct{}

// NewPerfilScraper creates a new instance of PerfilScraper.
func NewPerfilScraper() RealEstateScraper {
	return &PerfilScraper{}
}

// Run starts the scraping process for the given URL using the provided context.
func (p *PerfilScraper) Run(ctx context.Context, url string) ([]string, []string) {
	var (
		realStateurls = []string{}
		nextPages     = []string{}
	)

	c := collector.NewCollector()

	c.OnHTML("div#grid div.listing-item a[href]", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping collection due to context cancellation:", ctx.Err())
			return
		default:
			realStateurls = append(realStateurls, e.Attr("href"))
		}
	})

	select {
	case <-ctx.Done():
		fmt.Println("Stopping visit due to context cancellation:", ctx.Err())
		return nil, nil
	default:
		c.Visit(url)

		return realStateurls, nextPages
	}
}
