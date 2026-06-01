package unit

import (
	"testing"

	"github.com/mrbrowser/mrbrowser/intelligence/dom"
	"github.com/mrbrowser/mrbrowser/intelligence/resolver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// makeElement creates a test PageElement.
func makeElement(tag, text, role, elType, section string, visible, clickable bool) *dom.PageElement {
	return &dom.PageElement{
		Tag:       tag,
		Text:      text,
		Role:      role,
		Type:      elType,
		Section:   section,
		Visible:   visible,
		Clickable: clickable,
		Selector:  "#" + text,
		Attributes: map[string]string{
			"id":          text,
			"aria-label":  text,
			"placeholder": "",
		},
	}
}

var testElements = []*dom.PageElement{
	makeElement("button", "Login", "button", "button", "header", true, true),
	makeElement("button", "Sign In", "button", "button", "header", true, true),
	makeElement("input", "Email address", "textbox", "email-input", "main", true, false),
	makeElement("input", "Password", "textbox", "password-input", "main", true, false),
	makeElement("button", "Cancel", "button", "button", "main", true, true),
	makeElement("button", "Submit", "button", "button", "main", true, true),
	makeElement("a", "Forgot password?", "link", "link", "main", true, true),
	makeElement("button", "Download Invoice", "button", "button", "main", true, true),
	makeElement("input", "Search products", "textbox", "text-input", "main", true, false),
	makeElement("button", "Add to cart", "button", "button", "main", true, true),
}

func TestResolver_ExactTextMatch(t *testing.T) {
	r := resolver.New()
	el, candidates, err := r.Resolve("Login", testElements)
	require.NoError(t, err)
	assert.Equal(t, "Login", el.Text)
	assert.GreaterOrEqual(t, len(candidates), 1)
	assert.GreaterOrEqual(t, candidates[0].Confidence, 0.70)
}

func TestResolver_IntentMatch_Button(t *testing.T) {
	r := resolver.New()
	// "sign in button" should match "Sign In" button
	el, _, err := r.Resolve("sign in button", testElements)
	require.NoError(t, err)
	assert.Contains(t, el.Text, "Sign In")
}

func TestResolver_InputField(t *testing.T) {
	r := resolver.New()
	el, _, err := r.Resolve("email field", testElements)
	require.NoError(t, err)
	assert.Equal(t, "Email address", el.Text)
}

func TestResolver_PasswordField(t *testing.T) {
	r := resolver.New()
	el, _, err := r.Resolve("password input", testElements)
	require.NoError(t, err)
	assert.Equal(t, "Password", el.Text)
}

func TestResolver_SemanticSynonym(t *testing.T) {
	r := resolver.New()
	// "download invoice" → "Download Invoice"
	el, _, err := r.Resolve("download invoice", testElements)
	require.NoError(t, err)
	assert.Equal(t, "Download Invoice", el.Text)
}

func TestResolver_SearchInput(t *testing.T) {
	r := resolver.New()
	el, _, err := r.Resolve("search box", testElements)
	require.NoError(t, err)
	assert.Equal(t, "Search products", el.Text)
}

func TestResolver_CartButton(t *testing.T) {
	r := resolver.New()
	el, _, err := r.Resolve("add to cart", testElements)
	require.NoError(t, err)
	assert.Equal(t, "Add to cart", el.Text)
}

func TestResolver_NoMatch(t *testing.T) {
	r := resolver.New()
	_, _, err := r.Resolve("nonexistent xyz element", testElements)
	assert.Error(t, err)
}

func TestResolver_CandidatesRankedDescending(t *testing.T) {
	r := resolver.New()
	_, candidates, err := r.Resolve("login button", testElements)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(candidates), 2)

	for i := 1; i < len(candidates); i++ {
		assert.GreaterOrEqual(t, candidates[i-1].Confidence, candidates[i].Confidence,
			"candidates should be ranked descending by confidence")
	}
}

func TestResolver_EmptyElements(t *testing.T) {
	r := resolver.New()
	_, _, err := r.Resolve("login", []*dom.PageElement{})
	assert.Error(t, err)
}

// ─────────────────────────────────────────────────────────────
// Intent parser tests
// ─────────────────────────────────────────────────────────────

func TestIntentParse_Button(t *testing.T) {
	intent := resolver.Parse("login button")
	assert.Equal(t, "button", intent.InferredType)
	assert.Equal(t, "button", intent.InferredRole)
	assert.Contains(t, intent.Keywords, "login")
}

func TestIntentParse_InputField(t *testing.T) {
	intent := resolver.Parse("email field")
	assert.Equal(t, "input", intent.InferredType)
	assert.Contains(t, intent.Keywords, "email")
}

func TestIntentParse_NoNoise(t *testing.T) {
	intent := resolver.Parse("the login button")
	// "the" should be excluded as noise word
	assert.NotContains(t, intent.Keywords, "the")
	assert.Contains(t, intent.Keywords, "login")
}

func TestIntentParse_CaseInsensitive(t *testing.T) {
	a := resolver.Parse("LOGIN BUTTON")
	b := resolver.Parse("login button")
	assert.Equal(t, a.InferredType, b.InferredType)
}

// ─────────────────────────────────────────────────────────────
// Scorer tests
// ─────────────────────────────────────────────────────────────

func TestScorer_ExactTextHighScore(t *testing.T) {
	el := makeElement("button", "Login", "button", "button", "main", true, true)
	intent := resolver.Parse("Login")
	score := resolver.ScoreCandidate(el, intent, resolver.DefaultWeights())
	assert.GreaterOrEqual(t, score, 0.70, "exact text match should score >= 70%%")
}

func TestScorer_InvisiblePenalty(t *testing.T) {
	visible := makeElement("button", "Login", "button", "button", "main", true, true)
	invisible := makeElement("button", "Login", "button", "button", "main", false, false)
	intent := resolver.Parse("Login")
	w := resolver.DefaultWeights()

	visScore := resolver.ScoreCandidate(visible, intent, w)
	invScore := resolver.ScoreCandidate(invisible, intent, w)
	assert.Greater(t, visScore, invScore, "visible elements should score higher")
}

func TestScorer_RoleBonus(t *testing.T) {
	buttonEl := makeElement("button", "Submit", "button", "button", "main", true, true)
	divEl := makeElement("div", "Submit", "", "element", "main", true, false)
	intent := resolver.Parse("submit button")
	w := resolver.DefaultWeights()

	btnScore := resolver.ScoreCandidate(buttonEl, intent, w)
	divScore := resolver.ScoreCandidate(divEl, intent, w)
	assert.Greater(t, btnScore, divScore, "button element should score higher than div for button intent")
}
