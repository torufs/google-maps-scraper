package gmaps

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/playwright-community/playwright-go"
)

// Scraper holds configuration and state for scraping Google Maps.
type Scraper struct {
	// Concurrency controls how many browser tabs run in parallel.
	Concurrency int
	// MaxDepth limits how many result pages are followed (0 = unlimited).
	MaxDepth int
	// Lang is the language code passed to Google Maps (e.g. "en").
	Lang string
	// Logger is the structured logger used for diagnostic output.
	Logger *slog.Logger

	pw      *playwright.Playwright
	browser playwright.Browser
}

// NewScraper creates a Scraper with sensible defaults.
func NewScraper(opts ...Option) (*Scraper, error) {
	s := &Scraper{
		// Increased default concurrency to 3 for faster scraping on my machine.
		Concurrency: 3,
		MaxDepth:    0,
		Lang:        "en",
		Logger:      slog.Default(),
	}
	for _, o := range opts {
		o(s)
	}
	return s, nil
}

// Start initialises Playwright and launches a Chromium browser instance.
// Call Close when done to release resources.
func (s *Scraper) Start() error {
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("start playwright: %w", err)
	}
	s.pw = pw

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		_ = pw.Stop()
		return fmt.Errorf("launch browser: %w", err)
	}
	s.browser = browser
	return nil
}

// Close shuts down the browser and Playwright runtime.
func (s *Scraper) Close() {
	if s.browser != nil {
		_ = s.browser.Close()
	}
	if s.pw != nil {
		_ = s.pw.Stop()
	}
}

// Scrape searches Google Maps for the given query and sends each discovered
// Entry to the results channel. The channel is closed when scraping finishes
// or the context is cancelled.
func (s *Scraper) Scrape(ctx context.Context, query string, results chan<- Entry) error {
	defer close(results)

	page, err := s.browser.NewPage()
	if err != nil {
		return fmt.Errorf("new page: %w", err)
	}
	defer page.Close()

	url := fmt.Sprintf(
		"https://www.google.com/maps/search/%s/?hl=%s",
		playwright.String(query), s.Lang,
	)

	if _, err := page.Goto(*playwright.String(url), playwright.PageGotoOptions{
		Timeout:   playwright.Float(30_000),
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		return fmt.Errorf("navigate to %s: %w", url, err)
	}

	s.Logger.InfoContext(ctx, "page loaded", "query", query, "url", url)

	// Scroll and collect results until the end-of-list sentinel appears or
	// MaxDepth pages have been processed.
	depth := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		entries, err := extractEntries(page)
		if err != nil {
			s.Logger.WarnContext(ctx, "extract entries", "err", err)
		}
		for _, e := range entries {
			if e.IsValid() {
				results <- e
			}
		}

		depth++
		if s.MaxDepth > 0 && depth >= s.MaxDepth {
			break
		}

		reachedEnd, scrollErr := scrollResultsList(ctx, page)
		if scrollErr != nil {
			s.Logger.WarnContext(ctx, "scroll", "err", scrollErr)
			break
		}
		if r