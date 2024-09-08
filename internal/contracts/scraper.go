package contracts

import (
	"context"

	"github.com/gocolly/colly/v2"
)

type RealEstateScraper interface {
	GetRealEstates(ctx context.Context, url string) ([]string, []string)
	GetRealEstateData(ctx context.Context, ch chan RealEstate, re *RealEstate)
	SetRealEstateCode(ctx context.Context, c *colly.Collector, re *RealEstate)
	SetRealEstateName(ctx context.Context, c *colly.Collector, re *RealEstate)
	SetRealEstateDescription(ctx context.Context, c *colly.Collector, re *RealEstate)
	SetRealEstatePrice(ctx context.Context, c *colly.Collector, re *RealEstate)
	SetRealEstateBedrooms(ctx context.Context, c *colly.Collector, re *RealEstate)
	SetRealEstateBathrooms(ctx context.Context, c *colly.Collector, re *RealEstate)
	SetRealEstateArea(ctx context.Context, c *colly.Collector, re *RealEstate)
	SetRealEstateGarageSpaces(ctx context.Context, c *colly.Collector, re *RealEstate)
	SetRealEstateLocation(ctx context.Context, c *colly.Collector, re *RealEstate)
	SetRealEstateFurnished(ctx context.Context, c *colly.Collector, re *RealEstate)
	SetRealEstateYearBuilt(ctx context.Context, c *colly.Collector, re *RealEstate)
}
