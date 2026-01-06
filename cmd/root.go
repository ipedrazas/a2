package cmd

import (
	"fmt"
	"os"

	"github.com/ipedrazas/a2/pkg/checks"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/output"
	"github.com/ipedrazas/a2/pkg/runner"
	"github.com/spf13/cobra"
)

var (
	format string
)

var rootCmd = &cobra.Command{
	Use:   "a2",
	Short: "A2 - Application Analysis tool",
	Long: `A2 is a code quality checker that runs a suite of checks
against your repository and provides a health score.

It can be used as a CLI tool, GitHub Action, or pre-commit hook.`,
}

var checkCmd = &cobra.Command{
	Use:   "check [path]",
	Short: "Run checks on a directory",
	Long:  `Run all configured checks against the specified directory (defaults to current directory).`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCheck,
}

func init() {
	checkCmd.Flags().StringVarP(&format, "format", "f", "pretty", "Output format: pretty or json")
	rootCmd.AddCommand(checkCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runCheck(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// Load configuration
	cfg, err := config.Load(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Get the list of checks to run
	checkList := checks.GetChecks(cfg)

	// Run the suite with configured execution options
	opts := runner.RunSuiteOptions{Parallel: cfg.Execution.Parallel}
	result := runner.RunSuiteWithOptions(path, checkList, opts)

	// Output results
	switch format {
	case "json":
		return output.JSON(result)
	default:
		return output.Pretty(result, path)
	}
}
