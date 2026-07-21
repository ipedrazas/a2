package n8n

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// maxWorkflowFileSize caps how much of a single file the loader will read.
// Real exports are a few hundred KB at most; anything larger is suspect.
const maxWorkflowFileSize = 16 << 20 // 16 MiB

// Bundle is the set of n8n workflow files found under one directory.
type Bundle struct {
	Dir       string
	Workflows []*WorkflowFile // sorted by Path
	// ParseErrors lists files that look like workflow exports but could not
	// be parsed at all (truncated JSON, unreadable). They are part of the
	// bundle: checks must surface them, and they contribute to Hash.
	ParseErrors []ParseError
	// Hash is a SHA-256 over the sorted relative paths and contents of every
	// file in the bundle (workflows and parse-error files), so an identical
	// input bundle always hashes identically.
	Hash string
}

// WorkflowFile pairs a parsed workflow with its file location.
type WorkflowFile struct {
	Path     string // relative to Bundle.Dir, slash-separated
	Workflow *Workflow
}

// ParseError describes a bundle file that could not be parsed.
type ParseError struct {
	Path    string
	Message string
}

// skipDirs are directory names never scanned for workflow exports.
var skipDirs = map[string]bool{
	"node_modules": true,
	"vendor":       true,
	"dist":         true,
	"a2-cache":     true,
}

// IsWorkflowFile reports whether data has the n8n export shape used for
// domain detection: a JSON object with a top-level "nodes" array and a
// "connections" object.
func IsWorkflowFile(data []byte) bool {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(data, &top); err != nil {
		return false
	}
	nodes, ok := top["nodes"]
	if !ok {
		return false
	}
	var nodeList []json.RawMessage
	if err := json.Unmarshal(nodes, &nodeList); err != nil {
		return false
	}
	conns, ok := top["connections"]
	if !ok {
		return false
	}
	var connObj map[string]json.RawMessage
	return json.Unmarshal(conns, &connObj) == nil
}

// LoadDir finds and parses every n8n workflow file under dir (recursively,
// skipping hidden and dependency directories). JSON files without the export
// shape are ignored; files that look like exports but fail to parse are
// recorded as ParseErrors. The loader is stricter than IsWorkflowFile on
// purpose: a file with a "nodes" array but a missing/broken "connections"
// object is loaded as a degraded workflow rather than ignored, so a
// truncated estate can never look clean.
func LoadDir(dir string) (*Bundle, error) {
	bundle := &Bundle{Dir: dir}
	type bundleFile struct {
		path    string
		content []byte
	}
	var files []bundleFile

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if path != dir && (skipDirs[name] || strings.HasPrefix(name, ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.EqualFold(filepath.Ext(d.Name()), ".json") {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		info, err := d.Info()
		if err != nil {
			return err
		}
		if info.Size() > maxWorkflowFileSize {
			bundle.ParseErrors = append(bundle.ParseErrors, ParseError{
				Path:    rel,
				Message: fmt.Sprintf("file too large to scan (%d bytes, limit %d)", info.Size(), maxWorkflowFileSize),
			})
			return nil
		}
		data, err := os.ReadFile(path) // #nosec G304,G122 - path is from a controlled WalkDir over the scan directory
		if err != nil {
			return err
		}

		wf, perr := ParseWorkflow(data)
		switch {
		case perr == nil:
			bundle.Workflows = append(bundle.Workflows, &WorkflowFile{Path: rel, Workflow: wf})
			files = append(files, bundleFile{rel, data})
		case perr == ErrNotWorkflow:
			// Valid JSON, not an export (package.json, tsconfig, ...): ignore.
		case bytes.Contains(data, []byte(`"nodes"`)):
			// Broken JSON that mentions "nodes" is almost certainly a
			// truncated or corrupted export: it belongs to the bundle.
			bundle.ParseErrors = append(bundle.ParseErrors, ParseError{Path: rel, Message: perr.Error()})
			files = append(files, bundleFile{rel, data})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(bundle.Workflows, func(i, j int) bool { return bundle.Workflows[i].Path < bundle.Workflows[j].Path })
	sort.Slice(bundle.ParseErrors, func(i, j int) bool { return bundle.ParseErrors[i].Path < bundle.ParseErrors[j].Path })
	sort.Slice(files, func(i, j int) bool { return files[i].path < files[j].path })

	h := sha256.New()
	for _, f := range files {
		h.Write([]byte(f.path))
		h.Write([]byte{0})
		h.Write(f.content)
		h.Write([]byte{0})
	}
	bundle.Hash = hex.EncodeToString(h.Sum(nil))
	return bundle, nil
}

var (
	cacheMu sync.Mutex
	cache   = map[string]*Bundle{}
)

// LoadDirCached memoizes LoadDir per absolute path for the lifetime of the
// process, so the many independent n8n checks in one run parse the bundle
// once. Load errors are not cached.
func LoadDirCached(dir string) (*Bundle, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	cacheMu.Lock()
	defer cacheMu.Unlock()
	if b, ok := cache[abs]; ok {
		return b, nil
	}
	b, err := LoadDir(dir)
	if err != nil {
		return nil, err
	}
	cache[abs] = b
	return b, nil
}
