package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"baia/internal/contracts"
	"baia/internal/scraper/perfil"
	"baia/internal/utils"
	"baia/pkg/database"

	"github.com/joho/godotenv"
	"github.com/ricardocastanho/scrapify"
)

func main() {
	ctx, cancel := utils.NewCancelableContext()
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	uri := os.Getenv("NEO4J_URI")
	username := os.Getenv("NEO4J_USERNAME")
	password := os.Getenv("NEO4J_PASSWORD")

	if uri == "" || username == "" || password == "" {
		log.Fatal("Wrong Neo4j credentials in .env")
	}

	client := database.NewNeo4jClient(uri, username, password)

	driver, err := client.GetDriver()
	if err != nil {
		log.Fatalf("Failed to get Neo4j driver: %v", err)
	}
	defer client.Close()

	if err := driver.VerifyConnectivity(context.Background()); err != nil {
		log.Fatalf("Neo4j connection failed: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	perfilScraper := perfil.NewPerfilScraper(logger)

	strategy := []scrapify.ScraperStrategy[contracts.RealEstate]{
		{
			Scraper: perfilScraper,
			Url:     "https://www.imobiliariaperfil.imb.br/comprar-imoveis/apartamentos-santo-angelo/&pg=1",
		},
	}

	callback := func(data contracts.RealEstate) {
		logger.Info("Saving data in database:", "data", data)
		data.Save(ctx, driver)
	}

	scraper := scrapify.NewScraper(strategy, callback, time.Second*2)
	scraper.Run(ctx)

	logger.Info("Scraping completed.")
}
