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
	ch       chan contracts.RealState
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
		ch:       make(chan contracts.RealState),
	}
}

func (s *Scraper) getRealStateData(ctx context.Context) {
	for job := range s.jobs {
		for _, url := range job.urls {
			go func(url string) {
				defer s.wg.Done()
				job.scraper.GetRealStateData(ctx, s.ch, url)
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
		realStateUrls, _ := scraper.GetRealStates(ctx, scraperMap[scraper])

		s.wg.Add(len(realStateUrls))

		s.jobs <- ScraperJob{
			scraper: scraper,
			urls:    realStateUrls,
		}
	}
}

func (s *Scraper) Run(ctx context.Context) {
	s.wg.Add(len(s.scrapers))

	go s.getRealStateData(ctx)

	for i := range s.scrapers {
		scraperMap := s.scrapers[i]
		go s.runScraper(ctx, scraperMap)
	}

	s.wg.Wait()

	close(s.jobs)
}
