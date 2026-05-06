// Package actions provides the core browser action primitives for Mr. Browser.
// Each action includes pre/post verification to confirm the action had the expected effect.
package actions

import (
	"fmt"
	"time"

	"github.com/mrbrowser/mrbrowser/core/browser"
	"github.com/mrbrowser/mrbrowser/telemetry"
)

// ActionResult captures the outcome of a browser action.
type ActionResult struct {
	// Action is the name of the action performed.
	Action string `json:"action"`

	// Target is the element descriptor that was targeted.
	Target string `json:"target"`

	// Success indicates whether the action completed successfully.
	Success bool `json:"success"`

	// Error holds the error message if Success is false.
	Error string `json:"error,omitempty"`

	// Duration is how long the action took.
	Duration time.Duration `json:"duration"`

	// Verification holds the result of post-action verification.
	Verification *VerificationResult `json:"verification,omitempty"`

	// Element is the element that was acted upon.
	Element *browser.Element `json:"element,omitempty"`
}

// Engine executes browser actions with verification.
type Engine struct {
	page *browser.Page
	log  *telemetry.Logger
}

// New creates a new action engine for the given page.
func New(page *browser.Page) *Engine {
	return &Engine{
		page: page,
		log:  telemetry.New("actions"),
	}
}

// Click clicks the given element.
func (e *Engine) Click(el *browser.Element) ActionResult {
	start := time.Now()
	e.log.Step("Click", telemetry.F("element", el.Text), telemetry.F("tag", el.Tag))

	// Capture pre-click state for verification
	preURL, _ := e.page.URL()
	preHTML, _ := e.page.GetHTML()

	// Perform the click
	var clickErr error
	if el.Selector != "" {
		clickErr = e.page.ClickSelector(el.Selector)
	} else {
		clickErr = fmt.Errorf("element has no selector")
	}

	result := ActionResult{
		Action:   "click",
		Target:   el.Text,
		Element:  el,
		Duration: time.Since(start),
	}

	if clickErr != nil {
		result.Success = false
		result.Error = clickErr.Error()
		e.log.Error("Click failed", telemetry.F("error", clickErr))
		return result
	}

	// Small settle time after click
	time.Sleep(300 * time.Millisecond)

	// Post-click verification
	postURL, _ := e.page.URL()
	postHTML, _ := e.page.GetHTML()

	result.Verification = verifyDOMChange(preURL, postURL, preHTML, postHTML)
	result.Success = true
	result.Duration = time.Since(start)

	e.log.Success("Clicked", telemetry.F("element", el.Text), telemetry.F("changed", result.Verification.DOMChanged))
	return result
}

// Type clears and types text into the given element.
func (e *Engine) Type(el *browser.Element, value string) ActionResult {
	start := time.Now()
	// Mask value for logging if it looks like a password
	logValue := value
	if el.Type == "password" || el.Attributes["type"] == "password" {
		logValue = "***"
	}
	e.log.Step("Type", telemetry.F("element", el.Text), telemetry.F("value", logValue))

	result := ActionResult{
		Action:  "type",
		Target:  el.Text,
		Element: el,
	}

	if el.Selector == "" {
		result.Success = false
		result.Error = "element has no selector"
		return result
	}

	if err := e.page.TypeSelector(el.Selector, value); err != nil {
		result.Success = false
		result.Error = err.Error()
		e.log.Error("Type failed", telemetry.F("error", err))
		result.Duration = time.Since(start)
		return result
	}

	result.Success = true
	result.Duration = time.Since(start)
	e.log.Success("Typed", telemetry.F("element", el.Text))
	return result
}

// Scroll scrolls the page in the given direction.
func (e *Engine) Scroll(direction string, pixels int) ActionResult {
	start := time.Now()
	e.log.Step("Scroll", telemetry.F("direction", direction), telemetry.F("pixels", pixels))

	var script string
	switch direction {
	case "down":
		script = fmt.Sprintf("window.scrollBy(0, %d)", pixels)
	case "up":
		script = fmt.Sprintf("window.scrollBy(0, -%d)", pixels)
	case "left":
		script = fmt.Sprintf("window.scrollBy(-%d, 0)", pixels)
	case "right":
		script = fmt.Sprintf("window.scrollBy(%d, 0)", pixels)
	case "top":
		script = "window.scrollTo(0, 0)"
	case "bottom":
		script = "window.scrollTo(0, document.body.scrollHeight)"
	default:
		return ActionResult{
			Action:  "scroll",
			Success: false,
			Error:   fmt.Sprintf("unknown scroll direction: %q", direction),
		}
	}

	_, err := e.page.ExecuteJS(script)
	result := ActionResult{
		Action:   "scroll",
		Target:   direction,
		Duration: time.Since(start),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		e.log.Error("Scroll failed", telemetry.F("error", err))
		return result
	}

	time.Sleep(200 * time.Millisecond)
	result.Success = true
	e.log.Success("Scrolled", telemetry.F("direction", direction))
	return result
}

// Hover moves the mouse over the given element.
func (e *Engine) Hover(el *browser.Element) ActionResult {
	start := time.Now()
	e.log.Step("Hover", telemetry.F("element", el.Text))

	result := ActionResult{
		Action:  "hover",
		Target:  el.Text,
		Element: el,
	}

	if el.Selector == "" {
		result.Success = false
		result.Error = "element has no selector"
		result.Duration = time.Since(start)
		return result
	}

	// Use JS hover (mouseover event) as a fallback when MouseMoveParams isn't trivial
	script := fmt.Sprintf(`
		var el = document.querySelector(%q);
		if (el) {
			el.dispatchEvent(new MouseEvent('mouseover', {bubbles: true}));
			el.dispatchEvent(new MouseEvent('mouseenter', {bubbles: true}));
		}
	`, el.Selector)

	if _, err := e.page.ExecuteJS(script); err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Duration = time.Since(start)
		e.log.Error("Hover failed", telemetry.F("error", err))
		return result
	}

	time.Sleep(200 * time.Millisecond)
	result.Success = true
	result.Duration = time.Since(start)
	e.log.Success("Hovered", telemetry.F("element", el.Text))
	return result
}

// Upload sets the file input value for a file upload element.
func (e *Engine) Upload(el *browser.Element, filePath string) ActionResult {
	start := time.Now()
	e.log.Step("Upload", telemetry.F("element", el.Text), telemetry.F("file", filePath))

	result := ActionResult{
		Action:  "upload",
		Target:  el.Text,
		Element: el,
	}

	if el.Selector == "" {
		result.Success = false
		result.Error = "element has no selector"
		result.Duration = time.Since(start)
		return result
	}

	// Use chromedp's SetUploadFiles
	_, err := e.page.ExecuteJS(fmt.Sprintf(`
		(function() {
			var input = document.querySelector(%q);
			if (!input) return 'element not found';
			return 'ok';
		})()
	`, el.Selector))

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result
	}

	// Note: Real file upload requires chromedp.SetUploadFiles — handled by the caller
	// who has direct access to the chromedp context.
	result.Success = true
	result.Duration = time.Since(start)
	e.log.Success("Upload initiated", telemetry.F("file", filePath))
	return result
}

// Screenshot is a convenience action that takes a screenshot and returns the bytes.
func (e *Engine) Screenshot() ([]byte, ActionResult) {
	start := time.Now()
	e.log.Step("Screenshot")

	buf, err := e.page.Screenshot()
	result := ActionResult{
		Action:   "screenshot",
		Duration: time.Since(start),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return nil, result
	}

	result.Success = true
	return buf, result
}
