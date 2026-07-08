package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrbrowser/mrbrowser/core/browser"
	"github.com/mrbrowser/mrbrowser/core/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_LoginWorkflow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<body>
	<h1>Welcome Back</h1>
	<form>
		<input type="email" placeholder="Email Address" />
		<input type="password" placeholder="Password" />
		<button type="button" onclick="document.body.innerHTML = '<h2>Dashboard</h2>'">Sign In</button>
	</form>
</body>
</html>
		`))
	}))
	defer ts.Close()

	yamlData := []byte(`
task:
  name: e2e-login
  steps:
    - open:
        url: ` + ts.URL + `
    - type:
        target: Email Address
        value: test@example.com
    - type:
        target: Password
        value: secret
    - click:
        target: Sign In
    - assert:
        text_visible: Dashboard
`)

	task, err := runtime.Load(yamlData)
	require.NoError(t, err)

	session, err := browser.NewSession("e2e", browser.DefaultOptions())
	require.NoError(t, err)
	defer session.Close()

	exec, err := runtime.NewExecutor(session, runtime.ExecutorOptions{
		StopOnError: true,
		DBPath:      ":memory:",
	})
	require.NoError(t, err)
	defer exec.Close()

	result, err := exec.Run(task)
	require.NoError(t, err)
	assert.True(t, result.Success, "Task failed: %v", result.Error)
}

func TestE2E_DynamicSelectorSelfHealing(t *testing.T) {
	// A page where the element ID changes
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<body>
	<button id="random-btn-12345" onclick="document.body.innerHTML = '<h2>Downloaded</h2>'">Download Report</button>
</body>
</html>
		`))
	}))
	defer ts.Close()

	yamlData := []byte(`
task:
  name: e2e-healing
  steps:
    - open:
        url: ` + ts.URL + `
    - click:
        target: Download Report
    - assert:
        text_visible: Downloaded
`)

	task, err := runtime.Load(yamlData)
	require.NoError(t, err)

	session, err := browser.NewSession("e2e-healing", browser.DefaultOptions())
	require.NoError(t, err)
	defer session.Close()

	exec, err := runtime.NewExecutor(session, runtime.ExecutorOptions{
		StopOnError: true,
		DBPath:      ":memory:",
	})
	require.NoError(t, err)
	defer exec.Close()

	result, err := exec.Run(task)
	require.NoError(t, err)
	assert.True(t, result.Success, "Task failed: %v", result.Error)
}
