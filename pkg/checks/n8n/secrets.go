package n8ncheck

import (
	"context"
	"fmt"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/n8n"
)

// SecretsCheck scans node parameters, Code node source, and pinned data for
// hardcoded secret material. Values are redacted at detection time: results
// carry workflow, node, and field path only — never the matched value.
type SecretsCheck struct{}

func (c *SecretsCheck) ID() string   { return "n8n:secrets" }
func (c *SecretsCheck) Name() string { return "Hardcoded Secrets" }

func (c *SecretsCheck) Run(ctx context.Context, path string) (checker.Result, error) {
	bundle, early := loadBundle(c, path)
	if early != nil {
		return *early, nil
	}
	_, degraded := splitDegraded(bundle)

	// Scan every workflow, degraded ones included: whatever parsed is
	// scannable, and more coverage only helps here.
	var confirmed, candidates []finding
	for _, sf := range n8n.ScanBundleSecrets(bundle) {
		f := finding{
			File:     sf.File,
			Workflow: sf.Workflow,
			Node:     sf.Node,
			Note:     fmt.Sprintf("%s at %s (value redacted)", sf.Rule, sf.FieldPath),
		}
		if sf.Confirmed {
			confirmed = append(confirmed, f)
		} else {
			candidates = append(candidates, f)
		}
	}

	var res checker.Result
	switch {
	case len(confirmed) > 0:
		detail := formatFindings(append(confirmed, candidates...))
		res = failResult(c,
			fmt.Sprintf("%d hardcoded secret(s) detected in workflow exports", len(confirmed)),
			detail)
	case len(candidates) > 0:
		res = warnResult(c,
			fmt.Sprintf("%d high-entropy candidate(s) in auth fields need human review", len(candidates)),
			formatFindings(candidates))
	default:
		res = passResult(c, "no hardcoded secrets detected in node parameters, code, or pinned data")
	}
	return finish(c, res, degraded), nil
}
