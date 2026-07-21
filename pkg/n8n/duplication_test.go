package n8n

import (
	"strings"
	"testing"
)

func TestFindDuplicateSequencesFixture(t *testing.T) {
	b := loadFixtureDir(t, "testdata/duplication")
	clusters := FindDuplicateSequences(b.Workflows, 5)

	if len(clusters) != 1 {
		t.Fatalf("expected exactly 1 cluster, got %d: %+v", len(clusters), clusters)
	}
	cluster := clusters[0]

	// The copied run: scheduleTrigger → httpRequest → code → if → set → slack.
	wantTypes := []string{
		"n8n-nodes-base.scheduleTrigger",
		"n8n-nodes-base.httpRequest",
		"n8n-nodes-base.code",
		"n8n-nodes-base.if",
		"n8n-nodes-base.set",
		"n8n-nodes-base.slack",
	}
	if strings.Join(cluster.Types, "→") != strings.Join(wantTypes, "→") {
		t.Errorf("types = %v, want %v", cluster.Types, wantTypes)
	}

	if len(cluster.Occurrences) != 2 {
		t.Fatalf("expected 2 occurrences, got %+v", cluster.Occurrences)
	}
	files := []string{cluster.Occurrences[0].File, cluster.Occurrences[1].File}
	if files[0] != "invoice-alerts.json" || files[1] != "payment-alerts.json" {
		t.Errorf("occurrence files = %v", files)
	}
	if cluster.Occurrences[0].Nodes[len(cluster.Occurrences[0].Nodes)-1] != "Alert Finance" {
		t.Errorf("occurrence nodes = %v", cluster.Occurrences[0].Nodes)
	}
	if cluster.ParamSimilarity <= 0.9 {
		t.Errorf("near-identical copies should have high param-key similarity, got %v", cluster.ParamSimilarity)
	}
}

func TestFindDuplicateSequencesNoFalsePositives(t *testing.T) {
	for _, dir := range []string{"testdata/clean", "testdata/vulnerable", "testdata/mitigated"} {
		b := loadFixtureDir(t, dir)
		if clusters := FindDuplicateSequences(b.Workflows, 5); len(clusters) != 0 {
			t.Errorf("%s should have no duplication clusters, got %+v", dir, clusters)
		}
	}
}

func TestMainFlowPathsBranching(t *testing.T) {
	wf := loadFixtureWorkflow(t, "testdata/clean/webhook-intake.json")
	g := BuildGraph(wf)
	paths := mainFlowPaths(wf, g)

	// Roots: the webhook trigger (and any unconnected nodes). Every path
	// starts at a root and branches produce separate paths.
	if len(paths) < 2 {
		t.Fatalf("expected branch paths, got %v", paths)
	}
	joined := make([]string, len(paths))
	for i, p := range paths {
		joined[i] = strings.Join(p, "→")
	}
	all := strings.Join(joined, "\n")
	if !strings.Contains(all, "Lead Intake Webhook→Validate Payload→Save Lead→Respond Success") {
		t.Errorf("missing success path:\n%s", all)
	}
	if !strings.Contains(all, "Save Lead→Notify Failures") {
		t.Errorf("missing error-branch path:\n%s", all)
	}
}
