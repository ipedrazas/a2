package common

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// FileExistsCheck verifies that required files exist in the project.
type FileExistsCheck struct {
	Files []string // List of files to check for
}

func (c *FileExistsCheck) ID() string   { return "file_exists" }
func (c *FileExistsCheck) Name() string { return "Required Files" }

func (c *FileExistsCheck) Run(path string) (checker.Result, error) {
	var missing []string

	for _, file := range c.Files {
		// Use safepath to prevent directory traversal attacks
		if !safepath.Exists(path, file) {
			missing = append(missing, file)
		}
	}

	if len(missing) > 0 {
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   false,
			Status:   checker.Warn,
			Message:  "Missing files: " + strings.Join(missing, ", "),
			Language: checker.LangCommon,
		}, nil
	}

	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   true,
		Status:   checker.Pass,
		Message:  "All required files present",
		Language: checker.LangCommon,
	}, nil
}

// DefaultFileExistsCheck returns a FileExistsCheck with common project files.
func DefaultFileExistsCheck() *FileExistsCheck {
	return &FileExistsCheck{
		Files: []string{"README.md", "LICENSE"},
	}
}
