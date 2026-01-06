package checks

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
)

// GetChecks returns the checks to run based on configuration.
// Checks are ordered with critical (Fail) checks first.
func GetChecks(cfg *config.Config) []checker.Checker {
	allChecks := []checker.Checker{
		// Critical checks (Fail status - veto power)
		&ModuleCheck{},
		&BuildCheck{},
		&TestRunnerCheck{},

		// Warning checks
		&FileExistsCheck{Files: cfg.Files.Required},
		&GofmtCheck{},
		&GoVetCheck{},
		&CoverageCheck{Threshold: cfg.Coverage.Threshold},
		&DependencyCheck{},
	}

	// Add external checks from config
	for _, ext := range cfg.External {
		allChecks = append(allChecks, &ExternalCheck{
			CheckID:   ext.ID,
			CheckName: ext.Name,
			Command:   ext.Command,
			Args:      ext.Args,
			Severity:  ext.Severity,
		})
	}

	// Filter out disabled checks
	var enabled []checker.Checker
	for _, check := range allChecks {
		if !cfg.IsCheckDisabled(check.ID()) {
			enabled = append(enabled, check)
		}
	}

	return enabled
}

// DefaultChecks returns the default set of checks (for backward compatibility).
func DefaultChecks() []checker.Checker {
	return GetChecks(config.DefaultConfig())
}
