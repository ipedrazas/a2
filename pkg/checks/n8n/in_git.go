package n8ncheck

import (
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/n8n"
)

// InGitCheck verifies the workflow exports live in git: automation that only
// exists inside the live n8n instance has no history, no review, and no way
// back. Unlike the content checks, unparseable files still count here — a
// broken export in git is still versioned.
type InGitCheck struct{}

func (c *InGitCheck) ID() string   { return "n8n:in-git" }
func (c *InGitCheck) Name() string { return "Workflows Versioned In Git" }

func (c *InGitCheck) Run(path string) (checker.Result, error) {
	bundle, err := n8n.LoadDirCached(path)
	if err != nil {
		return erroredResult(c, "could not scan for n8n workflows: "+err.Error()), nil
	}
	var files []string
	for _, wf := range bundle.Workflows {
		files = append(files, wf.Path)
	}
	for _, pe := range bundle.ParseErrors {
		files = append(files, pe.Path)
	}
	if len(files) == 0 {
		return infoResult(c, "not applicable: no n8n workflow exports found"), nil
	}

	rb := checkutil.NewResultBuilder(c, checker.LangN8n)
	if !checkutil.ToolAvailable("git") {
		return rb.ToolNotInstalled("git", "install git to verify workflow version control"), nil
	}

	inside := rb.RunCommand(path, "git", "rev-parse", "--is-inside-work-tree")
	if !inside.Success() {
		return warnResult(c,
			fmt.Sprintf("%d workflow file(s) are not under version control (not a git repository)", len(files)),
			"exports live outside any git repository — commit them to a repo so changes have history and review"), nil
	}

	tracked := rb.RunCommand(path, "git", append([]string{"ls-files", "--"}, files...)...)
	if !tracked.Success() {
		return erroredResult(c, "git ls-files failed: "+strings.TrimSpace(tracked.CombinedOutput())), nil
	}
	trackedSet := map[string]bool{}
	for line := range strings.SplitSeq(tracked.Output(), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			trackedSet[line] = true
		}
	}

	var untracked []string
	for _, f := range files {
		if !trackedSet[f] {
			untracked = append(untracked, f+" — not tracked in git")
		}
	}
	if len(untracked) > 0 {
		res := warnResult(c,
			fmt.Sprintf("%d of %d workflow file(s) are not tracked in git", len(untracked), len(files)),
			strings.Join(untracked, "\n"))
		res.Command = tracked.DisplayCommand()
		return res, nil
	}

	res := passResult(c, fmt.Sprintf("all %d workflow file(s) are tracked in git", len(files)))
	res.Command = tracked.DisplayCommand()
	return res, nil
}
