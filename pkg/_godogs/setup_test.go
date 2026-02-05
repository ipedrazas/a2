package godogs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// a2BinaryPath is set by TestMain when the a2 binary is built for tests.
var a2BinaryPath string

// findRepoRoot walks up from the current directory to find the module root (containing go.mod).
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found in any parent directory")
		}
		dir = parent
	}
}

// buildA2 builds the a2 binary for use in godogs tests.
func buildA2() (string, error) {
	if path := os.Getenv("A2_BINARY"); path != "" {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	repoRoot, err := findRepoRoot()
	if err != nil {
		return "", fmt.Errorf("find repo root: %w", err)
	}

	// Build to dist/a2 in repo root so it can be reused across test runs
	distDir := filepath.Join(repoRoot, "dist")
	if err := os.MkdirAll(distDir, 0755); err != nil {
		return "", fmt.Errorf("create dist dir: %w", err)
	}
	a2Path := filepath.Join(distDir, "a2")

	cmd := exec.Command("go", "build", "-o", a2Path, ".")
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("build a2: %w", err)
	}

	return a2Path, nil
}

func TestMain(m *testing.M) {
	path, err := buildA2()
	if err != nil {
		fmt.Fprintf(os.Stderr, "godogs: build a2 for tests: %v\n", err)
		os.Exit(1)
	}
	a2BinaryPath = path
	os.Setenv("A2_BINARY", path)

	code := m.Run()
	os.Exit(code)
}
