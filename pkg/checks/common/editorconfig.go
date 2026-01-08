package common

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// EditorconfigCheck verifies editor configuration exists.
type EditorconfigCheck struct{}

func (c *EditorconfigCheck) ID() string   { return "common:editorconfig" }
func (c *EditorconfigCheck) Name() string { return "Editor Config" }

// Run checks for editor configuration files.
func (c *EditorconfigCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangCommon,
	}

	var found []string

	// Check for .editorconfig
	if safepath.Exists(path, ".editorconfig") {
		found = append(found, ".editorconfig")
	}

	// Check for VS Code settings
	if safepath.Exists(path, ".vscode/settings.json") {
		found = append(found, "VS Code settings")
	}
	if safepath.Exists(path, ".vscode/extensions.json") {
		found = append(found, "VS Code extensions")
	}

	// Check for JetBrains IDE config
	if safepath.Exists(path, ".idea") {
		found = append(found, "JetBrains IDE")
	}
	if safepath.Exists(path, ".idea/codeStyles") {
		found = append(found, "JetBrains code styles")
	}

	// Check for Vim/Neovim config
	vimConfigs := []string{".vimrc", ".nvimrc", ".vim", ".nvim"}
	for _, cfg := range vimConfigs {
		if safepath.Exists(path, cfg) {
			found = append(found, "Vim/Neovim config")
			break
		}
	}

	// Check for devcontainer (standardized dev environment)
	if safepath.Exists(path, ".devcontainer/devcontainer.json") ||
		safepath.Exists(path, ".devcontainer.json") {
		found = append(found, "Dev Container")
	}

	// Check for workspace config
	workspaceConfigs := []string{
		"*.code-workspace",
		".project",
		".classpath",
	}
	for _, cfg := range workspaceConfigs {
		if safepath.Exists(path, cfg) {
			found = append(found, "workspace config")
			break
		}
	}

	// Check .editorconfig content for quality
	if safepath.Exists(path, ".editorconfig") {
		if content, err := safepath.ReadFile(path, ".editorconfig"); err == nil {
			contentStr := strings.ToLower(string(content))
			var settings []string
			if strings.Contains(contentStr, "indent_style") {
				settings = append(settings, "indent")
			}
			if strings.Contains(contentStr, "end_of_line") {
				settings = append(settings, "line endings")
			}
			if strings.Contains(contentStr, "charset") {
				settings = append(settings, "charset")
			}
			if strings.Contains(contentStr, "trim_trailing_whitespace") {
				settings = append(settings, "whitespace")
			}
			if len(settings) > 0 {
				found = append(found, "configures: "+strings.Join(settings, ", "))
			}
		}
	}

	// Build result
	if len(found) > 0 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "Editor config: " + strings.Join(found, ", ")
	} else {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No editor config found (consider adding .editorconfig)"
	}

	return result, nil
}
