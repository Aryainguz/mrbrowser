package memory

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/mrbrowser/mrbrowser/core/browser"
	"github.com/mrbrowser/mrbrowser/intelligence/dom"
	"github.com/mrbrowser/mrbrowser/telemetry"
)

// MinRecoveryConfidence is the minimum score to accept a recovered element.
const MinRecoveryConfidence = 0.65

// SelfHealer uses stored fingerprints to recover from broken element selectors.
type SelfHealer struct {
	store *Store
	log   *telemetry.Logger
}

// NewSelfHealer creates a SelfHealer backed by the given store.
func NewSelfHealer(store *Store) *SelfHealer {
	return &SelfHealer{
		store: store,
		log:   telemetry.New("healer"),
	}
}

// RecoveryResult describes the outcome of a self-healing attempt.
type RecoveryResult struct {
	// Recovered indicates self-healing was successful.
	Recovered bool

	// Element is the recovered element.
	Element *browser.Element

	// Confidence is the similarity score [0,1].
	Confidence float64

	// Note describes how recovery was achieved.
	Note string
}

// TryRecover attempts to find a replacement element when the original selector fails.
//
// Algorithm:
//  1. Load stored fingerprints for this workflow+target.
//  2. Score each current page element against the fingerprint.
//  3. Return the highest-scoring candidate if it exceeds MinRecoveryConfidence.
func (h *SelfHealer) TryRecover(
	workflowID, target string,
	currentElements []*dom.PageElement,
) (*RecoveryResult, error) {
	h.log.Warn("Attempting self-healing",
		telemetry.F("workflow", workflowID),
		telemetry.F("target", target),
	)

	fps, err := h.store.GetFingerprints(workflowID, target)
	if err != nil {
		return nil, fmt.Errorf("load fingerprints: %w", err)
	}

	if len(fps) == 0 {
		return &RecoveryResult{Recovered: false, Note: "no fingerprints stored for this target"}, nil
	}

	// Use the most successful fingerprint
	fp := fps[0]

	type scoredEl struct {
		el    *dom.PageElement
		score float64
	}

	var best *scoredEl

	for _, el := range currentElements {
		score := scoreAgainstFingerprint(el, fp)
		if best == nil || score > best.score {
			best = &scoredEl{el: el, score: score}
		}
	}

	if best == nil || best.score < MinRecoveryConfidence {
		confidence := 0.0
		if best != nil {
			confidence = best.score
		}
		return &RecoveryResult{
			Recovered:  false,
			Confidence: confidence,
			Note:       fmt.Sprintf("best candidate scored %.0f%% (threshold: %.0f%%)", confidence*100, MinRecoveryConfidence*100),
		}, nil
	}

	recoveredEl := best.el.ToBrowserElement()
	recoveredEl.Confidence = best.score

	h.log.Recover("Self-healing succeeded",
		telemetry.F("target", target),
		telemetry.F("old_text", fp.Text),
		telemetry.F("new_text", recoveredEl.Text),
		telemetry.F("confidence", fmt.Sprintf("%.0f%%", best.score*100)),
	)

	return &RecoveryResult{
		Recovered:  true,
		Element:    recoveredEl,
		Confidence: best.score,
		Note: fmt.Sprintf("recovered: %q → %q (%.0f%% confidence)",
			fp.Text, recoveredEl.Text, best.score*100),
	}, nil
}

// scoreAgainstFingerprint scores a DOM element against a stored fingerprint.
// Uses text similarity, role match, section match, and nearby text overlap.
func scoreAgainstFingerprint(el *dom.PageElement, fp *ElementFingerprint) float64 {
	var score float64

	// 1. Text similarity (40% weight)
	textSim := normalizedSimilarity(
		strings.ToLower(el.Text),
		strings.ToLower(fp.Text),
	)
	score += 0.40 * textSim

	// 2. Role match (20% weight)
	if fp.Role != "" && strings.EqualFold(el.Role, fp.Role) {
		score += 0.20
	} else if fp.Role != "" && strings.Contains(strings.ToLower(el.Role), strings.ToLower(fp.Role)) {
		score += 0.10
	}

	// 3. Section match (15% weight)
	if fp.Section != "" && strings.EqualFold(el.Section, fp.Section) {
		score += 0.15
	}

	// 4. Type match (10% weight)
	if fp.Type != "" && strings.EqualFold(el.Type, fp.Type) {
		score += 0.10
	}

	// 5. Nearby text overlap (15% weight)
	nearbyOverlap := nearbyTextOverlap(el.NearbyText, fp.NearbyText)
	score += 0.15 * nearbyOverlap

	// Visibility bonus
	if el.Visible {
		score = math.Min(1.0, score+0.02)
	}

	return math.Min(1.0, math.Max(0.0, score))
}

// nearbyTextOverlap computes the Jaccard similarity of two nearby-text slices.
func nearbyTextOverlap(a, b []string) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	setA := make(map[string]bool)
	for _, s := range a {
		setA[strings.ToLower(strings.TrimSpace(s))] = true
	}
	intersection := 0
	for _, s := range b {
		if setA[strings.ToLower(strings.TrimSpace(s))] {
			intersection++
		}
	}
	union := len(a) + len(b) - intersection
	if union == 0 {
		return 0
	}
	return float64(intersection) / float64(union)
}

// normalizedSimilarity returns a similarity score in [0,1] between two strings.
func normalizedSimilarity(a, b string) float64 {
	if a == b {
		return 1.0
	}
	if a == "" || b == "" {
		return 0
	}
	if strings.Contains(b, a) || strings.Contains(a, b) {
		return 0.85
	}
	maxLen := math.Max(float64(len(a)), float64(len(b)))
	dist := float64(levenshteinDist(a, b))
	return math.Max(0, 1.0-dist/maxLen)
}

func levenshteinDist(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
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
			curr[j] = minI(curr[j-1]+1, minI(prev[j]+1, prev[j-1]+cost))
		}
		prev, curr = curr, prev
	}
	return prev[lb]
}

func minI(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// StoreFingerprint creates and saves a fingerprint for a successfully resolved element.
func (h *SelfHealer) StoreFingerprint(workflowID, stepName, target string, el *browser.Element) error {
	fp := &ElementFingerprint{
		ID:           generateID(),
		WorkflowID:   workflowID,
		StepName:     stepName,
		Target:       target,
		Selector:     el.Selector,
		Text:         el.Text,
		Role:         el.Role,
		Type:         el.Type,
		Section:      el.Section,
		NearbyText:   el.NearbyText,
		Attributes:   el.Attributes,
		LastSeen:     time.Now(),
		UseCount:     1,
		SuccessCount: 1,
		SuccessRate:  1.0,
	}
	return h.store.SaveFingerprint(fp)
}

// generateID creates a simple time-based unique ID.
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
