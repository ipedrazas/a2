package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/language"
	"github.com/ipedrazas/a2/pkg/tools"
	"github.com/spf13/cobra"
)

var (
	doctorFormat string
	doctorAll    bool
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system for required tools",
	Long: `Check if external tools required by a2 checks are installed.

Shows which tools are installed and provides install commands for missing ones.
By default, only shows tools relevant to detected languages in the current project.

Use --all to show tools for all supported languages.`,
	Run: runDoctor,
}

func init() {
	doctorCmd.Flags().StringVarP(&doctorFormat, "format", "f", "pretty", "Output format: pretty, json, or toon")
	doctorCmd.Flags().BoolVar(&doctorAll, "all", false, "Show tools for all languages, not just detected ones")
	rootCmd.AddCommand(doctorCmd)
}

// DoctorResult represents the doctor command output.
type DoctorResult struct {
	Environment tools.Environment `json:"environment"`
	Tools       []ToolResult      `json:"tools"`
	Summary     DoctorSummary     `json:"summary"`
}

// ToolResult represents a single tool's status.
type ToolResult struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Installed   bool   `json:"installed"`
	Version     string `json:"version,omitempty"`
	Install     string `json:"install,omitempty"`
}

// DoctorSummary summarizes the doctor results.
type DoctorSummary struct {
	Total     int `json:"total"`
	Installed int `json:"installed"`
	Missing   int `json:"missing"`
}

func runDoctor(cmd *cobra.Command, args []string) {
	// Detect environment
	env := tools.DetectEnvironment()

	// Determine which languages to check
	var langs []checker.Language
	if doctorAll {
		langs = checker.AllLanguages()
	} else {
		// Load config and detect languages
		cfg, _ := config.Load(".")
		if cfg == nil {
			cfg = config.DefaultConfig()
		}
		detected := language.DetectWithSourceDirs(".", cfg.GetSourceDirs())
		langs = detected.Languages
	}

	// Get tools for these languages
	toolList := tools.ForLanguages(langs)

	// Check installation status
	statuses := tools.CheckAllInstalled(toolList)

	// Build result
	result := DoctorResult{
		Environment: env,
	}

	for _, status := range statuses {
		tr := ToolResult{
			Name:        status.Tool.Name,
			Description: status.Tool.Description,
			Language:    string(status.Tool.Language),
			Installed:   status.Installed,
			Version:     status.Version,
		}
		if !status.Installed {
			tr.Install = tools.GetInstallCommand(status.Tool, env)
		}
		result.Tools = append(result.Tools, tr)

		result.Summary.Total++
		if status.Installed {
			result.Summary.Installed++
		} else {
			result.Summary.Missing++
		}
	}

	// Output based on format
	switch doctorFormat {
	case "json":
		outputDoctorJSON(result)
	case "toon":
		outputDoctorTOON(result, env, langs)
	default:
		outputDoctorPretty(result, env, langs)
	}
}

func outputDoctorPretty(result DoctorResult, env tools.Environment, langs []checker.Language) {
	fmt.Printf("Environment: %s/%s\n", env.OSName(), env.Arch)
	fmt.Printf("Package managers: %s\n", strings.Join(env.PackageManagers, ", "))
	fmt.Println()

	if len(langs) == 0 {
		fmt.Println("No languages detected. Use --all to show all tools.")
		return
	}

	// Group by language
	byLang := make(map[string][]ToolResult)
	for _, t := range result.Tools {
		byLang[t.Language] = append(byLang[t.Language], t)
	}

	// Display order
	langOrder := []string{"go", "python", "node", "typescript", "java", "rust", "swift", "common"}

	for _, lang := range langOrder {
		toolResults, ok := byLang[lang]
		if !ok || len(toolResults) == 0 {
			continue
		}

		fmt.Printf("%s:\n", formatLangHeader(lang))
		for _, t := range toolResults {
			if t.Installed {
				version := ""
				if t.Version != "" {
					// Truncate long version strings
					v := t.Version
					if len(v) > 30 {
						v = v[:30] + "..."
					}
					version = fmt.Sprintf(" (%s)", v)
				}
				fmt.Printf("  \033[32m✓\033[0m %-20s %s%s\n", t.Name, t.Description, version)
			} else {
				fmt.Printf("  \033[31m✗\033[0m %-20s %s\n", t.Name, t.Description)
				if t.Install != "" {
					fmt.Printf("      \033[33m→\033[0m %s\n", t.Install)
				}
			}
		}
		fmt.Println()
	}

	// Summary
	if result.Summary.Missing > 0 {
		fmt.Printf("Missing: %d tool(s)\n", result.Summary.Missing)
	} else {
		fmt.Printf("\033[32mAll %d tools installed!\033[0m\n", result.Summary.Total)
	}
}

func outputDoctorJSON(result DoctorResult) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(result)
}

func outputDoctorTOON(result DoctorResult, env tools.Environment, _ []checker.Language) {
	fmt.Printf("env %s/%s\n", env.OS, env.Arch)
	fmt.Printf("managers %s\n", strings.Join(env.PackageManagers, ","))

	if result.Summary.Missing > 0 {
		fmt.Println("missing")
		for _, t := range result.Tools {
			if !t.Installed {
				fmt.Printf("  %s %s\n", t.Name, t.Install)
			}
		}
	}

	fmt.Printf("summary %d/%d installed\n", result.Summary.Installed, result.Summary.Total)
}

func formatLangHeader(lang string) string {
	headers := map[string]string{
		"go":         "Go",
		"python":     "Python",
		"node":       "Node.js",
		"typescript": "TypeScript",
		"java":       "Java",
		"rust":       "Rust",
		"swift":      "Swift",
		"common":     "Common",
	}
	if h, ok := headers[lang]; ok {
		return h
	}
	return lang
}
