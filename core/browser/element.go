package browser

// Element represents a resolved DOM element with positional and semantic metadata.
// It is the core data structure passed between the intelligence engine and the action engine.
type Element struct {
	// NodeID is the CDP node ID.
	NodeID int `json:"node_id"`

	// Tag is the HTML tag name (e.g., "button", "input", "a").
	Tag string `json:"tag"`

	// Type is the element's semantic type (e.g., "button", "text-input", "link", "select").
	Type string `json:"type"`

	// Text is the visible text content of the element.
	Text string `json:"text"`

	// Role is the ARIA role or inferred role.
	Role string `json:"role"`

	// Visible indicates whether the element is visible in the viewport.
	Visible bool `json:"visible"`

	// Clickable indicates whether the element is interactable.
	Clickable bool `json:"clickable"`

	// Position is the bounding box in viewport coordinates.
	Position Rect `json:"position"`

	// Attributes are the HTML attributes of the element.
	Attributes map[string]string `json:"attributes"`

	// Placeholder is the input placeholder text, if any.
	Placeholder string `json:"placeholder,omitempty"`

	// Value is the current input value, if any.
	Value string `json:"value,omitempty"`

	// Selector is the CSS selector that uniquely identifies this element.
	// Used as fallback when we need to target the element via CDP.
	Selector string `json:"selector"`

	// XPath is an XPath expression that uniquely identifies this element.
	XPath string `json:"xpath,omitempty"`

	// NearbyText contains text from nearby elements (siblings, parent).
	// Used for context-aware self-healing.
	NearbyText []string `json:"nearby_text,omitempty"`

	// Section identifies the page section (header, nav, main, footer, etc.).
	Section string `json:"section,omitempty"`

	// Confidence is the resolver's confidence score (0.0–1.0).
	// Set by the ElementResolver when returning a match.
	Confidence float64 `json:"confidence,omitempty"`
}

// Rect represents a bounding box in viewport coordinates (pixels).
type Rect struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// Center returns the center point of the rectangle.
func (r Rect) Center() (float64, float64) {
	return r.X + r.Width/2, r.Y + r.Height/2
}

// IsZero returns true if the rect has no area (element not positioned).
func (r Rect) IsZero() bool {
	return r.Width == 0 && r.Height == 0
}

// Cookie represents an HTTP cookie.
type Cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Domain   string `json:"domain"`
	Path     string `json:"path"`
	Secure   bool   `json:"secure"`
	HTTPOnly bool   `json:"http_only"`
	SameSite string `json:"same_site"`
}
