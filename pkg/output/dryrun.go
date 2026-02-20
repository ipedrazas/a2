package output

import (
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/language"
)

// DryRun prints the list of checks that would run without executing them.
func DryRun(registrations []checker.CheckRegistration, detected language.DetectionResult, skipped []SkipInfo) {
	langs := make([]string, len(detected.Languages))
	for i, l := range detected.Languages {
		langs[i] = string(l)
	}

	fmt.Printf("Dry run: %d checks would run", len(registrations))
	if len(skipped) > 0 {
		fmt.Printf(" (%d skipped)", len(skipped))
	}
	fmt.Println()
	fmt.Println()
	fmt.Printf("Languages: %s\n", strings.Join(langs, ", "))
	fmt.Println()

	// Find max widths for alignment
	maxID := len("ID")
	maxName := len("Name")
	for _, reg := range registrations {
		if len(reg.Meta.ID) > maxID {
			maxID = len(reg.Meta.ID)
		}
		if len(reg.Meta.Name) > maxName {
			maxName = len(reg.Meta.Name)
		}
	}

	// Header
	fmt.Printf("  %-3s  %-*s  %-*s  %s\n", "#", maxID, "ID", maxName, "Name", "Critical")
	for i, reg := range registrations {
		critical := "no"
		if reg.Meta.Critical {
			critical = "yes"
		}
		fmt.Printf("  %-3d  %-*s  %-*s  %s\n",
			i+1, maxID, reg.Meta.ID, maxName, reg.Meta.Name, critical)
	}

	if len(skipped) > 0 {
		fmt.Println()
		fmt.Println("Skipped:")
		for _, s := range skipped {
			reason := s.Reason
			if s.Pattern != "" {
				reason = fmt.Sprintf("%s (%s)", s.Reason, s.Pattern)
			}
			fmt.Printf("  %-*s  %-*s  %s\n", maxID, s.ID, maxName, s.Name, reason)
		}
	}
}
