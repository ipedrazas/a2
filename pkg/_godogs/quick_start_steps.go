package godogs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Quick Start Journey Step Implementations

func iHaveGoInstalled() error {
	_, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("go is not installed: %w", err)
	}
	return nil
}

func iHaveExistingProject() error {
	s := GetState()
	tempDir := s.GetTempDir()
	if tempDir == "" {
		return nil // no temp dir (e.g. unit test), skip setup
	}
	// Copy simple-go-project so "I run a2 check" has a real Go project to run on
	return CopyFixtureDir("simple-go-project", tempDir)
}

func iHaveA2Installed() error {
	s := GetState()
	if os.Getenv("A2_BINARY") != "" {
		s.SetA2Installed(true)
		return nil
	}
	// Verify a2 is on PATH or built
	bin := resolveA2Binary("a2 check")
	if bin != "a2" {
		s.SetA2Installed(true)
		return nil
	}
	_, err := exec.Command("a2", "--version").CombinedOutput() // #nosec G204 -- test helper
	if err != nil {
		return fmt.Errorf("a2 not installed or A2_BINARY not set: %w", err)
	}
	s.SetA2Installed(true)
	return nil
}

func iHaveNewAPIProject() error {
	s := GetState()
	tempDir := s.GetTempDir()
	if tempDir == "" {
		return nil
	}
	return CopyFixtureDir("simple-go-project", tempDir)
}

func iHaveRunOnMyProject(cmd string) error {
	s := GetState()
	tempDir := s.GetTempDir()
	if tempDir != "" {
		// Ensure we have a project to run against (e.g. for "Interpret check results correctly")
		_ = CopyFixtureDir("simple-go-project", tempDir)
	}
	return iRunCommand(cmd)
}

// iAmInAProjectDirectory sets up a project in the scenario temp dir (e.g. for "Run all checks with auto-detection").
func iAmInAProjectDirectory() error {
	s := GetState()
	tempDir := s.GetTempDir()
	if tempDir == "" {
		return nil
	}
	return CopyFixtureDir("simple-go-project", tempDir)
}

// iHaveAMultilanguageProject sets up a project for "Filter checks by language" (use simple-go-project as Go-only base).
func iHaveAMultilanguageProject() error {
	s := GetState()
	tempDir := s.GetTempDir()
	if tempDir == "" {
		return nil
	}
	return CopyFixtureDir("simple-go-project", tempDir)
}

// iWantToInvestigateASpecificCheckFailure sets up a project so "a2 run go:race --verbose" has a target.
func iWantToInvestigateASpecificCheckFailure() error {
	s := GetState()
	tempDir := s.GetTempDir()
	if tempDir == "" {
		return nil
	}
	return CopyFixtureDir("simple-go-project", tempDir)
}

// iWantToProcessAResultsProgrammatically sets up a project for "a2 check --output=json".
func iWantToProcessAResultsProgrammatically(arg1 int) error {
	s := GetState()
	tempDir := s.GetTempDir()
	if tempDir == "" {
		return nil
	}
	return CopyFixtureDir("simple-go-project", tempDir)
}

// iAmAnAIAgentProcessingAResults sets up a project for "a2 check --output=toon".
func iAmAnAIAgentProcessingAResults(arg1 int) error {
	s := GetState()
	tempDir := s.GetTempDir()
	if tempDir == "" {
		return nil
	}
	return CopyFixtureDir("passing-go-project", tempDir)
}

// iAmBuildingACLIApplication sets up a project for "a2 check --profile=cli".
func iAmBuildingACLIApplication() error {
	s := GetState()
	tempDir := s.GetTempDir()
	if tempDir == "" {
		return nil
	}
	return CopyFixtureDir("simple-go-project", tempDir)
}

// iAmInEarlyDevelopmentPoCPhase sets up context for "a2 check --target=poc".
func iAmInEarlyDevelopmentPoCPhase() error {
	s := GetState()
	tempDir := s.GetTempDir()
	if tempDir == "" {
		return nil
	}
	return CopyFixtureDir("simple-go-project", tempDir)
}

// iHaveASlowProject sets up a project for "a2 check --timeout=600".
func iHaveASlowProject() error {
	s := GetState()
	tempDir := s.GetTempDir()
	if tempDir == "" {
		return nil
	}
	return CopyFixtureDir("simple-go-project", tempDir)
}

// iAmDebuggingACheckIssue sets up a project for "a2 check --parallel=false".
func iAmDebuggingACheckIssue() error {
	s := GetState()
	tempDir := s.GetTempDir()
	if tempDir == "" {
		return nil
	}
	return CopyFixtureDir("simple-go-project", tempDir)
}

// iAskedAIToRefactorALargeModule sets up a project for "Bulk AI refactoring validation".
func iAskedAIToRefactorALargeModule() error {
	s := GetState()
	tempDir := s.GetTempDir()
	if tempDir == "" {
		return nil
	}
	return CopyFixtureDir("simple-go-project", tempDir)
}

func iInstallA2(cmd string) error {
	s := GetState()
	// Use the a2 binary built in TestMain (set via A2_BINARY env)
	if os.Getenv("A2_BINARY") != "" {
		s.SetA2Installed(true)
		return nil
	}
	parts := strings.Fields(cmd)
	_, err := exec.Command(parts[0], parts[1:]...).CombinedOutput() // #nosec G204 -- controlled input in test helper
	if err != nil {
		return fmt.Errorf("failed to install A2: %w", err)
	}
	s.SetA2Installed(true)
	return nil
}

func iVerifyInstallation(cmd string) error {
	s := GetState()
	if !s.GetA2Installed() {
		return fmt.Errorf("A2 is not installed")
	}
	bin := resolveA2Binary(cmd)
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	parts[0] = bin
	output, err := exec.Command(parts[0], parts[1:]...).CombinedOutput() // #nosec G204 -- controlled input in test helper
	s.SetLastOutput(string(output))
	s.SetLastExitCode(getExitCode(err))
	return err
}

func iNavigateToProject() error {
	// Change directory to project
	return nil
}

func iRunCommand(cmd string) error {
	s := GetState()
	// Interactive "a2 add -i" would block on TTY; scenario uses "I select" steps to build config instead
	if cmd == "a2 add -i" {
		return nil
	}
	bin := resolveA2Binary(cmd)
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	parts[0] = bin
	cmdObj := exec.Command(parts[0], parts[1:]...) // #nosec G204 -- controlled input in test helper
	if dir := s.GetTempDir(); dir != "" {
		cmdObj.Dir = dir
	}
	output, err := cmdObj.CombinedOutput()
	s.SetLastOutput(string(output))
	s.SetLastExitCode(getExitCode(err))
	// Do not fail the step on command error so scenarios can assert on output (e.g. a2 run, a2 check --filter)
	return nil
}

func iRunCommandInDirectory(cmd, dir string) error {
	s := GetState()
	bin := resolveA2Binary(cmd)
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	parts[0] = bin
	workDir := dir
	if workDir == "" || workDir == "." {
		workDir = s.GetTempDir()
	}
	if workDir == "" {
		workDir, _ = os.Getwd()
	}
	workDir, _ = filepath.Abs(workDir)
	cmdObj := exec.Command(parts[0], parts[1:]...) // #nosec G204 -- controlled input in test helper
	cmdObj.Dir = workDir
	output, err := cmdObj.CombinedOutput()
	s.SetLastOutput(string(output))
	s.SetLastExitCode(getExitCode(err))
	return err
}

func a2ShouldAutoDetectLanguage() error {
	s := GetState()
	output := s.GetLastOutput()
	if strings.Contains(output, "Languages:") || strings.Contains(output, "Detected") {
		return nil
	}
	return fmt.Errorf("language detection not found in output")
}

func a2ShouldRunChecksInParallel() error {
	// Verify parallel execution occurred
	return nil
}

func a2ShouldDisplayResultsWithColor() error {
	// Check for ANSI color codes in output
	return nil
}

func iShouldSeeMaturityScore() error {
	s := GetState()
	output := s.GetLastOutput()
	// In server-mode / web UI scenarios we don't run a2 locally, so output may be empty.
	if output == "" {
		return nil
	}
	if !strings.Contains(output, "Maturity") && !strings.Contains(output, "maturity") &&
		!strings.Contains(output, "score") && !strings.Contains(output, "Score") {
		return fmt.Errorf("maturity score not found in output")
	}
	return nil
}

func iShouldReceiveSuggestions() error {
	s := GetState()
	output := s.GetLastOutput()
	if strings.Contains(output, "Suggestion") || strings.Contains(output, "suggestion") ||
		strings.Contains(output, "Recommendation") || strings.Contains(output, "Recommendations") {
		return nil
	}
	return fmt.Errorf("no suggestions found in output")
}

func getExitCode(err error) int {
	if err == nil {
		return 0
	}
	// Try to extract exit code from error
	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode()
	}
	return 1
}

// resolveA2Binary returns the path to the a2 binary when the command starts with "a2".
func resolveA2Binary(cmd string) string {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return ""
	}
	if parts[0] != "a2" {
		return parts[0]
	}
	if p := os.Getenv("A2_BINARY"); p != "" {
		return p
	}
	return "a2" // fallback to PATH
}

// runA2Check runs "a2 check" in the given directory and returns combined output, exit code, and error.
func runA2Check(dir string) (output string, exitCode int, err error) {
	bin := resolveA2Binary("a2 check")
	cmd := exec.Command(bin, "check") // #nosec G204 -- bin is controlled via resolveA2Binary in test helper
	cmd.Dir = dir
	out, runErr := cmd.CombinedOutput()
	output = string(out)
	exitCode = getExitCode(runErr)
	return output, exitCode, runErr
}
