package n8ncheck

import (
	"context"
	"fmt"
	"regexp"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/n8n"
)

// authParamNamePattern matches header/query parameter names that carry
// credentials. Findings report the parameter name and location only — never
// the value.
var authParamNamePattern = regexp.MustCompile(`(?i)(auth|token|secret|api[-_]?key|access[-_]?key|password|bearer|credential)`)

// InlineAuthCheck flags HTTP nodes that carry literal auth material in
// header or query parameters instead of referencing an n8n credential. It is
// a weaker signal than n8n:secrets (the literal may be a placeholder), so it
// only warns.
type InlineAuthCheck struct{}

func (c *InlineAuthCheck) ID() string   { return "n8n:inline-auth" }
func (c *InlineAuthCheck) Name() string { return "Inline Auth Vs Credential Store" }

func (c *InlineAuthCheck) Run(ctx context.Context, path string) (checker.Result, error) {
	bundle, early := loadBundle(c, path)
	if early != nil {
		return *early, nil
	}
	clean, degraded := splitDegraded(bundle)

	var findings []finding
	for _, wf := range clean {
		for i := range wf.Workflow.Nodes {
			node := &wf.Workflow.Nodes[i]
			if node.Disabled || !isHTTPish(node) {
				continue
			}
			for _, section := range []string{"headerParameters", "queryParameters"} {
				for _, name := range literalAuthParams(node, section) {
					findings = append(findings, finding{
						File:     wf.Path,
						Workflow: workflowLabel(wf),
						Node:     node.Name,
						Note:     fmt.Sprintf("literal value in %s %q — use an n8n credential instead", sectionLabel(section), name),
					})
				}
			}
		}
	}

	var res checker.Result
	if len(findings) == 0 {
		res = passResult(c, "no HTTP nodes carry literal auth parameters")
	} else {
		res = warnResult(c,
			fmt.Sprintf("%d HTTP parameter(s) carry literal auth material instead of a credential reference", len(findings)),
			formatFindings(findings))
	}
	return finish(c, res, degraded), nil
}

func isHTTPish(node *n8n.Node) bool {
	return n8n.IsHTTPRequestNode(node) || node.Type == "@n8n/n8n-nodes-langchain.toolHttpRequest"
}

func sectionLabel(section string) string {
	if section == "queryParameters" {
		return "query parameter"
	}
	return "header"
}

// literalAuthParams returns the names of auth-looking parameters in the
// given section ({header,query}Parameters.parameters[]) whose value is a
// non-empty literal rather than an expression.
func literalAuthParams(node *n8n.Node, section string) []string {
	sectionMap, ok := node.Parameters[section].(map[string]any)
	if !ok {
		return nil
	}
	params, ok := sectionMap["parameters"].([]any)
	if !ok {
		return nil
	}
	var names []string
	for _, p := range params {
		entry, ok := p.(map[string]any)
		if !ok {
			continue
		}
		name, _ := entry["name"].(string)
		value, valueIsString := entry["value"].(string)
		if name == "" || !valueIsString || value == "" || isExpression(value) {
			continue
		}
		if authParamNamePattern.MatchString(name) {
			names = append(names, name)
		}
	}
	return names
}
