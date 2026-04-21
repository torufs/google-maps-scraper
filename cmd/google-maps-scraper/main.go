// Package main is the entry point for the google-maps-scraper CLI tool.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gosom/google-maps-scraper/pkg/gmaps"
)

const defaultConcurrency = 5

func main() {
	var (
		queries     string
		outputFile  string
		outputFmt   string
		concurrency int
		lang        string
		depth       int
		dsn         string
	)

	flag.StringVar(&queries, "queries", "", "Comma-separated list of search queries (required)")
	flag.StringVar(&outputFile, "output", "output", "Output file path (without extension for most formats)")
	flag.StringVar(&outputFmt, "format", "csv", "Output format: csv, json, jsonl, tsv, excel, sqlite, postgres")
	flag.IntVar(&concurrency, "concurrency", defaultConcurrency, "Number of concurrent scrapers")
	flag.StringVar(&lang, "lang", "en", "Language for results (e.g. en, de, fr)")
	flag.IntVar(&depth, "depth", 10, "Max scroll depth per query")
	flag.StringVar(&dsn, "dsn", "", "DSN for postgres output (required when format=postgres)")
	flag.Parse()

	if queries == "" {
		fmt.Fprintln(os.Stderr, "error: -queries flag is required")
		flag.Usage()
		os.Exit(1)
	}

	if outputFmt == "postgres" && dsn == "" {
		fmt.Fprintln(os.Stderr, "error: -dsn flag is required when format=postgres")
		os.Exit(1)
	}

	queryList := splitAndTrim(queries, ",")
	if len(queryList) == 0 {
		fmt.Fprintln(os.Stderr, "error: no valid queries provided")
		os.Exit(1)
	}

	writer, err := gmaps.NewWriter(outputFmt, outputFile, dsn)
	if err != nil {
		log.Fatalf("failed to create writer: %v", err)
	}
	defer func() {
		if c, ok := writer.(interface{ Close() error }); ok {
			if cerr := c.Close(); cerr != nil {
				log.Printf("warning: error closing writer: %v", cerr)
			}
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	scraper, err := gmaps.NewScraper(gmaps.ScraperOptions{
		Concurrency: concurrency,
		Lang:        lang,
		Depth:       depth,
	})
	if err != nil {
		log.Fatalf("failed to create scraper: %v", err)
	}

	results, err := scraper.Scrape(ctx, queryList)
	if err != nil {
		log.Fatalf("scraping failed: %v", err)
	}

	if err := gmaps.WriteAll(writer, results); err != nil {
		log.Fatalf("failed to write results: %v", err)
	}

	log.Printf("done: wrote %d entries to %s", len(results), outputFile)
}

// splitAndTrim splits s by sep and trims whitespace from each element,
// omitting empty strings.
func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
