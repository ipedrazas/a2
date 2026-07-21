package n8n

import (
	"os"
	"path/filepath"
	"testing"
)

func loadFixtureDir(t *testing.T, dir string) *Bundle {
	t.Helper()
	b, err := LoadDir(dir)
	if err != nil {
		t.Fatalf("LoadDir(%s): %v", dir, err)
	}
	return b
}

func findWorkflow(t *testing.T, b *Bundle, path string) *Workflow {
	t.Helper()
	for _, wf := range b.Workflows {
		if wf.Path == path {
			return wf.Workflow
		}
	}
	t.Fatalf("workflow %s not found in bundle (have %d workflows)", path, len(b.Workflows))
	return nil
}

func TestLoadDirClean(t *testing.T) {
	b := loadFixtureDir(t, "testdata/clean")

	if len(b.Workflows) != 2 {
		t.Fatalf("expected 2 workflows, got %d", len(b.Workflows))
	}
	if len(b.ParseErrors) != 0 {
		t.Fatalf("expected no parse errors, got %v", b.ParseErrors)
	}

	wf := findWorkflow(t, b, "webhook-intake.json")
	if wf.Name != "Lead Intake" {
		t.Errorf("name = %q, want Lead Intake", wf.Name)
	}
	if !wf.Active {
		t.Error("expected workflow to be active")
	}
	if wf.Settings.ErrorWorkflow != "error-handler-wf-001" {
		t.Errorf("errorWorkflow = %q", wf.Settings.ErrorWorkflow)
	}
	if wf.Degraded() {
		t.Errorf("unexpected parse issues: %v", wf.ParseIssues)
	}
	if len(wf.PinData) != 0 {
		t.Errorf("expected empty pinData, got %d entries", len(wf.PinData))
	}

	webhook := wf.NodeByName("Lead Intake Webhook")
	if webhook == nil {
		t.Fatal("webhook node not found")
	}
	if auth, _ := webhook.StringParam("authentication"); auth != "headerAuth" {
		t.Errorf("webhook authentication = %q, want headerAuth", auth)
	}
	if len(webhook.Credentials) == 0 {
		t.Error("webhook should reference a credential")
	}

	save := wf.NodeByName("Save Lead")
	if save == nil {
		t.Fatal("Save Lead node not found")
	}
	if !save.RetryOnFail || save.MaxTries != 4 || save.WaitBetweenTries != 2000 {
		t.Errorf("retry settings = %v/%d/%d, want true/4/2000", save.RetryOnFail, save.MaxTries, save.WaitBetweenTries)
	}
	if save.OnError != "continueErrorOutput" {
		t.Errorf("onError = %q, want continueErrorOutput", save.OnError)
	}
}

func TestLoadDirVulnerable(t *testing.T) {
	b := loadFixtureDir(t, "testdata/vulnerable")

	if len(b.Workflows) != 4 {
		t.Fatalf("expected 4 workflows, got %d", len(b.Workflows))
	}
	if len(b.ParseErrors) != 0 {
		t.Fatalf("expected no parse errors, got %v", b.ParseErrors)
	}

	lead := findWorkflow(t, b, "lead-webhook.json")
	if lead.Settings.ErrorWorkflow != "" {
		t.Errorf("lead-webhook should have no error workflow, got %q", lead.Settings.ErrorWorkflow)
	}
	if _, ok := lead.PinData["Webhook"]; !ok {
		t.Error("lead-webhook should have pinData on the Webhook node")
	}
	if !lead.Active {
		t.Error("lead-webhook should be active")
	}
	webhook := lead.NodeByName("Webhook")
	if webhook == nil {
		t.Fatal("Webhook node not found")
	}
	if _, ok := webhook.StringParam("authentication"); ok {
		t.Error("vulnerable webhook should have no authentication parameter (defaults to none)")
	}

	bulk := findWorkflow(t, b, "bulk-code.json")
	send := bulk.NodeByName("Send Email")
	if send == nil || !send.ContinueOnFail {
		t.Error("bulk-code Send Email should have continueOnFail set")
	}

	report := findWorkflow(t, b, "daily-report.json")
	if report.Settings.ErrorWorkflow != "err-handler-wf-vuln" {
		t.Errorf("daily-report errorWorkflow = %q", report.Settings.ErrorWorkflow)
	}
}

func TestLoadDirMalformed(t *testing.T) {
	b := loadFixtureDir(t, "testdata/malformed")

	// truncated.json is broken JSON mentioning "nodes": a bundle parse error.
	if len(b.ParseErrors) != 1 {
		t.Fatalf("expected 1 parse error, got %v", b.ParseErrors)
	}
	if b.ParseErrors[0].Path != "truncated.json" {
		t.Errorf("parse error path = %q, want truncated.json", b.ParseErrors[0].Path)
	}

	if len(b.Workflows) != 3 {
		t.Fatalf("expected 3 workflows, got %d", len(b.Workflows))
	}

	// Unknown node types are retained with raw parameters; the wrong-typed
	// typeVersion degrades the workflow but salvages the other fields.
	unknown := findWorkflow(t, b, "unknown-nodes.json")
	if !unknown.Degraded() {
		t.Error("unknown-nodes should be degraded (string typeVersion)")
	}
	mystery := unknown.NodeByName("Mystery")
	if mystery == nil {
		t.Fatal("Mystery node was dropped; unknown node types must be retained")
	}
	if mystery.Type != "acme.futureWidget" {
		t.Errorf("mystery type = %q", mystery.Type)
	}
	if mystery.Parameters["gadget"] == nil {
		t.Error("mystery raw parameters were dropped")
	}

	missing := findWorkflow(t, b, "missing-connections.json")
	if !missing.Degraded() {
		t.Error("missing-connections should be degraded")
	}
	if missing.Connections != nil {
		t.Error("missing connections should stay nil")
	}

	bad := findWorkflow(t, b, "bad-connections.json")
	if !bad.Degraded() {
		t.Error("bad-connections should be degraded")
	}
	if bad.Connections != nil {
		t.Error("wrong-typed connections should stay nil")
	}
}

func TestLoadDirIgnoresNonWorkflowJSON(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "package.json", `{"name": "app", "dependencies": {"left-pad": "^1.0.0"}}`)
	writeFile(t, dir, "tsconfig.json", `{"compilerOptions": {"strict": true}}`)
	writeFile(t, dir, "broken.json", `{"foo": `)

	b := loadFixtureDir(t, dir)
	if len(b.Workflows) != 0 {
		t.Errorf("expected no workflows, got %d", len(b.Workflows))
	}
	if len(b.ParseErrors) != 0 {
		t.Errorf("non-workflow JSON must not produce parse errors, got %v", b.ParseErrors)
	}
}

func TestBundleHashDeterministic(t *testing.T) {
	first := loadFixtureDir(t, "testdata/clean")
	second := loadFixtureDir(t, "testdata/clean")
	if first.Hash == "" {
		t.Fatal("hash is empty")
	}
	if first.Hash != second.Hash {
		t.Errorf("same directory hashed differently: %s vs %s", first.Hash, second.Hash)
	}

	other := loadFixtureDir(t, "testdata/vulnerable")
	if other.Hash == first.Hash {
		t.Error("different bundles must not share a hash")
	}
}

func TestBundleHashChangesWithContent(t *testing.T) {
	dir := t.TempDir()
	src, err := os.ReadFile("testdata/clean/crm-sync.json")
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, dir, "crm-sync.json", string(src))
	before := loadFixtureDir(t, dir)

	writeFile(t, dir, "crm-sync.json", string(src)+"\n")
	after := loadFixtureDir(t, dir)

	if before.Hash == after.Hash {
		t.Error("hash must change when file content changes")
	}
}

func TestIsWorkflowFile(t *testing.T) {
	readFixture := func(path string) []byte {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{"clean export", readFixture("testdata/clean/crm-sync.json"), true},
		{"agent export", readFixture("testdata/vulnerable/email-agent.json"), true},
		{"missing connections", readFixture("testdata/malformed/missing-connections.json"), false},
		{"connections wrong type", readFixture("testdata/malformed/bad-connections.json"), false},
		{"truncated", readFixture("testdata/malformed/truncated.json"), false},
		{"package.json", []byte(`{"name": "app", "dependencies": {}}`), false},
		{"nodes not an array", []byte(`{"nodes": {}, "connections": {}}`), false},
		{"not json", []byte(`hello`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsWorkflowFile(tt.data); got != tt.want {
				t.Errorf("IsWorkflowFile = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadDirCached(t *testing.T) {
	first, err := LoadDirCached("testdata/clean")
	if err != nil {
		t.Fatal(err)
	}
	second, err := LoadDirCached("testdata/clean")
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Error("expected the cached bundle pointer on the second load")
	}
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
