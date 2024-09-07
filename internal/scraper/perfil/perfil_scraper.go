package perfil

import (
	"baia/internal/contracts"
	"baia/pkg/collector"
	"context"
	"fmt"

	"github.com/gocolly/colly/v2"
)

// PerfilScraper implements the RealEstateScraperInterface for the "Perfil" real estate website.
type PerfilScraper struct{}

// NewPerfilScraper creates a new instance of PerfilScraper.
func NewPerfilScraper() contracts.RealEstateScraper {
	return &PerfilScraper{}
}

// GetRealStates starts the scraping process for the given URL using the provided context.
func (p *PerfilScraper) GetRealStates(ctx context.Context, url string) ([]string, []string) {
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

// GetRealStateData gets all the data from a given url
func (p *PerfilScraper) GetRealStateData(ctx context.Context, ch chan contracts.RealState, url string) {
	realState := contracts.RealState{}
	realState.Url = url

	c := collector.NewCollector()

	c.OnHTML("div.property-title", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping collection due to context cancellation:", ctx.Err())
			return
		default:
			h2 := e.DOM.Find("h2")
			span := h2.Find("span")

			realState.Cod = span.Text()

			span.Remove()

			realState.Name = h2.Text()
		}
	})

	c.OnHTML("div.valor-imovel span", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping collection due to context cancellation:", ctx.Err())
			return
		default:
			err := realState.SetPrice(e.Text)

			if err != nil {
				fmt.Println("Error while trying to parse real state price:", err)
				return
			}
		}
	})

	select {
	case <-ctx.Done():
		fmt.Println("Stopping visit due to context cancellation:", ctx.Err())
	default:
		c.Visit(url)
		ch <- realState
	}
}
