package pythoncheck

import (
	"path/filepath"

	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// runPythonCommand runs a command within the project's Python environment.
// It detects uv, poetry, or local virtualenvs and wraps the command accordingly.
func runPythonCommand(projectPath, command string, args ...string) *checkutil.CommandResult {
	name, prefix := resolvePythonEnv(projectPath, command)
	allArgs := append(prefix, args...)
	return checkutil.RunCommand(projectPath, name, allArgs...)
}

// pythonToolAvailable checks if a Python tool is available for a project,
// considering virtualenvs and package managers.
func pythonToolAvailable(projectPath, tool string) bool {
	for _, venvDir := range []string{".venv", "venv"} {
		binPath := filepath.Join(venvDir, "bin", tool)
		if safepath.Exists(projectPath, binPath) {
			return true
		}
	}
	return checkutil.ToolAvailable(tool)
}

// resolvePythonEnv determines how to run a command within the project's Python environment.
// Returns the resolved executable name and any prefix arguments.
//
// Detection order:
//  1. uv project (uv.lock exists, uv available) → uv run <command>
//  2. poetry project (poetry.lock exists, poetry available) → poetry run <command>
//  3. local venv (.venv/bin/<command> or venv/bin/<command>) → direct path
//  4. bare command (system PATH)
func resolvePythonEnv(projectPath, command string) (name string, prefixArgs []string) {
	// uv project
	if safepath.Exists(projectPath, "uv.lock") && checkutil.ToolAvailable("uv") {
		return "uv", []string{"run", command}
	}

	// poetry project
	if safepath.Exists(projectPath, "poetry.lock") && checkutil.ToolAvailable("poetry") {
		return "poetry", []string{"run", command}
	}

	// local venv
	for _, venvDir := range []string{".venv", "venv"} {
		binPath := filepath.Join(venvDir, "bin", command)
		if safepath.Exists(projectPath, binPath) {
			return filepath.Join(projectPath, venvDir, "bin", command), nil
		}
	}

	return command, nil
}
