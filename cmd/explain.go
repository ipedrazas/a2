package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ipedrazas/a2/pkg/checks"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain CHECK_ID",
	Short: "Show detailed explanation of a check",
	Long: `Display comprehensive information about what a check does, why it matters,
and how to fix issues when they arise.

Example:
  a2 explain go:race
  a2 explain common:health`,
	Args: cobra.ExactArgs(1),
	Run:  runExplain,
}

func init() {
	rootCmd.AddCommand(explainCmd)
}

func runExplain(cmd *cobra.Command, args []string) {
	checkID := args[0]
	cfg := config.DefaultConfig()
	allRegs := checks.GetAllCheckRegistrations(cfg)

	for _, reg := range allRegs {
		if reg.Meta.ID == checkID {
			fmt.Printf("Check ID:     %s\n", reg.Meta.ID)
			fmt.Printf("Name:         %s\n", reg.Meta.Name)

			if reg.Meta.Description != "" {
				fmt.Printf("Description:  %s\n", reg.Meta.Description)
			}

			// Format languages
			langs := make([]string, len(reg.Meta.Languages))
			for i, l := range reg.Meta.Languages {
				langs[i] = string(l)
			}
			fmt.Printf("Languages:    %s\n", strings.Join(langs, ", "))

			if reg.Meta.Critical {
				fmt.Printf("Critical:     Yes (failure stops execution)\n")
			} else {
				fmt.Printf("Critical:     No\n")
			}

			if reg.Meta.Suggestion != "" {
				fmt.Printf("Suggestion:   %s\n", reg.Meta.Suggestion)
			}

			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown check ID: %s\n", checkID)
	fmt.Fprintf(os.Stderr, "Use 'a2 list checks' to see all available check IDs.\n")
	os.Exit(1)
}
