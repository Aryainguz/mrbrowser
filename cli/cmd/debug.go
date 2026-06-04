package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mrbrowser/mrbrowser/core/browser"
	"github.com/mrbrowser/mrbrowser/core/runtime"
	"github.com/mrbrowser/mrbrowser/telemetry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var debugCmd = &cobra.Command{
	Use:   "debug <workflow.yaml>",
	Short: "Step through a workflow with interactive prompts",
	Long: `Execute a workflow step by step, pausing after each step for user confirmation.

Example:
  mr-browser debug examples/login.yaml --headless=false`,
	Args: cobra.ExactArgs(1),
	RunE: debugWorkflow,
}

func debugWorkflow(cmd *cobra.Command, args []string) error {
	log := telemetry.New("debug")
	path := args[0]

	task, err := runtime.LoadFile(path)
	if err != nil {
		return fmt.Errorf("load task: %w", err)
	}

	log.Info("Debug mode", telemetry.F("task", task.Name), telemetry.F("steps", len(task.Steps)))
	fmt.Printf("\n  🐛 Debug: %s (%d steps)\n\n", task.Name, len(task.Steps))

	opts := browser.Options{
		Headless:     viper.GetBool("headless"),
		ChromiumPath: viper.GetString("chromium_path"),
		NoSandbox:    viper.GetBool("no_sandbox"),
	}

	session, err := browser.NewSession(task.Name, opts)
	if err != nil {
		return fmt.Errorf("launch browser: %w", err)
	}
	defer func() { _ = session.Close() }()

	executor, err := runtime.NewExecutor(session, runtime.ExecutorOptions{
		StopOnError: false,
		DBPath:      viper.GetString("db_path"),
	})
	if err != nil {
		return fmt.Errorf("create executor: %w", err)
	}
	defer func() { _ = executor.Close() }()

	scanner := bufio.NewScanner(os.Stdin)

	for i, step := range task.Steps {
		stepNum := i + 1
		fmt.Printf("  ── Step %d/%d [%s] ─────────────\n", stepNum, len(task.Steps), step.Kind())

		// Print step details
		printStepDetails(&step)

		fmt.Printf("  Press ENTER to execute, 's' to skip, 'q' to quit: ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		switch input {
		case "q", "quit":
			fmt.Println("  Aborted.")
			return nil
		case "s", "skip":
			fmt.Printf("  ⏭  Skipped step %d\n\n", stepNum)
			continue
		}

		// Run just this step
		singleTask := &runtime.Task{
			Name:  task.Name,
			Steps: []runtime.Step{step},
		}
		result, _ := executor.Run(singleTask)
		if result != nil && len(result.StepResults) > 0 {
			sr := result.StepResults[0]
			if sr.Success {
				fmt.Printf("  ✓ Step %d completed (%s)\n\n", stepNum, sr.Duration.Round(1e6))
			} else {
				fmt.Printf("  ✗ Step %d failed: %s\n\n", stepNum, sr.Error)
			}
		}
	}

	fmt.Println("  ✓ Debug session complete.")
	return nil
}

func printStepDetails(step *runtime.Step) {
	switch step.Kind() {
	case "open":
		fmt.Printf("    URL: %s\n", step.Open.URL)
	case "click":
		fmt.Printf("    Target: %q\n", step.Click.Target)
	case "type":
		fmt.Printf("    Target: %q  Value: %q\n", step.Type.Target, step.Type.Value)
	case "scroll":
		fmt.Printf("    Direction: %s  Pixels: %d\n", step.Scroll.Direction, step.Scroll.Pixels)
	case "screenshot":
		out := step.Screenshot.Output
		if out == "" {
			out = "(auto)"
		}
		fmt.Printf("    Output: %s\n", out)
	case "wait":
		if step.Wait.Selector != "" {
			fmt.Printf("    Wait for selector: %q\n", step.Wait.Selector)
		} else {
			fmt.Printf("    Wait: %.1fs\n", step.Wait.Seconds)
		}
	}
}
