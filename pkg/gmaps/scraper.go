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
		Concurrency: 1,
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
		if reachedEnd {
			break
		}

		// Brief pause to allow dynamic content to render.
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}

	return nil
}

// extractEntries reads result cards currently visible on the page.
func extractEntries(page playwright.Page) ([]Entry, error) {
	// Placeholder: real implementation would evaluate JS / parse DOM nodes.
	_ = page
	return nil, nil
}

// scrollResultsList scrolls the sidebar result list down one viewport.
// It returns true when the end-of-results sentinel is detected.
func scrollResultsList(ctx context.Context, page playwright.Page) (bool, error) {
	_ = ctx
	_, err := page.Evaluate(`() => {
		const list = document.querySelector('[role="feed"]');
		if (list) list.scrollBy(0, list.clientHeight);
	}`)
	if err != nil {
		return false, fmt.Errorf("scroll eval: %w", err)
	}

	// Check for the "You've reached the end of the list" element.
	endEl, _ := page.QuerySelector(".HlvSq")
	return endEl != nil, nil
}
