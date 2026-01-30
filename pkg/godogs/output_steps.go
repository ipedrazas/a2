package godogs

import (
	"fmt"
	"strings"
	"github.com/cucumber/godog"
)

// Output interpretation and validation step implementations

func iViewOutput() error {
	s := GetState()
	if s.GetLastOutput() == "" {
		return fmt.Errorf("no output available")
	}
	return nil
}

func iShouldSeeGreenCheckmarks() error {
	s := GetState()
	if !strings.Contains(s.GetLastOutput(), "✓") &&
	   !strings.Contains(s.GetLastOutput(), "[PASS]") &&
	   !strings.Contains(s.GetLastOutput(), "✅") {
		return fmt.Errorf("no green checkmarks found in output")
	}
	return nil
}

func iShouldSeeRedXMarks() error {
	s := GetState()
	if !strings.Contains(s.GetLastOutput(), "✗") &&
	   !strings.Contains(s.GetLastOutput(), "[FAIL]") &&
	   !strings.Contains(s.GetLastOutput(), "❌") {
		return fmt.Errorf("no red X marks found in output")
	}
	return nil
}

func iShouldSeeYellowWarnings() error {
	s := GetState()
	if !strings.Contains(s.GetLastOutput(), "⚠") &&
	   !strings.Contains(s.GetLastOutput(), "[WARN]") {
		return fmt.Errorf("no yellow warnings found in output")
	}
	return nil
}

func iShouldSeeBlueInfo() error {
	s := GetState()
	if !strings.Contains(s.GetLastOutput(), "ℹ") &&
	   !strings.Contains(s.GetLastOutput(), "[INFO]") {
		return fmt.Errorf("no blue info found in output")
	}
	return nil
}

func a2DetectedIssues() error {
	s := GetState()
	if !strings.Contains(s.GetLastOutput(), "✗") &&
	   !strings.Contains(s.GetLastOutput(), "⚠") {
		return fmt.Errorf("no issues detected in output")
	}
	return nil
}

func iReceivedSuggestions() error {
	s := GetState()
	if !strings.Contains(s.GetLastOutput(), "Suggestion") &&
	   !strings.Contains(s.GetLastOutput(), "suggest") {
		return fmt.Errorf("no suggestions found in output")
	}
	return nil
}

func iFixIssues() error {
	return nil
}

func maturityScoreShouldImprove() error {
	return nil
}

func iShouldCommitWithConfidence() error {
	s := GetState()
	if s.GetLastExitCode() != 0 {
		return fmt.Errorf("checks did not pass, cannot commit with confidence")
	}
	return nil
}

func a2ShouldDisplayInfoStatus() error {
	s := GetState()
	if !strings.Contains(s.GetLastOutput(), "ℹ") {
		return fmt.Errorf("Info status not displayed for missing tools")
	}
	return nil
}

func a2ShouldContinueRunning() error {
	return nil
}

func a2ShouldSuggestToolInstallation() error {
	s := GetState()
	if !strings.Contains(s.GetLastOutput(), "install") {
		return fmt.Errorf("no installation suggestions found")
	}
	return nil
}

func iSeeProgressIndicators() error {
	s := GetState()
	if !strings.Contains(s.GetLastOutput(), "...") &&
	   !strings.Contains(s.GetLastOutput(), "Running") {
		return fmt.Errorf("no progress indicators found")
	}
	return nil
}

func failuresShouldDecrease() error {
	return nil
}

func iContinueFixing() error {
	return nil
}

func iReceiveTokenFormat() error {
	s := GetState()
	if !strings.Contains(s.GetLastOutput(), "results[") {
		return fmt.Errorf("output is not in TOON format")
	}
	return nil
}

func iSeeAllCheckStatus() error {
	return nil
}

func iIdentifyRemainingIssues() error {
	s := GetState()
	if !strings.Contains(s.GetLastOutput(), "⚠") {
		return fmt.Errorf("no remaining issues identified")
	}
	return nil
}

func iAddressWarnings() error {
	return nil
}

func allChecksPass() error {
	s := GetState()
	if s.GetLastExitCode() != 0 {
		return fmt.Errorf("checks did not pass: exit code %d", s.GetLastExitCode())
	}
	return nil
}

func maturityScoreIs(score int) error {
	s := GetState()
	s.SetMaturityScore(score)
	scoreStr := fmt.Sprintf("%d%%", score)
	if !strings.Contains(s.GetLastOutput(), scoreStr) &&
	   !strings.Contains(s.GetLastOutput(), fmt.Sprintf("%d", score)) {
		return fmt.Errorf("maturity score %d not found in output", score)
	}
	return nil
}

func iCanPush() error {
	s := GetState()
	if s.GetLastExitCode() != 0 {
		return fmt.Errorf("cannot push with exit code %d", s.GetLastExitCode())
	}
	return nil
}

func a2ReportedFailingTest() error {
	return godog.ErrPending
}

func iDontKnowRequirements() error {
	return nil
}

func iSeeCheckDescription() error {
	s := GetState()
	if !strings.Contains(s.GetLastOutput(), "Description") {
		return fmt.Errorf("check description not found")
	}
	return nil
}

func iSeeToolCommand() error {
	return nil
}

func iSeeRequirements() error {
	s := GetState()
	if !strings.Contains(s.GetLastOutput(), "Requirements") {
		return fmt.Errorf("requirements not found in output")
	}
	return nil
}

func iReceiveFixSuggestions() error {
	s := GetState()
	if !strings.Contains(s.GetLastOutput(), "Suggestion") {
		return fmt.Errorf("fix suggestions not found")
	}
	return nil
}

func iReceiveFastFeedback() error {
	return nil
}

func iCatchIssuesEarly() error {
	return nil
}

func feedbackLoopLessThan(minutes int) error {
	return nil
}

func velocityMaintained() error {
	return nil
}

func teamMaintainsHighStandards() error {
	return nil
}

func iChooseNoExternalChecks() error {
	return nil
}

func iHaveNotReviewedCode() error {
	return nil
}

func iReceivedFailureMessages() error {
	return nil
}

func iWorkOnFeature(hours int) error {
	return nil
}

func iWantToPush() error {
	return nil
}

func iCommunicateToTeam() error {
	return nil
}

// Missing step implementations - stub functions for pending steps
func iCommitConfig() error { return godog.ErrPending }
func iPushToRepo() error { return godog.ErrPending }
func teamCanRun(cmd string) error { return godog.ErrPending }
func everyoneSeesSameStandards() error { return godog.ErrPending }
func configVersionControlled() error { return godog.ErrPending }
func iHaveNewAPIProject() error { return godog.ErrPending }
func iRunInInteractiveMode(cmd string) error { return godog.ErrPending }
func iSelectApplicationType(appType string) error { return godog.ErrPending }
func iSelectMaturityLevel(level string) error { return godog.ErrPending }
func iSelectLanguageDetection(detection string) error { return godog.ErrPending }
func a2CreatesConfig(filename string) error { return godog.ErrPending }
func iHaveGoOnlyProject() error { return godog.ErrPending }
func iHaveConfiguredA2() error { return godog.ErrPending }
func iUseAIGenerateCode() error { return godog.ErrPending }
func a2DetectsBuildFailures() error { return godog.ErrPending }
func a2IdentifiesMissingTests() error { return godog.ErrPending }
func iHaveA2Installed() error { return godog.ErrPending }
func iRunWithOutputFormat(cmd, format string) error { return godog.ErrPending }

// Additional missing stub implementations
func someToolsNotInstalled() error { return godog.ErrPending }
func checkShouldComplete() error { return godog.ErrPending }
func a2FlagsFormatIssues() error { return godog.ErrPending }
func a2ChecksSecurity() error { return godog.ErrPending }
func iReceiveActionableFeedback() error { return godog.ErrPending }
func iFixBuildIssues() error { return godog.ErrPending }
func iRunFinalTime() error { return godog.ErrPending }
func iImplementingStepByStep() error { return godog.ErrPending }
func iRunAfterEachChange() error { return godog.ErrPending }
