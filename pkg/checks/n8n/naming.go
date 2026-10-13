package n8ncheck

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
)

// n8nDefaultNames are default names n8n assigns to workflows and freshly
// dropped nodes. A name matching one of these (optionally with a numeric
// suffix, e.g. "Webhook1") means nobody renamed it.
var n8nDefaultNames = []string{
	"My workflow",
	"Webhook",
	"HTTP Request",
	"Code",
	"Function",
	"Function Item",
	"Set",
	"Edit Fields",
	"If",
	"Switch",
	"Merge",
	"No Operation, do nothing",
	"Schedule Trigger",
	"Manual Trigger",
	"When clicking 'Test workflow'",
	"When clicking ‘Test workflow’", // typographic-quote variant
	"When Executed by Another Workflow",
	"AI Agent",
	"OpenAI Chat Model",
	"Email Trigger (IMAP)",
	"Send Email",
	"Respond to Webhook",
	"Slack",
	"Gmail",
	"Google Sheets",
	"Postgres",
	"MySQL",
	"Telegram",
	"Discord",
	"Notion",
	"Airtable",
}

// defaultNamePattern matches any default name with an optional numeric
// suffix ("Webhook", "Webhook1", "Code 2").
var defaultNamePattern = buildDefaultNamePattern()

func buildDefaultNamePattern() *regexp.Regexp {
	quoted := make([]string, len(n8nDefaultNames))
	for i, name := range n8nDefaultNames {
		quoted[i] = regexp.QuoteMeta(name)
	}
	return regexp.MustCompile(`^(?:` + strings.Join(quoted, "|") + `) ?\d*$`)
}

// NamingCheck measures workflow and node name compliance: against a
// configured convention regex when one is set, otherwise by detecting
// untouched n8n default names.
type NamingCheck struct {
	// Regex is the naming convention names must match (n8n.naming_regex).
	// Empty enables defaults-detection mode.
	Regex string
}

func (c *NamingCheck) ID() string   { return "n8n:naming" }
func (c *NamingCheck) Name() string { return "Naming Convention" }

func (c *NamingCheck) Run(ctx context.Context, path string) (checker.Result, error) {
	var convention *regexp.Regexp
	if c.Regex != "" {
		var err error
		convention, err = regexp.Compile(c.Regex)
		if err != nil {
			return erroredResult(c, "invalid n8n.naming_regex: "+err.Error()), nil
		}
	}

	bundle, early := loadBundle(c, path)
	if early != nil {
		return *early, nil
	}
	clean, degraded := splitDegraded(bundle)

	var total, violations int
	var findings []finding
	checkName := func(wfFile, wfLabel, node, name string) {
		if name == "" {
			return
		}
		total++
		var bad bool
		var note string
		if convention != nil {
			bad = !convention.MatchString(name)
			note = "does not match the configured naming convention"
		} else {
			bad = defaultNamePattern.MatchString(name)
			note = "untouched default name"
		}
		if bad {
			violations++
			findings = append(findings, finding{File: wfFile, Workflow: wfLabel, Node: node, Note: note})
		}
	}

	for _, wf := range clean {
		label := workflowLabel(wf)
		checkName(wf.Path, label, "", wf.Workflow.Name)
		for i := range wf.Workflow.Nodes {
			checkName(wf.Path, label, wf.Workflow.Nodes[i].Name, wf.Workflow.Nodes[i].Name)
		}
	}

	if total == 0 {
		return finish(c, infoResult(c, "not applicable: no named workflows or nodes"), degraded), nil
	}

	compliance := float64(total-violations) / float64(total) * 100
	mode := "defaults-detection (set n8n.naming_regex to check a real convention)"
	if convention != nil {
		mode = "convention " + c.Regex
	}

	var res checker.Result
	if violations == 0 {
		res = passResult(c, fmt.Sprintf("all %d names comply (%s)", total, mode))
	} else {
		res = warnResult(c,
			fmt.Sprintf("%d of %d names non-compliant (%.0f%% compliance, %s)", violations, total, compliance, mode),
			formatFindings(findings))
	}
	return finish(c, res, degraded), nil
}
