// Package memory provides the SQLite-backed storage layer for Mr. Browser.
// It stores element fingerprints, workflow results, and action history
// to enable self-healing automation.
package memory

import (
	"time"
)

// ElementFingerprint stores a durable record of a DOM element's identity.
// Used by the self-healing system to recover when selectors break.
type ElementFingerprint struct {
	// ID is the unique identifier for this fingerprint (UUID).
	ID string `json:"id" db:"id"`

	// WorkflowID is the name of the workflow that created this fingerprint.
	WorkflowID string `json:"workflow_id" db:"workflow_id"`

	// StepName is the step within the workflow.
	StepName string `json:"step_name" db:"step_name"`

	// Target is the original user-provided intent string.
	Target string `json:"target" db:"target"`

	// Selector is the last-known working CSS selector.
	Selector string `json:"selector" db:"selector"`

	// Text is the visible text of the element.
	Text string `json:"text" db:"text"`

	// Role is the ARIA role or semantic role.
	Role string `json:"role" db:"role"`

	// Type is the semantic element type.
	Type string `json:"type" db:"type"`

	// Section is the page section (header, nav, main, footer).
	Section string `json:"section" db:"section"`

	// NearbyText is a JSON array of nearby text snippets.
	NearbyText []string `json:"nearby_text" db:"-"`

	// Attributes is a JSON object of key element attributes.
	Attributes map[string]string `json:"attributes" db:"-"`

	// LastSeen is when this fingerprint was last successfully used.
	LastSeen time.Time `json:"last_seen" db:"last_seen"`

	// UseCount is how many times this fingerprint has been used.
	UseCount int `json:"use_count" db:"use_count"`

	// SuccessCount is the number of successful uses.
	SuccessCount int `json:"success_count" db:"success_count"`

	// SuccessRate is SuccessCount / UseCount.
	SuccessRate float64 `json:"success_rate" db:"success_rate"`
}

// WorkflowResult records the outcome of a complete workflow run.
type WorkflowResult struct {
	ID           string        `json:"id" db:"id"`
	WorkflowName string        `json:"workflow_name" db:"workflow_name"`
	StartedAt    time.Time     `json:"started_at" db:"started_at"`
	CompletedAt  time.Time     `json:"completed_at" db:"completed_at"`
	Duration     time.Duration `json:"duration" db:"duration_ms"`
	Success      bool          `json:"success" db:"success"`
	StepsTotal   int           `json:"steps_total" db:"steps_total"`
	StepsPassed  int           `json:"steps_passed" db:"steps_passed"`
	ErrorMessage string        `json:"error_message,omitempty" db:"error_message"`
	URL          string        `json:"url" db:"url"`
}

// ActionRecord stores a single action's outcome.
type ActionRecord struct {
	ID           string        `json:"id" db:"id"`
	WorkflowID   string        `json:"workflow_id" db:"workflow_id"`
	StepName     string        `json:"step_name" db:"step_name"`
	Action       string        `json:"action" db:"action"`
	Target       string        `json:"target" db:"target"`
	Selector     string        `json:"selector" db:"selector"`
	Success      bool          `json:"success" db:"success"`
	Duration     time.Duration `json:"duration" db:"duration_ms"`
	ErrorMessage string        `json:"error_message,omitempty" db:"error_message"`
	Recovered    bool          `json:"recovered" db:"recovered"`
	RecoveryNote string        `json:"recovery_note,omitempty" db:"recovery_note"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`
}

// RecoveryEvent records when the self-healer successfully recovered from a broken selector.
type RecoveryEvent struct {
	ID          string    `json:"id" db:"id"`
	WorkflowID  string    `json:"workflow_id" db:"workflow_id"`
	Target      string    `json:"target" db:"target"`
	OldSelector string    `json:"old_selector" db:"old_selector"`
	NewSelector string    `json:"new_selector" db:"new_selector"`
	OldText     string    `json:"old_text" db:"old_text"`
	NewText     string    `json:"new_text" db:"new_text"`
	Confidence  float64   `json:"confidence" db:"confidence"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
