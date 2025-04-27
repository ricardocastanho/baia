package perfil

import (
	"baia/internal/contracts"
	"baia/pkg/collector"
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

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
func (p *PerfilScraper) GetData(ctx context.Context, ch chan<- contracts.RealEstate, re *contracts.RealEstate, url string) {
	c := collector.NewCollector(p.logger)

	p.SetRealEstateCode(ctx, c, re)
	p.SetRealEstateName(ctx, c, re)
	p.SetRealEstateDescription(ctx, c, re)
	p.SetRealEstatePrice(ctx, c, re)
	p.SetRealEstateBedrooms(ctx, c, re)
	p.SetRealEstateBathrooms(ctx, c, re)
	p.SetRealEstateArea(ctx, c, re)
	p.SetRealEstateGarageSpaces(ctx, c, re)
	p.SetRealEstateDistrict(ctx, c, re)
	p.SetRealEstateCity(ctx, c, re)
	p.SetRealEstateFurnished(ctx, c, re)
	p.SetRealEstateYearBuilt(ctx, c, re)
	p.SetRealEstatePhotos(ctx, c, re)
	p.SetRealEstateTags(ctx, c, re)

	c.OnScraped(func(c *colly.Response) {
		ch <- *re
	})

	select {
	case <-ctx.Done():
		p.logger.Debug(fmt.Sprint("Stopping visit due to context cancellation:", ctx.Err()))
	default:
		re.Url = url
		re.Agency = "Perfil"

		if strings.Contains(url, "alugar") || strings.Contains(url, "locacao") {
			re.ForRent = true
		} else {
			re.ForSale = true
		}

		if strings.Contains(url, "apartamento") {
			re.Type = contracts.Apartment
		} else if strings.Contains(url, "casa") {
			re.Type = contracts.House
		} else if strings.Contains(url, "terreno") {
			re.Type = contracts.Land
		}

		c.Visit(url)
	}
}

func (p *PerfilScraper) SetRealEstateCode(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div.property-title h2 span.imovel-codigo", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			p.logger.Debug(fmt.Sprint("Stopping collection due to context cancellation:", ctx.Err()))
			return
		default:
			code := strings.TrimSpace(strings.Replace(e.Text, "Cód.", "", 1))
			r.SetCode(code)
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

func (p *PerfilScraper) SetRealEstateDistrict(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
}

func (p *PerfilScraper) SetRealEstateCity(ctx context.Context, c *colly.Collector, r *contracts.RealEstate) {
	c.OnHTML("div.property-title span a span[data-tag='address']", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			p.logger.Debug(fmt.Sprint("Stopping collection due to context cancellation:", ctx.Err()))
			return
		default:
			span := e.DOM

			i := span.Find("i")

			i.Remove()

			r.SetCity(strings.Split(span.Text(), " / ")[0])
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
