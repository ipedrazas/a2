package n8ncheck

import (
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/n8n"
)

// AgentToolsCheck emits the per-agent tool inventory with read/write
// classification. Whether each write was justified is a human judgment —
// this check produces the inventory that makes that conversation concrete,
// and warns when write-capable (or unclassifiable) tools are attached.
type AgentToolsCheck struct {
	classifier *n8n.Classifier
}

func (c *AgentToolsCheck) ID() string   { return "n8n:agent-tools" }
func (c *AgentToolsCheck) Name() string { return "Agent Tool Inventory" }

func (c *AgentToolsCheck) Run(path string) (checker.Result, error) {
	bundle, early := loadBundle(c, path)
	if early != nil {
		return *early, nil
	}
	clean, degraded := splitDegraded(bundle)

	var agents, agentsWithWrite, agentsWithUnknown int
	var inventory []string
	for _, wf := range clean {
		g := n8n.BuildGraph(wf.Workflow)
		for i := range wf.Workflow.Nodes {
			agent := &wf.Workflow.Nodes[i]
			if agent.Disabled || !n8n.IsAgentNode(agent) {
				continue
			}
			agents++
			tools := g.AttachedTools(agent.Name)
			var lines []string
			var writes, unknowns int
			for _, tool := range tools {
				capability, reason := c.classifier.ClassifyToolCapability(tool)
				switch capability {
				case n8n.CapWrite:
					writes++
				case n8n.CapUnknown:
					unknowns++
				}
				lines = append(lines, fmt.Sprintf("  %s — %s — %s (%s)",
					tool.Name, shortNodeType(tool.Type), capability, reason))
			}
			if writes > 0 {
				agentsWithWrite++
			}
			if unknowns > 0 {
				agentsWithUnknown++
			}
			header := fmt.Sprintf("%s · %s · %s — %d tool(s): %d write-capable, %d unknown, %d read-only",
				wf.Path, workflowLabel(wf), agent.Name, len(tools), writes, unknowns, len(tools)-writes-unknowns)
			if len(tools) == 0 {
				header = fmt.Sprintf("%s · %s · %s — no tools attached", wf.Path, workflowLabel(wf), agent.Name)
			}
			inventory = append(inventory, header)
			inventory = append(inventory, lines...)
		}
	}

	if agents == 0 {
		return finish(c, infoResult(c, "not applicable: no AI agent nodes"), degraded), nil
	}

	detail := strings.Join(inventory, "\n")
	var res checker.Result
	switch {
	case agentsWithWrite > 0:
		res = warnResult(c,
			fmt.Sprintf("%d of %d agent(s) hold write-capable tools — review each write against its blast radius", agentsWithWrite, agents),
			detail)
	case agentsWithUnknown > 0:
		res = warnResult(c,
			fmt.Sprintf("%d of %d agent(s) hold tools of unknown capability — classify them (n8n.write_tools_extra)", agentsWithUnknown, agents),
			detail)
	default:
		res = passResult(c, fmt.Sprintf("%d agent(s), read-only tools throughout", agents))
		res.RawOutput = detail
	}
	return finish(c, res, degraded), nil
}
