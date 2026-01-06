package cmd

import (
	"fmt"
	"os"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checks"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/language"
	"github.com/ipedrazas/a2/pkg/output"
	"github.com/ipedrazas/a2/pkg/runner"
	"github.com/ipedrazas/a2/pkg/version"
	"github.com/spf13/cobra"
)

var (
	format    string
	languages []string // Explicit language selection
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

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print version, git SHA, and build date information.`,
	Run:   runVersion,
}

func init() {
	checkCmd.Flags().StringVarP(&format, "format", "f", "pretty", "Output format: pretty or json")
	checkCmd.Flags().StringSliceVarP(&languages, "lang", "l", nil, "Languages to check (go, python). Auto-detects if not specified.")
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(versionCmd)
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

	// Detect or use explicit languages
	var detected language.DetectionResult
	if len(languages) > 0 {
		// Convert string flags to Language types
		langs := make([]checker.Language, len(languages))
		for i, l := range languages {
			langs[i] = checker.Language(l)
		}
		detected = language.DetectWithOverride(path, langs)
	} else if len(cfg.Language.Explicit) > 0 {
		// Use explicit languages from config
		langs := make([]checker.Language, len(cfg.Language.Explicit))
		for i, l := range cfg.Language.Explicit {
			langs[i] = checker.Language(l)
		}
		detected = language.DetectWithOverride(path, langs)
	} else {
		// Auto-detect languages
		detected = language.Detect(path)
	}

	// Fallback to Go if nothing detected (backward compatibility)
	if len(detected.Languages) == 0 {
		detected.Languages = []checker.Language{checker.LangGo}
		detected.Primary = checker.LangGo
	}

	// Get the list of checks to run
	checkList := checks.GetChecks(cfg, detected)

	// Run the suite with configured execution options
	opts := runner.RunSuiteOptions{Parallel: cfg.Execution.Parallel}
	result := runner.RunSuiteWithOptions(path, checkList, opts)

	// Output results
	switch format {
	case "json":
		return output.JSON(result, detected)
	default:
		return output.Pretty(result, path, detected)
	}
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("Version:   %s\n", version.Version)
	fmt.Printf("Git SHA:   %s\n", version.GitSHA)
	fmt.Printf("Build Date: %s\n", version.BuildDate)
}
