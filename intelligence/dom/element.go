// Package dom provides DOM extraction and element metadata for the Mr. Browser intelligence engine.
package dom

import (
	"github.com/mrbrowser/mrbrowser/core/browser"
)

// PageElement is the internal representation of a DOM node with full metadata.
// It mirrors browser.Element but includes tree-structure information.
type PageElement struct {
	// NodeID is the CDP node ID.
	NodeID int `json:"node_id"`

	// Tag is the HTML tag name.
	Tag string `json:"tag"`

	// Type is the semantic type (button, text-input, link, etc.).
	Type string `json:"type"`

	// Text is the visible text content (trimmed).
	Text string `json:"text"`

	// Role is the ARIA role or inferred semantic role.
	Role string `json:"role"`

	// Visible indicates the element is visible in the viewport.
	Visible bool `json:"visible"`

	// Clickable indicates the element can receive click events.
	Clickable bool `json:"clickable"`

	// Position is the bounding box.
	Position browser.Rect `json:"position"`

	// Attributes are all HTML attributes.
	Attributes map[string]string `json:"attributes"`

	// Placeholder is the input placeholder text.
	Placeholder string `json:"placeholder,omitempty"`

	// Value is the current form input value.
	Value string `json:"value,omitempty"`

	// Selector is a generated unique CSS selector for this element.
	Selector string `json:"selector"`

	// XPath is an XPath expression for this element.
	XPath string `json:"xpath,omitempty"`

	// Section identifies the semantic page section.
	Section string `json:"section,omitempty"`

	// NearbyText contains text from adjacent siblings and parent.
	NearbyText []string `json:"nearby_text,omitempty"`

	// Children contains child PageElements.
	Children []*PageElement `json:"children,omitempty"`

	// Depth is the nesting depth in the DOM tree.
	Depth int `json:"depth"`
}

// ToBrowserElement converts a PageElement to a browser.Element for use in actions.
func (pe *PageElement) ToBrowserElement() *browser.Element {
	return &browser.Element{
		NodeID:      pe.NodeID,
		Tag:         pe.Tag,
		Type:        pe.Type,
		Text:        pe.Text,
		Role:        pe.Role,
		Visible:     pe.Visible,
		Clickable:   pe.Clickable,
		Position:    pe.Position,
		Attributes:  pe.Attributes,
		Placeholder: pe.Placeholder,
		Value:       pe.Value,
		Selector:    pe.Selector,
		XPath:       pe.XPath,
		NearbyText:  pe.NearbyText,
		Section:     pe.Section,
	}
}

// semanticTypeMap maps HTML tags to their semantic type.
var semanticTypeMap = map[string]string{
	"button":   "button",
	"a":        "link",
	"input":    "input",
	"textarea": "textarea",
	"select":   "select",
	"form":     "form",
	"img":      "image",
	"label":    "label",
	"h1":       "heading",
	"h2":       "heading",
	"h3":       "heading",
	"h4":       "heading",
	"h5":       "heading",
	"h6":       "heading",
	"p":        "paragraph",
	"nav":      "navigation",
	"header":   "header",
	"footer":   "footer",
	"main":     "main",
	"section":  "section",
	"article":  "article",
	"aside":    "aside",
	"table":    "table",
	"th":       "table-header",
	"td":       "table-cell",
	"li":       "list-item",
	"ul":       "list",
	"ol":       "list",
}

// SemanticType returns the semantic type for a given tag and input type.
func SemanticType(tag, inputType string) string {
	if tag == "input" {
		switch inputType {
		case "submit", "button", "reset":
			return "button"
		case "text", "email", "password", "search", "tel", "url", "number":
			return inputType + "-input"
		case "checkbox":
			return "checkbox"
		case "radio":
			return "radio"
		case "file":
			return "file-input"
		}
	}
	if t, ok := semanticTypeMap[tag]; ok {
		return t
	}
	return "element"
}

// ARIARole returns the effective ARIA role for a tag.
func ARIARole(tag, ariaRole string) string {
	if ariaRole != "" {
		return ariaRole
	}
	implicit := map[string]string{
		"button":  "button",
		"a":       "link",
		"input":   "textbox",
		"select":  "listbox",
		"nav":     "navigation",
		"header":  "banner",
		"footer":  "contentinfo",
		"main":    "main",
		"section": "region",
		"article": "article",
		"aside":   "complementary",
		"form":    "form",
		"h1":      "heading",
		"h2":      "heading",
		"h3":      "heading",
		"img":     "img",
		"table":   "table",
		"ul":      "list",
		"ol":      "list",
		"li":      "listitem",
	}
	if r, ok := implicit[tag]; ok {
		return r
	}
	return ""
}
