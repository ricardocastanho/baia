package contracts

import (
	"context"

	"github.com/gocolly/colly/v2"
)

type RealEstateScraper interface {
	GetRealStates(ctx context.Context, url string) ([]string, []string)
	GetRealStateData(ctx context.Context, ch chan RealState, url string)
	SetRealStateCode(ctx context.Context, c *colly.Collector)
	SetRealStateName(ctx context.Context, c *colly.Collector)
	SetRealStateDescription(ctx context.Context, c *colly.Collector)
	SetRealStatePrice(ctx context.Context, c *colly.Collector)
	SetRealStateBedrooms(ctx context.Context, c *colly.Collector)
}
