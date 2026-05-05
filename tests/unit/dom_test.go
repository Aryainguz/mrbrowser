package unit

import (
	"testing"

	"github.com/mrbrowser/mrbrowser/intelligence/dom"
	"github.com/stretchr/testify/assert"
)

func TestDOM_SemanticType(t *testing.T) {
	tests := []struct {
		tag       string
		inputType string
		expected  string
	}{
		{"button", "", "button"},
		{"a", "", "link"},
		{"input", "text", "text-input"},
		{"input", "password", "password-input"},
		{"input", "submit", "button"},
		{"input", "checkbox", "checkbox"},
		{"input", "radio", "radio"},
		{"select", "", "select"},
		{"textarea", "", "textarea"},
		{"h1", "", "heading"},
		{"nav", "", "navigation"},
		{"form", "", "form"},
		{"div", "", "element"}, // unknown falls back to element
	}

	for _, tc := range tests {
		actual := dom.SemanticType(tc.tag, tc.inputType)
		assert.Equal(t, tc.expected, actual, "tag:%s type:%s", tc.tag, tc.inputType)
	}
}

func TestDOM_ARIARole(t *testing.T) {
	tests := []struct {
		tag      string
		role     string
		expected string
	}{
		{"div", "button", "button"},   // explicit role overrides
		{"nav", "banner", "banner"},   // explicit role overrides
		{"button", "", "button"},      // implicit
		{"a", "", "link"},             // implicit
		{"header", "", "banner"},      // implicit
		{"footer", "", "contentinfo"}, // implicit
		{"main", "", "main"},          // implicit
		{"input", "", "textbox"},      // implicit
		{"div", "", ""},               // no implicit role
	}

	for _, tc := range tests {
		actual := dom.ARIARole(tc.tag, tc.role)
		assert.Equal(t, tc.expected, actual, "tag:%s role:%s", tc.tag, tc.role)
	}
}
