package scraper

import (
	"baia/internal/contracts"
	"context"
	"fmt"
)

type Scraper struct {
	scrapers []map[contracts.RealEstateScraper]string
}

func NewScraper(s []map[contracts.RealEstateScraper]string) *Scraper {
	return &Scraper{
		scrapers: s,
	}
}

func (s *Scraper) Run(ctx context.Context) {
	for i := range s.scrapers {
		scraperMap := s.scrapers[i]

		for scraper := range scraperMap {
			realStateUrls, nextPages := scraper.GetRealStates(ctx, scraperMap[scraper])

			fmt.Println("Urls:", realStateUrls)
			fmt.Println("Next pages:", nextPages)
		}
	}
}
