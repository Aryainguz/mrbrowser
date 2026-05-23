package resolver

import (
	"fmt"
	"sort"

	"github.com/mrbrowser/mrbrowser/core/browser"
	"github.com/mrbrowser/mrbrowser/intelligence/dom"
	"github.com/mrbrowser/mrbrowser/telemetry"
)

// MinConfidence is the minimum score a candidate must achieve to be returned.
const MinConfidence = 0.40

// ScoredCandidate pairs a DOM element with its resolver confidence score.
type ScoredCandidate struct {
	Element    *browser.Element `json:"element"`
	Confidence float64          `json:"confidence"`
}

// Resolver scores DOM candidates against a user intent and returns the best match.
type Resolver struct {
	weights Weights
	log     *telemetry.Logger
}

// New creates a Resolver with default scoring weights.
func New() *Resolver {
	return &Resolver{
		weights: DefaultWeights(),
		log:     telemetry.New("resolver"),
	}
}

// WithWeights creates a Resolver with custom scoring weights.
func WithWeights(w Weights) *Resolver {
	return &Resolver{
		weights: w,
		log:     telemetry.New("resolver"),
	}
}

// Resolve finds the best-matching element for the given intent string
// from a list of DOM elements.
//
// Returns:
//   - The best-matching browser.Element
//   - All scored candidates (sorted descending by confidence)
//   - An error if no candidate exceeds MinConfidence
func (r *Resolver) Resolve(target string, elements []*dom.PageElement) (*browser.Element, []ScoredCandidate, error) {
	if len(elements) == 0 {
		return nil, nil, fmt.Errorf("resolve %q: no elements to search", target)
	}

	intent := Parse(target)
	r.log.Debug("Resolving intent",
		telemetry.F("target", target),
		telemetry.F("keywords", intent.Keywords),
		telemetry.F("role", intent.InferredRole),
		telemetry.F("type", intent.InferredType),
	)

	candidates := make([]ScoredCandidate, 0, len(elements))
	for _, el := range elements {
		score := ScoreCandidate(el, intent, r.weights)
		if score >= MinConfidence {
			browserEl := el.ToBrowserElement()
			browserEl.Confidence = score
			candidates = append(candidates, ScoredCandidate{
				Element:    browserEl,
				Confidence: score,
			})
		}
	}

	if len(candidates) == 0 {
		return nil, nil, fmt.Errorf("resolve %q: no element matched with confidence >= %.0f%%", target, MinConfidence*100)
	}

	// Sort descending by confidence
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Confidence > candidates[j].Confidence
	})

	best := candidates[0]

	r.log.Info("Element resolved",
		telemetry.F("target", target),
		telemetry.F("matched", best.Element.Text),
		telemetry.F("confidence", fmt.Sprintf("%.0f%%", best.Confidence*100)),
		telemetry.F("selector", best.Element.Selector),
	)

	// Log top alternatives for debugging
	if len(candidates) > 1 {
		for i := 1; i < len(candidates) && i < 3; i++ {
			c := candidates[i]
			r.log.Debug("Alternative candidate",
				telemetry.F("rank", i+1),
				telemetry.F("text", c.Element.Text),
				telemetry.F("confidence", fmt.Sprintf("%.0f%%", c.Confidence*100)),
			)
		}
	}

	return best.Element, candidates, nil
}

// ResolveAll returns all candidates above MinConfidence, sorted by confidence.
// Useful for inspection/debugging.
func (r *Resolver) ResolveAll(target string, elements []*dom.PageElement) ([]ScoredCandidate, error) {
	_, candidates, err := r.Resolve(target, elements)
	return candidates, err
}
