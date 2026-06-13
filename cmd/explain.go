package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checks"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/spf13/cobra"
)

// docsBaseURL is the canonical location of the checks reference docs.
const docsBaseURL = "https://github.com/ipedrazas/a2/blob/main/docs"

// languagesWithDocs are the languages that have a dedicated docs/checks/<lang>.md page.
var languagesWithDocs = map[string]bool{
	"go": true, "python": true, "node": true,
	"java": true, "rust": true, "typescript": true,
}

// docReference returns a URL pointing at the documentation for a check ID.
// Language checks (e.g. "go:race") link to docs/checks/<lang>.md#<anchor>;
// everything else links to the top-level docs/CHECKS.md.
func docReference(id string) string {
	lang, _, ok := strings.Cut(id, ":")
	if ok && languagesWithDocs[lang] {
		// GitHub heading anchors lowercase the text and drop the colon:
		// "## go:race" -> "#gorace".
		anchor := strings.ReplaceAll(id, ":", "")
		return fmt.Sprintf("%s/checks/%s.md#%s", docsBaseURL, lang, anchor)
	}
	return docsBaseURL + "/CHECKS.md"
}

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

			if reg.Meta.Speed == checker.SpeedSlow {
				fmt.Printf("Speed:        Slow (skipped by 'a2 check --quick')\n")
			} else {
				fmt.Printf("Speed:        Fast (runs in 'a2 check --quick')\n")
			}

			fmt.Printf("Docs:         %s\n", docReference(reg.Meta.ID))

			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown check ID: %s\n", checkID)
	fmt.Fprintf(os.Stderr, "Use 'a2 list checks' to see all available check IDs.\n")
	os.Exit(1)
}
