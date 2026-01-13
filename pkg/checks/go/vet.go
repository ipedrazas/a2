package gocheck

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
)

// VetCheck runs go vet to find suspicious constructs.
type VetCheck struct{}

func (c *VetCheck) ID() string   { return "go:vet" }
func (c *VetCheck) Name() string { return "Go Vet" }

func (c *VetCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangGo)

	result := checkutil.RunCommand(path, "go", "vet", "./...")
	if !result.Success() {
		return rb.Warn(result.Output()), nil
	}

	return rb.Pass("No issues found"), nil
}
