package main

import (
	"fmt"

	"baia/internal/contracts"
	"baia/internal/scraper"
	"baia/internal/scraper/perfil"
	"baia/internal/utils"
)

func main() {
	ctx, cancel := utils.NewCancelableContext()
	defer cancel()

	scrapers := make([]scraper.ScraperStrategy, 0)

	perfilScraper := perfil.NewPerfilScraper()

	scrapers = append(scrapers, scraper.ScraperStrategy{
		Scraper: perfilScraper,
		Url:     "https://www.imobiliariaperfil.imb.br/comprar-imoveis/apartamentos-santo-angelo/",
		Type:    contracts.Apartment,
		ForSale: true,
	})

	s := scraper.NewScraper(scrapers)

	s.Run(ctx)

	fmt.Println("Scraping completed.")
}
