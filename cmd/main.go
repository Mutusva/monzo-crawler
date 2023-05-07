package main

import (
	"flag"
	"github.com/Mutusva/monzo-webcrawler/crawler"
)

func main() {
	seedUrl := flag.String("seed_url", "https://monzo.com/", "The seed url to start crawling")
	shouldCrawlerExternal := flag.Bool("ext", false, "Can crawl external sites")
	flag.Parse()

	seed := []string{*seedUrl}
	cr := crawler.New(seed)
	err := cr.Start(*shouldCrawlerExternal)
	_ = err
}
