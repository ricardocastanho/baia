package contracts

import (
	"context"
)

type RealEstateScraper interface {
	GetRealStates(ctx context.Context, url string) ([]string, []string)
	GetRealStateData(ctx context.Context, ch chan RealState, url string)
}
