package gocheck

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
)

// VetCheck runs go vet to find suspicious constructs.
type VetCheck struct{}

func (c *VetCheck) ID() string   { return "go:vet" }
func (c *VetCheck) Name() string { return "Go Vet" }

func (c *VetCheck) Run(path string) (checker.Result, error) {
	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// go vet returns non-zero if it finds issues
		output := strings.TrimSpace(stderr.String())
		if output == "" {
			output = strings.TrimSpace(stdout.String())
		}

		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   false,
			Status:   checker.Warn,
			Message:  output,
			Language: checker.LangGo,
		}, nil
	}

	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   true,
		Status:   checker.Pass,
		Message:  "No issues found",
		Language: checker.LangGo,
	}, nil
}
