package n8ncheck

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/n8n"
)

// n8n's built-in retry defaults; a whole estate sitting exactly on these
// values suggests nobody consciously chose a retry policy.
const (
	defaultMaxTries         = 3
	defaultWaitBetweenTries = 1000
)

// RetryCheck verifies that external-call nodes (HTTP, databases, SaaS
// integrations) have retry configured, and flags estates where every retry
// uses the untouched defaults.
type RetryCheck struct{}

func (c *RetryCheck) ID() string   { return "n8n:retry-config" }
func (c *RetryCheck) Name() string { return "Retry Configuration" }

func (c *RetryCheck) Run(path string) (checker.Result, error) {
	bundle, early := loadBundle(c, path)
	if early != nil {
		return *early, nil
	}
	clean, degraded := splitDegraded(bundle)

	var external, withRetry, allDefaults int
	var findings []finding
	distribution := map[string]int{}
	for _, wf := range clean {
		for i := range wf.Workflow.Nodes {
			node := &wf.Workflow.Nodes[i]
			if node.Disabled || !n8n.IsExternalCallNode(node) {
				continue
			}
			external++
			if !node.RetryOnFail {
				findings = append(findings, finding{
					File:     wf.Path,
					Workflow: workflowLabel(wf),
					Node:     node.Name,
					Note:     "no retry configured",
				})
				continue
			}
			withRetry++
			maxTries, wait := node.MaxTries, node.WaitBetweenTries
			if maxTries == 0 {
				maxTries = defaultMaxTries
			}
			if wait == 0 {
				wait = defaultWaitBetweenTries
			}
			if maxTries == defaultMaxTries && wait == defaultWaitBetweenTries {
				allDefaults++
			}
			distribution[fmt.Sprintf("maxTries=%d wait=%dms", maxTries, wait)]++
		}
	}

	if external == 0 {
		return finish(c, infoResult(c, "not applicable: no external-call nodes"), degraded), nil
	}

	detail := formatFindings(findings)
	if len(distribution) > 0 {
		var lines []string
		for k, v := range distribution {
			lines = append(lines, fmt.Sprintf("%s × %d", k, v))
		}
		sort.Strings(lines)
		detail = strings.TrimSpace(detail + "\n\nretry distribution:\n" + strings.Join(lines, "\n"))
	}

	var res checker.Result
	switch {
	case len(findings) > 0:
		res = warnResult(c,
			fmt.Sprintf("%d of %d external-call nodes have no retry configured", len(findings), external),
			detail)
	case withRetry >= 3 && allDefaults == withRetry:
		res = warnResult(c,
			fmt.Sprintf("all %d retry configurations use the untouched defaults (maxTries=%d, wait=%dms) — nobody chose a retry policy",
				withRetry, defaultMaxTries, defaultWaitBetweenTries),
			detail)
	default:
		res = passResult(c, fmt.Sprintf("all %d external-call nodes have retry configured", external))
	}
	return finish(c, res, degraded), nil
}
