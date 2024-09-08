package scraper

import (
	"baia/internal/contracts"
	"context"
	"fmt"
	"sync"
)

type Scraper struct {
	scrapers []map[contracts.RealEstateScraper]string
	jobs     chan ScraperJob
	ch       chan contracts.RealEstate
	wg       sync.WaitGroup
}

type ScraperJob struct {
	scraper contracts.RealEstateScraper
	urls    []string
}

func NewScraper(s []map[contracts.RealEstateScraper]string) *Scraper {
	return &Scraper{
		scrapers: s,
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

func (s *Scraper) runScraper(ctx context.Context, scraperMap map[contracts.RealEstateScraper]string) {
	defer s.wg.Done()

	for scraper := range scraperMap {
		realEstateUrls, _ := scraper.GetRealEstates(ctx, scraperMap[scraper])

		s.wg.Add(len(realEstateUrls))

		s.jobs <- ScraperJob{
			scraper: scraper,
			urls:    realEstateUrls,
		}
	}
}

func (s *Scraper) Run(ctx context.Context) {
	s.wg.Add(len(s.scrapers))

	go s.getRealEstateData(ctx)

	for i := range s.scrapers {
		scraperMap := s.scrapers[i]
		go s.runScraper(ctx, scraperMap)
	}

	s.wg.Wait()

	close(s.jobs)
}
