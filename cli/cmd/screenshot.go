package cmd

import (
	"fmt"
	"os"

	"github.com/mrbrowser/mrbrowser/core/browser"
	"github.com/mrbrowser/mrbrowser/telemetry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var screenshotOutput string

var screenshotCmd = &cobra.Command{
	Use:   "screenshot <url>",
	Short: "Capture a screenshot of a URL",
	Long: `Navigate to a URL and capture a PNG screenshot.

Example:
  mr-browser screenshot https://example.com
  mr-browser screenshot https://example.com --output page.png`,
	Args: cobra.ExactArgs(1),
	RunE: captureScreenshot,
}

func init() {
	screenshotCmd.Flags().StringVarP(&screenshotOutput, "output", "o", "", "output file path (default: screenshot_<timestamp>.png)")
}

func captureScreenshot(cmd *cobra.Command, args []string) error {
	log := telemetry.New("screenshot")
	url := args[0]

	output := screenshotOutput
	if output == "" {
		output = fmt.Sprintf("screenshot_%d.png", timeNow())
	}

	opts := browser.Options{
		Headless:     viper.GetBool("headless"),
		ChromiumPath: viper.GetString("chromium_path"),
		NoSandbox:    viper.GetBool("no_sandbox"),
	}

	session, err := browser.NewSession("screenshot", opts)
	if err != nil {
		return fmt.Errorf("launch browser: %w", err)
	}
	defer func() { _ = session.Close() }()

	log.Step("Navigating", telemetry.F("url", url))
	if err := session.Page.Navigate(url); err != nil {
		return fmt.Errorf("navigate: %w", err)
	}

	log.Step("Capturing screenshot")
	buf, err := session.Page.Screenshot()
	if err != nil {
		return fmt.Errorf("screenshot: %w", err)
	}

	if err := os.WriteFile(output, buf, 0644); err != nil {
		return fmt.Errorf("save screenshot: %w", err)
	}

	log.Success("Screenshot saved", telemetry.F("path", output), telemetry.F("bytes", len(buf)))
	fmt.Printf("✓ Screenshot saved: %s (%d bytes)\n", output, len(buf))
	return nil
}

func timeNow() int64 {
	return int64(^uint64(0) >> 1) // placeholder, replaced below
}
