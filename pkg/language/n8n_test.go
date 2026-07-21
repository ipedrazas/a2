package language

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
)

func TestDetectN8n(t *testing.T) {
	t.Run("workflow exports in root", func(t *testing.T) {
		result := Detect("../n8n/testdata/clean")
		if !result.HasLanguage(checker.LangN8n) {
			t.Fatalf("n8n not detected; got %v", result.Languages)
		}
		if len(result.Indicators[checker.LangN8n]) == 0 {
			t.Error("expected workflow files as indicators")
		}
		if result.Primary != checker.LangN8n {
			t.Errorf("primary = %v, want n8n in a pure workflow directory", result.Primary)
		}
	})

	t.Run("workflow exports one level down", func(t *testing.T) {
		dir := t.TempDir()
		sub := filepath.Join(dir, "flows")
		if err := os.MkdirAll(sub, 0o755); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile("../n8n/testdata/clean/crm-sync.json")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(sub, "crm-sync.json"), data, 0o644); err != nil {
			t.Fatal(err)
		}
		result := Detect(dir)
		if !result.HasLanguage(checker.LangN8n) {
			t.Fatalf("n8n not detected one level down; got %v", result.Languages)
		}
	})

	t.Run("plain json is not n8n", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"name": "app", "dependencies": {}}`), 0o644); err != nil {
			t.Fatal(err)
		}
		result := Detect(dir)
		if result.HasLanguage(checker.LangN8n) {
			t.Error("package.json must not detect as n8n")
		}
		if !result.HasLanguage(checker.LangNode) {
			t.Error("package.json should still detect node")
		}
	})

	t.Run("detection respects n8n source_dir", func(t *testing.T) {
		dir := t.TempDir()
		sub := filepath.Join(dir, "exports", "prod")
		if err := os.MkdirAll(sub, 0o755); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile("../n8n/testdata/clean/crm-sync.json")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(sub, "crm-sync.json"), data, 0o644); err != nil {
			t.Fatal(err)
		}
		// Two levels deep: undetected without a source_dir, detected with one.
		if Detect(dir).HasLanguage(checker.LangN8n) {
			t.Fatal("two-level nesting should not be detected without source_dir")
		}
		result := DetectWithSourceDirs(dir, map[string][]string{"n8n": {"exports/prod"}})
		if !result.HasLanguage(checker.LangN8n) {
			t.Errorf("n8n not detected via source_dir; got %v", result.Languages)
		}
	})
}
