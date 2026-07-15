package report

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// GenerateMarkdown creates a Markdown report and writes it to disk.
func (t *Tracker) GenerateMarkdown(taskName, outputFile string) error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# QA Automation Report: %s\n\n", taskName))
	sb.WriteString(fmt.Sprintf("**Date:** %s\n", time.Now().Format(time.RFC1123)))
	sb.WriteString(fmt.Sprintf("**Total Steps:** %d\n", len(t.Steps)))
	sb.WriteString(fmt.Sprintf("**Failed Steps:** %d\n", t.GetErrorCount()))
	
	failedAPIs := t.GetFailedAPIs()
	sb.WriteString(fmt.Sprintf("**Failing APIs (4xx/5xx):** %d\n\n", len(failedAPIs)))

	// Executive Summary (Failing APIs)
	if len(failedAPIs) > 0 {
		sb.WriteString("## 🚨 Failing APIs\n\n")
		sb.WriteString("| Method | Status | Duration | URL |\n")
		sb.WriteString("|--------|--------|----------|-----|\n")
		for _, api := range failedAPIs {
			sb.WriteString(fmt.Sprintf("| `%s` | **%d** | %s | `%s` |\n", api.Method, api.StatusCode, api.Duration.Round(time.Millisecond), api.URL))
		}
		sb.WriteString("\n")
	}

	// Execution Table
	sb.WriteString("## 📋 Execution Steps\n\n")
	sb.WriteString("| Step | Action | Target | Duration | Status | Screenshot/Details |\n")
	sb.WriteString("|------|--------|--------|----------|--------|-------------------|\n")

	for _, s := range t.Steps {
		status := "✅ PASS"
		if s.Error != nil {
			status = "❌ FAIL"
		}
		
		target := s.Target
		if target == "" {
			target = "-"
		}

		details := ""
		if s.Screenshot != "" {
			details = fmt.Sprintf("📸 [%s](./%s)", s.Screenshot, s.Screenshot)
		}
		if s.Error != nil {
			if details != "" {
				details += "<br>"
			}
			details += fmt.Sprintf("⚠️ `%s`", s.Error.Error())
		}

		sb.WriteString(fmt.Sprintf("| %d | `%s` | %s | %s | %s | %s |\n",
			s.Index, s.Action, target, s.Duration.Round(time.Millisecond), status, details))
	}

	sb.WriteString("\n\n---\n*Report generated natively by MrBrowser*")

	return os.WriteFile(outputFile, []byte(sb.String()), 0644)
}
