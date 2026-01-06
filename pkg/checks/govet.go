package checks

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
)

// GoVetCheck runs go vet to find suspicious constructs.
type GoVetCheck struct{}

func (c *GoVetCheck) ID() string   { return "govet" }
func (c *GoVetCheck) Name() string { return "Go Vet" }

func (c *GoVetCheck) Run(path string) (checker.Result, error) {
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

		// Count the number of issues (each line is typically one issue)
		lines := strings.Split(output, "\n")
		issueCount := 0
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				issueCount++
			}
		}

		return checker.Result{
			Name:    c.Name(),
			ID:      c.ID(),
			Passed:  false,
			Status:  checker.Warn,
			Message: output,
		}, nil
	}

	return checker.Result{
		Name:    c.Name(),
		ID:      c.ID(),
		Passed:  true,
		Status:  checker.Pass,
		Message: "No issues found",
	}, nil
}
