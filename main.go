package main

import (
	"fmt"
	"sync"

	"baia/internal/scraper/perfil"
	"baia/internal/utils"
)

// runScraper runs a given scraper function with a specific URL and prints scraped URLs.
func runScraper(wg *sync.WaitGroup, scraperFunc func() ([]string, []string)) {
	defer wg.Done()

	realStateUrls, nextPages := scraperFunc()

	fmt.Println("Urls:", realStateUrls)
	fmt.Println("next pages:", nextPages)
}

func main() {
	ctx, cancel := utils.NewCancelableContext()
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	perfilScraperFunc := func() ([]string, []string) {
		perfilScraper := perfil.NewPerfilScraper()
		return perfilScraper.GetRealStates(ctx, "https://www.imobiliariaperfil.imb.br/comprar-imoveis/apartamentos-santo-angelo/")
	}

	go runScraper(&wg, perfilScraperFunc)

	wg.Wait()

	fmt.Println("Scraping completed.")
}
