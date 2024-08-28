package main

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly/v2"
)

func getRealStateUrls(ch chan string, wg *sync.WaitGroup, url string) {
	c := colly.NewCollector()

	c.OnHTML("div#grid div.listing-item a[href]", func(e *colly.HTMLElement) {
		ch <- e.Attr("href")
	})

	c.OnScraped(func(r *colly.Response) {
		wg.Done()
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(url)
}

func main() {
	numGoroutines := 1

	ch := make(chan string)
	defer close(ch)

	var wg sync.WaitGroup

	wg.Add(numGoroutines)

	go getRealStateUrls(ch, &wg, "https://www.imobiliariaperfil.imb.br/comprar-imoveis/apartamentos-santo-angelo/")

	for link := range ch {
		fmt.Println(link)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()
}
