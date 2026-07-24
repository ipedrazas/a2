package gocheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LoggingCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *LoggingCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "logging-test-*")
	s.Require().NoError(err)
}

func (s *LoggingCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *LoggingCheckTestSuite) write(name, content string) {
	s.Require().NoError(os.WriteFile(filepath.Join(s.tempDir, name), []byte(content), 0644))
}

func (s *LoggingCheckTestSuite) TestDetectsSlog() {
	s.write("main.go", "package main\n\nimport \"log/slog\"\n\nfunc main() { slog.Info(\"hi\") }\n")

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed, "slog import must count as structured logging: %s", result.Reason)
}

func (s *LoggingCheckTestSuite) TestWarnsOnPrints() {
	s.write("main.go", "package main\n\nimport \"fmt\"\n\nfunc main() { fmt.Println(\"hi\") }\n")

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Contains(result.Reason, "fmt.Print")
}

// TestDotPathScansTree is the regression test for the root-skip bug: with
// path "." the walk root's Name() is ".", which matched the hidden-dir skip
// and silently skipped the entire tree — no files were ever scanned, so
// projects using slog were reported as having no structured logging.
func (s *LoggingCheckTestSuite) TestDotPathScansTree() {
	s.Require().NoError(os.Mkdir(filepath.Join(s.tempDir, "internal"), 0755))
	s.write(filepath.Join("internal", "app.go"),
		"package app\n\nimport \"log/slog\"\n\nfunc Log() { slog.Info(\"hi\") }\n")

	wd, err := os.Getwd()
	s.Require().NoError(err)
	s.Require().NoError(os.Chdir(s.tempDir))
	defer func() { s.Require().NoError(os.Chdir(wd)) }()

	check := &LoggingCheck{}
	result, err := check.Run(".")

	s.NoError(err)
	s.True(result.Passed, "run with path '.' must scan the tree: %s", result.Reason)
}

// TestZeroFilesScannedWarns: a verdict from a walk that saw no files is
// meaningless — the check must say so instead of passing or reporting
// "no structured logging".
func (s *LoggingCheckTestSuite) TestZeroFilesScannedWarns() {
	s.write("notes.txt", "not go code")

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Contains(result.Reason, "No Go files scanned")
}

func TestLoggingCheckTestSuite(t *testing.T) {
	suite.Run(t, new(LoggingCheckTestSuite))
}
