package n8ncheck

import (
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/n8n"
)

// WebhookAuthCheck finds webhook triggers with no authentication and
// traverses the graph behind them: an open endpoint that can reach
// side-effectful nodes (sends, writes, payments) means anyone on the
// internet can trigger those effects.
type WebhookAuthCheck struct {
	classifier *n8n.Classifier
}

func (c *WebhookAuthCheck) ID() string   { return "n8n:webhook-auth" }
func (c *WebhookAuthCheck) Name() string { return "Unauthenticated Webhooks" }

func (c *WebhookAuthCheck) Run(path string) (checker.Result, error) {
	bundle, early := loadBundle(c, path)
	if early != nil {
		return *early, nil
	}
	clean, degraded := splitDegraded(bundle)

	var webhooks int
	var critical, warnings []finding
	for _, wf := range clean {
		g := n8n.BuildGraph(wf.Workflow)
		for i := range wf.Workflow.Nodes {
			node := &wf.Workflow.Nodes[i]
			if node.Disabled || node.Type != "n8n-nodes-base.webhook" {
				continue
			}
			webhooks++
			if auth, ok := node.StringParam("authentication"); ok && auth != "none" {
				continue
			}

			// Unauthenticated: what can a caller reach?
			var writePaths, unknownNodes []string
			for _, reached := range g.Reachable(node.Name, n8n.ConnMain) {
				target := g.Node(reached)
				if target == nil || target.Disabled {
					continue
				}
				effect, reason := c.classifier.ClassifySideEffect(target)
				switch effect {
				case n8n.EffectWrite:
					writePaths = append(writePaths, fmt.Sprintf("%s [%s]", renderPath(g, node.Name, reached), reason))
				case n8n.EffectUnknown:
					unknownNodes = append(unknownNodes, fmt.Sprintf("%s (%s)", reached, reason))
				}
			}

			switch {
			case len(writePaths) > 0:
				critical = append(critical, finding{
					File:     wf.Path,
					Workflow: workflowLabel(wf),
					Node:     node.Name,
					Note: "no authentication, side effects reachable:\n    " +
						strings.Join(writePaths, "\n    "),
				})
			case len(unknownNodes) > 0:
				warnings = append(warnings, finding{
					File:     wf.Path,
					Workflow: workflowLabel(wf),
					Node:     node.Name,
					Note:     "no authentication; reaches nodes of unknown effect: " + strings.Join(unknownNodes, ", "),
				})
			default:
				warnings = append(warnings, finding{
					File:     wf.Path,
					Workflow: workflowLabel(wf),
					Node:     node.Name,
					Note:     "no authentication (read-only path — still an open endpoint)",
				})
			}
		}
	}

	if webhooks == 0 {
		return finish(c, infoResult(c, "not applicable: no webhook triggers"), degraded), nil
	}

	var res checker.Result
	switch {
	case len(critical) > 0:
		res = failResult(c,
			fmt.Sprintf("%d unauthenticated webhook(s) can trigger external writes", len(critical)),
			formatFindings(append(critical, warnings...)))
	case len(warnings) > 0:
		res = warnResult(c,
			fmt.Sprintf("%d unauthenticated webhook(s) exposed (no side effects reachable)", len(warnings)),
			formatFindings(warnings))
	default:
		res = passResult(c, fmt.Sprintf("all %d webhook trigger(s) require authentication", webhooks))
	}
	return finish(c, res, degraded), nil
}
