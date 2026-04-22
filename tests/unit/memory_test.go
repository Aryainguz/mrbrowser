package unit

import (
	"testing"
	"time"

	"github.com/mrbrowser/mrbrowser/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *memory.Store {
	store, err := memory.Open(":memory:")
	require.NoError(t, err)
	return store
}

func TestStore_SaveAndGetFingerprint(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	fp := &memory.ElementFingerprint{
		ID:           "test-id-1",
		WorkflowID:   "test-wf",
		StepName:     "click-login",
		Target:       "login button",
		Selector:     "#login-btn",
		Text:         "Login",
		Role:         "button",
		Type:         "button",
		Section:      "header",
		NearbyText:   []string{"Forgot Password", "Sign Up"},
		Attributes:   map[string]string{"id": "login-btn", "class": "btn primary"},
		LastSeen:     time.Now().Round(time.Second),
		UseCount:     5,
		SuccessCount: 4,
		SuccessRate:  0.8,
	}

	err := store.SaveFingerprint(fp)
	require.NoError(t, err)

	fps, err := store.GetFingerprints("test-wf", "login button")
	require.NoError(t, err)
	require.Len(t, fps, 1)

	saved := fps[0]
	assert.Equal(t, "test-id-1", saved.ID)
	assert.Equal(t, "login button", saved.Target)
	assert.Equal(t, "Login", saved.Text)
	assert.Equal(t, []string{"Forgot Password", "Sign Up"}, saved.NearbyText)
	assert.Equal(t, "login-btn", saved.Attributes["id"])
	assert.Equal(t, fp.LastSeen.UTC(), saved.LastSeen.UTC())
}

func TestStore_UpdateFingerprint(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	fp := &memory.ElementFingerprint{
		ID:         "test-id-1",
		WorkflowID: "test-wf",
		Target:     "login",
		Selector:   "#old",
	}
	require.NoError(t, store.SaveFingerprint(fp))

	fp.Selector = "#new"
	fp.SuccessCount = 10
	require.NoError(t, store.SaveFingerprint(fp))

	fps, _ := store.GetFingerprints("test-wf", "login")
	assert.Equal(t, "#new", fps[0].Selector)
	assert.Equal(t, 10, fps[0].SuccessCount)
}

func TestStore_WorkflowResults(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	res := &memory.WorkflowResult{
		ID:           "run-1",
		WorkflowName: "daily-check",
		StartedAt:    time.Now().Add(-time.Minute).Round(time.Second),
		CompletedAt:  time.Now().Round(time.Second),
		Duration:     time.Minute,
		Success:      true,
		StepsTotal:   5,
		StepsPassed:  5,
	}

	require.NoError(t, store.SaveWorkflowResult(res))

	history, err := store.GetWorkflowHistory("daily-check", 10)
	require.NoError(t, err)
	require.Len(t, history, 1)

	assert.Equal(t, "run-1", history[0].ID)
	assert.True(t, history[0].Success)
	assert.Equal(t, time.Minute, history[0].Duration)
}
