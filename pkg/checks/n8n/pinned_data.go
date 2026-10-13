package n8ncheck

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
)

// PinnedDataCheck flags pinData on active workflows: test fixtures left
// live, a marker that production and testing share an environment.
type PinnedDataCheck struct{}

func (c *PinnedDataCheck) ID() string   { return "n8n:pinned-data" }
func (c *PinnedDataCheck) Name() string { return "Pinned Data In Production" }

func (c *PinnedDataCheck) Run(ctx context.Context, path string) (checker.Result, error) {
	bundle, early := loadBundle(c, path)
	if early != nil {
		return *early, nil
	}
	clean, degraded := splitDegraded(bundle)

	var findings []finding
	for _, wf := range clean {
		if !wf.Workflow.Active || len(wf.Workflow.PinData) == 0 {
			continue
		}
		nodes := make([]string, 0, len(wf.Workflow.PinData))
		for name := range wf.Workflow.PinData {
			nodes = append(nodes, name)
		}
		sort.Strings(nodes)
		findings = append(findings, finding{
			File:     wf.Path,
			Workflow: workflowLabel(wf),
			Note:     fmt.Sprintf("pinned test data on %d node(s): %s", len(nodes), strings.Join(nodes, ", ")),
		})
	}

	var res checker.Result
	if len(findings) == 0 {
		res = passResult(c, "no active workflows carry pinned test data")
	} else {
		res = warnResult(c,
			fmt.Sprintf("%d active %s carry pinned test data", len(findings), pluralWorkflows(len(findings))),
			formatFindings(findings))
	}
	return finish(c, res, degraded), nil
}
