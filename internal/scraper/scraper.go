package scraper

import (
	"context"
)

type RealEstateScraper interface {
	Run(ctx context.Context, url string)
}

type GetRealEstateScraper interface {
	Run(ctx context.Context, urls []string)
}
