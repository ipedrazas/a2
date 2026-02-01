package godogs

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

// getGodogsDir returns the directory containing the godogs package (pkg/godogs).
func getGodogsDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil
	}
	return filepath.Dir(file), nil
}

// FixturesDir returns the path to pkg/godogs/fixtures.
func FixturesDir() (string, error) {
	dir, err := getGodogsDir()
	if err != nil || dir == "" {
		return "", err
	}
	return filepath.Join(dir, "fixtures"), nil
}

// ClearDir removes all contents of dir (files and subdirs) but not dir itself.
func ClearDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		path := filepath.Join(dir, e.Name())
		if e.IsDir() {
			if err := os.RemoveAll(path); err != nil {
				return err
			}
		} else {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}
	return nil
}

// CopyFixtureDir copies a fixture directory (e.g. "with-issues", "simple-go-project") into dest.
func CopyFixtureDir(fixtureName, dest string) error {
	fixtures, err := FixturesDir()
	if err != nil {
		return err
	}
	src := filepath.Join(fixtures, fixtureName)
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dest, rel)
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		data, err := os.ReadFile(path) // #nosec G304 -- path from controlled fixtures directory in test helper
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(destPath), 0750); err != nil {
			return err
		}
		return os.WriteFile(destPath, data, info.Mode())
	})
}
