package perfil

import (
	"baia/internal/contracts"
	"baia/pkg/collector"
	"context"
	"fmt"
	"regexp"

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
func (p *PerfilScraper) GetRealStateData(ctx context.Context, ch chan contracts.RealEstate, re *contracts.RealEstate) {
	c := collector.NewCollector()

	p.SetRealStateCode(ctx, c, re)
	p.SetRealStateName(ctx, c, re)
	p.SetRealStateDescription(ctx, c, re)
	p.SetRealStatePrice(ctx, c, re)
	p.SetRealStateBedrooms(ctx, c, re)
	p.SetRealStateBathrooms(ctx, c, re)
	p.SetRealStateArea(ctx, c, re)
	p.SetRealStateGarageSpaces(ctx, c, re)

	c.OnScraped(func(c *colly.Response) {
		ch <- *re
	})

	select {
	case <-ctx.Done():
		fmt.Println("Stopping visit due to context cancellation:", ctx.Err())
	default:
		c.Visit(re.Url)
		return
	}
}

func (p *PerfilScraper) SetRealStateCode(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div.property-title h2 span.imovel-codigo", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping collection due to context cancellation:", ctx.Err())
			return
		default:
			r.SetCode(e.Text)
		}
	})
}

func (p *PerfilScraper) SetRealStateName(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div.property-title", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping collection due to context cancellation:", ctx.Err())
			return
		default:
			h2 := e.DOM.Find("h2")
			span := h2.Find("span")

			span.Remove()

			r.SetName(h2.Text())
		}
	})
}

func (p *PerfilScraper) SetRealStateDescription(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div#text-0 div p", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping collection due to context cancellation:", ctx.Err())
			return
		default:
			r.SetDescription(e.Text)
		}
	})
}

func (p *PerfilScraper) SetRealStatePrice(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div.valor-imovel span", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping collection due to context cancellation:", ctx.Err())
			return
		default:
			err := r.SetPrice(e.Text)

			if err != nil {
				fmt.Println("Error while trying to parse real state price:", err)
				return
			}
		}
	})
}

func (p *PerfilScraper) SetRealStateBedrooms(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div.property-title span a span:nth-child(1)", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping collection due to context cancellation:", ctx.Err())
			return
		default:
			re := regexp.MustCompile(`\d+`)
			match := re.FindString(e.Text)

			if match != "" {
				r.SetBedrooms(match)
			}
		}
	})
}

func (p *PerfilScraper) SetRealStateBathrooms(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("#conteudo div.container div div.col-lg-8.col-md-7 div ul li.banheiros span", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping collection due to context cancellation:", ctx.Err())
			return
		default:
			r.SetBathrooms(e.Text)
		}
	})
}

func (p *PerfilScraper) SetRealStateArea(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div.property-description ul.listing-features li.area span", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping collection due to context cancellation:", ctx.Err())
			return
		default:
			r.SetArea(e.Text)
		}
	})
}

func (p *PerfilScraper) SetRealStateGarageSpaces(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div.property-description ul.listing-features li.vagas span", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping collection due to context cancellation:", ctx.Err())
			return
		default:
			r.SetGarageSpaces(e.Text)
		}
	})
}
