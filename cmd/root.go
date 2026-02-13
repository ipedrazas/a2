package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checks"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/language"
	"github.com/ipedrazas/a2/pkg/output"
	"github.com/ipedrazas/a2/pkg/profiles"
	"github.com/ipedrazas/a2/pkg/runner"
	"github.com/ipedrazas/a2/pkg/targets"
	"github.com/ipedrazas/a2/pkg/userconfig"
	"github.com/ipedrazas/a2/pkg/version"
	"github.com/spf13/cobra"
)

var (
	format        string
	outputFormat  string        // Alias for format (e.g. --output=toon)
	languages     []string      // Explicit language selection
	skippedChecks []string      // Checks to skip via CLI
	profile       string        // Application type profile (cli, api, library, desktop)
	target        string        // Maturity target (poc, production)
	timeout       time.Duration // Timeout for each individual check
	verbosity     int           // Verbosity level (0=normal, 1=failures, 2=all)
	failFast      bool          // Cancel remaining checks on first critical failure (parallel mode)
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

var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "List available application profiles",
	Long:  `List all built-in and user-defined profiles that can be used with the --profile flag.`,
	Run:   runProfiles,
}

var profilesInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize user profiles directory with built-in profiles",
	Long:  `Creates the user profiles directory (~/.config/a2/profiles) and writes built-in profiles as editable YAML files.`,
	RunE:  runProfilesInit,
}

var targetsCmd = &cobra.Command{
	Use:   "targets",
	Short: "List available maturity targets",
	Long:  `List all built-in and user-defined targets that can be used with the --target flag.`,
	Run:   runTargets,
}

var targetsInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize user targets directory with built-in targets",
	Long:  `Creates the user targets directory (~/.config/a2/targets) and writes built-in targets as editable YAML files.`,
	RunE:  runTargetsInit,
}

var profilesValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate user-defined profiles",
	Long: `Validate all user-defined profiles for correctness.

Checks for:
- Unknown check IDs (typos or invalid references)
- Duplicate disabled check IDs
- Profiles that override built-in profiles`,
	Run: runProfilesValidate,
}

var targetsValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate user-defined targets",
	Long: `Validate all user-defined targets for correctness.

Checks for:
- Unknown check IDs (typos or invalid references)
- Duplicate disabled check IDs
- Targets that override built-in targets`,
	Run: runTargetsValidate,
}

func init() {
	rootCmd.Version = version.Version
	checkCmd.Flags().StringVarP(&format, "format", "f", "pretty", "Output format: pretty, json, or toon")
	checkCmd.Flags().StringVar(&outputFormat, "output", "", "Output format (alias for --format): pretty, json, or toon")
	checkCmd.Flags().StringSliceVarP(&languages, "lang", "l", nil, "Languages to check (go, python). Auto-detects if not specified.")
	checkCmd.Flags().StringSliceVar(&skippedChecks, "skip", nil, "Checks to skip (e.g., --skip=license,k8s)")
	checkCmd.Flags().StringVar(&profile, "profile", "", "Application profile (cli, api, library, desktop)")
	checkCmd.Flags().StringVar(&target, "target", "", "Maturity target (poc, production)")
	checkCmd.Flags().DurationVar(&timeout, "timeout", 0, "Timeout for each individual check (e.g., 30s, 1m). 0 means no timeout")
	checkCmd.Flags().CountVarP(&verbosity, "verbose", "v", "Increase verbosity (-v for failures, -vv for all)")
	checkCmd.Flags().BoolVar(&failFast, "fail-fast", false, "Cancel remaining checks on first critical failure (parallel mode only)")
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(profilesCmd)
	rootCmd.AddCommand(targetsCmd)
	rootCmd.AddCommand(addCmd)

	// Add init subcommands
	profilesCmd.AddCommand(profilesInitCmd)
	targetsCmd.AddCommand(targetsInitCmd)

	// Add validate subcommands
	profilesCmd.AddCommand(profilesValidateCmd)
	targetsCmd.AddCommand(targetsValidateCmd)
}

func Execute() {
	// Initialize user profiles and targets
	// Errors are non-fatal; we'll just use built-in definitions
	_ = profiles.Init()
	_ = targets.Init()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runCheck(cmd *cobra.Command, args []string) error {
	if outputFormat != "" {
		format = outputFormat
	}
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
	baseDisabled := append([]string{}, cfg.Checks.Disabled...)

	// Apply target if specified (maturity level)
	var targetDisabled []string
	if target != "" {
		t, ok := targets.Get(target)
		if !ok {
			return fmt.Errorf("unknown target: %s (available: %s)", target, strings.Join(targets.Names(), ", "))
		}
		targetDisabled = append(targetDisabled, t.Disabled...)
		cfg.Checks.Disabled = append(cfg.Checks.Disabled, targetDisabled...)
	}

	// Apply profile if specified (application type)
	var profileDisabled []string
	if profile != "" {
		p, ok := profiles.Get(profile)
		if !ok {
			return fmt.Errorf("unknown profile: %s (available: %s)", profile, strings.Join(profiles.Names(), ", "))
		}
		profileDisabled = append(profileDisabled, p.Disabled...)
		cfg.Checks.Disabled = append(cfg.Checks.Disabled, profileDisabled...)
	}

	// Apply CLI skip flags
	if len(skippedChecks) > 0 {
		cfg.Checks.Disabled = append(cfg.Checks.Disabled, skippedChecks...)
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
		// Auto-detect languages, checking configured source directories
		detected = language.DetectWithSourceDirs(path, cfg.GetSourceDirs())
	}

	// Exit with error if no language detected
	if len(detected.Languages) == 0 {
		return fmt.Errorf("no supported language detected. Supported languages: go, python, node, java, rust, typescript, swift. Use --lang to specify explicitly")
	}

	// Get the list of checks to run
	registrations := checks.GetChecks(cfg, detected)
	skipped := buildSkippedChecks(cfg, detected, registrations, baseDisabled, targetDisabled, profileDisabled, skippedChecks, target, profile)

	// Set up progress callback for Pretty format only
	var progress *output.ProgressReporter
	var progressFunc runner.ProgressFunc

	if format == "pretty" && len(registrations) > 0 {
		progress = output.NewProgressReporter()
		progressFunc = progress.Update
	}

	// Run the suite with configured execution options
	opts := runner.RunSuiteOptions{
		Parallel:   cfg.Execution.Parallel,
		FailFast:   failFast,
		Timeout:    timeout,
		OnProgress: progressFunc,
	}
	result := runner.RunSuiteWithOptions(path, registrations, opts)

	// Clear progress display before showing results
	if progress != nil {
		progress.Done()
	}

	// Output results and handle exit code
	var success bool
	var outputErr error
	verbosityLevel := output.VerbosityLevel(verbosity)

	switch format {
	case "json":
		success, outputErr = output.JSON(result, detected, verbosityLevel)
	case "toon":
		success, outputErr = output.TOON(result, detected, verbosityLevel)
	default:
		success, outputErr = output.Pretty(result, path, detected, verbosityLevel, skipped)
	}

	if outputErr != nil {
		return outputErr
	}

	// Exit with code 1 if checks failed
	if !success {
		os.Exit(1)
	}

	return nil
}

type disabledSource struct {
	pattern string
	reason  string
}

func buildSkippedChecks(cfg *config.Config, detected language.DetectionResult, enabled []checker.CheckRegistration, baseDisabled, targetDisabled, profileDisabled, cliDisabled []string, targetName, profileName string) []output.SkipInfo {
	enabledIDs := make(map[string]struct{}, len(enabled))
	for _, reg := range enabled {
		enabledIDs[reg.Meta.ID] = struct{}{}
	}

	allRegs := checks.GetAllCheckRegistrations(cfg)
	var sources []disabledSource
	for _, p := range cliDisabled {
		sources = append(sources, disabledSource{pattern: p, reason: "cli --skip"})
	}
	if profileName != "" {
		for _, p := range profileDisabled {
			sources = append(sources, disabledSource{pattern: p, reason: "profile: " + profileName})
		}
	}
	if targetName != "" {
		for _, p := range targetDisabled {
			sources = append(sources, disabledSource{pattern: p, reason: "target: " + targetName})
		}
	}
	for _, p := range baseDisabled {
		sources = append(sources, disabledSource{pattern: p, reason: "config"})
	}

	var skipped []output.SkipInfo
	for _, reg := range allRegs {
		if !isApplicableForDetected(reg.Meta.Languages, detected.Languages) {
			continue
		}
		if _, ok := enabledIDs[reg.Meta.ID]; ok {
			continue
		}
		if !cfg.IsCheckDisabled(reg.Meta.ID) {
			continue
		}
		reason := "disabled"
		pattern := ""
		for _, src := range sources {
			if config.MatchDisabled(reg.Meta.ID, src.pattern) {
				reason = src.reason
				pattern = src.pattern
				break
			}
		}
		skipped = append(skipped, output.SkipInfo{
			ID:      reg.Meta.ID,
			Name:    reg.Meta.Name,
			Reason:  reason,
			Pattern: pattern,
		})
	}

	return skipped
}

func isApplicableForDetected(checkLangs []checker.Language, detected []checker.Language) bool {
	if len(checkLangs) == 0 {
		return true
	}
	for _, lang := range checkLangs {
		if lang == checker.LangCommon {
			return true
		}
		for _, d := range detected {
			if lang == d {
				return true
			}
		}
	}
	return false
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("Version:   %s\n", version.Version)
	fmt.Printf("Git SHA:   %s\n", version.GitSHA)
	fmt.Printf("Build Date: %s\n", version.BuildDate)
}

func runProfiles(cmd *cobra.Command, args []string) {
	fmt.Println("Application Profiles:")
	fmt.Println()
	fmt.Println("Profiles define which checks are relevant for your application type.")
	if dir, err := userconfig.GetSubDir("profiles"); err == nil {
		fmt.Printf("User profiles are stored in %s and override built-in profiles.\n", dir)
	} else {
		fmt.Println("User profiles override built-in profiles.")
	}
	fmt.Println()
	for _, p := range profiles.List() {
		fmt.Printf("  %s (%s)\n", p.Name, p.Source)
		fmt.Printf("    %s\n", p.Description)
		if len(p.Disabled) > 0 {
			fmt.Printf("    Skips %d checks\n", len(p.Disabled))
		}
		fmt.Println()
	}
	fmt.Println("Usage: a2 check --profile=<name>")
	fmt.Println("       a2 profiles init  # Initialize user profiles directory")
}

func runProfilesInit(cmd *cobra.Command, args []string) error {
	fmt.Println("Initializing user profiles directory...")
	return profiles.WriteBuiltInProfiles()
}

func runTargets(cmd *cobra.Command, args []string) {
	fmt.Println("Maturity Targets:")
	fmt.Println()
	fmt.Println("Targets control the strictness level of checks for your project stage.")
	if dir, err := userconfig.GetSubDir("targets"); err == nil {
		fmt.Printf("User targets are stored in %s and override built-in targets.\n", dir)
	} else {
		fmt.Println("User targets override built-in targets.")
	}
	fmt.Println()
	for _, t := range targets.List() {
		fmt.Printf("  %s (%s)\n", t.Name, t.Source)
		fmt.Printf("    %s\n", t.Description)
		if len(t.Disabled) > 0 {
			fmt.Printf("    Skips %d checks\n", len(t.Disabled))
		}
		fmt.Println()
	}
	fmt.Println("Usage: a2 check --target=<name>")
	fmt.Println("       a2 targets init  # Initialize user targets directory")
}

func runTargetsInit(cmd *cobra.Command, args []string) error {
	fmt.Println("Initializing user targets directory...")
	return targets.WriteBuiltInTargets()
}

func runProfilesValidate(cmd *cobra.Command, args []string) {
	results := profiles.ValidateAllUserProfiles()

	if len(results) == 0 {
		fmt.Println("No user profiles found to validate.")
		fmt.Println()
		fmt.Println("User profiles are stored in ~/.config/a2/profiles/")
		fmt.Println("Run 'a2 profiles init' to create sample profiles.")
		return
	}

	hasErrors := false
	for name, result := range results {
		if name == "_error" {
			fmt.Printf("Error: %s\n", result.Errors[0])
			os.Exit(1)
		}

		fmt.Printf("Profile: %s\n", name)
		if result.Valid {
			fmt.Printf("  Status: VALID\n")
		} else {
			fmt.Printf("  Status: INVALID\n")
			hasErrors = true
		}
		for _, err := range result.Errors {
			fmt.Printf("  ERROR: %s\n", err)
		}
		for _, warn := range result.Warnings {
			fmt.Printf("  WARNING: %s\n", warn)
		}
		fmt.Println()
	}

	if hasErrors {
		os.Exit(1)
	}
}

func runTargetsValidate(cmd *cobra.Command, args []string) {
	results := targets.ValidateAllUserTargets()

	if len(results) == 0 {
		fmt.Println("No user targets found to validate.")
		fmt.Println()
		fmt.Println("User targets are stored in ~/.config/a2/targets/")
		fmt.Println("Run 'a2 targets init' to create sample targets.")
		return
	}

	hasErrors := false
	for name, result := range results {
		if name == "_error" {
			fmt.Printf("Error: %s\n", result.Errors[0])
			os.Exit(1)
		}

		fmt.Printf("Target: %s\n", name)
		if result.Valid {
			fmt.Printf("  Status: VALID\n")
		} else {
			fmt.Printf("  Status: INVALID\n")
			hasErrors = true
		}
		for _, err := range result.Errors {
			fmt.Printf("  ERROR: %s\n", err)
		}
		for _, warn := range result.Warnings {
			fmt.Printf("  WARNING: %s\n", warn)
		}
		fmt.Println()
	}

	if hasErrors {
		os.Exit(1)
	}
}
