package scraper

import (
	"baia/internal/contracts"
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ScraperStrategy struct {
	Scraper contracts.RealEstateScraper
	Type    string
	Url     string
	ForSale bool
	ForRent bool
}

type Scraper struct {
	driver      neo4j.DriverWithContext
	logger      *slog.Logger
	strategy    []ScraperStrategy
	jobs        chan ScraperJob
	ch          chan contracts.RealEstate
	wg          sync.WaitGroup
	scrapedUrls map[string]bool
}

type ScraperJob struct {
	scraper contracts.RealEstateScraper
	urls    []string
	Type    string
	ForSale bool
	ForRent bool
}

func NewScraper(driver neo4j.DriverWithContext, logger *slog.Logger, s []ScraperStrategy) *Scraper {
	return &Scraper{
		driver:      driver,
		logger:      logger,
		strategy:    s,
		jobs:        make(chan ScraperJob),
		ch:          make(chan contracts.RealEstate),
		scrapedUrls: make(map[string]bool),
	}
}

func (s *Scraper) getRealEstate(ctx context.Context) {
	for job := range s.jobs {
		for _, url := range job.urls {
			go func(url string) {
				defer s.wg.Done()

				_, ok := s.scrapedUrls[url]
				if ok {
					return
				}

				realEstate := contracts.RealEstate{Url: url, Type: job.Type, ForSale: job.ForSale, ForRent: job.ForRent}
				job.scraper.GetRealEstate(ctx, s.ch, &realEstate)
				s.scrapedUrls[url] = true
			}(url)
			time.Sleep(time.Second * 3)
		}

		select {
		case <-ctx.Done():
			return
		default:
			for re := range s.ch {
				err := re.Save(ctx, s.driver)
				if err != nil {
					s.logger.Error(fmt.Sprintf("Error while trying to save property: %v", err))
				}
			}
		}
	}

	go func() {
		s.wg.Wait()
		close(s.ch)
	}()
}

func (s *Scraper) getNextPages(ctx context.Context, strategy ScraperStrategy, nextPages []string) {
	for _, newUrl := range nextPages {
		_, ok := s.scrapedUrls[newUrl]

		if ok {
			continue
		}

		s.wg.Add(1)
		s.scrapedUrls[newUrl] = true

		go s.runScraper(ctx, ScraperStrategy{
			Scraper: strategy.Scraper,
			Type:    strategy.Type,
			ForSale: strategy.ForSale,
			ForRent: strategy.ForRent,
			Url:     newUrl,
		})
	}
}

func (s *Scraper) runScraper(ctx context.Context, strategy ScraperStrategy) {
	defer func() {
		s.wg.Done()
	}()

	realEstateUrls, nextPages := strategy.Scraper.GetRealEstateUrls(ctx, strategy.Url)
	s.scrapedUrls[strategy.Url] = true

	s.wg.Add(len(realEstateUrls))

	s.jobs <- ScraperJob{
		scraper: strategy.Scraper,
		urls:    realEstateUrls,
		Type:    strategy.Type,
		ForSale: strategy.ForSale,
		ForRent: strategy.ForRent,
	}

	s.getNextPages(ctx, strategy, nextPages)
}

func (s *Scraper) Run(ctx context.Context) {
	s.wg.Add(len(s.strategy))

	go s.getRealEstate(ctx)

	for i := range s.strategy {
		strategy := s.strategy[i]
		go s.runScraper(ctx, strategy)
	}

	s.wg.Wait()

	close(s.jobs)
}
