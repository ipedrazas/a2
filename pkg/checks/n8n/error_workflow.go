package n8ncheck

import (
	"context"
	"fmt"

	"github.com/ipedrazas/a2/pkg/checker"
)

// ErrorWorkflowCheck verifies that active workflows have an error workflow
// configured (settings.errorWorkflow), so failures alert someone instead of
// dying silently.
type ErrorWorkflowCheck struct {
	// FailRatio is the coverage ratio below which the check fails instead
	// of warns (default 0.5).
	FailRatio float64
}

func (c *ErrorWorkflowCheck) ID() string   { return "n8n:error-workflow" }
func (c *ErrorWorkflowCheck) Name() string { return "Error Workflow Coverage" }

func (c *ErrorWorkflowCheck) Run(ctx context.Context, path string) (checker.Result, error) {
	bundle, early := loadBundle(c, path)
	if early != nil {
		return *early, nil
	}
	clean, degraded := splitDegraded(bundle)

	var active, uncovered int
	var findings []finding
	for _, wf := range clean {
		if !wf.Workflow.Active {
			continue
		}
		active++
		if wf.Workflow.Settings.ErrorWorkflow == "" {
			uncovered++
			findings = append(findings, finding{
				File:     wf.Path,
				Workflow: workflowLabel(wf),
				Note:     "no error workflow configured",
			})
		}
	}

	if active == 0 {
		return finish(c, infoResult(c, "not applicable: no active workflows"), degraded), nil
	}

	coverage := float64(active-uncovered) / float64(active)
	msg := fmt.Sprintf("%d of %d active %s have no error workflow (%.0f%% coverage)",
		uncovered, active, pluralWorkflows(active), coverage*100)

	var res checker.Result
	switch {
	case uncovered == 0:
		res = passResult(c, fmt.Sprintf("all %d active %s have an error workflow", active, pluralWorkflows(active)))
	case coverage < c.FailRatio:
		res = failResult(c, msg, formatFindings(findings))
	default:
		res = warnResult(c, msg, formatFindings(findings))
	}
	return finish(c, res, degraded), nil
}
