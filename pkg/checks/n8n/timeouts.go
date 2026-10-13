package n8ncheck

import (
	"context"
	"fmt"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/n8n"
)

// callerWaitingTriggers are trigger types where a caller is synchronously
// waiting on the workflow, making unbounded HTTP calls a hang risk.
var callerWaitingTriggers = map[string]bool{
	"n8n-nodes-base.webhook":               true,
	"n8n-nodes-base.formTrigger":           true,
	"@n8n/n8n-nodes-langchain.chatTrigger": true,
}

// TimeoutsCheck flags HTTP Request nodes without an explicit timeout in
// workflows whose trigger implies someone is waiting (webhook/form/chat).
type TimeoutsCheck struct{}

func (c *TimeoutsCheck) ID() string   { return "n8n:timeouts" }
func (c *TimeoutsCheck) Name() string { return "HTTP Timeout Discipline" }

func (c *TimeoutsCheck) Run(ctx context.Context, path string) (checker.Result, error) {
	bundle, early := loadBundle(c, path)
	if early != nil {
		return *early, nil
	}
	clean, degraded := splitDegraded(bundle)

	var inScope int
	var findings []finding
	for _, wf := range clean {
		if !hasCallerWaitingTrigger(wf) {
			continue
		}
		for i := range wf.Workflow.Nodes {
			node := &wf.Workflow.Nodes[i]
			if node.Disabled || !n8n.IsHTTPRequestNode(node) {
				continue
			}
			inScope++
			if !hasExplicitTimeout(node) {
				findings = append(findings, finding{
					File:     wf.Path,
					Workflow: workflowLabel(wf),
					Node:     node.Name,
					Note:     "no explicit timeout while a caller is waiting on this workflow",
				})
			}
		}
	}

	if inScope == 0 {
		return finish(c, infoResult(c, "not applicable: no HTTP Request nodes in caller-waiting workflows"), degraded), nil
	}

	var res checker.Result
	if len(findings) == 0 {
		res = passResult(c, fmt.Sprintf("all %d HTTP Request nodes in caller-waiting workflows set a timeout", inScope))
	} else {
		res = warnResult(c,
			fmt.Sprintf("%d of %d HTTP Request nodes in caller-waiting workflows have no timeout", len(findings), inScope),
			formatFindings(findings))
	}
	return finish(c, res, degraded), nil
}

func hasCallerWaitingTrigger(wf *n8n.WorkflowFile) bool {
	for i := range wf.Workflow.Nodes {
		if callerWaitingTriggers[wf.Workflow.Nodes[i].Type] {
			return true
		}
	}
	return false
}

// hasExplicitTimeout reports whether the node's options.timeout is set (a
// positive number, or an expression — an expression means someone thought
// about it).
func hasExplicitTimeout(node *n8n.Node) bool {
	options, ok := node.Parameters["options"].(map[string]any)
	if !ok {
		return false
	}
	switch v := options["timeout"].(type) {
	case float64:
		return v > 0
	case string:
		return isExpression(v)
	default:
		return false
	}
}
