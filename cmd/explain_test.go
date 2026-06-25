package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
)

// captureStdout runs fn and returns everything it wrote to os.Stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	fn()

	w.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("copy: %v", err)
	}
	return buf.String()
}

func TestRunExplain_ExactMatch(t *testing.T) {
	var out string
	err := func() error {
		var runErr error
		out = captureStdout(t, func() { runErr = runExplain(explainCmd, []string{"security:shell_injection"}) })
		return runErr
	}()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Check ID:     security:shell_injection") {
		t.Errorf("expected explanation for security:shell_injection, got:\n%s", out)
	}
	if strings.Count(out, "Check ID:") != 1 {
		t.Errorf("expected exactly one check, got:\n%s", out)
	}
}

func TestRunExplain_Wildcard(t *testing.T) {
	var out string
	var runErr error
	out = captureStdout(t, func() { runErr = runExplain(explainCmd, []string{"security:*"}) })
	if runErr != nil {
		t.Fatalf("unexpected error: %v", runErr)
	}
	// security:* must expand to more than one check, including shell_injection.
	if n := strings.Count(out, "Check ID:"); n < 2 {
		t.Errorf("expected multiple checks for security:*, got %d:\n%s", n, out)
	}
	if !strings.Contains(out, "security:shell_injection") {
		t.Errorf("expected security:shell_injection in wildcard output, got:\n%s", out)
	}
	// Every matched check must belong to the security namespace.
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "Check ID:") && !strings.Contains(line, "security:") {
			t.Errorf("wildcard matched a non-security check: %q", line)
		}
	}
}

// TestExplainCheck_CommandBased asserts a check that shells out renders its
// Command and not the scan-check "Where" pointer.
func TestExplainCheck_CommandBased(t *testing.T) {
	reg := checker.CheckRegistration{
		Meta: checker.CheckMeta{
			ID:        "go:vet",
			Name:      "Go Vet",
			Languages: []checker.Language{checker.LangGo},
			Command:   "go vet ./...",
		},
	}
	out := captureStdout(t, func() { explainCheck(reg) })

	if !strings.Contains(out, "Command:      go vet ./...") {
		t.Errorf("expected Command line, got:\n%s", out)
	}
	if strings.Contains(out, "Where:") {
		t.Errorf("command-based check must not render a Where pointer, got:\n%s", out)
	}
}

// TestExplainCheck_ScanCheck asserts a scan check (no Command) renders Detail,
// the 'a2 run <id>' pointer, and FixPrompt.
func TestExplainCheck_ScanCheck(t *testing.T) {
	reg := checker.CheckRegistration{
		Meta: checker.CheckMeta{
			ID:        "security:obfuscation",
			Name:      "Code Obfuscation Detection",
			Languages: []checker.Language{checker.LangCommon},
			Critical:  true,
			Detail:    "Flags base64-encoded payloads and eval-of-decoded-string patterns.",
			FixPrompt: "Review the flagged lines and rewrite any dynamically decoded code.",
		},
	}
	out := captureStdout(t, func() { explainCheck(reg) })

	if !strings.Contains(out, "Detail:") || !strings.Contains(out, "base64-encoded payloads") {
		t.Errorf("expected Detail line, got:\n%s", out)
	}
	if !strings.Contains(out, "a2 run security:obfuscation") {
		t.Errorf("expected 'a2 run' pointer, got:\n%s", out)
	}
	if !strings.Contains(out, "Fix prompt:") || !strings.Contains(out, "dynamically decoded code") {
		t.Errorf("expected Fix prompt line, got:\n%s", out)
	}
}

// TestExplainCheck_MultilineDetailIndented asserts multi-line Detail/FixPrompt
// fields are indented to align under their label.
func TestExplainCheck_MultilineDetailIndented(t *testing.T) {
	reg := checker.CheckRegistration{
		Meta: checker.CheckMeta{
			ID:        "security:network",
			Name:      "Network",
			Languages: []checker.Language{checker.LangCommon},
			Detail:    "line one\nline two",
		},
	}
	out := captureStdout(t, func() { explainCheck(reg) })

	if !strings.Contains(out, "\n              line two") {
		t.Errorf("expected continuation line indented under label, got:\n%s", out)
	}
}

func TestRunExplain_NoMatch(t *testing.T) {
	captureStdout(t, func() {
		if err := runExplain(explainCmd, []string{"nonexistent:*"}); err == nil {
			t.Error("expected error for non-matching pattern, got nil")
		}
	})
}
