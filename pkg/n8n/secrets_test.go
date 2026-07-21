package n8n

import (
	"fmt"
	"strings"
	"testing"
)

// fake builds obviously-fake key material at runtime so no secret-shaped
// literal sits in the source (which would trip repo-level secret scanners).
func fake(prefix string, filler rune, n int) string {
	return prefix + strings.Repeat(string(filler), n)
}

func TestScanValuePatterns(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantRule  string
		confirmed bool
	}{
		{"openai style", "Bearer " + fake("sk-", 'A', 40), "openai-style key (sk-...)", true},
		{"github pat", fake("ghp_", 'a', 36), "github token (gh?_...)", true},
		{"slack bot token", fake("xoxb-", '1', 24), "slack token (xox...)", true},
		{"aws key id", "AKIA" + strings.Repeat("A", 16), "aws access key id (AKIA/ASIA...)", true},
		{"google api key", fake("AIza", 'B', 35), "google api key (AIza...)", true},
		{"jwt", fake("eyJ", 'x', 12) + "." + strings.Repeat("y", 20) + "." + strings.Repeat("z", 10), "jwt", true},
		{"private key", "-----BEGIN RSA PRIVATE KEY-----", "private key block", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := scanValue("wf", "node", "parameters.value", "value", tt.value)
			if len(findings) != 1 {
				t.Fatalf("got %d findings, want 1", len(findings))
			}
			if findings[0].Rule != tt.wantRule || findings[0].Confirmed != tt.confirmed {
				t.Errorf("rule=%q confirmed=%v, want %q/%v", findings[0].Rule, findings[0].Confirmed, tt.wantRule, tt.confirmed)
			}
		})
	}
}

func TestScanValueEntropyCandidates(t *testing.T) {
	random := "q7Zp3vNb8Rt2Lw9XmK4Jh6Fd"

	// High-entropy literal in an auth-labelled field: candidate.
	findings := scanValue("wf", "node", "parameters.headerParameters.parameters[0].value", "X-Api-Key", random)
	if len(findings) != 1 || findings[0].Confirmed {
		t.Fatalf("expected one unconfirmed candidate, got %+v", findings)
	}

	// Same value in a non-auth field: silent.
	if got := scanValue("wf", "node", "parameters.url", "url", random); len(got) != 0 {
		t.Errorf("non-auth field should not produce candidates, got %+v", got)
	}

	// Expressions are credential references, not literals.
	if got := scanValue("wf", "node", "parameters.p.value", "Authorization", "={{ $credentials.apiKey }}"); len(got) != 0 {
		t.Errorf("expressions should not produce candidates, got %+v", got)
	}

	// Low-entropy placeholder in an auth field: silent.
	if got := scanValue("wf", "node", "parameters.p.value", "Authorization", "YOUR_API_KEY_HERE_PLEASE"); len(got) != 0 {
		t.Errorf("placeholders should not produce candidates, got %+v", got)
	}
}

func TestScanBundleSecretsVulnerableFixture(t *testing.T) {
	b := loadFixtureDir(t, "testdata/vulnerable")
	findings := ScanBundleSecrets(b)

	var confirmed []SecretFinding
	for _, f := range findings {
		if f.Confirmed {
			confirmed = append(confirmed, f)
		}
	}
	if len(confirmed) != 2 {
		t.Fatalf("expected 2 confirmed findings (sk- and AKIA), got %d: %+v", len(confirmed), findings)
	}
	for _, f := range confirmed {
		if f.File != "lead-webhook.json" || f.Node != "HTTP Request" {
			t.Errorf("finding at unexpected location: %+v", f)
		}
		if !strings.Contains(f.FieldPath, "headerParameters") {
			t.Errorf("field path should locate the header, got %q", f.FieldPath)
		}
	}

	// The finding structs themselves must not carry the matched values.
	serialized := fmt.Sprintf("%+v", findings)
	for _, secret := range []string{"sk-FAKE", "AKIAFAKE"} {
		if strings.Contains(serialized, secret) {
			t.Errorf("findings leak secret material %q", secret)
		}
	}
}

func TestScanCleanAndMitigatedFixtures(t *testing.T) {
	for _, dir := range []string{"testdata/clean", "testdata/mitigated"} {
		b := loadFixtureDir(t, dir)
		if findings := ScanBundleSecrets(b); len(findings) != 0 {
			t.Errorf("%s should have no secret findings, got %+v", dir, findings)
		}
	}
}

func TestScanPinData(t *testing.T) {
	data := []byte(`{
		"name": "Pinned Secret",
		"nodes": [{"name": "Webhook", "type": "n8n-nodes-base.webhook", "typeVersion": 2, "parameters": {}}],
		"connections": {},
		"pinData": {"Webhook": [{"json": {"api_token": "` + fake("sk-", 'C', 32) + `"}}]}
	}`)
	wf, err := ParseWorkflow(data)
	if err != nil {
		t.Fatal(err)
	}
	findings := ScanWorkflowSecrets(wf)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding in pinData, got %+v", findings)
	}
	if !strings.HasPrefix(findings[0].FieldPath, "pinData.Webhook") {
		t.Errorf("field path = %q, want pinData.Webhook...", findings[0].FieldPath)
	}
}
