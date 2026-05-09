// Package server provides the HTTP REST API for Mr. Browser.
// It allows the Python SDK and remote clients to control the engine.
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/mrbrowser/mrbrowser/core/browser"
	"github.com/mrbrowser/mrbrowser/core/runtime"
	"github.com/mrbrowser/mrbrowser/memory"
	"github.com/mrbrowser/mrbrowser/telemetry"
)

// Server is the Mr. Browser REST API server.
type Server struct {
	port       int
	httpServer *http.Server
	opts       browser.Options
	store      *memory.Store
	log        *telemetry.Logger

	mu       sync.Mutex
	sessions map[string]*SessionContext
}

// SessionContext holds the state for a single API-driven browser session.
type SessionContext struct {
	ID       string
	Session  *browser.Session
	Executor *runtime.Executor
	LastUsed time.Time
}

// New creates a new API server.
func New(port int, opts browser.Options, dbPath string) (*Server, error) {
	store, err := memory.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("open store: %w", err)
	}

	s := &Server{
		port:     port,
		opts:     opts,
		store:    store,
		log:      telemetry.New("server"),
		sessions: make(map[string]*SessionContext),
	}

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /version", s.handleVersion)

	mux.HandleFunc("POST /api/v1/sessions", s.handleCreateSession)
	mux.HandleFunc("POST /api/v1/sessions/close-all", s.handleCloseAll)
	mux.HandleFunc("DELETE /api/v1/sessions/{id}", s.handleCloseSession)

	// Page actions
	mux.HandleFunc("POST /api/v1/sessions/{id}/navigate", s.handleNavigate)
	mux.HandleFunc("GET /api/v1/sessions/{id}/url", s.handleGetURL)
	mux.HandleFunc("GET /api/v1/sessions/{id}/title", s.handleGetTitle)
	mux.HandleFunc("POST /api/v1/sessions/{id}/reload", s.handleReload)

	// Intent actions
	mux.HandleFunc("POST /api/v1/sessions/{id}/click", s.handleClick)
	mux.HandleFunc("POST /api/v1/sessions/{id}/type", s.handleType)
	mux.HandleFunc("POST /api/v1/sessions/{id}/hover", s.handleHover)
	mux.HandleFunc("POST /api/v1/sessions/{id}/scroll", s.handleScroll)
	mux.HandleFunc("POST /api/v1/sessions/{id}/upload", s.handleUpload)

	// Extraction & Data
	mux.HandleFunc("POST /api/v1/sessions/{id}/screenshot", s.handleScreenshot)
	mux.HandleFunc("GET /api/v1/sessions/{id}/html", s.handleGetHTML)
	mux.HandleFunc("GET /api/v1/sessions/{id}/elements", s.handleGetElements)
	mux.HandleFunc("POST /api/v1/sessions/{id}/js", s.handleExecuteJS)

	// Wait
	mux.HandleFunc("POST /api/v1/sessions/{id}/wait", s.handleWait)

	// Workflows
	mux.HandleFunc("POST /api/v1/workflows/run", s.handleRunWorkflow)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.loggingMiddleware(mux),
	}

	return s, nil
}

// Start runs the HTTP server in a goroutine.
func (s *Server) Start() error {
	s.log.Info("Starting API server", telemetry.F("port", s.port))
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("Server failed", telemetry.F("error", err))
		}
	}()

	// Start session reaper
	go s.reapIdleSessions()
	return nil
}

// Stop gracefully shuts down the server.
func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("Stopping API server")
	s.closeAllSessions()
	_ = s.store.Close()
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) getSession(id string) (*SessionContext, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, ok := s.sessions[id]
	if !ok {
		return nil, fmt.Errorf("session not found: %s", id)
	}
	ctx.LastUsed = time.Now()
	return ctx, nil
}

func (s *Server) closeAllSessions() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, ctx := range s.sessions {
		_ = ctx.Executor.Close()
		_ = ctx.Session.Close()
		delete(s.sessions, id)
	}
}

// reapIdleSessions closes sessions that haven't been used in 30 minutes.
func (s *Server) reapIdleSessions() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		var stale []string
		for id, ctx := range s.sessions {
			if now.Sub(ctx.LastUsed) > 30*time.Minute {
				stale = append(stale, id)
			}
		}

		for _, id := range stale {
			ctx := s.sessions[id]
			s.log.Info("Reaping idle session", telemetry.F("id", id))
			_ = ctx.Executor.Close()
			_ = ctx.Session.Close()
			delete(s.sessions, id)
		}
		s.mu.Unlock()
	}
}

// loggingMiddleware adds basic request logging.
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.log.Debug("HTTP",
			telemetry.F("method", r.Method),
			telemetry.F("path", r.URL.Path),
			telemetry.F("duration", time.Since(start).Round(time.Millisecond)),
		)
	})
}
