package browser

import (
	"time"

	"github.com/mrbrowser/mrbrowser/telemetry"
)

// Session ties together a Browser, a Page, and runtime context for a single automation run.
// It is the primary object handed to a task executor.
type Session struct {
	ID      string
	Browser *Browser
	Page    *Page
	log     *telemetry.Logger
	started time.Time
	meta    map[string]interface{}
}

// NewSession creates a new Session by launching a browser and opening a page.
func NewSession(id string, opts Options) (*Session, error) {
	b, err := Launch(opts)
	if err != nil {
		return nil, err
	}

	pg, err := b.NewPage()
	if err != nil {
		_ = b.Close()
		return nil, err
	}

	return &Session{
		ID:      id,
		Browser: b,
		Page:    pg,
		log:     telemetry.New("session"),
		started: time.Now(),
		meta:    make(map[string]interface{}),
	}, nil
}

// Set stores a key-value pair in the session context.
// Useful for passing data between task steps (e.g., extracted values).
func (s *Session) Set(key string, value interface{}) {
	s.meta[key] = value
}

// Get retrieves a value from the session context.
func (s *Session) Get(key string) (interface{}, bool) {
	v, ok := s.meta[key]
	return v, ok
}

// Duration returns how long this session has been running.
func (s *Session) Duration() time.Duration {
	return time.Since(s.started)
}

// Close closes the page and browser for this session.
func (s *Session) Close() error {
	s.log.Info("Closing session", telemetry.F("id", s.ID), telemetry.F("duration", s.Duration()))
	return s.Browser.Close()
}
