// Package resolver implements the intent-based element resolver for Mr. Browser.
// It converts a natural-language target description into a scored list of DOM elements.
package resolver

import (
	"strings"
	"unicode"
)

// Intent represents a parsed user intent for finding an element.
type Intent struct {
	// Raw is the original user-provided string.
	Raw string

	// Tokens are the normalized, split words from Raw.
	Tokens []string

	// Noun is the primary noun (usually the element concept, e.g., "button", "field", "link").
	Noun string

	// Keywords are the significant content words (excluding noise words).
	Keywords []string

	// InferredRole is the ARIA/semantic role implied by the intent.
	InferredRole string

	// InferredType is the element type implied by the intent.
	InferredType string
}

// Parse converts a natural-language target string into a structured Intent.
func Parse(target string) Intent {
	raw := strings.TrimSpace(target)
	tokens := tokenize(raw)
	keywords, noun := extractKeywordsAndNoun(tokens)
	role, elType := inferRoleAndType(tokens)

	return Intent{
		Raw:          raw,
		Tokens:       tokens,
		Noun:         noun,
		Keywords:     keywords,
		InferredRole: role,
		InferredType: elType,
	}
}

// tokenize splits the target string into normalized lowercase tokens.
func tokenize(s string) []string {
	s = strings.ToLower(s)
	var tokens []string
	var buf strings.Builder

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			buf.WriteRune(r)
		} else {
			if buf.Len() > 0 {
				tokens = append(tokens, buf.String())
				buf.Reset()
			}
		}
	}
	if buf.Len() > 0 {
		tokens = append(tokens, buf.String())
	}
	return tokens
}

// noiseWords are articles and prepositions that carry no semantic value.
var noiseWords = map[string]bool{
	"the": true, "a": true, "an": true, "for": true, "of": true,
	"in": true, "on": true, "at": true, "to": true, "with": true,
	"and": true, "or": true, "my": true, "your": true, "their": true,
}

// elementNouns are words that describe what kind of DOM element is wanted.
var elementNouns = map[string]string{
	"button":   "button",
	"btn":      "button",
	"link":     "link",
	"input":    "input",
	"field":    "input",
	"box":      "input",
	"checkbox": "checkbox",
	"radio":    "radio",
	"dropdown": "select",
	"select":   "select",
	"menu":     "select",
	"textarea": "textarea",
	"form":     "form",
	"tab":      "tab",
	"option":   "option",
	"item":     "listitem",
	"heading":  "heading",
	"title":    "heading",
	"image":    "image",
	"icon":     "image",
	"submit":   "button",
	"search":   "input",
}

// roleMap maps element nouns to ARIA roles.
var roleMap = map[string]string{
	"button":   "button",
	"link":     "link",
	"input":    "textbox",
	"checkbox": "checkbox",
	"radio":    "radio",
	"select":   "listbox",
	"tab":      "tab",
	"heading":  "heading",
	"image":    "img",
	"listitem": "listitem",
}

// actionWords are verbs that hint at the element type.
var actionWords = map[string]string{
	"click":    "button",
	"press":    "button",
	"submit":   "button",
	"type":     "input",
	"enter":    "input",
	"fill":     "input",
	"search":   "input",
	"select":   "select",
	"choose":   "select",
	"check":    "checkbox",
	"navigate": "link",
	"go":       "link",
	"download": "link",
	"upload":   "input",
	"sign":     "button",
	"log":      "button",
}

func extractKeywordsAndNoun(tokens []string) (keywords []string, noun string) {
	for _, t := range tokens {
		if noiseWords[t] {
			continue
		}
		if n, ok := elementNouns[t]; ok {
			noun = n
			continue
		}
		if actionWords[t] != "" {
			// Action words inform the type but aren't content keywords
			continue
		}
		keywords = append(keywords, t)
	}
	return keywords, noun
}

func inferRoleAndType(tokens []string) (role, elType string) {
	for _, t := range tokens {
		// Direct element noun match
		if n, ok := elementNouns[t]; ok {
			elType = n
			role = roleMap[n]
			return
		}
		// Action word match
		if et, ok := actionWords[t]; ok {
			elType = et
			role = roleMap[et]
			return
		}
	}

	// Default: no strong signal
	return "", ""
}
