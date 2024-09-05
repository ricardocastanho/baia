package scraper

import (
	"context"
)

type RealEstateScraper interface {
	Run(ctx context.Context, url string) ([]string, []string)
}
