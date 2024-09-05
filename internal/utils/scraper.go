package utils

import (
	"context"
)

type RealEstateScraper interface {
	GetRealStates(ctx context.Context, url string) ([]string, []string)
}
