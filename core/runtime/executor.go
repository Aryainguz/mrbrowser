package runtime

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mrbrowser/mrbrowser/core/actions"
	"github.com/mrbrowser/mrbrowser/core/browser"
	"github.com/mrbrowser/mrbrowser/core/report"
	"github.com/mrbrowser/mrbrowser/intelligence/dom"
	"github.com/mrbrowser/mrbrowser/intelligence/resolver"
	"github.com/mrbrowser/mrbrowser/memory"
	"github.com/mrbrowser/mrbrowser/telemetry"
)

// StepResult captures the outcome of a single task step.
type StepResult struct {
	Step      int           `json:"step"`
	Kind      string        `json:"kind"`
	Success   bool          `json:"success"`
	Error     string        `json:"error,omitempty"`
	Duration  time.Duration `json:"duration"`
	Recovered bool          `json:"recovered,omitempty"`
}

// ExecutionResult captures the outcome of a complete task run.
type ExecutionResult struct {
	TaskName    string        `json:"task_name"`
	Success     bool          `json:"success"`
	StepResults []StepResult  `json:"step_results"`
	Duration    time.Duration `json:"duration"`
	Error       string        `json:"error,omitempty"`
}

// Executor runs a Task against a live browser session.
type Executor struct {
	session     *browser.Session
	store       *memory.Store
	healer      *memory.SelfHealer
	resolver    *resolver.Resolver
	engine      *actions.Engine
	extractor   *dom.Extractor
	log         *telemetry.Logger
	stopOnError bool
	tracker     *report.Tracker
}

// ExecutorOptions configures the Executor.
type ExecutorOptions struct {
	StopOnError bool
	DBPath      string
}

// DefaultExecutorOptions returns sensible defaults.
func DefaultExecutorOptions() ExecutorOptions {
	return ExecutorOptions{
		StopOnError: true,
		DBPath:      "./mrbrowser.db",
	}
}

// NewExecutor creates an Executor for the given session.
func NewExecutor(session *browser.Session, opts ExecutorOptions) (*Executor, error) {
	store, err := memory.Open(opts.DBPath)
	if err != nil {
		return nil, fmt.Errorf("open memory store: %w", err)
	}

	return &Executor{
		session:     session,
		store:       store,
		healer:      memory.NewSelfHealer(store),
		resolver:    resolver.New(),
		engine:      actions.New(session.Page),
		extractor:   dom.NewExtractor(session.Page),
		log:         telemetry.New("executor"),
		stopOnError: opts.StopOnError,
		tracker:     report.NewTracker(),
	}, nil
}

// Run executes all steps of the given task.
func (e *Executor) Run(task *Task) (*ExecutionResult, error) {
	started := time.Now()
	e.log.Info("Starting task", telemetry.F("name", task.Name), telemetry.F("steps", len(task.Steps)))

	e.session.Page.EnableNetworkTracking(e.tracker)

	result := &ExecutionResult{
		TaskName: task.Name,
	}

	for i, step := range task.Steps {
		stepNum := i + 1
		e.log.Step(fmt.Sprintf("Step %d/%d [%s]", stepNum, len(task.Steps), step.Kind()))

		e.tracker.StartStep(stepNum, step.Kind(), getStepTarget(&step))

		stepResult, err := e.runStep(task.Name, stepNum, &step)
		result.StepResults = append(result.StepResults, stepResult)
		
		e.tracker.EndStep(err)

		if err != nil {
			e.log.Error("Step failed",
				telemetry.F("step", stepNum),
				telemetry.F("kind", step.Kind()),
				telemetry.F("error", err),
			)
			if e.stopOnError {
				result.Success = false
				result.Error = fmt.Sprintf("step %d [%s]: %v", stepNum, step.Kind(), err)
				result.Duration = time.Since(started)
				
				if task.Config.GenerateReport != "" {
					_ = e.tracker.GenerateMarkdown(task.Name, task.Config.GenerateReport)
				}
				
				return result, nil
			}
		} else {
			e.log.Success(fmt.Sprintf("Step %d complete", stepNum))
		}
	}

	result.Success = true
	result.Duration = time.Since(started)
	e.log.Success("Task complete",
		telemetry.F("name", task.Name),
		telemetry.F("duration", result.Duration.Round(time.Millisecond)),
	)
	
	if task.Config.GenerateReport != "" {
		_ = e.tracker.GenerateMarkdown(task.Name, task.Config.GenerateReport)
	}
	
	return result, nil
}

func getStepTarget(s *Step) string {
	if s.Click != nil { return s.Click.Target }
	if s.Type != nil { return s.Type.Target }
	if s.Open != nil { return s.Open.URL }
	return ""
}

func (e *Executor) runStep(taskName string, stepNum int, step *Step) (StepResult, error) {
	start := time.Now()
	sr := StepResult{Step: stepNum, Kind: step.Kind()}

	var err error
	switch step.Kind() {
	case "open":
		err = e.stepOpen(step.Open)
	case "click":
		sr.Recovered, err = e.stepClick(taskName, stepNum, step.Click)
	case "type":
		sr.Recovered, err = e.stepType(taskName, stepNum, step.Type)
	case "scroll":
		err = e.stepScroll(step.Scroll)
	case "hover":
		_, err = e.stepHover(taskName, stepNum, step.Hover)
	case "screenshot":
		err = e.stepScreenshot(step.Screenshot)
	case "wait":
		err = e.stepWait(step.Wait)
	case "assert":
		err = e.stepAssert(step.Assert)
	case "reload":
		err = e.session.Page.Reload()
	default:
		err = fmt.Errorf("unknown step kind: %s", step.Kind())
	}

	sr.Duration = time.Since(start)
	if err != nil {
		sr.Success = false
		sr.Error = err.Error()
	} else {
		sr.Success = true
	}
	return sr, err
}

func (e *Executor) highlightAndScreenshot(el *browser.Element, taskName string, stepNum int) {
	if el == nil || el.Selector == "" {
		return
	}

	escapedSelector := strings.ReplaceAll(el.Selector, `"`, `\"`)
	script := fmt.Sprintf(`
		(function() {
			var el = document.querySelector("%s");
			if (!el) return "";
			var oldOutline = el.style.outline;
			var oldOffset = el.style.outlineOffset;
			el.style.outline = "4px solid red";
			el.style.outlineOffset = "2px";
			return oldOutline + "|||" + oldOffset;
		})()
	`, escapedSelector)

	result, err := e.session.Page.ExecuteJS(script)
	if err != nil {
		return
	}

	// Capture the screenshot with the highlight
	buf, _ := e.session.Page.Screenshot()
	if len(buf) > 0 {
		output := fmt.Sprintf("highlight_%s_step%d.png", taskName, stepNum)
		if os.WriteFile(output, buf, 0644) == nil {
			e.tracker.RecordScreenshot(output)
		}
	}

	// Restore the old outline
	if resStr, ok := result.(string); ok {
		parts := strings.SplitN(resStr, "|||", 2)
		if len(parts) == 2 {
			oldOutline := strings.ReplaceAll(parts[0], `"`, `\"`)
			oldOffset := strings.ReplaceAll(parts[1], `"`, `\"`)
			
			restoreScript := fmt.Sprintf(`
				(function() {
					var el = document.querySelector("%s");
					if (el) {
						el.style.outline = "%s";
						el.style.outlineOffset = "%s";
					}
				})()
			`, escapedSelector, oldOutline, oldOffset)
			e.session.Page.ExecuteJS(restoreScript)
		}
	}
}

// stepOpen navigates to a URL.
func (e *Executor) stepOpen(s *OpenStep) error {
	return e.session.Page.Navigate(s.URL)
}

// stepClick finds and clicks an element, with self-healing fallback.
func (e *Executor) stepClick(taskName string, stepNum int, s *ClickStep) (recovered bool, err error) {
	el, recovered, err := e.resolveElement(taskName, s.Target, s.Selector)
	if err != nil {
		return false, err
	}
	
	e.highlightAndScreenshot(el, taskName, stepNum)

	result := e.engine.Click(el)
	if !result.Success {
		return recovered, fmt.Errorf("click: %s", result.Error)
	}

	// Store fingerprint after success
	_ = e.healer.StoreFingerprint(taskName, fmt.Sprintf("click-%s", s.Target), s.Target, el)
	return recovered, nil
}

// stepType finds an element and types text, with self-healing fallback.
func (e *Executor) stepType(taskName string, stepNum int, s *TypeStep) (recovered bool, err error) {
	el, recovered, err := e.resolveElement(taskName, s.Target, s.Selector)
	if err != nil {
		return false, err
	}

	e.highlightAndScreenshot(el, taskName, stepNum)

	val := s.Value
	if s.Enter {
		val += "\r"
	}

	result := e.engine.Type(el, val)
	if !result.Success {
		return recovered, fmt.Errorf("type: %s", result.Error)
	}
	_ = e.healer.StoreFingerprint(taskName, fmt.Sprintf("type-%s", s.Target), s.Target, el)
	return recovered, nil
}

// stepScroll scrolls the page.
func (e *Executor) stepScroll(s *ScrollStep) error {
	px := s.Pixels
	if px == 0 {
		px = 300
	}
	result := e.engine.Scroll(s.Direction, px)
	if !result.Success {
		return fmt.Errorf("scroll: %s", result.Error)
	}
	return nil
}

// stepHover hovers over an element.
func (e *Executor) stepHover(taskName string, stepNum int, s *HoverStep) (recovered bool, err error) {
	el, recovered, err := e.resolveElement(taskName, s.Target, s.Selector)
	if err != nil {
		return false, err
	}

	e.highlightAndScreenshot(el, taskName, stepNum)

	result := e.engine.Hover(el)
	if !result.Success {
		return recovered, fmt.Errorf("hover: %s", result.Error)
	}
	return recovered, nil
}

// stepScreenshot captures a screenshot.
func (e *Executor) stepScreenshot(s *ScreenshotStep) error {
	buf, result := e.engine.Screenshot()
	if !result.Success {
		return fmt.Errorf("screenshot: %s", result.Error)
	}
	output := s.Output
	if output == "" {
		output = fmt.Sprintf("screenshot_%d.png", time.Now().Unix())
	}
	if err := os.WriteFile(output, buf, 0644); err != nil {
		return fmt.Errorf("save screenshot to %q: %w", output, err)
	}
	e.tracker.RecordScreenshot(output)
	e.log.Success("Screenshot saved", telemetry.F("path", output))
	return nil
}

// stepWait waits for a duration or condition.
func (e *Executor) stepWait(s *WaitStep) error {
	if s.Selector != "" {
		return e.session.Page.WaitForSelector(s.Selector)
	}
	if s.URL != "" {
		return e.session.Page.WaitForURL(s.URL)
	}
	if s.Seconds > 0 {
		time.Sleep(time.Duration(s.Seconds * float64(time.Second)))
	}
	return nil
}

// stepAssert checks a condition about the current page state.
func (e *Executor) stepAssert(s *AssertStep) error {
	if s.URLContains != "" {
		url, err := e.session.Page.URL()
		if err != nil {
			return fmt.Errorf("assert url: %w", err)
		}
		if !strings.Contains(url, s.URLContains) {
			return fmt.Errorf("assert url_contains %q: got %q", s.URLContains, url)
		}
	}
	if s.TextVisible != "" {
		html, err := e.session.Page.GetHTML()
		if err != nil {
			return fmt.Errorf("assert text: %w", err)
		}
		if !strings.Contains(html, s.TextVisible) {
			return fmt.Errorf("assert text_visible %q: not found on page", s.TextVisible)
		}
	}
	if s.ElementExists != "" {
		elements, err := e.extractor.Extract()
		if err != nil {
			return fmt.Errorf("assert element_exists: %w", err)
		}
		_, _, err = e.resolver.Resolve(s.ElementExists, elements)
		if err != nil {
			return fmt.Errorf("assert element_exists %q: %w", s.ElementExists, err)
		}
	}
	return nil
}

// resolveElement finds an element by intent, using explicit selector if provided,
// then falling back to the resolver, then to self-healing.
func (e *Executor) resolveElement(taskName, target, explicitSelector string) (*browser.Element, bool, error) {
	// If an explicit selector is given, use it directly
	if explicitSelector != "" {
		return &browser.Element{Selector: explicitSelector, Text: target}, false, nil
	}

	// Extract current page elements
	elements, err := e.extractor.Extract()
	if err != nil {
		return nil, false, fmt.Errorf("extract elements: %w", err)
	}

	// Try resolver
	el, _, err := e.resolver.Resolve(target, elements)
	if err == nil {
		return el, false, nil
	}

	// Resolver failed — try self-healing
	e.log.Warn("Resolver failed, attempting self-healing",
		telemetry.F("target", target),
		telemetry.F("error", err),
	)

	recovery, healErr := e.healer.TryRecover(taskName, target, elements)
	if healErr != nil {
		return nil, false, fmt.Errorf("resolve %q failed and healing errored: %v / %v", target, err, healErr)
	}
	if !recovery.Recovered {
		return nil, false, fmt.Errorf("resolve %q failed: %v; healing: %s", target, err, recovery.Note)
	}

	return recovery.Element, true, nil
}

// Close releases executor resources.
func (e *Executor) Close() error {
	return e.store.Close()
}
