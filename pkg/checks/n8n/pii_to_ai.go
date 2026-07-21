package n8ncheck

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/n8n"
)

// fullPassthroughPattern spots whole-payload handoffs into an AI node:
// {{ $json }} or JSON.stringify($json) sends everything, minimized nothing.
var fullPassthroughPattern = regexp.MustCompile(`\{\{\s*\$json\s*\}\}|JSON\.stringify\(\s*\$json\s*\)`)

// piiFieldPattern spots expressions mapping PII-suggestive fields into AI
// inputs. These are field NAMES (heuristic), safe to print.
var piiFieldPattern = regexp.MustCompile(`\$json(?:\.\w+)*\.(?i)(email|e_mail|phone|mobile|address|dob|date_of_birth|birthdate|salary|iban|ssn|social_security|passport|national_id|first_name|last_name|full_name|surname)\b`)

// PiiToAICheck flags AI node inputs that carry PII-suggestive fields or the
// full $json payload to an external model provider without evidence of
// minimization. Field names are heuristic; findings seed a review, so this
// check warns.
type PiiToAICheck struct{}

func (c *PiiToAICheck) ID() string   { return "n8n:pii-to-ai" }
func (c *PiiToAICheck) Name() string { return "PII Into AI Inputs" }

func (c *PiiToAICheck) Run(path string) (checker.Result, error) {
	bundle, early := loadBundle(c, path)
	if early != nil {
		return *early, nil
	}
	clean, degraded := splitDegraded(bundle)

	var aiNodes int
	var findings []finding
	for _, wf := range clean {
		for i := range wf.Workflow.Nodes {
			node := &wf.Workflow.Nodes[i]
			if node.Disabled || !n8n.IsAIProcessingNode(node) {
				continue
			}
			aiNodes++
			var passthroughPaths []string
			piiFields := map[string]bool{}
			node.WalkParams(func(paramPath, _, value string) {
				if !strings.Contains(value, "{{") && !strings.Contains(value, "$json") {
					return
				}
				if fullPassthroughPattern.MatchString(value) {
					passthroughPaths = append(passthroughPaths, paramPath)
				}
				for _, m := range piiFieldPattern.FindAllStringSubmatch(value, -1) {
					piiFields[strings.ToLower(m[1])] = true
				}
			})
			for _, p := range passthroughPaths {
				findings = append(findings, finding{
					File:     wf.Path,
					Workflow: workflowLabel(wf),
					Node:     node.Name,
					Note:     "full $json passthrough into AI input at " + p + " — no field minimization",
				})
			}
			if len(piiFields) > 0 {
				fields := make([]string, 0, len(piiFields))
				for f := range piiFields {
					fields = append(fields, f)
				}
				sort.Strings(fields)
				findings = append(findings, finding{
					File:     wf.Path,
					Workflow: workflowLabel(wf),
					Node:     node.Name,
					Note:     "PII-suggestive fields mapped into AI input: " + strings.Join(fields, ", "),
				})
			}
		}
	}

	if aiNodes == 0 {
		return finish(c, infoResult(c, "not applicable: no AI nodes"), degraded), nil
	}

	var res checker.Result
	if len(findings) == 0 {
		res = passResult(c, fmt.Sprintf("no PII-suggestive fields or full payloads mapped into %d AI node(s)", aiNodes))
	} else {
		res = warnResult(c,
			fmt.Sprintf("%d AI input(s) carry PII-suggestive fields or the full payload", len(findings)),
			formatFindings(findings))
	}
	return finish(c, res, degraded), nil
}
