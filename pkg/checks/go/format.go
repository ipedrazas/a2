package gocheck

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
)

// FormatCheck verifies that all Go code is properly formatted.
type FormatCheck struct{}

func (c *FormatCheck) ID() string   { return "go:format" }
func (c *FormatCheck) Name() string { return "Go Format" }

func (c *FormatCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangGo)

	// Run gofmt -l to list files that need formatting
	result := checkutil.RunCommand(path, "gofmt", "-l", ".")

	// gofmt returns non-zero if there's an error (not just unformatted files)
	if result.Err != nil && result.Stderr != "" {
		return rb.Warn("gofmt error: " + checkutil.TruncateMessage(result.Stderr, 150)), nil
	}

	output := strings.TrimSpace(result.Stdout)
	if output != "" {
		// Count unformatted files
		files := strings.Split(output, "\n")
		return rb.Warn(checkutil.TruncateMessage("Unformatted files: "+strings.Join(files, ", ")+". Run 'gofmt -w .' to fix.", 200)), nil
	}

	return rb.Pass("All Go files are properly formatted"), nil
}
