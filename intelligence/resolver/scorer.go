package resolver

import (
	"math"
	"strings"

	"github.com/mrbrowser/mrbrowser/intelligence/dom"
)

// Weights controls how much each scoring dimension contributes to the final score.
// These are tuned empirically and can be overridden.
type Weights struct {
	TextSimilarity float64 // Normalized edit-distance similarity of text
	KeywordMatch   float64 // Fraction of keywords found in element text/attrs
	RoleMatch      float64 // Bonus for ARIA role match
	TypeMatch      float64 // Bonus for semantic type match
	Visibility     float64 // Penalty for invisible elements
	Clickability   float64 // Bonus for clickable elements when intent implies action
	PositionScore  float64 // Prefer elements in the upper/central viewport
	AttributeMatch float64 // Keywords found in id, name, placeholder, aria-label
}

// DefaultWeights returns the default scoring weights.
func DefaultWeights() Weights {
	return Weights{
		TextSimilarity: 0.40,
		KeywordMatch:   0.25,
		RoleMatch:      0.10,
		TypeMatch:      0.10,
		Visibility:     0.05,
		Clickability:   0.05,
		PositionScore:  0.02,
		AttributeMatch: 0.03,
	}
}

// ScoreCandidate scores a single DOM element against a parsed intent.
// Returns a confidence score in [0.0, 1.0].
func ScoreCandidate(el *dom.PageElement, intent Intent, weights Weights) float64 {
	elText := strings.ToLower(strings.TrimSpace(el.Text))
	elPlaceholder := strings.ToLower(el.Placeholder)
	elID := strings.ToLower(el.Attributes["id"])
	elName := strings.ToLower(el.Attributes["name"])
	elAriaLabel := strings.ToLower(el.Attributes["aria-label"])
	elClass := strings.ToLower(el.Attributes["class"])

	// Combine all text signals
	allText := strings.Join([]string{elText, elPlaceholder, elAriaLabel}, " ")
	allAttrs := strings.Join([]string{elID, elName, elPlaceholder, elAriaLabel, elClass}, " ")

	intentLower := strings.ToLower(intent.Raw)

	// 1. Text similarity (normalized Levenshtein against the raw intent)
	textSim := textSimilarity(intentLower, elText)
	// Also try similarity against combined text
	combinedSim := textSimilarity(intentLower, allText)
	if combinedSim > textSim {
		textSim = combinedSim
	}

	// 2. Keyword match: what fraction of intent keywords appear in element text/attrs?
	kwScore := keywordMatchScore(intent.Keywords, allText)
	attrKwScore := keywordMatchScore(intent.Keywords, allAttrs)

	// 3. Role match
	roleScore := 0.0
	if intent.InferredRole != "" && strings.EqualFold(el.Role, intent.InferredRole) {
		roleScore = 1.0
	}

	// 4. Type match
	typeScore := 0.0
	if intent.InferredType != "" && strings.EqualFold(el.Type, intent.InferredType) {
		typeScore = 1.0
	} else if intent.InferredType != "" && strings.Contains(el.Type, intent.InferredType) {
		typeScore = 0.5
	}

	// 5. Visibility: invisible elements are penalized
	visScore := 0.0
	if el.Visible {
		visScore = 1.0
	}

	// 6. Clickability: if intent implies clicking, prefer clickable elements
	clickScore := 0.0
	if el.Clickable {
		clickScore = 1.0
	}

	// 7. Position: prefer elements in the upper 60% of the viewport (heuristic)
	posScore := 0.0
	if el.Position.Y < 800 && el.Position.Y >= 0 {
		// Normalize: y=0 → 1.0, y=800 → 0.0
		posScore = math.Max(0, 1.0-el.Position.Y/800.0)
	}

	// 8. Attribute match: keywords in id/name/placeholder/aria-label
	attrScore := attrKwScore

	score := weights.TextSimilarity*textSim +
		weights.KeywordMatch*kwScore +
		weights.RoleMatch*roleScore +
		weights.TypeMatch*typeScore +
		weights.Visibility*visScore +
		weights.Clickability*clickScore +
		weights.PositionScore*posScore +
		weights.AttributeMatch*attrScore

	// Hard penalty: invisible elements with no text similarity are very unlikely
	if !el.Visible && textSim < 0.3 {
		score *= 0.1
	}

	return math.Min(1.0, math.Max(0.0, score))
}

// textSimilarity returns a normalized similarity in [0,1] between two strings.
// Uses a combination of:
//   - Exact match (1.0)
//   - Contains match (0.8)
//   - Token overlap
//   - Normalized Levenshtein
func textSimilarity(a, b string) float64 {
	if a == "" || b == "" {
		return 0
	}
	if a == b {
		return 1.0
	}
	// Exact substring
	if strings.Contains(b, a) {
		return 0.85
	}
	if strings.Contains(a, b) {
		return 0.80
	}

	// Token overlap
	aTokens := strings.Fields(a)
	bTokens := strings.Fields(b)
	tokenOverlap := jaccardSimilarity(aTokens, bTokens)

	// Levenshtein similarity
	levSim := 1.0 - float64(levenshtein(a, b))/float64(max(len(a), len(b)))

	// Take the best signal
	best := math.Max(tokenOverlap, levSim)
	return best
}

// keywordMatchScore returns the fraction of keywords found in the text.
func keywordMatchScore(keywords []string, text string) float64 {
	if len(keywords) == 0 {
		return 0
	}
	hits := 0
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			hits++
		}
	}
	return float64(hits) / float64(len(keywords))
}

// jaccardSimilarity computes the Jaccard similarity between two token slices.
func jaccardSimilarity(a, b []string) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 1.0
	}
	setA := make(map[string]bool, len(a))
	for _, t := range a {
		setA[t] = true
	}
	setB := make(map[string]bool, len(b))
	for _, t := range b {
		setB[t] = true
	}
	intersection := 0
	for t := range setA {
		if setB[t] {
			intersection++
		}
	}
	union := len(setA) + len(setB) - intersection
	if union == 0 {
		return 0
	}
	return float64(intersection) / float64(union)
}

// levenshtein computes the edit distance between two strings.
func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)

	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	// Use two-row DP
	prev := make([]int, lb+1)
	curr := make([]int, lb+1)

	for j := 0; j <= lb; j++ {
		prev[j] = j
	}

	for i := 1; i <= la; i++ {
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			curr[j] = minInt(
				curr[j-1]+1,
				minInt(prev[j]+1, prev[j-1]+cost),
			)
		}
		prev, curr = curr, prev
	}

	return prev[lb]
}

func minInt(a, b int) int {
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
