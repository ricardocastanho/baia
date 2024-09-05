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

	scrapers := make([]map[contracts.RealEstateScraper]string, 0)

	perfilScraper := perfil.NewPerfilScraper()

	scrapers = append(scrapers, map[contracts.RealEstateScraper]string{
		perfilScraper: "https://www.imobiliariaperfil.imb.br/comprar-imoveis/apartamentos-santo-angelo/",
	})

	s := scraper.NewScraper(scrapers)

	s.Run(ctx)

	fmt.Println("Scraping completed.")
}
