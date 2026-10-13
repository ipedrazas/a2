package n8ncheck

import (
	"context"
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/n8n"
)

// DuplicationCheck reports node sequences copied across workflows — each
// cluster is a sub-workflow candidate. It informs a refactoring
// conversation; it does not gate.
type DuplicationCheck struct {
	// MinSequence is the minimum shared sequence length reported
	// (default 5).
	MinSequence int
}

func (c *DuplicationCheck) ID() string   { return "n8n:duplication" }
func (c *DuplicationCheck) Name() string { return "Cross-Workflow Duplication" }

func (c *DuplicationCheck) Run(ctx context.Context, path string) (checker.Result, error) {
	bundle, early := loadBundle(c, path)
	if early != nil {
		return *early, nil
	}
	clean, degraded := splitDegraded(bundle)

	clusters := n8n.FindDuplicateSequences(clean, c.MinSequence)
	if len(clusters) == 0 {
		return finish(c, passResult(c, "no duplicated node sequences across workflows"), degraded), nil
	}

	var lines []string
	for _, cluster := range clusters {
		types := make([]string, len(cluster.Types))
		for i, t := range cluster.Types {
			types[i] = shortNodeType(t)
		}
		lines = append(lines, fmt.Sprintf("shared %d-node sequence (%.0f%% parameter-shape similarity): %s",
			len(cluster.Types), cluster.ParamSimilarity*100, strings.Join(types, " → ")))
		for _, occ := range cluster.Occurrences {
			lines = append(lines, fmt.Sprintf("  %s · %s · %s → … → %s",
				occ.File, occ.Workflow, occ.Nodes[0], occ.Nodes[len(occ.Nodes)-1]))
		}
	}

	res := warnResult(c,
		fmt.Sprintf("%d duplicated node sequence(s) across workflows — sub-workflow candidates", len(clusters)),
		strings.Join(lines, "\n"))
	return finish(c, res, degraded), nil
}

// shortNodeType strips the package prefixes for readable sequences.
func shortNodeType(t string) string {
	t = strings.TrimPrefix(t, "n8n-nodes-base.")
	t = strings.TrimPrefix(t, "@n8n/n8n-nodes-langchain.")
	return t
}
