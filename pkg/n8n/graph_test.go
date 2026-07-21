package n8n

import (
	"os"
	"reflect"
	"testing"
)

func loadFixtureWorkflow(t *testing.T, path string) *Workflow {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	wf, err := ParseWorkflow(data)
	if err != nil {
		t.Fatalf("ParseWorkflow(%s): %v", path, err)
	}
	return wf
}

func TestReachableTransitive(t *testing.T) {
	wf := loadFixtureWorkflow(t, "testdata/vulnerable/lead-webhook.json")
	g := BuildGraph(wf)

	got := g.Reachable("Webhook", ConnMain)
	want := []string{"Create HubSpot Contact", "HTTP Request", "Set"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Reachable(Webhook) = %v, want %v", got, want)
	}

	// Reachability is transitive, not adjacency: the CRM write is three
	// hops from the webhook.
	if got := g.Reachable("Create HubSpot Contact", ConnMain); len(got) != 0 {
		t.Errorf("Reachable(sink) = %v, want empty", got)
	}
}

func TestPathTo(t *testing.T) {
	wf := loadFixtureWorkflow(t, "testdata/vulnerable/lead-webhook.json")
	g := BuildGraph(wf)

	got := g.PathTo("Webhook", "Create HubSpot Contact", ConnMain)
	want := []string{"Webhook", "Set", "HTTP Request", "Create HubSpot Contact"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PathTo = %v, want %v", got, want)
	}

	if got := g.PathTo("Create HubSpot Contact", "Webhook", ConnMain); got != nil {
		t.Errorf("reverse path should be nil, got %v", got)
	}
}

func TestErrorOutputTraversal(t *testing.T) {
	wf := loadFixtureWorkflow(t, "testdata/clean/webhook-intake.json")
	g := BuildGraph(wf)

	// Notify Failures hangs off Save Lead's error output (port 1); main
	// traversal must follow error branches too.
	reachable := g.Reachable("Lead Intake Webhook", ConnMain)
	found := false
	for _, name := range reachable {
		if name == "Notify Failures" {
			found = true
		}
	}
	if !found {
		t.Errorf("error-output branch not reachable from trigger; got %v", reachable)
	}
}

func TestAttachedTools(t *testing.T) {
	wf := loadFixtureWorkflow(t, "testdata/vulnerable/email-agent.json")
	g := BuildGraph(wf)

	tools := g.AttachedTools("AI Agent")
	var names []string
	for _, n := range tools {
		names = append(names, n.Name)
	}
	want := []string{"Look Up CRM Sheet", "Send Email"}
	if !reflect.DeepEqual(names, want) {
		t.Errorf("AttachedTools = %v, want %v", names, want)
	}

	// The untrusted trigger reaches the agent over the main flow.
	reachable := g.Reachable("Email Trigger (IMAP)", ConnMain)
	if !reflect.DeepEqual(reachable, []string{"AI Agent"}) {
		t.Errorf("Reachable(IMAP trigger) = %v", reachable)
	}

	// Without a type filter, traversal crosses AI connection types too.
	all := g.Reachable("OpenAI Chat Model")
	if !reflect.DeepEqual(all, []string{"AI Agent"}) {
		t.Errorf("Reachable(model, all types) = %v", all)
	}
}

func TestFlattenedConnectionsParse(t *testing.T) {
	// Some template exports flatten the ports array; the parser accepts both.
	data := []byte(`{
		"name": "Flat",
		"nodes": [
			{"name": "A", "type": "n8n-nodes-base.manualTrigger", "typeVersion": 1, "parameters": {}},
			{"name": "B", "type": "n8n-nodes-base.set", "typeVersion": 3, "parameters": {}}
		],
		"connections": {
			"A": {"main": [{"node": "B", "type": "main", "index": 0}]}
		}
	}`)
	wf, err := ParseWorkflow(data)
	if err != nil {
		t.Fatal(err)
	}
	if wf.Degraded() {
		t.Fatalf("unexpected parse issues: %v", wf.ParseIssues)
	}
	g := BuildGraph(wf)
	if got := g.Reachable("A", ConnMain); !reflect.DeepEqual(got, []string{"B"}) {
		t.Errorf("Reachable(A) = %v, want [B]", got)
	}
}

func TestDanglingConnectionTarget(t *testing.T) {
	data := []byte(`{
		"name": "Dangling",
		"nodes": [
			{"name": "A", "type": "n8n-nodes-base.manualTrigger", "typeVersion": 1, "parameters": {}}
		],
		"connections": {
			"A": {"main": [[{"node": "Ghost", "type": "main", "index": 0}, null]]},
			"Phantom": {"main": [[{"node": "A", "type": "main", "index": 0}]]}
		}
	}`)
	wf, err := ParseWorkflow(data)
	if err != nil {
		t.Fatal(err)
	}
	g := BuildGraph(wf)
	if got := g.Reachable("A"); len(got) != 0 {
		t.Errorf("edges to unknown nodes must be dropped, got %v", got)
	}
	if got := g.InEdges("A"); len(got) != 0 {
		t.Errorf("edges from unknown nodes must be dropped, got %v", got)
	}
}
