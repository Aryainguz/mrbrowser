// Package browser provides the core browser abstraction for Mr. Browser.
// It manages Chromium lifecycle, tabs, and the Page abstraction via CDP.
package browser

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/mrbrowser/mrbrowser/telemetry"
)

// Browser manages a Chromium instance and its pages.
type Browser struct {
	opts        Options
	allocCtx    context.Context
	allocCancel context.CancelFunc
	pages       []*Page
	mu          sync.Mutex
	log         *telemetry.Logger
	closed      bool
}

// Launch creates and starts a new Chromium browser instance.
func Launch(opts Options) (*Browser, error) {
	opts = opts.withDefaults()

	log := telemetry.New("browser")
	log.Info("Launching Chromium", telemetry.F("headless", opts.Headless), telemetry.F("path", opts.ChromiumPath))

	allocOpts := chromedp.DefaultExecAllocatorOptions[:]

	if opts.Headless {
		allocOpts = append(allocOpts,
			chromedp.Headless,
			chromedp.Flag("disable-gpu", true),
		)
	} else {
		allocOpts = append(allocOpts, chromedp.Flag("headless", false))
	}

	if opts.ChromiumPath != "" {
		allocOpts = append(allocOpts, chromedp.ExecPath(opts.ChromiumPath))
	}

	if opts.NoSandbox {
		allocOpts = append(allocOpts, chromedp.NoSandbox)
	}

	if opts.DisableExtensions {
		allocOpts = append(allocOpts, chromedp.Flag("disable-extensions", true))
	}

	allocOpts = append(allocOpts,
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-features", "TranslateUI"),
		chromedp.WindowSize(opts.WindowWidth, opts.WindowHeight),
	)

	if opts.UserAgent != "" {
		allocOpts = append(allocOpts, chromedp.UserAgent(opts.UserAgent))
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), allocOpts...)

	b := &Browser{
		opts:        opts,
		allocCtx:    allocCtx,
		allocCancel: allocCancel,
		log:         log,
	}

	log.Success("Chromium launched")
	return b, nil
}

// NewPage opens a new browser tab and returns a Page.
func (b *Browser) NewPage() (*Page, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil, fmt.Errorf("browser is closed")
	}

	ctx, cancel := chromedp.NewContext(b.allocCtx)

	// Set a default timeout for all operations
	ctx, timeoutCancel := context.WithTimeout(ctx, b.opts.DefaultTimeout)
	_ = timeoutCancel // page.Close() will cancel this

	// Trigger the browser to actually connect by running a no-op
	if err := chromedp.Run(ctx); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to open new page: %w", err)
	}

	p := &Page{
		ctx:     ctx,
		cancel:  cancel,
		browser: b,
		log:     telemetry.New("page"),
		created: time.Now(),
	}

	b.pages = append(b.pages, p)
	b.log.Info("New page opened")
	return p, nil
}

// Close shuts down the browser and all open pages.
func (b *Browser) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil
	}

	b.log.Info("Closing browser")
	b.closed = true

	for _, p := range b.pages {
		_ = p.closeInternal()
	}

	b.allocCancel()
	b.log.Success("Browser closed")
	return nil
}

// PageCount returns the number of open pages.
func (b *Browser) PageCount() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.pages)
}

// removePage removes a closed page from the browser's page list.
func (b *Browser) removePage(p *Page) {
	b.mu.Lock()
	defer b.mu.Unlock()
	filtered := make([]*Page, 0, len(b.pages))
	for _, pg := range b.pages {
		if pg != p {
			filtered = append(filtered, pg)
		}
	}
	b.pages = filtered
}
