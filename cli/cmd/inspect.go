package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mrbrowser/mrbrowser/core/browser"
	"github.com/mrbrowser/mrbrowser/intelligence/dom"
	"github.com/mrbrowser/mrbrowser/telemetry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	inspectJSON    bool
	inspectVisible bool
	inspectTarget  string
)

var inspectCmd = &cobra.Command{
	Use:   "inspect <url>",
	Short: "Inspect page elements at a URL",
	Long: `Navigate to a URL and print all detected interactive elements.

Example:
  mr-browser inspect https://example.com
  mr-browser inspect https://example.com --visible-only
  mr-browser inspect https://example.com --target "login button"
  mr-browser inspect https://example.com --json`,
	Args: cobra.ExactArgs(1),
	RunE: inspectPage,
}

func init() {
	inspectCmd.Flags().BoolVar(&inspectJSON, "json", false, "output as JSON")
	inspectCmd.Flags().BoolVar(&inspectVisible, "visible-only", false, "show only visible elements")
	inspectCmd.Flags().StringVar(&inspectTarget, "target", "", "resolve a specific intent and show ranked candidates")
}

func inspectPage(cmd *cobra.Command, args []string) error {
	log := telemetry.New("inspect")
	url := args[0]

	opts := browser.Options{
		Headless:     viper.GetBool("headless"),
		ChromiumPath: viper.GetString("chromium_path"),
		NoSandbox:    viper.GetBool("no_sandbox"),
	}

	session, err := browser.NewSession("inspect", opts)
	if err != nil {
		return fmt.Errorf("launch browser: %w", err)
	}
	defer func() { _ = session.Close() }()

	log.Step("Navigating", telemetry.F("url", url))
	if err := session.Page.Navigate(url); err != nil {
		return fmt.Errorf("navigate: %w", err)
	}

	// Brief settle time for JS-heavy pages
	time.Sleep(500 * time.Millisecond)

	extractor := dom.NewExtractor(session.Page)
	elements, err := extractor.Extract()
	if err != nil {
		return fmt.Errorf("extract elements: %w", err)
	}

	// Filter if needed
	if inspectVisible {
		filtered := make([]*dom.PageElement, 0, len(elements))
		for _, el := range elements {
			if el.Visible {
				filtered = append(filtered, el)
			}
		}
		elements = filtered
	}

	if inspectJSON {
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(elements)
	}

	// Pretty print
	title, _ := session.Page.Title()
	currentURL, _ := session.Page.URL()

	fmt.Printf("\n  🌐 %s\n", currentURL)
	if title != "" {
		fmt.Printf("  📄 %s\n", title)
	}
	fmt.Printf("  Found %d elements\n\n", len(elements))

	// Section headers
	currentSection := ""
	for _, el := range elements {
		if el.Section != currentSection {
			currentSection = el.Section
			fmt.Printf("  ── %s ─────────────────────\n", strings.ToUpper(el.Section))
		}

		visibility := "👁 "
		if !el.Visible {
			visibility = "   "
		}
		clickable := ""
		if el.Clickable {
			clickable = " [clickable]"
		}

		text := el.Text
		if text == "" {
			text = el.Placeholder
		}
		if text == "" {
			text = el.Attributes["aria-label"]
		}
		if len(text) > 60 {
			text = text[:57] + "..."
		}

		fmt.Printf("  %s %-12s  %-10s  %s%s\n",
			visibility,
			tagLabel(el.Tag, el.Type),
			el.Role,
			text,
			clickable,
		)
	}
	fmt.Println()
	return nil
}

func tagLabel(tag, semType string) string {
	if semType != "" && semType != "element" {
		return fmt.Sprintf("<%s>", semType)
	}
	return fmt.Sprintf("<%s>", tag)
}
