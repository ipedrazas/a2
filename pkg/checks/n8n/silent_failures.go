package n8ncheck

import (
	"fmt"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/n8n"
)

// SilentFailuresCheck finds nodes whose failures are swallowed inline:
// errors merged into the regular output (continueOnFail /
// onError=continueRegularOutput), or routed to an error output that is not
// connected to anything.
type SilentFailuresCheck struct{}

func (c *SilentFailuresCheck) ID() string   { return "n8n:silent-failures" }
func (c *SilentFailuresCheck) Name() string { return "Silent Failure Settings" }

func (c *SilentFailuresCheck) Run(path string) (checker.Result, error) {
	bundle, early := loadBundle(c, path)
	if early != nil {
		return *early, nil
	}
	clean, degraded := splitDegraded(bundle)

	var findings []finding
	for _, wf := range clean {
		g := n8n.BuildGraph(wf.Workflow)
		for i := range wf.Workflow.Nodes {
			node := &wf.Workflow.Nodes[i]
			if node.Disabled {
				continue
			}
			switch {
			case node.OnError == "continueRegularOutput" || (node.ContinueOnFail && node.OnError == ""):
				findings = append(findings, finding{
					File:     wf.Path,
					Workflow: workflowLabel(wf),
					Node:     node.Name,
					Note:     "failures continue into the regular output (no error branch)",
				})
			case node.OnError == "continueErrorOutput" && !hasErrorBranch(g, node.Name):
				findings = append(findings, finding{
					File:     wf.Path,
					Workflow: workflowLabel(wf),
					Node:     node.Name,
					Note:     "error output is not connected to anything",
				})
			}
		}
	}

	var res checker.Result
	if len(findings) == 0 {
		res = passResult(c, "no nodes swallow failures inline")
	} else {
		res = warnResult(c,
			fmt.Sprintf("%d node(s) swallow failures without a connected error branch", len(findings)),
			formatFindings(findings))
	}
	return finish(c, res, degraded), nil
}

// hasErrorBranch reports whether the node has an outgoing main connection
// from a non-zero output port (where the error output lives for nodes with
// onError=continueErrorOutput).
func hasErrorBranch(g *n8n.Graph, nodeName string) bool {
	for _, e := range g.OutEdges(nodeName) {
		if e.Type == n8n.ConnMain && e.OutputIndex > 0 {
			return true
		}
	}
	return false
}
