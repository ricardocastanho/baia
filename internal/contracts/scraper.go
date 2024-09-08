package contracts

import (
	"context"

	"github.com/gocolly/colly/v2"
)

type RealEstateScraper interface {
	GetRealStates(ctx context.Context, url string) ([]string, []string)
	GetRealStateData(ctx context.Context, ch chan RealEstate, re *RealEstate)
	SetRealStateCode(ctx context.Context, c *colly.Collector, re *RealEstate)
	SetRealStateName(ctx context.Context, c *colly.Collector, re *RealEstate)
	SetRealStateDescription(ctx context.Context, c *colly.Collector, re *RealEstate)
	SetRealStatePrice(ctx context.Context, c *colly.Collector, re *RealEstate)
	SetRealStateBedrooms(ctx context.Context, c *colly.Collector, re *RealEstate)
	SetRealStateBathrooms(ctx context.Context, c *colly.Collector, re *RealEstate)
	SetRealStateArea(ctx context.Context, c *colly.Collector, re *RealEstate)
	SetRealStateGarageSpaces(ctx context.Context, c *colly.Collector, re *RealEstate)
}
