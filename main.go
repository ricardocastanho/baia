package main

import (
	"fmt"
	"sync"
	"time"

	"baia/internal/scraper"
	"baia/internal/utils"
)

// runScraper runs a given scraper function with a specific URL and prints scraped URLs.
func runScraper(wg *sync.WaitGroup, scraperFunc func(chan string)) {
	defer wg.Done()

	ch := make(chan string)

	go scraperFunc(ch)

	for link := range ch {
		fmt.Println("Scraped URL:", link)
	}
}

func main() {
	ctx, cancel := utils.NewTimeoutContext(10 * time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	perfilScraperFunc := func(ch chan string) {
		perfilScraper := scraper.NewPerfilScraper(ch)
		perfilScraper.Run(ctx, "https://www.imobiliariaperfil.imb.br/comprar-imoveis/apartamentos-santo-angelo/")
	}

	go runScraper(&wg, perfilScraperFunc)

	wg.Wait()

	fmt.Println("Scraping completed.")
}
