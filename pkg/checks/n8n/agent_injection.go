package n8ncheck

import (
	"context"
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/n8n"
)

// AgentInjectionCheck traces untrusted input (webhooks, inbound email,
// forms, chat, RSS) through the graph to AI nodes. Untrusted text reaching
// an LLM is a prompt-injection path; if the LLM is an agent holding
// write-capable tools, the attacker's text can drive real-world actions —
// send email, write to the CRM, post messages — with the workflow's
// credentials.
type AgentInjectionCheck struct {
	classifier *n8n.Classifier
}

func (c *AgentInjectionCheck) ID() string   { return "n8n:agent-injection-path" }
func (c *AgentInjectionCheck) Name() string { return "Agent Injection Path" }

func (c *AgentInjectionCheck) Run(ctx context.Context, path string) (checker.Result, error) {
	bundle, early := loadBundle(c, path)
	if early != nil {
		return *early, nil
	}
	clean, degraded := splitDegraded(bundle)

	var critical, warnings []finding
	for _, wf := range clean {
		g := n8n.BuildGraph(wf.Workflow)
		for i := range wf.Workflow.Nodes {
			trigger := &wf.Workflow.Nodes[i]
			if trigger.Disabled || !c.classifier.IsUntrustedTrigger(trigger) {
				continue
			}
			for _, reached := range g.Reachable(trigger.Name, n8n.ConnMain) {
				target := g.Node(reached)
				if target == nil || target.Disabled || !n8n.IsAIProcessingNode(target) {
					continue
				}
				pathStr := renderPath(g, trigger.Name, reached)
				writeTools, unknownTools, readTools := c.classifyAgentTools(g, target)

				switch {
				case len(writeTools) > 0:
					critical = append(critical, finding{
						File:     wf.Path,
						Workflow: workflowLabel(wf),
						Note: fmt.Sprintf("untrusted input path: %s — agent holds write-capable tools: %s",
							pathStr, strings.Join(writeTools, ", ")),
					})
				case len(unknownTools) > 0:
					warnings = append(warnings, finding{
						File:     wf.Path,
						Workflow: workflowLabel(wf),
						Note: fmt.Sprintf("untrusted input path: %s — agent tools of unknown capability: %s",
							pathStr, strings.Join(unknownTools, ", ")),
					})
				case len(readTools) == 0:
					// A tool-less AI node: injection cannot act, but it can
					// still corrupt whatever consumes the output.
					warnings = append(warnings, finding{
						File:     wf.Path,
						Workflow: workflowLabel(wf),
						Note:     fmt.Sprintf("untrusted input path: %s — injection can corrupt this AI node's output", pathStr),
					})
				}
				// Read-only-tooled agents pass here; their inventory is
				// still surfaced by n8n:agent-tools.
			}
		}
	}

	var res checker.Result
	switch {
	case len(critical) > 0:
		res = failResult(c,
			fmt.Sprintf("%d untrusted-input path(s) reach agents with write-capable tools", len(critical)),
			formatFindings(append(critical, warnings...)))
	case len(warnings) > 0:
		res = warnResult(c,
			fmt.Sprintf("%d untrusted-input path(s) reach AI nodes", len(warnings)),
			formatFindings(warnings))
	default:
		res = passResult(c, "no untrusted input reaches AI nodes with write-capable tools")
	}
	return finish(c, res, degraded), nil
}

// classifyAgentTools splits an AI node's attached tools by capability,
// rendered as "Name (type, reason)" strings.
func (c *AgentInjectionCheck) classifyAgentTools(g *n8n.Graph, agent *n8n.Node) (write, unknown, read []string) {
	for _, tool := range g.AttachedTools(agent.Name) {
		capability, reason := c.classifier.ClassifyToolCapability(tool)
		rendered := fmt.Sprintf("%s (%s — %s)", tool.Name, shortNodeType(tool.Type), reason)
		switch capability {
		case n8n.CapWrite:
			write = append(write, rendered)
		case n8n.CapUnknown:
			unknown = append(unknown, rendered)
		default:
			read = append(read, rendered)
		}
	}
	return write, unknown, read
}
