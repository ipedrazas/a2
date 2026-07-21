package n8ncheck

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/config"
)

const (
	cleanDir       = "../../n8n/testdata/clean"
	vulnerableDir  = "../../n8n/testdata/vulnerable"
	mitigatedDir   = "../../n8n/testdata/mitigated"
	malformedDir   = "../../n8n/testdata/malformed"
	duplicationDir = "../../n8n/testdata/duplication"
)

// allCheckers builds every n8n check with default config.
func allCheckers(t *testing.T) []checker.Checker {
	t.Helper()
	regs := Register(config.DefaultConfig())
	checkers := make([]checker.Checker, len(regs))
	for i, reg := range regs {
		checkers[i] = reg.Checker
	}
	return checkers
}

func runCheck(t *testing.T, c checker.Checker, dir string) checker.Result {
	t.Helper()
	res, err := c.Run(dir)
	if err != nil {
		t.Fatalf("%s.Run(%s) returned error: %v", c.ID(), dir, err)
	}
	return res
}

// TestChecksAcrossFixtures is the per-check status matrix over the fixture
// estates.
func TestChecksAcrossFixtures(t *testing.T) {
	tests := []struct {
		check       checker.Checker
		dir         string
		wantStatus  checker.Status
		wantMessage string // substring of Message; empty = no assertion
		wantDetail  string // substring of Reason+RawOutput; empty = no assertion
	}{
		// clean estate: everything passes (or is N/A).
		{&ErrorWorkflowCheck{FailRatio: 0.5}, cleanDir, checker.Pass, "all 2 active workflows", ""},
		{&SilentFailuresCheck{}, cleanDir, checker.Pass, "", ""},
		{&RetryCheck{}, cleanDir, checker.Pass, "all 4 external-call nodes", ""},
		{&TimeoutsCheck{}, cleanDir, checker.Info, "not applicable", ""},
		{&CodeBulkCheck{LineLimit: 80}, cleanDir, checker.Pass, "no Code nodes", ""},
		{&PinnedDataCheck{}, cleanDir, checker.Pass, "", ""},
		{&NamingCheck{}, cleanDir, checker.Pass, "comply", ""},
		{&InlineAuthCheck{}, cleanDir, checker.Pass, "", ""},
		{&SecretsCheck{}, cleanDir, checker.Pass, "no hardcoded secrets", ""},
		{&DuplicationCheck{MinSequence: 5}, cleanDir, checker.Pass, "", ""},

		// vulnerable estate: the demo failures.
		{&ErrorWorkflowCheck{FailRatio: 0.5}, vulnerableDir, checker.Fail, "2 of 3 active workflows have no error workflow", "lead-webhook.json · Lead Capture"},
		{&SilentFailuresCheck{}, vulnerableDir, checker.Warn, "1 node(s) swallow failures", "Send Email"},
		{&RetryCheck{}, vulnerableDir, checker.Warn, "no retry configured", "HTTP Request"},
		{&TimeoutsCheck{}, vulnerableDir, checker.Warn, "1 of 1 HTTP Request nodes", "lead-webhook.json · Lead Capture · HTTP Request"},
		{&CodeBulkCheck{LineLimit: 80}, vulnerableDir, checker.Warn, "Code node issue", "122 lines (limit 80)"},
		{&CodeBulkCheck{LineLimit: 80}, vulnerableDir, checker.Warn, "", "embeds HTTP calls"},
		{&PinnedDataCheck{}, vulnerableDir, checker.Warn, "1 active workflow carry pinned test data", "pinned test data on 1 node(s): Webhook"},
		{&NamingCheck{}, vulnerableDir, checker.Warn, "non-compliant", "untouched default name"},
		{&InlineAuthCheck{}, vulnerableDir, checker.Warn, "2 HTTP parameter(s)", `literal value in header "Authorization"`},
		{&SecretsCheck{}, vulnerableDir, checker.Fail, "2 hardcoded secret(s)", "value redacted"},
		{&SecretsCheck{}, vulnerableDir, checker.Fail, "", "lead-webhook.json · Lead Capture · HTTP Request — openai-style key"},

		// duplication estate: one copied sequence, reported once.
		{&DuplicationCheck{MinSequence: 5}, duplicationDir, checker.Warn, "1 duplicated node sequence(s)", "httpRequest → code → if → set → slack"},
		{&DuplicationCheck{MinSequence: 5}, duplicationDir, checker.Warn, "", "invoice-alerts.json · Invoice Alerts"},

		// mitigated estate: the vulnerable shapes, fixed.
		{&ErrorWorkflowCheck{FailRatio: 0.5}, mitigatedDir, checker.Pass, "", ""},
		{&TimeoutsCheck{}, mitigatedDir, checker.Pass, "set a timeout", ""},
		{&RetryCheck{}, mitigatedDir, checker.Pass, "", ""},
		{&PinnedDataCheck{}, mitigatedDir, checker.Pass, "", ""},
		{&InlineAuthCheck{}, mitigatedDir, checker.Pass, "", ""},
		{&SilentFailuresCheck{}, mitigatedDir, checker.Pass, "", ""},
		{&SecretsCheck{}, mitigatedDir, checker.Pass, "", ""},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s/%s", tt.check.ID(), filepath.Base(tt.dir)), func(t *testing.T) {
			res := runCheck(t, tt.check, tt.dir)
			if res.Status != tt.wantStatus {
				t.Errorf("status = %v, want %v (message: %s / reason: %s)", res.Status, tt.wantStatus, res.Message, res.Reason)
			}
			if tt.wantMessage != "" && !strings.Contains(res.Message, tt.wantMessage) {
				t.Errorf("message %q does not contain %q", res.Message, tt.wantMessage)
			}
			if tt.wantDetail != "" && !strings.Contains(res.Reason+res.RawOutput, tt.wantDetail) {
				t.Errorf("detail does not contain %q:\n%s\n%s", tt.wantDetail, res.Reason, res.RawOutput)
			}
		})
	}
}

// TestMalformedBundleDegradesToError: file-level parse errors must degrade
// every content check to an error status — never a silent pass.
func TestMalformedBundleDegradesToError(t *testing.T) {
	for _, c := range allCheckers(t) {
		if c.ID() == "n8n:in-git" {
			continue // in-git is about files, not content; covered separately
		}
		t.Run(c.ID(), func(t *testing.T) {
			res := runCheck(t, c, malformedDir)
			if res.Status != checker.Errored {
				t.Errorf("status = %v, want ERROR (message: %s)", res.Status, res.Message)
			}
			if !strings.Contains(res.RawOutput, "truncated.json") {
				t.Errorf("detail should name the unparseable file, got: %s", res.RawOutput)
			}
		})
	}
}

// TestDegradedWorkflowCapsPass: a workflow with parse issues must cap an
// otherwise passing check at Warn.
func TestDegradedWorkflowCapsPass(t *testing.T) {
	dir := t.TempDir()
	copyFixture(t, "../../n8n/testdata/clean/crm-sync.json", filepath.Join(dir, "crm-sync.json"))
	copyFixture(t, "../../n8n/testdata/malformed/missing-connections.json", filepath.Join(dir, "orphan.json"))

	res := runCheck(t, &ErrorWorkflowCheck{FailRatio: 0.5}, dir)
	if res.Status != checker.Warn {
		t.Fatalf("status = %v, want WARN (message: %s)", res.Status, res.Message)
	}
	if !strings.Contains(res.Message, "could not be fully parsed") {
		t.Errorf("message should mention unparsed workflows, got: %s", res.Message)
	}
	if !strings.Contains(res.RawOutput, "orphan.json") {
		t.Errorf("detail should name the degraded file, got: %s", res.RawOutput)
	}
}

// TestNoSecretValuesInResults: the fake secrets planted in the vulnerable
// fixture must never appear in any check's serialized result.
func TestNoSecretValuesInResults(t *testing.T) {
	secrets := []string{"sk-FAKE", "AKIAFAKE"}
	for _, c := range allCheckers(t) {
		res := runCheck(t, c, vulnerableDir)
		serialized := fmt.Sprintf("%+v", res)
		for _, secret := range secrets {
			if strings.Contains(serialized, secret) {
				t.Errorf("%s leaks secret material (%s...) in its result", c.ID(), secret)
			}
		}
	}
}

func TestInGit(t *testing.T) {
	if !checkutil.ToolAvailable("git") {
		t.Skip("git not available")
	}

	t.Run("tracked fixtures pass", func(t *testing.T) {
		res := runCheck(t, &InGitCheck{}, cleanDir)
		if res.Status != checker.Pass {
			t.Errorf("status = %v, want PASS (message: %s)", res.Status, res.Message)
		}
	})

	t.Run("outside a repository warns", func(t *testing.T) {
		dir := t.TempDir()
		copyFixture(t, "../../n8n/testdata/clean/crm-sync.json", filepath.Join(dir, "crm-sync.json"))
		res := runCheck(t, &InGitCheck{}, dir)
		// A temp dir may still sit inside some parent repository; both
		// outcomes are a warn (untracked or no repo), never a pass.
		if res.Status != checker.Warn {
			t.Errorf("status = %v, want WARN (message: %s)", res.Status, res.Message)
		}
	})

	t.Run("malformed files still counted", func(t *testing.T) {
		res := runCheck(t, &InGitCheck{}, malformedDir)
		if res.Status != checker.Pass {
			t.Errorf("status = %v, want PASS — parse errors are irrelevant to version control (message: %s)", res.Status, res.Message)
		}
		if !strings.Contains(res.Message, "4 workflow file(s)") {
			t.Errorf("expected all 4 files counted, got: %s", res.Message)
		}
	})
}

func TestNamingWithConventionRegex(t *testing.T) {
	// Clean fixture names are Title Case; a lowercase-only convention makes
	// them all violations.
	res := runCheck(t, &NamingCheck{Regex: `^[a-z0-9-]+$`}, cleanDir)
	if res.Status != checker.Warn {
		t.Fatalf("status = %v, want WARN", res.Status)
	}
	if !strings.Contains(res.Message, "0% compliance") {
		t.Errorf("expected 0%% compliance, got: %s", res.Message)
	}

	invalid := runCheck(t, &NamingCheck{Regex: `([`}, cleanDir)
	if invalid.Status != checker.Errored {
		t.Errorf("invalid regex should error, got %v", invalid.Status)
	}
}

func TestRetryUniformDefaults(t *testing.T) {
	dir := t.TempDir()
	// Three external-call nodes, all retrying with untouched defaults.
	wf := `{
		"name": "Uniform Retries",
		"active": true,
		"nodes": [
			{"name": "Fetch A", "type": "n8n-nodes-base.httpRequest", "typeVersion": 4, "parameters": {"url": "https://a.test"}, "retryOnFail": true},
			{"name": "Fetch B", "type": "n8n-nodes-base.httpRequest", "typeVersion": 4, "parameters": {"url": "https://b.test"}, "retryOnFail": true, "maxTries": 3, "waitBetweenTries": 1000},
			{"name": "Fetch C", "type": "n8n-nodes-base.httpRequest", "typeVersion": 4, "parameters": {"url": "https://c.test"}, "retryOnFail": true, "maxTries": 3}
		],
		"connections": {},
		"settings": {"errorWorkflow": "err-1"}
	}`
	if err := os.WriteFile(filepath.Join(dir, "uniform.json"), []byte(wf), 0o644); err != nil {
		t.Fatal(err)
	}

	res := runCheck(t, &RetryCheck{}, dir)
	if res.Status != checker.Warn {
		t.Fatalf("status = %v, want WARN (message: %s)", res.Status, res.Message)
	}
	if !strings.Contains(res.Message, "untouched defaults") {
		t.Errorf("expected uniform-defaults warning, got: %s", res.Message)
	}
}

func copyFixture(t *testing.T, src, dst string) {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		t.Fatal(err)
	}
}
