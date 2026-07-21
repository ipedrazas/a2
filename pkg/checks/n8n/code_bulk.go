package n8ncheck

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
)

// codeNodeTypes are the node types carrying user-written code, mapped to the
// parameter holding the source.
var codeNodeTypes = map[string][]string{
	"n8n-nodes-base.code":         {"jsCode", "pythonCode"},
	"n8n-nodes-base.function":     {"functionCode"},
	"n8n-nodes-base.functionItem": {"functionCode"},
}

// embeddedHTTPPattern spots network calls inside Code nodes: business logic
// and integration entangled in one untestable blob.
var embeddedHTTPPattern = regexp.MustCompile(`\bfetch\s*\(|\baxios\b|http\.request|XMLHttpRequest|\burllib\b|\brequests\.(get|post|put|patch|delete|request)\b`)

// CodeBulkCheck flags Code/Function nodes that exceed a line threshold or
// embed HTTP calls.
type CodeBulkCheck struct {
	// LineLimit is the line count above which a Code node is flagged
	// (default 80).
	LineLimit int
}

func (c *CodeBulkCheck) ID() string   { return "n8n:code-node-bulk" }
func (c *CodeBulkCheck) Name() string { return "Code Node Bulk" }

func (c *CodeBulkCheck) Run(path string) (checker.Result, error) {
	bundle, early := loadBundle(c, path)
	if early != nil {
		return *early, nil
	}
	clean, degraded := splitDegraded(bundle)

	var codeNodes int
	var findings []finding
	lineCounts := map[string]int{}
	for _, wf := range clean {
		for i := range wf.Workflow.Nodes {
			node := &wf.Workflow.Nodes[i]
			sourceParams, isCode := codeNodeTypes[node.Type]
			if !isCode || node.Disabled {
				continue
			}
			source := ""
			for _, param := range sourceParams {
				if s, ok := node.StringParam(param); ok && s != "" {
					source = s
					break
				}
			}
			codeNodes++
			lines := strings.Count(source, "\n") + 1
			lineCounts[fmt.Sprintf("%s · %s · %s", wf.Path, workflowLabel(wf), node.Name)] = lines
			if lines > c.LineLimit {
				findings = append(findings, finding{
					File:     wf.Path,
					Workflow: workflowLabel(wf),
					Node:     node.Name,
					Note:     fmt.Sprintf("%d lines (limit %d)", lines, c.LineLimit),
				})
			}
			if embeddedHTTPPattern.MatchString(source) {
				findings = append(findings, finding{
					File:     wf.Path,
					Workflow: workflowLabel(wf),
					Node:     node.Name,
					Note:     "embeds HTTP calls in code (use an HTTP Request node so retries/timeouts apply)",
				})
			}
		}
	}

	if codeNodes == 0 {
		return finish(c, passResult(c, "no Code nodes"), degraded), nil
	}

	detail := formatFindings(findings)
	var countLines []string
	for k, v := range lineCounts {
		countLines = append(countLines, fmt.Sprintf("%s — %d lines", k, v))
	}
	sort.Strings(countLines)
	detail = strings.TrimSpace(detail + "\n\nper-node line counts:\n" + strings.Join(countLines, "\n"))

	var res checker.Result
	if len(findings) == 0 {
		res = passResult(c, fmt.Sprintf("all %d Code nodes are under %d lines with no embedded HTTP calls", codeNodes, c.LineLimit))
	} else {
		res = warnResult(c,
			fmt.Sprintf("%d Code node issue(s): oversized nodes or embedded HTTP calls", len(findings)),
			detail)
	}
	return finish(c, res, degraded), nil
}
