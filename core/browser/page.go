package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/mrbrowser/mrbrowser/telemetry"
)

// Page represents a single browser tab/page.
// All navigation, querying, and interaction happen through a Page.
type Page struct {
	ctx     context.Context
	cancel  context.CancelFunc
	browser *Browser
	log     *telemetry.Logger
	created time.Time
	url     string
	closed  bool
}

// Navigate navigates to the given URL and waits for the page to load.
func (p *Page) Navigate(url string) error {
	p.log.Step("Navigating", telemetry.F("url", url))

	if err := chromedp.Run(p.ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
	); err != nil {
		return fmt.Errorf("navigate to %q: %w", url, err)
	}

	p.url = url
	p.log.Success("Navigated", telemetry.F("url", url))
	return nil
}

// URL returns the current page URL.
func (p *Page) URL() (string, error) {
	var u string
	if err := chromedp.Run(p.ctx,
		chromedp.Location(&u),
	); err != nil {
		return "", fmt.Errorf("get url: %w", err)
	}
	return u, nil
}

// Title returns the current page title.
func (p *Page) Title() (string, error) {
	var t string
	if err := chromedp.Run(p.ctx,
		chromedp.Title(&t),
	); err != nil {
		return "", fmt.Errorf("get title: %w", err)
	}
	return t, nil
}

// Screenshot captures a full-page screenshot and returns the PNG bytes.
func (p *Page) Screenshot() ([]byte, error) {
	p.log.Debug("Taking screenshot")
	var buf []byte
	if err := chromedp.Run(p.ctx,
		chromedp.FullScreenshot(&buf, 90),
	); err != nil {
		return nil, fmt.Errorf("screenshot: %w", err)
	}
	p.log.Debug("Screenshot captured", telemetry.F("bytes", len(buf)))
	return buf, nil
}

// GetHTML returns the full HTML content of the current page.
func (p *Page) GetHTML() (string, error) {
	var html string
	if err := chromedp.Run(p.ctx,
		chromedp.OuterHTML("html", &html, chromedp.ByQuery),
	); err != nil {
		return "", fmt.Errorf("get html: %w", err)
	}
	return html, nil
}

// ExecuteJS runs a JavaScript expression in the page context and returns the result.
func (p *Page) ExecuteJS(script string) (interface{}, error) {
	var result interface{}
	if err := chromedp.Run(p.ctx,
		chromedp.Evaluate(script, &result),
	); err != nil {
		return nil, fmt.Errorf("execute js: %w", err)
	}
	return result, nil
}

// WaitForSelector waits until a CSS selector matches an element.
func (p *Page) WaitForSelector(selector string) error {
	return chromedp.Run(p.ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
	)
}

// WaitForURL waits until the page URL matches (or contains) the given string.
func (p *Page) WaitForURL(urlContains string) error {
	return chromedp.Run(p.ctx,
		chromedp.Poll(fmt.Sprintf("window.location.href.includes(%q)", urlContains), nil, chromedp.WithPollingInterval(100*time.Millisecond)),
	)
}

// Reload reloads the current page.
func (p *Page) Reload() error {
	return chromedp.Run(p.ctx,
		page.Reload(),
		chromedp.WaitReady("body", chromedp.ByQuery),
	)
}

// SetCookies sets cookies on the page.
func (p *Page) SetCookies(cookies []*Cookie) error {
	for _, c := range cookies {
		params := network.SetCookieParams{
			Name:     c.Name,
			Value:    c.Value,
			Domain:   c.Domain,
			Path:     c.Path,
			Secure:   c.Secure,
			HTTPOnly: c.HTTPOnly,
		}
		if err := chromedp.Run(p.ctx,
			chromedp.ActionFunc(func(ctx context.Context) error {
				err := network.SetCookie(params.Name, params.Value).
					WithDomain(params.Domain).
					WithPath(params.Path).
					WithSecure(params.Secure).
					WithHTTPOnly(params.HTTPOnly).
					Do(ctx)
				return err
			}),
		); err != nil {
			return fmt.Errorf("set cookie %q: %w", c.Name, err)
		}
	}
	return nil
}

// GetCookies returns all cookies for the current page.
func (p *Page) GetCookies() ([]*Cookie, error) {
	var rawCookies []*network.Cookie
	if err := chromedp.Run(p.ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			rawCookies, err = network.GetCookies().Do(ctx)
			return err
		}),
	); err != nil {
		return nil, fmt.Errorf("get cookies: %w", err)
	}

	cookies := make([]*Cookie, len(rawCookies))
	for i, c := range rawCookies {
		cookies[i] = &Cookie{
			Name:     c.Name,
			Value:    c.Value,
			Domain:   c.Domain,
			Path:     c.Path,
			Secure:   c.Secure,
			HTTPOnly: c.HTTPOnly,
		}
	}
	return cookies, nil
}

// ClickSelector clicks an element by CSS selector.
// Prefer using actions.Click with an Element for intent-based interaction.
func (p *Page) ClickSelector(selector string) error {
	return chromedp.Run(p.ctx,
		chromedp.Click(selector, chromedp.ByQuery),
	)
}

// TypeSelector types text into an element identified by CSS selector.
func (p *Page) TypeSelector(selector, value string) error {
	return chromedp.Run(p.ctx,
		chromedp.Clear(selector, chromedp.ByQuery),
		chromedp.SendKeys(selector, value, chromedp.ByQuery),
	)
}

// ScrollTo scrolls the page to the given element selector.
func (p *Page) ScrollTo(selector string) error {
	return chromedp.Run(p.ctx,
		chromedp.ScrollIntoView(selector, chromedp.ByQuery),
	)
}

// Ctx returns the underlying chromedp context. Used internally by the action engine.
func (p *Page) Ctx() context.Context {
	return p.ctx
}

// Close closes this page.
func (p *Page) Close() error {
	if err := p.closeInternal(); err != nil {
		return err
	}
	p.browser.removePage(p)
	return nil
}

func (p *Page) closeInternal() error {
	if p.closed {
		return nil
	}
	p.closed = true
	p.cancel()
	return nil
}
