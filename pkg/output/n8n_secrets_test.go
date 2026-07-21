package output_test

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ipedrazas/a2/pkg/checker"
	n8ncheck "github.com/ipedrazas/a2/pkg/checks/n8n"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/language"
	"github.com/ipedrazas/a2/pkg/output"
	"github.com/ipedrazas/a2/pkg/runner"
)

const vulnerableFixtures = "../n8n/testdata/vulnerable"

// TestN8nSecretsAbsentFromJSONOutput is the Phase 3 acceptance test for the
// n8n domain: run the full n8n suite over the vulnerable fixture estate,
// render it through the real `-f json` writer at maximum verbosity, and
// assert the planted fake secrets never appear while the secrets check
// itself fails.
func TestN8nSecretsAbsentFromJSONOutput(t *testing.T) {
	regs := n8ncheck.Register(config.DefaultConfig())
	// Parallel matches how `a2 check` runs (config default) and avoids the
	// sequential runner's veto-abort on the first critical failure, so every
	// check's output is present to scan.
	result := runner.RunSuiteWithOptions(vulnerableFixtures, regs, runner.RunSuiteOptions{
		Parallel: true,
		Timeout:  time.Minute,
	})

	detected := language.DetectionResult{
		Languages: []checker.Language{checker.LangN8n},
		Primary:   checker.LangN8n,
	}

	// output.JSON writes to stdout; capture it.
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	success, outputErr := output.JSON(result, detected, output.VerbosityAll)
	w.Close()
	os.Stdout = old
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if outputErr != nil {
		t.Fatalf("output.JSON: %v", outputErr)
	}
	if success {
		t.Error("vulnerable estate must not report success (n8n:secrets and n8n:error-workflow fail)")
	}

	out := string(data)
	for _, secret := range []string{"sk-FAKE", "AKIAFAKE"} {
		if strings.Contains(out, secret) {
			t.Errorf("JSON output leaks secret material %q", secret)
		}
	}
	if !strings.Contains(out, "n8n:secrets") {
		t.Error("JSON output missing the n8n:secrets result")
	}
	if !strings.Contains(out, "value redacted") {
		t.Error("JSON output should carry the redacted finding locations")
	}
}
