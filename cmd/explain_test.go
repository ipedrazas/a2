package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
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

func TestRunExplain_NoMatch(t *testing.T) {
	captureStdout(t, func() {
		if err := runExplain(explainCmd, []string{"nonexistent:*"}); err == nil {
			t.Error("expected error for non-matching pattern, got nil")
		}
	})
}
