package language

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/n8n"
)

// Bounds for n8n detection so it stays cheap on large repositories: only the
// root and one level of subdirectories are scanned, only this many JSON
// files are inspected, and only files up to this size are read.
const (
	n8nDetectMaxFiles    = 128
	n8nDetectMaxFileSize = 4 << 20 // 4 MiB
	n8nDetectMaxHits     = 3
)

// n8nSkipDirs are directory names skipped during n8n detection.
var n8nSkipDirs = map[string]bool{
	"node_modules": true,
	"vendor":       true,
	"dist":         true,
	"a2-cache":     true,
}

// detectN8nIndicators reports files with the n8n workflow export shape (a
// JSON object with a top-level "nodes" array and "connections" object) in
// path or one level of subdirectories. It returns up to n8nDetectMaxHits
// relative paths to use as detection indicators; deeper nesting is left to
// the n8n checks' own recursive loader.
func detectN8nIndicators(path string) []string {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	var hits []string
	inspected := 0

	scanDir := func(dir, prefix string) {
		dirEntries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, e := range dirEntries {
			if len(hits) >= n8nDetectMaxHits || inspected >= n8nDetectMaxFiles {
				return
			}
			if e.IsDir() || !strings.EqualFold(filepath.Ext(e.Name()), ".json") {
				continue
			}
			info, err := e.Info()
			if err != nil || info.Size() > n8nDetectMaxFileSize {
				continue
			}
			inspected++
			data, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				continue
			}
			if n8n.IsWorkflowFile(data) {
				hits = append(hits, filepath.ToSlash(filepath.Join(prefix, e.Name())))
			}
		}
	}

	scanDir(path, "")
	for _, e := range entries {
		if len(hits) >= n8nDetectMaxHits || inspected >= n8nDetectMaxFiles {
			break
		}
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		if n8nSkipDirs[name] || strings.HasPrefix(name, ".") {
			continue
		}
		scanDir(filepath.Join(path, name), name)
	}
	return hits
}
