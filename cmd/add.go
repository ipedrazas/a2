package cmd

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/language"
	"github.com/ipedrazas/a2/pkg/profiles"
	"github.com/ipedrazas/a2/pkg/safepath"
	"github.com/ipedrazas/a2/pkg/targets"
	"github.com/spf13/cobra"
)

var (
	addInteractive bool
	addProfile     string
	addTarget      string
	addLanguages   []string
	addFiles       []string
	addCoverage    float64
	addOutput      string
	addForce       bool
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Generate a .a2.yaml configuration file",
	Long: `Generate a .a2.yaml configuration file for your project.

Two modes are available:
  - Interactive mode (-i): Guided prompts to configure all options
  - Non-interactive mode: Pass flags directly to generate the config

Examples:
  a2 add -i                                    # Interactive mode
  a2 add --profile cli --target poc            # CLI app for PoC
  a2 add --lang go,python --coverage 90        # Multi-language with 90% coverage`,
	RunE: runAdd,
}

func init() {
	addCmd.Flags().BoolVarP(&addInteractive, "interactive", "i", false, "Run in interactive mode")
	addCmd.Flags().StringVar(&addProfile, "profile", "", "Application profile (cli, api, library, desktop)")
	addCmd.Flags().StringVar(&addTarget, "target", "", "Maturity target (poc, production)")
	addCmd.Flags().StringSliceVar(&addLanguages, "lang", nil, "Languages (go, python, node, java, rust, typescript)")
	addCmd.Flags().StringSliceVar(&addFiles, "files", nil, "Required files (default: README.md,LICENSE)")
	addCmd.Flags().Float64Var(&addCoverage, "coverage", 0, "Coverage threshold (default: 80)")
	addCmd.Flags().StringVarP(&addOutput, "output", "o", ".a2.yaml", "Output file path")
	addCmd.Flags().BoolVarP(&addForce, "force", "f", false, "Overwrite existing file")
}

func runAdd(cmd *cobra.Command, args []string) error {
	// Check if output file already exists
	if !addForce && safepath.Exists(".", addOutput) {
		return fmt.Errorf("file %s already exists. Use --force to overwrite", addOutput)
	}

	var opts config.GeneratorOptions

	if addInteractive {
		var err error
		opts, err = runInteractiveMode()
		if err != nil {
			return err
		}
	} else {
		opts = config.GeneratorOptions{
			Profile:   addProfile,
			Target:    addTarget,
			Languages: addLanguages,
			Files:     addFiles,
			Coverage:  addCoverage,
		}
	}

	// Validate profile if set
	if opts.Profile != "" {
		if _, ok := profiles.Get(opts.Profile); !ok {
			return fmt.Errorf("unknown profile: %s (available: %s)", opts.Profile, strings.Join(profiles.Names(), ", "))
		}
	}

	// Validate target if set
	if opts.Target != "" {
		if _, ok := targets.Get(opts.Target); !ok {
			return fmt.Errorf("unknown target: %s (available: %s)", opts.Target, strings.Join(targets.Names(), ", "))
		}
	}

	// Generate config
	cfg := config.Generate(opts)

	// Convert to YAML
	yamlData, err := config.ToYAML(cfg, opts)
	if err != nil {
		return fmt.Errorf("failed to generate YAML: %w", err)
	}

	// Write the file
	if err := os.WriteFile(addOutput, yamlData, 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Created %s\n", addOutput)
	if opts.Profile != "" || opts.Target != "" {
		fmt.Print("Run checks with: a2 check")
		if opts.Profile != "" {
			fmt.Printf(" --profile %s", opts.Profile)
		}
		if opts.Target != "" {
			fmt.Printf(" --target %s", opts.Target)
		}
		fmt.Println()
	} else {
		fmt.Println("Run checks with: a2 check")
	}

	return nil
}

func runInteractiveMode() (config.GeneratorOptions, error) {
	var opts config.GeneratorOptions

	// Step 1: Profile selection
	profileOptions := []huh.Option[string]{
		huh.NewOption("None (use defaults)", ""),
	}
	for _, p := range profiles.List() {
		profileOptions = append(profileOptions, huh.NewOption(
			fmt.Sprintf("%s - %s", p.Name, p.Description),
			p.Name,
		))
	}

	// Step 2: Target selection
	targetOptions := []huh.Option[string]{
		huh.NewOption("None (use defaults)", ""),
	}
	for _, t := range targets.List() {
		targetOptions = append(targetOptions, huh.NewOption(
			fmt.Sprintf("%s - %s", t.Name, t.Description),
			t.Name,
		))
	}

	// Step 3: Language detection
	detected := language.Detect(".")
	var detectedLangs []string
	for _, lang := range detected.Languages {
		detectedLangs = append(detectedLangs, string(lang))
	}

	allLanguages := []string{"go", "python", "node", "java", "rust", "typescript"}
	languageOptions := make([]huh.Option[string], len(allLanguages))
	for i, lang := range allLanguages {
		label := lang
		if slices.Contains(detectedLangs, lang) {
			label = fmt.Sprintf("%s (detected)", lang)
		}
		languageOptions[i] = huh.NewOption(label, lang)
	}

	// If nothing detected, default to empty (auto-detect)
	selectedLangs := detectedLangs

	// Step 4: Required files
	fileOptions := []huh.Option[string]{
		huh.NewOption("README.md", "README.md"),
		huh.NewOption("LICENSE", "LICENSE"),
		huh.NewOption("CHANGELOG.md", "CHANGELOG.md"),
		huh.NewOption("CONTRIBUTING.md", "CONTRIBUTING.md"),
		huh.NewOption(".editorconfig", ".editorconfig"),
	}
	selectedFiles := []string{"README.md", "LICENSE"}

	// Step 5: Coverage threshold
	coverageStr := "80"

	// Build the form
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Application Profile").
				Description("Choose a profile that matches your application type").
				Options(profileOptions...).
				Value(&opts.Profile),

			huh.NewSelect[string]().
				Title("Maturity Target").
				Description("Choose the maturity level for your project").
				Options(targetOptions...).
				Value(&opts.Target),
		),

		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Languages").
				Description(buildLanguageDescription(detected)).
				Options(languageOptions...).
				Value(&selectedLangs),
		),

		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Required Files").
				Description("Select files that must exist in your project").
				Options(fileOptions...).
				Value(&selectedFiles),

			huh.NewInput().
				Title("Coverage Threshold").
				Description("Minimum code coverage percentage (0-100)").
				Value(&coverageStr).
				Validate(validateCoverage),
		),
	)

	if err := form.Run(); err != nil {
		return opts, err
	}

	// Parse coverage
	if cov, err := strconv.ParseFloat(coverageStr, 64); err == nil {
		opts.Coverage = cov
	}

	// Only set languages if explicitly selected (not auto-detect)
	if len(selectedLangs) > 0 {
		opts.Languages = selectedLangs
	}

	// Only set files if different from default
	if !stringSliceEqual(selectedFiles, []string{"README.md", "LICENSE"}) {
		opts.Files = selectedFiles
	}

	// Show preview and confirm
	cfg := config.Generate(opts)
	yamlData, err := config.ToYAML(cfg, opts)
	if err != nil {
		return opts, fmt.Errorf("failed to generate preview: %w", err)
	}

	fmt.Println("\n--- Preview ---")
	fmt.Println(string(yamlData))
	fmt.Println("---------------")

	var confirm bool
	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Create .a2.yaml with these settings?").
				Value(&confirm),
		),
	)

	if err := confirmForm.Run(); err != nil {
		return opts, err
	}

	if !confirm {
		return opts, fmt.Errorf("cancelled by user")
	}

	return opts, nil
}

func buildLanguageDescription(detected language.DetectionResult) string {
	if len(detected.Languages) == 0 {
		return "No languages detected. Select languages to check."
	}

	var parts []string
	for _, lang := range detected.Languages {
		indicators := detected.Indicators[lang]
		if len(indicators) > 0 {
			parts = append(parts, fmt.Sprintf("%s (%s)", lang, strings.Join(indicators, ", ")))
		}
	}
	return "Detected: " + strings.Join(parts, ", ")
}

func validateCoverage(s string) error {
	if s == "" {
		return nil
	}
	cov, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("must be a number")
	}
	if cov < 0 || cov > 100 {
		return fmt.Errorf("must be between 0 and 100")
	}
	return nil
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
