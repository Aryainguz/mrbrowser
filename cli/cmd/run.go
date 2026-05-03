package cmd

import (
	"fmt"
	"os"

	"github.com/mrbrowser/mrbrowser/core/browser"
	"github.com/mrbrowser/mrbrowser/core/runtime"
	"github.com/mrbrowser/mrbrowser/telemetry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:   "run <workflow.yaml>",
	Short: "Execute a YAML workflow",
	Long: `Execute a YAML task workflow.

Example:
  mr-browser run examples/login.yaml
  mr-browser run workflow.yaml --headless=false`,
	Args: cobra.ExactArgs(1),
	RunE: runWorkflow,
}

func runWorkflow(cmd *cobra.Command, args []string) error {
	log := telemetry.New("run")
	path := args[0]

	task, err := runtime.LoadFile(path)
	if err != nil {
		return fmt.Errorf("load task: %w", err)
	}

	log.Info("Loaded task",
		telemetry.F("name", task.Name),
		telemetry.F("steps", len(task.Steps)),
	)

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
		StopOnError: true,
		DBPath:      viper.GetString("db_path"),
	})
	if err != nil {
		return fmt.Errorf("create executor: %w", err)
	}
	defer func() { _ = executor.Close() }()

	result, err := executor.Run(task)
	if err != nil {
		return fmt.Errorf("run task: %w", err)
	}

	// Print summary
	fmt.Println()
	fmt.Printf("  Task: %s\n", result.TaskName)
	fmt.Printf("  Steps: %d/%d passed\n", countPassed(result.StepResults), len(result.StepResults))
	fmt.Printf("  Duration: %s\n", result.Duration.Round(1e6))

	if !result.Success {
		fmt.Printf("  Status: FAILED — %s\n", result.Error)
		os.Exit(1)
	}
	fmt.Printf("  Status: SUCCESS\n")
	return nil
}

func countPassed(steps []runtime.StepResult) int {
	n := 0
	for _, s := range steps {
		if s.Success {
			n++
		}
	}
	return n
}
