// Package n8ncheck implements the n8n:* checks: deterministic static
// analysis of exported n8n workflow JSON. Checks operate on the parsed
// bundle from pkg/n8n and never execute workflows or call the network.
//
// Degradation policy (never false-pass): a bundle with file-level parse
// errors makes every content check report an error status; a workflow with
// structural parse issues is excluded from evaluation and caps an otherwise
// passing check at warn.
package n8ncheck

import (
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/n8n"
)

// loadBundle loads the workflow bundle for a check. When the check cannot or
// need not proceed it returns a non-nil Result to return as-is: an Errored
// result for unreadable/partially unreadable bundles, or an Info result when
// the directory contains no n8n workflows at all.
func loadBundle(c checker.Checker, path string) (*n8n.Bundle, *checker.Result) {
	bundle, err := n8n.LoadDirCached(path)
	if err != nil {
		res := erroredResult(c, "could not scan for n8n workflows: "+err.Error())
		return nil, &res
	}
	if len(bundle.ParseErrors) > 0 {
		var lines []string
		for _, pe := range bundle.ParseErrors {
			lines = append(lines, pe.Path+" — "+pe.Message)
		}
		res := erroredResult(c, fmt.Sprintf("%d workflow file(s) could not be parsed; fix the export and re-run", len(bundle.ParseErrors)))
		res.RawOutput = strings.Join(lines, "\n")
		return nil, &res
	}
	if len(bundle.Workflows) == 0 {
		res := infoResult(c, "not applicable: no n8n workflow exports found")
		return nil, &res
	}
	return bundle, nil
}

// erroredResult builds a Result for a check a2 could not evaluate. The
// runner treats Errored as "surfaced but unscored", which is the safe
// degradation for unparseable input — never a silent pass.
func erroredResult(c checker.Checker, msg string) checker.Result {
	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   false,
		Status:   checker.Errored,
		Message:  msg,
		Reason:   msg,
		Language: checker.LangN8n,
	}
}

func infoResult(c checker.Checker, msg string) checker.Result {
	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   true,
		Status:   checker.Info,
		Message:  msg,
		Reason:   msg,
		Language: checker.LangN8n,
	}
}

func passResult(c checker.Checker, msg string) checker.Result {
	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   true,
		Status:   checker.Pass,
		Message:  msg,
		Reason:   msg,
		Language: checker.LangN8n,
	}
}

func warnResult(c checker.Checker, msg, detail string) checker.Result {
	return checker.Result{
		Name:      c.Name(),
		ID:        c.ID(),
		Passed:    false,
		Status:    checker.Warn,
		Message:   msg,
		Reason:    detail,
		RawOutput: detail,
		Language:  checker.LangN8n,
	}
}

func failResult(c checker.Checker, msg, detail string) checker.Result {
	return checker.Result{
		Name:      c.Name(),
		ID:        c.ID(),
		Passed:    false,
		Status:    checker.Fail,
		Message:   msg,
		Reason:    detail,
		RawOutput: detail,
		Language:  checker.LangN8n,
	}
}

// finding is one located issue. The rendered line format
// "<file> · <workflow> · <node> — <note>" is a stable contract: downstream
// tooling parses it, so changes must be treated as breaking.
type finding struct {
	File     string
	Workflow string
	Node     string // empty for workflow-level findings
	Note     string // empty when the line is self-explanatory
}

func (f finding) String() string {
	line := f.File + " · " + f.Workflow
	if f.Node != "" {
		line += " · " + f.Node
	}
	if f.Note != "" {
		line += " — " + f.Note
	}
	return line
}

// formatFindings renders findings one per line, in bundle order (bundle
// workflows are sorted by path, nodes keep export order), so output is
// deterministic for identical input.
func formatFindings(findings []finding) string {
	lines := make([]string, len(findings))
	for i, f := range findings {
		lines[i] = f.String()
	}
	return strings.Join(lines, "\n")
}

// splitDegraded separates cleanly parsed workflows from degraded ones.
// Degraded workflows cannot be trusted for evaluation: callers evaluate the
// clean set and pass the degraded set to finish().
func splitDegraded(b *n8n.Bundle) (clean, degraded []*n8n.WorkflowFile) {
	for _, wf := range b.Workflows {
		if wf.Workflow.Degraded() {
			degraded = append(degraded, wf)
		} else {
			clean = append(clean, wf)
		}
	}
	return clean, degraded
}

// finish applies the degradation policy to a computed result: when degraded
// workflows exist, a passing result is capped at Warn (parser uncertainty is
// never a pass) and the degraded files are appended to the detail in every
// case.
func finish(c checker.Checker, res checker.Result, degraded []*n8n.WorkflowFile) checker.Result {
	if len(degraded) == 0 {
		return res
	}
	var lines []string
	for _, wf := range degraded {
		lines = append(lines, fmt.Sprintf("%s · %s — could not be fully parsed: %s",
			wf.Path, workflowLabel(wf), strings.Join(wf.Workflow.ParseIssues, "; ")))
	}
	note := fmt.Sprintf("%d workflow(s) could not be fully parsed and were not evaluated", len(degraded))
	if res.Status == checker.Pass {
		res = warnResult(c, res.Message+"; "+note, strings.Join(lines, "\n"))
	} else {
		res.Message += "; " + note
		res.Reason = strings.TrimSpace(res.Reason + "\n" + strings.Join(lines, "\n"))
		res.RawOutput = strings.TrimSpace(res.RawOutput + "\n" + strings.Join(lines, "\n"))
	}
	return res
}

// workflowLabel names a workflow for output, falling back to the file path.
func workflowLabel(wf *n8n.WorkflowFile) string {
	if wf.Workflow.Name != "" {
		return wf.Workflow.Name
	}
	return wf.Path
}

// renderPath renders the main-flow path between two nodes as
// "A → B → C" — the readable injection/exposure path shown in findings.
func renderPath(g *n8n.Graph, from, to string) string {
	path := g.PathTo(from, to, n8n.ConnMain)
	if path == nil {
		return from + " → … → " + to
	}
	return strings.Join(path, " → ")
}

// isExpression reports whether an n8n parameter value is an expression
// (starts with "=") rather than a literal.
func isExpression(value string) bool {
	return strings.HasPrefix(value, "=")
}

// pluralWorkflows is a tiny formatting helper for messages.
func pluralWorkflows(n int) string {
	if n == 1 {
		return "workflow"
	}
	return "workflows"
}
