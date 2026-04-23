package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrbrowser/mrbrowser/core/browser"
	"github.com/mrbrowser/mrbrowser/core/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBrowser_LaunchAndNavigate(t *testing.T) {
	b, err := browser.Launch(browser.DefaultOptions())
	require.NoError(t, err)
	defer b.Close()

	page, err := b.NewPage()
	require.NoError(t, err)

	err = page.Navigate("about:blank")
	require.NoError(t, err)

	url, err := page.URL()
	require.NoError(t, err)
	assert.Equal(t, "about:blank", url)
}

func TestBrowser_Screenshot(t *testing.T) {
	b, err := browser.Launch(browser.DefaultOptions())
	require.NoError(t, err)
	defer b.Close()

	page, err := b.NewPage()
	require.NoError(t, err)

	err = page.Navigate("about:blank")
	require.NoError(t, err)

	buf, err := page.Screenshot()
	require.NoError(t, err)
	assert.Greater(t, len(buf), 100, "screenshot should contain data")
}

func TestBrowser_ExecuteJS(t *testing.T) {
	b, err := browser.Launch(browser.DefaultOptions())
	require.NoError(t, err)
	defer b.Close()

	page, err := b.NewPage()
	require.NoError(t, err)

	err = page.Navigate("about:blank")
	require.NoError(t, err)

	res, err := page.ExecuteJS("1 + 2")
	require.NoError(t, err)

	// Chromedp returns float64 for numbers
	val, ok := res.(float64)
	require.True(t, ok)
	assert.Equal(t, 3.0, val)
}

func TestRuntime_ExecuteTask(t *testing.T) {
	// Start a local test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<body>
				<h1 id="title">Login Page</h1>
				<input type="text" placeholder="Username" id="user" />
				<input type="password" placeholder="Password" id="pass" />
				<button id="login-btn" onclick="document.getElementById('title').innerText = 'Success!'">Sign In</button>
			</body>
			</html>
		`))
	}))
	defer ts.Close()

	yamlData := []byte(`
task:
  name: integration-test
  steps:
    - open:
        url: ` + ts.URL + `
    - assert:
        text_visible: Login Page
    - type:
        target: Username
        value: testuser
    - type:
        target: Password
        value: secret
    - click:
        target: Sign In
    - assert:
        text_visible: Success!
`)

	task, err := runtime.Load(yamlData)
	require.NoError(t, err)

	session, err := browser.NewSession("int-test", browser.DefaultOptions())
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

	assert.True(t, result.Success, "Task execution failed: %v", result.Error)
	assert.Len(t, result.StepResults, 6)

	for _, sr := range result.StepResults {
		assert.True(t, sr.Success, "step %d failed: %v", sr.Step, sr.Error)
	}
}
