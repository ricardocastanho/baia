package scraper

import (
	"baia/internal/contracts"
	"context"
	"fmt"
	"sync"
)

type ScraperStrategy struct {
	Scraper contracts.RealEstateScraper
	Type    string
	Url     string
}

type Scraper struct {
	strategy []ScraperStrategy
	jobs     chan ScraperJob
	ch       chan contracts.RealEstate
	wg       sync.WaitGroup
}

type ScraperJob struct {
	scraper contracts.RealEstateScraper
	urls    []string
}

func NewScraper(s []ScraperStrategy) *Scraper {
	return &Scraper{
		strategy: s,
		jobs:     make(chan ScraperJob),
		ch:       make(chan contracts.RealEstate),
	}
}

func (s *Scraper) getRealEstateData(ctx context.Context) {
	for job := range s.jobs {
		for _, url := range job.urls {
			go func(url string) {
				defer s.wg.Done()
				realEstate := contracts.RealEstate{Url: url}
				job.scraper.GetRealEstateData(ctx, s.ch, &realEstate)
			}(url)
		}

		select {
		case <-ctx.Done():
			return
		default:
			for data := range s.ch {
				fmt.Println("Real State: ", data)
			}
		}
	}

	go func() {
		s.wg.Wait()
		close(s.ch)
	}()
}

func (s *Scraper) runScraper(ctx context.Context, strategy ScraperStrategy) {
	defer s.wg.Done()

	realEstateUrls, _ := strategy.Scraper.GetRealEstates(ctx, strategy.Url)

	s.wg.Add(len(realEstateUrls))

	s.jobs <- ScraperJob{
		scraper: strategy.Scraper,
		urls:    realEstateUrls,
	}
}

func (s *Scraper) Run(ctx context.Context) {
	s.wg.Add(len(s.strategy))

	go s.getRealEstateData(ctx)

	for i := range s.strategy {
		strategy := s.strategy[i]
		go s.runScraper(ctx, strategy)
	}

	s.wg.Wait()

	close(s.jobs)
}
