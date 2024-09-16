package perfil

import (
	"baia/internal/contracts"
	"baia/pkg/collector"
	"context"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/gocolly/colly/v2"
	"github.com/ricardocastanho/scrapify"
)

// PerfilScraper implements the RealEstateScraperInterface for the "Perfil" real estate website.
type PerfilScraper struct {
	logger *slog.Logger
}

// NewPerfilScraper creates a new instance of PerfilScraper.
func NewPerfilScraper(logger *slog.Logger) scrapify.IScraper[contracts.RealEstate] {
	return &PerfilScraper{
		logger: logger,
	}
}

// GetRealEstateUrls starts the scraping process for the given URL using the provided context.
func (p *PerfilScraper) GetUrls(ctx context.Context, url string) ([]string, []string) {
	var (
		realEstateurls = []string{}
		nextPages      = []string{}
	)

	c := collector.NewCollector(p.logger)

	c.OnHTML("div#grid div.listing-item a[href]", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			p.logger.Debug(fmt.Sprint("Stopping collection due to context cancellation:", ctx.Err()))
			return
		default:
			realEstateurls = append(realEstateurls, e.Attr("href"))
		}
	})

	c.OnHTML("#conteudo > div.container.lista-imoveis-container > div.pagination-container > nav > ul > li > a", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			p.logger.Debug(fmt.Sprint("Stopping collection due to context cancellation:", ctx.Err()))
			return
		default:
			nextPages = append(nextPages, e.Attr("href"))
		}
	})

	select {
	case <-ctx.Done():
		p.logger.Debug(fmt.Sprint("Stopping visit due to context cancellation:", ctx.Err()))
		return nil, nil
	default:
		c.Visit(url)

		return realEstateurls, nextPages
	}
}

// GetRealEstate gets all the data from a given url
func (p *PerfilScraper) GetData(ctx context.Context, ch chan<- contracts.RealEstate, data *contracts.RealEstate, url string) {
	c := collector.NewCollector(p.logger)

	p.SetRealEstateCode(ctx, c, data)
	p.SetRealEstateName(ctx, c, data)
	p.SetRealEstateDescription(ctx, c, data)
	p.SetRealEstatePrice(ctx, c, data)
	p.SetRealEstateBedrooms(ctx, c, data)
	p.SetRealEstateBathrooms(ctx, c, data)
	p.SetRealEstateArea(ctx, c, data)
	p.SetRealEstateGarageSpaces(ctx, c, data)
	p.SetRealEstateLocation(ctx, c, data)
	p.SetRealEstateFurnished(ctx, c, data)
	p.SetRealEstateYearBuilt(ctx, c, data)
	p.SetRealEstatePhotos(ctx, c, data)
	p.SetRealEstateTags(ctx, c, data)

	data.Url = url

	c.OnScraped(func(c *colly.Response) {
		ch <- *data
	})

	select {
	case <-ctx.Done():
		p.logger.Debug(fmt.Sprint("Stopping visit due to context cancellation:", ctx.Err()))
	default:
		c.Visit(url)
		return
	}
}

func (p *PerfilScraper) SetRealEstateCode(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div.property-title h2 span.imovel-codigo", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			p.logger.Debug(fmt.Sprint("Stopping collection due to context cancellation:", ctx.Err()))
			return
		default:
			r.SetCode(e.Text)
		}
	})
}

func (p *PerfilScraper) SetRealEstateName(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div.property-title", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			p.logger.Debug(fmt.Sprint("Stopping collection due to context cancellation:", ctx.Err()))
			return
		default:
			h2 := e.DOM.Find("h2")
			span := h2.Find("span")

			span.Remove()

			r.SetName(h2.Text())
		}
	})
}

func (p *PerfilScraper) SetRealEstateDescription(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div#text-0 div p", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			p.logger.Debug(fmt.Sprint("Stopping collection due to context cancellation:", ctx.Err()))
			return
		default:
			r.SetDescription(e.Text)
		}
	})
}

func (p *PerfilScraper) SetRealEstatePrice(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div.valor-imovel span", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			p.logger.Debug(fmt.Sprint("Stopping collection due to context cancellation:", ctx.Err()))
			return
		default:
			err := r.SetPrice(e.Text)

			if err != nil {
				p.logger.Error(fmt.Sprint("Error while trying to parse real state price:", err))
				return
			}
		}
	})
}

func (p *PerfilScraper) SetRealEstateBedrooms(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div.property-title span a span:nth-child(1)", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			p.logger.Debug(fmt.Sprint("Stopping collection due to context cancellation:", ctx.Err()))
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

func (p *PerfilScraper) SetRealEstateBathrooms(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("#conteudo div.container div div.col-lg-8.col-md-7 div ul li.banheiros span", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			p.logger.Debug(fmt.Sprint("Stopping collection due to context cancellation:", ctx.Err()))
			return
		default:
			r.SetBathrooms(e.Text)
		}
	})
}

func (p *PerfilScraper) SetRealEstateArea(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div.property-description ul.listing-features li.area span", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			p.logger.Debug(fmt.Sprint("Stopping collection due to context cancellation:", ctx.Err()))
			return
		default:
			r.SetArea(e.Text)
		}
	})
}

func (p *PerfilScraper) SetRealEstateGarageSpaces(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div.property-description ul.listing-features li.vagas span", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			p.logger.Debug(fmt.Sprint("Stopping collection due to context cancellation:", ctx.Err()))
			return
		default:
			r.SetGarageSpaces(e.Text)
		}
	})
}

func (p *PerfilScraper) SetRealEstateLocation(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div.property-title span a span[data-tag='address']", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			p.logger.Debug(fmt.Sprint("Stopping collection due to context cancellation:", ctx.Err()))
			return
		default:
			span := e.DOM

			i := span.Find("i")

			i.Remove()

			r.SetLocation(span.Text())
		}
	})
}

func (p *PerfilScraper) SetRealEstateFurnished(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div.property-description ul.listing-features li.mobilia span", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			p.logger.Debug(fmt.Sprint("Stopping collection due to context cancellation:", ctx.Err()))
			return
		default:
			isFurnished := e.Text == "Semi" || e.Text == "Sim"
			r.SetFurnished(isFurnished)
		}
	})
}

func (p *PerfilScraper) SetRealEstateYearBuilt(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
}

func (p *PerfilScraper) SetRealEstatePhotos(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("img.sp-image", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			p.logger.Debug(fmt.Sprint("Stopping collection due to context cancellation:", ctx.Err()))
			return
		default:
			src := e.Attr("src")
			r.SetPhoto(src)
		}
	})
}

func (p *PerfilScraper) SetRealEstateTags(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("ul.property-features li", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			p.logger.Debug(fmt.Sprint("Stopping collection due to context cancellation:", ctx.Err()))
			return
		default:
			r.SetTag(e.Text)
		}
	})
}
