package perfil

import (
	"baia/internal/contracts"
	"baia/pkg/collector"
	"context"
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

// PerfilScraper implements the RealEstateScraperInterface for the "Perfil" real estate website.
type PerfilScraper struct {
	realState contracts.RealState
}

// NewPerfilScraper creates a new instance of PerfilScraper.
func NewPerfilScraper() contracts.RealEstateScraper {
	return &PerfilScraper{
		realState: contracts.RealState{},
	}
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
	p.realState.Url = url

	c := collector.NewCollector()

	p.SetRealStateCode(ctx, c)
	p.SetRealStateName(ctx, c)
	p.SetRealStatePrice(ctx, c)

	select {
	case <-ctx.Done():
		fmt.Println("Stopping visit due to context cancellation:", ctx.Err())
	default:
		c.Visit(url)
		ch <- p.realState
	}
}

func (p *PerfilScraper) SetRealStateCode(ctx context.Context, c *colly.Collector) {
	c.OnHTML("div.property-title h2 span.imovel-codigo", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping collection due to context cancellation:", ctx.Err())
			return
		default:
			p.realState.SetCode(e.Text)
		}
	})
}

func (p *PerfilScraper) SetRealStateName(ctx context.Context, c *colly.Collector) {
	c.OnHTML("div.property-title", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping collection due to context cancellation:", ctx.Err())
			return
		default:
			h2 := e.DOM.Find("h2")
			span := h2.Find("span")

			span.Remove()

			p.realState.Name = strings.TrimSpace(h2.Text())
		}
	})
}

func (p *PerfilScraper) SetRealStatePrice(ctx context.Context, c *colly.Collector) {
	c.OnHTML("div.valor-imovel span", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping collection due to context cancellation:", ctx.Err())
			return
		default:
			err := p.realState.SetPrice(e.Text)

			if err != nil {
				fmt.Println("Error while trying to parse real state price:", err)
				return
			}
		}
	})
}
