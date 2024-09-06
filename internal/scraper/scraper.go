package scraper

import (
	"baia/internal/contracts"
	"context"
	"fmt"
	"sync"
)

type Scraper struct {
	scrapers []map[contracts.RealEstateScraper]string
}

type ScraperJob struct {
	scraper contracts.RealEstateScraper
	urls    []string
}

func NewScraper(s []map[contracts.RealEstateScraper]string) *Scraper {
	return &Scraper{
		scrapers: s,
	}
}

func (s *Scraper) getRealStateData(ctx context.Context, jobs <-chan ScraperJob) {
	ch := make(chan contracts.RealState)
	var wg sync.WaitGroup

	for job := range jobs {
		wg.Add(len(job.urls))

		for _, url := range job.urls {
			go func(url string) {
				defer wg.Done()
				job.scraper.GetRealStateData(ctx, ch, url)
			}(url)
		}

		select {
		case <-ctx.Done():
			return
		default:
			for data := range ch {
				fmt.Println("Real State data:", data)
			}
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()
}

func (s *Scraper) runScraper(ctx context.Context, wg *sync.WaitGroup, scraperMap map[contracts.RealEstateScraper]string, jobs chan<- ScraperJob) {
	defer wg.Done()

	for scraper := range scraperMap {
		realStateUrls, nextPages := scraper.GetRealStates(ctx, scraperMap[scraper])

		fmt.Println("Urls:", realStateUrls)
		fmt.Println("Next pages:", nextPages)

		jobs <- ScraperJob{
			scraper: scraper,
			urls:    realStateUrls,
		}
	}
}

func (s *Scraper) Run(ctx context.Context) {
	var wg sync.WaitGroup

	jobs := make(chan ScraperJob)

	go s.getRealStateData(ctx, jobs)

	wg.Add(len(s.scrapers))

	for i := range s.scrapers {
		scraperMap := s.scrapers[i]
		go s.runScraper(ctx, &wg, scraperMap, jobs)
	}

	wg.Wait()

	close(jobs)
}
