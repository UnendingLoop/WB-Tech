// Package cmd is an entry point to wGet
package cmd

import (
	"flag"
	"log"
	"net/url"
	"os"

	"miniWget/crawler"
)

func InitWget() {
	parser := flag.NewFlagSet("wGet-clone with website mirroring function only", flag.ExitOnError)

	t := parser.Int("timeout", 0, "--timeout <N> - устанавливает таймаут для каждого URL-запроса")
	r := parser.Int("r", 0, "-r <N> - задает глубину рекурсии скачивания страниц")

	if err := parser.Parse(os.Args[1:]); err != nil {
		log.Fatalf("Failed to parse os.Args: %v", err)
	}

	if parser.NArg() < 1 {
		log.Fatal("Usage: wget-clone [options] <URL>")
	}
	rawURL := parser.Arg(0)

	webURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
		return
	}

	config := &crawler.Crawler{
		InitURL:     rawURL,
		Domain:      webURL.Hostname(),
		MaxDepth:    *r,
		Timeout:     *t,
		Downloaded:  make(map[string]bool),
		RobotAllows: make(map[string]bool),
		Queue:       []crawler.QueueItem{},
	}

	config.LoadRobots()
	config.Crawl()
}
