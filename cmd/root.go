package cmd

import (
	"fmt"
	"os"
	"strings"

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
	languages     []string // Explicit language selection
	skippedChecks []string // Checks to skip via CLI
	profile       string   // Application type profile (cli, api, library, desktop)
	target        string   // Maturity target (poc, production)
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

func init() {
	checkCmd.Flags().StringVarP(&format, "format", "f", "pretty", "Output format: pretty or json")
	checkCmd.Flags().StringSliceVarP(&languages, "lang", "l", nil, "Languages to check (go, python). Auto-detects if not specified.")
	checkCmd.Flags().StringSliceVar(&skippedChecks, "skip", nil, "Checks to skip (e.g., --skip=license,k8s)")
	checkCmd.Flags().StringVar(&profile, "profile", "", "Application profile (cli, api, library, desktop)")
	checkCmd.Flags().StringVar(&target, "target", "", "Maturity target (poc, production)")
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(profilesCmd)
	rootCmd.AddCommand(targetsCmd)
	rootCmd.AddCommand(addCmd)

	// Add init subcommands
	profilesCmd.AddCommand(profilesInitCmd)
	targetsCmd.AddCommand(targetsInitCmd)
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

	// Apply target if specified (maturity level)
	if target != "" {
		t, ok := targets.Get(target)
		if !ok {
			return fmt.Errorf("unknown target: %s (available: %s)", target, strings.Join(targets.Names(), ", "))
		}
		cfg.Checks.Disabled = append(cfg.Checks.Disabled, t.Disabled...)
	}

	// Apply profile if specified (application type)
	if profile != "" {
		p, ok := profiles.Get(profile)
		if !ok {
			return fmt.Errorf("unknown profile: %s (available: %s)", profile, strings.Join(profiles.Names(), ", "))
		}
		cfg.Checks.Disabled = append(cfg.Checks.Disabled, p.Disabled...)
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
