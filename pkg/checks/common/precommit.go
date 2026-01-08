package common

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// PrecommitCheck verifies that pre-commit hooks are configured.
type PrecommitCheck struct{}

func (c *PrecommitCheck) ID() string   { return "common:precommit" }
func (c *PrecommitCheck) Name() string { return "Pre-commit Hooks" }

// Run checks for pre-commit hook configurations.
func (c *PrecommitCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangCommon,
	}

	var foundTools []string

	// Check for pre-commit framework (Python-based)
	if safepath.Exists(path, ".pre-commit-config.yaml") || safepath.Exists(path, ".pre-commit-config.yml") {
		foundTools = append(foundTools, "pre-commit")
	}

	// Check for Husky (Node.js)
	if c.hasHusky(path) {
		foundTools = append(foundTools, "Husky")
	}

	// Check for Lefthook
	if safepath.Exists(path, "lefthook.yml") || safepath.Exists(path, "lefthook.yaml") || safepath.Exists(path, ".lefthook.yml") {
		foundTools = append(foundTools, "Lefthook")
	}

	// Check for Overcommit (Ruby)
	if safepath.Exists(path, ".overcommit.yml") {
		foundTools = append(foundTools, "Overcommit")
	}

	// Check for commitlint
	if c.hasCommitlint(path) {
		foundTools = append(foundTools, "commitlint")
	}

	// Check for lint-staged (often used with Husky)
	if c.hasLintStaged(path) {
		foundTools = append(foundTools, "lint-staged")
	}

	// Check for git hooks directory with executable hooks
	if c.hasGitHooks(path) {
		foundTools = append(foundTools, "git hooks")
	}

	if len(foundTools) > 0 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "Pre-commit hooks configured: " + strings.Join(foundTools, ", ")
		return result, nil
	}

	result.Passed = false
	result.Status = checker.Warn
	result.Message = "No pre-commit hooks configured (consider adding pre-commit, Husky, or Lefthook)"
	return result, nil
}

// hasHusky checks if Husky is configured.
func (c *PrecommitCheck) hasHusky(path string) bool {
	// Check for .husky directory
	huskyDir := filepath.Join(path, ".husky")
	if info, err := os.Stat(huskyDir); err == nil && info.IsDir() {
		// Check if there are hook files in .husky
		entries, err := os.ReadDir(huskyDir)
		if err == nil {
			for _, entry := range entries {
				name := entry.Name()
				// Look for common hook files (not starting with _ which are husky internals)
				if !strings.HasPrefix(name, "_") && !strings.HasPrefix(name, ".") {
					if name == "pre-commit" || name == "pre-push" || name == "commit-msg" {
						return true
					}
				}
			}
		}
	}

	// Check for husky config in package.json
	content, err := safepath.ReadFile(path, "package.json")
	if err != nil {
		return false
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(content, &pkg); err != nil {
		return false
	}

	// Check for "husky" key in package.json
	if _, ok := pkg["husky"]; ok {
		return true
	}

	// Check for husky in devDependencies
	if devDeps, ok := pkg["devDependencies"].(map[string]interface{}); ok {
		if _, hasHusky := devDeps["husky"]; hasHusky {
			return true
		}
	}

	return false
}

// hasCommitlint checks if commitlint is configured.
func (c *PrecommitCheck) hasCommitlint(path string) bool {
	commitlintConfigs := []string{
		"commitlint.config.js",
		"commitlint.config.cjs",
		"commitlint.config.mjs",
		"commitlint.config.ts",
		".commitlintrc",
		".commitlintrc.json",
		".commitlintrc.yaml",
		".commitlintrc.yml",
		".commitlintrc.js",
		".commitlintrc.cjs",
	}

	for _, config := range commitlintConfigs {
		if safepath.Exists(path, config) {
			return true
		}
	}

	// Check for commitlint config in package.json
	content, err := safepath.ReadFile(path, "package.json")
	if err != nil {
		return false
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(content, &pkg); err != nil {
		return false
	}

	if _, ok := pkg["commitlint"]; ok {
		return true
	}

	return false
}

// hasLintStaged checks if lint-staged is configured.
func (c *PrecommitCheck) hasLintStaged(path string) bool {
	lintStagedConfigs := []string{
		"lint-staged.config.js",
		"lint-staged.config.cjs",
		"lint-staged.config.mjs",
		".lintstagedrc",
		".lintstagedrc.json",
		".lintstagedrc.yaml",
		".lintstagedrc.yml",
		".lintstagedrc.js",
		".lintstagedrc.cjs",
	}

	for _, config := range lintStagedConfigs {
		if safepath.Exists(path, config) {
			return true
		}
	}

	// Check for lint-staged config in package.json
	content, err := safepath.ReadFile(path, "package.json")
	if err != nil {
		return false
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(content, &pkg); err != nil {
		return false
	}

	if _, ok := pkg["lint-staged"]; ok {
		return true
	}

	return false
}

// hasGitHooks checks if there are custom git hooks in .git/hooks.
func (c *PrecommitCheck) hasGitHooks(path string) bool {
	hooksDir := filepath.Join(path, ".git", "hooks")
	entries, err := os.ReadDir(hooksDir)
	if err != nil {
		return false
	}

	// Look for executable hook files (not .sample files)
	hookNames := []string{"pre-commit", "pre-push", "commit-msg", "prepare-commit-msg"}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Skip sample files
		if strings.HasSuffix(name, ".sample") {
			continue
		}
		for _, hookName := range hookNames {
			if name == hookName {
				// Check if file is executable
				info, err := entry.Info()
				if err != nil {
					continue
				}
				if info.Mode()&0111 != 0 {
					return true
				}
			}
		}
	}

	return false
}
