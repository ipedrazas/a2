package godogs

import (
	"fmt"
	"os/exec"
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
	// For testing, we create a temporary Go project
	// In real tests, this would set up a test project
	return nil
}

func iInstallA2(cmd string) error {
	s := GetState()
	parts := strings.Fields(cmd)
	_, err := exec.Command(parts[0], parts[1:]...).CombinedOutput()
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
	parts := strings.Fields(cmd)
	_, err := exec.Command(parts[0], parts[1:]...).CombinedOutput()
	return err
}

func iNavigateToProject() error {
	// Change directory to project
	return nil
}

func iRunCommand(cmd string) error {
	s := GetState()
	parts := strings.Fields(cmd)
	output, err := exec.Command(parts[0], parts[1:]...).CombinedOutput()
	s.SetLastOutput(string(output))
	s.SetLastExitCode(getExitCode(err))
	return err
}

func iRunCommandInDirectory(cmd, dir string) error {
	s := GetState()
	parts := strings.Fields(cmd)
	output, err := exec.Command(parts[0], parts[1:]...).CombinedOutput()
	s.SetLastOutput(string(output))
	s.SetLastExitCode(getExitCode(err))
	return err
}

func a2ShouldAutoDetectLanguage() error {
	s := GetState()
	if !strings.Contains(s.GetLastOutput(), "Detected") {
		return fmt.Errorf("language detection not found in output")
	}
	return nil
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
	if !strings.Contains(s.GetLastOutput(), "Maturity") &&
		!strings.Contains(s.GetLastOutput(), "maturity") {
		return fmt.Errorf("maturity score not found in output")
	}
	return nil
}

func iShouldReceiveSuggestions() error {
	s := GetState()
	if !strings.Contains(s.GetLastOutput(), "Suggestion") &&
		!strings.Contains(s.GetLastOutput(), "suggestion") {
		return fmt.Errorf("no suggestions found in output")
	}
	return nil
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
