package gocheck

import (
	"context"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
)

// VetCheck runs go vet to find suspicious constructs.
type VetCheck struct{}

func (c *VetCheck) ID() string   { return "go:vet" }
func (c *VetCheck) Name() string { return "Go Vet" }

func (c *VetCheck) Run(ctx context.Context, path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangGo)

	result := rb.RunCommand(ctx, path, "go", "vet", "./...")
	output := result.CombinedOutput()
	if !result.Success() {
		return rb.WarnWithOutput("go vet issues: "+checkutil.TruncateMessage(result.Output(), 200), output), nil
	}

	return rb.Pass("No issues found"), nil
}
