package cmd

import (
	"fmt"
	"sort"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checks"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/spf13/cobra"
)

var listExplain bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available resources",
	Long:  `List available resources like checks, profiles, and targets.`,
}

var listChecksCmd = &cobra.Command{
	Use:     "checks",
	Aliases: []string{"c"},
	Short:   "List all available checks",
	Long: `List all available checks that a2 can run.

Shows check IDs (for use with --skip), names, and which language they apply to.
Checks marked as [critical] will cause the run to fail if they don't pass.`,
	Run: runListChecks,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listChecksCmd)
	listChecksCmd.Flags().BoolVar(&listExplain, "explain", false, "Show detailed descriptions for each check")
}

func runListChecks(cmd *cobra.Command, args []string) {
	cfg := config.DefaultConfig()
	allRegs := checks.GetAllCheckRegistrations(cfg)

	// Group checks by language
	byLanguage := make(map[checker.Language][]checker.CheckRegistration)
	for _, reg := range allRegs {
		for _, lang := range reg.Meta.Languages {
			byLanguage[lang] = append(byLanguage[lang], reg)
		}
	}

	// Define language display order
	languageOrder := []checker.Language{
		checker.LangGo,
		checker.LangPython,
		checker.LangNode,
		checker.LangTypeScript,
		checker.LangJava,
		checker.LangRust,
		checker.LangSwift,
		checker.LangCommon,
	}

	fmt.Println("Available Checks:")
	fmt.Println()
	fmt.Println("Use --skip=<id> to skip specific checks when running 'a2 check'.")
	fmt.Println()

	for _, lang := range languageOrder {
		regs, ok := byLanguage[lang]
		if !ok || len(regs) == 0 {
			continue
		}

		// Sort by order within the language
		sort.Slice(regs, func(i, j int) bool {
			return regs[i].Meta.Order < regs[j].Meta.Order
		})

		// Display language header
		langName := formatLanguageName(lang)
		fmt.Printf("%s:\n", langName)

		for _, reg := range regs {
			critical := ""
			if reg.Meta.Critical {
				critical = " [critical]"
			}
			fmt.Printf("  %-25s %s%s\n", reg.Meta.ID, reg.Meta.Name, critical)
			if listExplain && reg.Meta.Description != "" {
				fmt.Printf("    %s\n", reg.Meta.Description)
			}
		}
		fmt.Println()
	}

	// Count total checks
	total := len(allRegs)
	fmt.Printf("Total: %d checks\n", total)
	fmt.Println()
	fmt.Println("Example: a2 check --skip=common:k8s,common:health")
}

func formatLanguageName(lang checker.Language) string {
	names := map[checker.Language]string{
		checker.LangGo:         "Go",
		checker.LangPython:     "Python",
		checker.LangNode:       "Node.js",
		checker.LangTypeScript: "TypeScript",
		checker.LangJava:       "Java",
		checker.LangRust:       "Rust",
		checker.LangSwift:      "Swift",
		checker.LangCommon:     "Common (language-agnostic)",
	}
	if name, ok := names[lang]; ok {
		return name
	}
	// Fallback for unknown languages
	return string(lang)
}
