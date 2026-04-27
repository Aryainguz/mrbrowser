package browser

import "time"

// Options configures a Chromium browser instance.
type Options struct {
	// Headless runs the browser without a visible window (default: true).
	Headless bool

	// ChromiumPath is the path to the Chromium/Chrome executable.
	// If empty, chromedp will find it automatically.
	ChromiumPath string

	// NoSandbox disables the Chrome sandbox (required in Docker).
	NoSandbox bool

	// DisableExtensions disables Chrome extensions.
	DisableExtensions bool

	// WindowWidth and WindowHeight set the browser window size.
	WindowWidth  int
	WindowHeight int

	// UserAgent overrides the browser user agent string.
	UserAgent string

	// DefaultTimeout is the maximum time for any single browser operation.
	DefaultTimeout time.Duration

	// PageLoadTimeout is the maximum time to wait for a page to load.
	PageLoadTimeout time.Duration

	// DownloadDir is the directory where downloaded files are saved.
	DownloadDir string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Headless:          true,
		NoSandbox:         false,
		DisableExtensions: true,
		WindowWidth:       1920,
		WindowHeight:      1080,
		DefaultTimeout:    30 * time.Second,
		PageLoadTimeout:   30 * time.Second,
	}
}

// withDefaults fills in zero-value fields with defaults.
func (o Options) withDefaults() Options {
	d := DefaultOptions()
	if o.WindowWidth == 0 {
		o.WindowWidth = d.WindowWidth
	}
	if o.WindowHeight == 0 {
		o.WindowHeight = d.WindowHeight
	}
	if o.DefaultTimeout == 0 {
		o.DefaultTimeout = d.DefaultTimeout
	}
	if o.PageLoadTimeout == 0 {
		o.PageLoadTimeout = d.PageLoadTimeout
	}
	return o
}
