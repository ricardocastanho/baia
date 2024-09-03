package scraper

import (
	"baia/pkg/collector"
	"context"
	"fmt"

	"sync"

	"github.com/gocolly/colly/v2"
)

// PerfilScraper implements the RealEstateScraperInterface for the "Perfil" real estate website.
type PerfilScraper struct {
	ch chan string
	wg *sync.WaitGroup
}

// NewPerfilScraper creates a new instance of PerfilScraper.
func NewPerfilScraper(ch chan string) RealEstateScraper {
	return &PerfilScraper{
		ch: ch,
		wg: &sync.WaitGroup{},
	}
}

// Run starts the scraping process for the given URL using the provided context.
func (p *PerfilScraper) Run(ctx context.Context, url string) {
	numGoroutines := 1
	p.wg.Add(numGoroutines)

	go p.getRealEstateUrls(ctx, url)

	go func() {
		p.wg.Wait()
		close(p.ch)
	}()
}

// getRealEstateUrls scrapes the "Perfil" website for real estate URLs.
func (p *PerfilScraper) getRealEstateUrls(ctx context.Context, url string) {
	c := collector.NewCollector()

	c.OnHTML("div#grid div.listing-item a[href]", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping collection due to context cancellation:", ctx.Err())
			return
		case p.ch <- e.Attr("href"):
		}
	})

	c.OnScraped(func(r *colly.Response) {
		p.wg.Done()
	})

	select {
	case <-ctx.Done():
		fmt.Println("Stopping visit due to context cancellation:", ctx.Err())
		return
	default:
		c.Visit(url)
	}
}
