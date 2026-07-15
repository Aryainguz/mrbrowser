package report

import (
	"time"
)

// APIRequest represents a network request made during the workflow.
type APIRequest struct {
	Method     string
	URL        string
	StatusCode int
	Duration   time.Duration
}

// StepRecord records the execution details of a single step.
type StepRecord struct {
	Index        int
	Action       string
	Target       string
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
	Screenshot   string
	Error        error
	APIRequests  []APIRequest
}

// Tracker stores the execution state for reporting.
type Tracker struct {
	Steps       []*StepRecord
	CurrentStep *StepRecord
}

// NewTracker initializes a new report tracker.
func NewTracker() *Tracker {
	return &Tracker{
		Steps: make([]*StepRecord, 0),
	}
}

// StartStep marks the beginning of a step.
func (t *Tracker) StartStep(index int, action, target string) {
	record := &StepRecord{
		Index:     index,
		Action:    action,
		Target:    target,
		StartTime: time.Now(),
	}
	t.Steps = append(t.Steps, record)
	t.CurrentStep = record
}

// EndStep marks the end of a step.
func (t *Tracker) EndStep(err error) {
	if t.CurrentStep == nil {
		return
	}
	t.CurrentStep.EndTime = time.Now()
	t.CurrentStep.Duration = t.CurrentStep.EndTime.Sub(t.CurrentStep.StartTime)
	t.CurrentStep.Error = err
}

// RecordScreenshot attaches a screenshot to the current step.
func (t *Tracker) RecordScreenshot(path string) {
	if t.CurrentStep != nil {
		t.CurrentStep.Screenshot = path
	}
}

// RecordAPIRequest logs a network request against the current step.
func (t *Tracker) RecordAPIRequest(req APIRequest) {
	if t.CurrentStep != nil {
		t.CurrentStep.APIRequests = append(t.CurrentStep.APIRequests, req)
	}
}

// GetErrorCount returns the number of failed steps.
func (t *Tracker) GetErrorCount() int {
	count := 0
	for _, s := range t.Steps {
		if s.Error != nil {
			count++
		}
	}
	return count
}

// GetSlowAPIs returns all APIs taking longer than the threshold.
func (t *Tracker) GetSlowAPIs(threshold time.Duration) []APIRequest {
	var slow []APIRequest
	for _, s := range t.Steps {
		for _, api := range s.APIRequests {
			if api.Duration >= threshold {
				slow = append(slow, api)
			}
		}
	}
	return slow
}

// GetFailedAPIs returns all APIs with 4xx or 5xx status codes.
func (t *Tracker) GetFailedAPIs() []APIRequest {
	var failed []APIRequest
	for _, s := range t.Steps {
		for _, api := range s.APIRequests {
			if api.StatusCode >= 400 {
				failed = append(failed, api)
			}
		}
	}
	return failed
}
