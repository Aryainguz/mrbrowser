// Package runtime provides the YAML task definition and execution engine.
package runtime

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Task is a complete automation workflow loaded from a YAML file.
type Task struct {
	// Name uniquely identifies this workflow.
	Name string `yaml:"name"`

	// Description is an optional human-readable description.
	Description string `yaml:"description,omitempty"`

	// Steps is the ordered list of actions to execute.
	Steps []Step `yaml:"steps"`

	// Config holds optional per-task configuration.
	Config TaskConfig `yaml:"config,omitempty"`
}

// TaskConfig holds optional per-task overrides.
type TaskConfig struct {
	// Headless overrides the global headless setting.
	Headless *bool `yaml:"headless,omitempty"`

	// Timeout overrides the default step timeout (in seconds).
	TimeoutSeconds int `yaml:"timeout_seconds,omitempty"`

	// StopOnError controls whether to stop on first error (default: true).
	StopOnError *bool `yaml:"stop_on_error,omitempty"`
}

// Step represents a single step in a workflow.
// Exactly one field should be non-nil.
type Step struct {
	Open       *OpenStep       `yaml:"open,omitempty"`
	Click      *ClickStep      `yaml:"click,omitempty"`
	Type       *TypeStep       `yaml:"type,omitempty"`
	Scroll     *ScrollStep     `yaml:"scroll,omitempty"`
	Hover      *HoverStep      `yaml:"hover,omitempty"`
	Screenshot *ScreenshotStep `yaml:"screenshot,omitempty"`
	Wait       *WaitStep       `yaml:"wait,omitempty"`
	Upload     *UploadStep     `yaml:"upload,omitempty"`
	Extract    *ExtractStep    `yaml:"extract,omitempty"`
	Assert     *AssertStep     `yaml:"assert,omitempty"`
	Reload     *ReloadStep     `yaml:"reload,omitempty"`
}

// Kind returns the step type name.
func (s *Step) Kind() string {
	switch {
	case s.Open != nil:
		return "open"
	case s.Click != nil:
		return "click"
	case s.Type != nil:
		return "type"
	case s.Scroll != nil:
		return "scroll"
	case s.Hover != nil:
		return "hover"
	case s.Screenshot != nil:
		return "screenshot"
	case s.Wait != nil:
		return "wait"
	case s.Upload != nil:
		return "upload"
	case s.Extract != nil:
		return "extract"
	case s.Assert != nil:
		return "assert"
	case s.Reload != nil:
		return "reload"
	}
	return "unknown"
}

// OpenStep navigates to a URL.
type OpenStep struct {
	URL string `yaml:"url"`
}

// ClickStep clicks an element identified by intent.
type ClickStep struct {
	// Target is the natural-language element description.
	Target string `yaml:"target"`
	// Selector is an optional explicit CSS selector (overrides resolver).
	Selector string `yaml:"selector,omitempty"`
}

// TypeStep types text into an element.
type TypeStep struct {
	Target   string `yaml:"target"`
	Value    string `yaml:"value"`
	// Clear clears the field before typing (default: true).
	Clear    *bool  `yaml:"clear,omitempty"`
	Selector string `yaml:"selector,omitempty"`
	Enter    bool   `yaml:"enter,omitempty"`
}

// ScrollStep scrolls the page.
type ScrollStep struct {
	// Direction: up, down, left, right, top, bottom.
	Direction string `yaml:"direction"`
	// Pixels to scroll (ignored for top/bottom).
	Pixels int `yaml:"pixels,omitempty"`
}

// HoverStep hovers over an element.
type HoverStep struct {
	Target   string `yaml:"target"`
	Selector string `yaml:"selector,omitempty"`
}

// ScreenshotStep captures a screenshot.
type ScreenshotStep struct {
	// Output is the file path to save the screenshot PNG.
	Output string `yaml:"output,omitempty"`
}

// WaitStep pauses execution.
type WaitStep struct {
	// Seconds to wait.
	Seconds float64 `yaml:"seconds,omitempty"`
	// Selector to wait for (visibility).
	Selector string `yaml:"selector,omitempty"`
	// URL to wait for (substring match).
	URL string `yaml:"url,omitempty"`
}

// UploadStep uploads a file.
type UploadStep struct {
	Target   string `yaml:"target"`
	FilePath string `yaml:"file"`
	Selector string `yaml:"selector,omitempty"`
}

// ExtractStep extracts text from an element and saves it to the session.
type ExtractStep struct {
	Target string `yaml:"target"`
	SaveAs string `yaml:"save_as"`
}

// AssertStep asserts a condition about the current page.
type AssertStep struct {
	// URLContains asserts the current URL contains this string.
	URLContains string `yaml:"url_contains,omitempty"`
	// TextVisible asserts the given text is visible on the page.
	TextVisible string `yaml:"text_visible,omitempty"`
	// ElementExists asserts an element matching the target exists.
	ElementExists string `yaml:"element_exists,omitempty"`
}

// ReloadStep reloads the current page.
type ReloadStep struct{}

// LoadFile parses a YAML task file.
func LoadFile(path string) (*Task, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read task file %q: %w", path, err)
	}
	return Load(data)
}

// Load parses YAML task bytes.
func Load(data []byte) (*Task, error) {
	// Support both bare YAML and `task:` wrapper
	var wrapper struct {
		Task *Task `yaml:"task"`
	}
	if err := yaml.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("parse task yaml: %w", err)
	}

	if wrapper.Task != nil {
		return validate(wrapper.Task)
	}

	// Try direct parse
	var task Task
	if err := yaml.Unmarshal(data, &task); err != nil {
		return nil, fmt.Errorf("parse task yaml: %w", err)
	}
	return validate(&task)
}

func validate(t *Task) (*Task, error) {
	if t.Name == "" {
		t.Name = "unnamed"
	}
	if len(t.Steps) == 0 {
		return nil, fmt.Errorf("task %q has no steps", t.Name)
	}
	for i, s := range t.Steps {
		if s.Kind() == "unknown" {
			return nil, fmt.Errorf("task %q step %d: no recognized action", t.Name, i+1)
		}
	}
	return t, nil
}
