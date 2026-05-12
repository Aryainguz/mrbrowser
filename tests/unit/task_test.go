package unit

import (
	"testing"

	"github.com/mrbrowser/mrbrowser/core/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTask_LoadYAML(t *testing.T) {
	yamlData := []byte(`
task:
  name: test-workflow
  steps:
    - open:
        url: https://example.com
    - click:
        target: login button
    - type:
        target: email
        value: user@test.com
    - assert:
        text_visible: Welcome
`)

	task, err := runtime.Load(yamlData)
	require.NoError(t, err)

	assert.Equal(t, "test-workflow", task.Name)
	require.Len(t, task.Steps, 4)

	assert.Equal(t, "open", task.Steps[0].Kind())
	assert.Equal(t, "https://example.com", task.Steps[0].Open.URL)

	assert.Equal(t, "click", task.Steps[1].Kind())
	assert.Equal(t, "login button", task.Steps[1].Click.Target)

	assert.Equal(t, "type", task.Steps[2].Kind())
	assert.Equal(t, "email", task.Steps[2].Type.Target)
	assert.Equal(t, "user@test.com", task.Steps[2].Type.Value)

	assert.Equal(t, "assert", task.Steps[3].Kind())
	assert.Equal(t, "Welcome", task.Steps[3].Assert.TextVisible)
}

func TestTask_LoadYAMLDirect(t *testing.T) {
	// Test parsing without the top-level 'task:' wrapper
	yamlData := []byte(`
name: direct-workflow
steps:
  - open:
      url: https://example.com
`)

	task, err := runtime.Load(yamlData)
	require.NoError(t, err)
	assert.Equal(t, "direct-workflow", task.Name)
	require.Len(t, task.Steps, 1)
}

func TestTask_InvalidStep(t *testing.T) {
	yamlData := []byte(`
name: bad
steps:
  - unknown_action:
      foo: bar
`)
	_, err := runtime.Load(yamlData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no recognized action")
}
