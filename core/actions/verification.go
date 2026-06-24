package actions

import "strings"

// VerificationResult describes what changed after an action was performed.
type VerificationResult struct {
	// URLChanged indicates the page URL changed after the action.
	URLChanged bool `json:"url_changed"`

	// DOMChanged indicates the page DOM changed after the action.
	DOMChanged bool `json:"dom_changed"`

	// PreviousURL is the URL before the action.
	PreviousURL string `json:"previous_url,omitempty"`

	// CurrentURL is the URL after the action.
	CurrentURL string `json:"current_url,omitempty"`
}

// verifyDOMChange compares before/after state and returns a VerificationResult.
func verifyDOMChange(preURL, postURL, preHTML, postHTML string) *VerificationResult {
	v := &VerificationResult{
		PreviousURL: preURL,
		CurrentURL:  postURL,
		URLChanged:  preURL != postURL,
	}

	// For DOM change, we compare lengths and a quick diff.
	// Exact comparison would be too expensive; length delta + content sample is sufficient.
	if len(preHTML) != len(postHTML) {
		v.DOMChanged = true
	} else {
		// Sample 3 sections: beginning, middle, end
		sections := [][2]int{
			{0, min(500, len(preHTML))},
			{len(preHTML)/2 - 250, len(preHTML)/2 + 250},
			{max(0, len(preHTML)-500), len(preHTML)},
		}
		for _, s := range sections {
			start, end := s[0], s[1]
			if end > len(preHTML) {
				end = len(preHTML)
			}
			if end > len(postHTML) {
				end = len(postHTML)
			}
			if start >= end {
				continue
			}
			if preHTML[start:end] != postHTML[start:end] {
				v.DOMChanged = true
				break
			}
		}
	}

	return v
}

// containsAny returns true if s contains any of the substrings.
func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
