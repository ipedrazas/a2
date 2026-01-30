package godogs

import (
	"fmt"
)

// This file documents step implementations that live in other files.
// Only steps that are still defined here (no-op or delegation) remain.

// Steps implemented in config_steps.go: iCreateACustomProfileIn, iCreateImprovementPhasesIn,
// iDefineExternalChecksIn, iSelectProfile, iSelectTarget, iSetSeverity_modeTo, iSetTargetTo,
// iDisableCloudnativeChecksHealthMetricsTracing, iDisableContainerChecksDockerfileKS (no arg),
// iRelaxTestingRequirementsCoverageCyclomatic, iCreateForDevelopment, iCreateForMainBranch, iEdit, etc.

// Steps implemented in ci_steps.go: iHaveAGitHubRepository, iCreate, iHaveAConfiguredInTheProject,
// workflow/CI steps, theResultsShouldBeAvailableForDownload, etc.

// Steps implemented in output_steps.go: AI-assisted, detection, explanation, test result steps.

// Steps implemented in pending_steps.go: all remaining context/outcome steps (no-op).

func iHaveNotConfiguredA(_ int) error { return nil }

func iDontUnderstandWhatACheckDoes() error { return nil }

func iAmRunningAForTheFirstTime(_ int) error { return nil }

func iOnlyWantToCheckGoCode() error { return nil }

func iShouldSeeUpdatedResults() error {
	s := GetState()
	if s.GetLastOutput() == "" {
		return fmt.Errorf("no updated results in output")
	}
	return nil
}

func iRerun(arg1 string) error    { return iRunAgain(arg1) }
func iRunAgain(arg1 string) error { return iRunCommand(arg1) }
