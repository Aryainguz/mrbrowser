package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mrbrowser/mrbrowser/core/browser"
	"github.com/mrbrowser/mrbrowser/core/runtime"
	"github.com/mrbrowser/mrbrowser/intelligence/dom"
)

// APIError is a standard JSON error response.
type APIError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Target  string `json:"target,omitempty"`
		Action  string `json:"action,omitempty"`
	} `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, msg, target, action string) {
	errResp := APIError{}
	errResp.Error.Code = code
	errResp.Error.Message = msg
	errResp.Error.Target = target
	errResp.Error.Action = action
	writeJSON(w, status, errResp)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"version": "0.1.0"})
}

func (s *Server) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	id := fmt.Sprintf("sess_%d", time.Now().UnixNano())

	session, err := browser.NewSession(id, s.opts)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "SESSION_FAILED", err.Error(), "", "")
		return
	}

	executor, err := runtime.NewExecutor(session, runtime.ExecutorOptions{
		StopOnError: true,
		DBPath:      s.store.DBPath(), // Assuming we can get this, or pass it
	})
	if err != nil {
		_ = session.Close()
		writeError(w, http.StatusInternalServerError, "EXECUTOR_FAILED", err.Error(), "", "")
		return
	}

	if req.URL != "" {
		_ = session.Page.Navigate(req.URL)
	}

	s.mu.Lock()
	s.sessions[id] = &SessionContext{
		ID:       id,
		Session:  session,
		Executor: executor,
		LastUsed: time.Now(),
	}
	s.mu.Unlock()

	writeJSON(w, http.StatusOK, map[string]string{"session_id": id})
}

func (s *Server) handleCloseSession(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	s.mu.Lock()
	ctx, ok := s.sessions[id]
	if ok {
		_ = ctx.Executor.Close()
		_ = ctx.Session.Close()
		delete(s.sessions, id)
	}
	s.mu.Unlock()

	writeJSON(w, http.StatusOK, map[string]bool{"closed": ok})
}

func (s *Server) handleCloseAll(w http.ResponseWriter, r *http.Request) {
	s.closeAllSessions()
	writeJSON(w, http.StatusOK, map[string]bool{"closed": true})
}

func (s *Server) handleNavigate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx, err := s.getSession(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error(), "", "")
		return
	}

	var req struct {
		URL string `json:"url"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	if err := ctx.Session.Page.Navigate(req.URL); err != nil {
		writeError(w, http.StatusInternalServerError, "NAVIGATE_FAILED", err.Error(), "", "")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (s *Server) handleGetURL(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx, err := s.getSession(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error(), "", "")
		return
	}
	url, _ := ctx.Session.Page.URL()
	writeJSON(w, http.StatusOK, map[string]string{"url": url})
}

func (s *Server) handleGetTitle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx, err := s.getSession(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error(), "", "")
		return
	}
	title, _ := ctx.Session.Page.Title()
	writeJSON(w, http.StatusOK, map[string]string{"title": title})
}

func (s *Server) handleReload(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx, err := s.getSession(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error(), "", "")
		return
	}
	_ = ctx.Session.Page.Reload()
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// ── Action Handlers ──

type actionReq struct {
	Target    string `json:"target"`
	Selector  string `json:"selector"`
	Value     string `json:"value"`     // for type
	FilePath  string `json:"file_path"` // for upload
	Direction string `json:"direction"` // for scroll
	Pixels    int    `json:"pixels"`    // for scroll
}

func (s *Server) handleClick(w http.ResponseWriter, r *http.Request) {
	s.handleAction(w, r, "click")
}

func (s *Server) handleType(w http.ResponseWriter, r *http.Request) {
	s.handleAction(w, r, "type")
}

func (s *Server) handleHover(w http.ResponseWriter, r *http.Request) {
	s.handleAction(w, r, "hover")
}

func (s *Server) handleScroll(w http.ResponseWriter, r *http.Request) {
	s.handleAction(w, r, "scroll")
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	s.handleAction(w, r, "upload")
}

func (s *Server) handleAction(w http.ResponseWriter, r *http.Request, action string) {
	id := r.PathValue("id")
	ctx, err := s.getSession(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error(), "", "")
		return
	}

	var req actionReq
	_ = json.NewDecoder(r.Body).Decode(&req)

	// Build a one-step task to reuse the Executor's resolver + healer logic
	task := &runtime.Task{
		Name:  "api-action",
		Steps: []runtime.Step{{}},
	}

	switch action {
	case "click":
		task.Steps[0].Click = &runtime.ClickStep{Target: req.Target, Selector: req.Selector}
	case "type":
		task.Steps[0].Type = &runtime.TypeStep{Target: req.Target, Selector: req.Selector, Value: req.Value}
	case "hover":
		task.Steps[0].Hover = &runtime.HoverStep{Target: req.Target, Selector: req.Selector}
	case "scroll":
		task.Steps[0].Scroll = &runtime.ScrollStep{Direction: req.Direction, Pixels: req.Pixels}
	case "upload":
		task.Steps[0].Upload = &runtime.UploadStep{Target: req.Target, Selector: req.Selector, FilePath: req.FilePath}
	}

	res, err := ctx.Executor.Run(task)
	if err != nil || !res.Success {
		errMsg := res.Error
		if err != nil {
			errMsg = err.Error()
		}

		// Map errors to specific codes for the SDK
		code := "ACTION_FAILED"
		if errMsg != "" && (contains(errMsg, "no element matched") || contains(errMsg, "resolver failed")) {
			code = "ELEMENT_NOT_FOUND"
		}

		writeError(w, http.StatusBadRequest, code, errMsg, req.Target, action)
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (s *Server) handleScreenshot(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx, err := s.getSession(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error(), "", "")
		return
	}

	buf, err := ctx.Session.Page.Screenshot()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "SCREENSHOT_FAILED", err.Error(), "", "")
		return
	}

	// Send raw bytes inside JSON for simplicity with Python SDK (base64 could be better for large imgs)
	writeJSON(w, http.StatusOK, map[string][]byte{"data": buf})
}

func (s *Server) handleGetHTML(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx, err := s.getSession(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error(), "", "")
		return
	}
	html, _ := ctx.Session.Page.GetHTML()
	writeJSON(w, http.StatusOK, map[string]string{"html": html})
}

func (s *Server) handleExecuteJS(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx, err := s.getSession(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error(), "", "")
		return
	}

	var req struct {
		Script string `json:"script"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	res, err := ctx.Session.Page.ExecuteJS(req.Script)
	if err != nil {
		writeError(w, http.StatusBadRequest, "JS_ERROR", err.Error(), "", "")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"result": res})
}

func (s *Server) handleGetElements(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx, err := s.getSession(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error(), "", "")
		return
	}

	visibleOnly := r.URL.Query().Get("visible_only") == "true"

	ext := dom.NewExtractor(ctx.Session.Page)
	elements, err := ext.Extract()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "EXTRACT_FAILED", err.Error(), "", "")
		return
	}

	if visibleOnly {
		var filtered []*dom.PageElement
		for _, el := range elements {
			if el.Visible {
				filtered = append(filtered, el)
			}
		}
		elements = filtered
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"elements": elements})
}

func (s *Server) handleWait(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx, err := s.getSession(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error(), "", "")
		return
	}

	var req struct {
		Selector string  `json:"selector"`
		URL      string  `json:"url"`
		Timeout  float64 `json:"timeout"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	var waitErr error
	if req.Selector != "" {
		waitErr = ctx.Session.Page.WaitForSelector(req.Selector)
	} else if req.URL != "" {
		waitErr = ctx.Session.Page.WaitForURL(req.URL)
	} else {
		time.Sleep(time.Duration(req.Timeout * float64(time.Second)))
	}

	if waitErr != nil {
		writeError(w, http.StatusRequestTimeout, "TIMEOUT", waitErr.Error(), "", "")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (s *Server) handleRunWorkflow(w http.ResponseWriter, r *http.Request) {
	var req struct {
		YAML string `json:"yaml"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	task, err := runtime.Load([]byte(req.YAML))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_YAML", err.Error(), "", "")
		return
	}

	session, err := browser.NewSession(task.Name, s.opts)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "SESSION_FAILED", err.Error(), "", "")
		return
	}
	defer session.Close()

	executor, err := runtime.NewExecutor(session, runtime.ExecutorOptions{
		StopOnError: true,
		DBPath:      s.store.DBPath(),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "EXECUTOR_FAILED", err.Error(), "", "")
		return
	}
	defer executor.Close()

	res, err := executor.Run(task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "RUN_FAILED", err.Error(), "", "")
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
