package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checks"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/spf13/cobra"
)

var runFormat string

var runCmd = &cobra.Command{
	Use:   "run CHECK_ID [path]",
	Short: "Run a specific check with full output",
	Long: `Run a single check and display the complete stdout/stderr from the underlying tool.

This is useful for debugging or understanding why a check failed.

Example:
  a2 run go:race
  a2 run go:build ./path/to/project
  a2 run common:secrets`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runSingleCheck,
}

func init() {
	runCmd.Flags().StringVarP(&runFormat, "format", "f", "pretty", "Output format: pretty or json")
	rootCmd.AddCommand(runCmd)
}

func runSingleCheck(cmd *cobra.Command, args []string) error {
	checkID := args[0]
	path := "."
	if len(args) > 1 {
		path = args[1]
	}

	// Load configuration
	cfg, err := config.Load(path)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Get all check registrations
	allRegs := checks.GetAllCheckRegistrations(cfg)

	// Find the check by ID
	var found *checker.CheckRegistration
	for i := range allRegs {
		if allRegs[i].Meta.ID == checkID {
			found = &allRegs[i]
			break
		}
	}

	if found == nil {
		fmt.Fprintf(os.Stderr, "Unknown check ID: %s\n", checkID)
		fmt.Fprintf(os.Stderr, "Use 'a2 list checks' to see all available check IDs.\n")
		return fmt.Errorf("unknown check ID: %s", checkID)
	}

	// Run the check
	start := time.Now()
	result, err := found.Checker.Run(path)
	duration := time.Since(start)
	result.Duration = duration

	if err != nil {
		return fmt.Errorf("check execution error: %w", err)
	}

	// Output results based on format
	if runFormat == "json" {
		outputRunResultJSON(result)
	} else {
		outputRunResultPretty(result, found.Meta)
	}

	// Exit with error if check failed
	if result.Status == checker.Fail {
		os.Exit(1)
	}

	return nil
}

func outputRunResultPretty(result checker.Result, meta checker.CheckMeta) {
	// Status symbol
	var symbol string
	switch result.Status {
	case checker.Pass:
		symbol = "\u2713" // checkmark
	case checker.Warn:
		symbol = "!"
	case checker.Fail:
		symbol = "\u2717" // X
	case checker.Info:
		symbol = "\u2139" // info
	}

	// Status label
	statusLabel := result.Status.String()

	// Print header
	fmt.Printf("%s %s %s (%.1fs) - %s\n", symbol, statusLabel, result.Name, result.Duration.Seconds(), result.ID)
	if result.Message != "" {
		fmt.Printf("    %s\n", result.Message)
	}
	if result.Reason != "" {
		if result.Reason != result.Message {
			if result.Message != "" {
				fmt.Printf("    Reason: %s\n", result.Reason)
			} else {
				fmt.Printf("    %s\n", result.Reason)
			}
		}
	}

	// Print suggestion if available and check didn't pass
	if meta.Suggestion != "" && result.Status != checker.Pass {
		fmt.Printf("\nSuggestion: %s\n", meta.Suggestion)
	}

	// Print raw output if available
	if result.RawOutput != "" {
		fmt.Printf("\n--- Output ---\n")
		fmt.Println(result.RawOutput)
	}
}

func outputRunResultJSON(result checker.Result) {
	fmt.Printf("{\n")
	fmt.Printf("  \"id\": %q,\n", result.ID)
	fmt.Printf("  \"name\": %q,\n", result.Name)
	fmt.Printf("  \"status\": %q,\n", result.Status.String())
	fmt.Printf("  \"passed\": %t,\n", result.Passed)
	fmt.Printf("  \"message\": %q,\n", result.Message)
	fmt.Printf("  \"reason\": %q,\n", result.Reason)
	fmt.Printf("  \"language\": %q,\n", result.Language)
	fmt.Printf("  \"duration_ms\": %d", result.Duration.Milliseconds())
	if result.RawOutput != "" {
		fmt.Printf(",\n  \"raw_output\": %q", result.RawOutput)
	}
	fmt.Printf("\n}\n")
}
