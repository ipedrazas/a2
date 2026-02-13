package checkutil

import (
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// mockChecker is a simple mock for testing ResultBuilder.
type mockChecker struct {
	id   string
	name string
}

func (m *mockChecker) ID() string   { return m.id }
func (m *mockChecker) Name() string { return m.name }
func (m *mockChecker) Run(path string) (checker.Result, error) {
	return checker.Result{}, nil
}

type UtilTestSuite struct {
	suite.Suite
}

func TestUtilTestSuite(t *testing.T) {
	suite.Run(t, new(UtilTestSuite))
}

func (s *UtilTestSuite) TestTruncateMessage_ShortMessage() {
	msg := "Short message"
	result := TruncateMessage(msg, 100)
	assert.Equal(s.T(), "Short message", result)
}

func (s *UtilTestSuite) TestTruncateMessage_ExactLength() {
	msg := "12345"
	result := TruncateMessage(msg, 5)
	assert.Equal(s.T(), "12345", result)
}

func (s *UtilTestSuite) TestTruncateMessage_LongMessage() {
	msg := "This is a very long message that should be truncated"
	result := TruncateMessage(msg, 20)
	assert.Equal(s.T(), "This is a very long ...", result)
}

func (s *UtilTestSuite) TestTruncateMessage_WithWhitespace() {
	msg := "   Message with whitespace   "
	result := TruncateMessage(msg, 100)
	assert.Equal(s.T(), "Message with whitespace", result)
}

func (s *UtilTestSuite) TestTruncateMessage_EmptyMessage() {
	result := TruncateMessage("", 10)
	assert.Equal(s.T(), "", result)
}

func (s *UtilTestSuite) TestTruncateMessage_WhitespaceOnly() {
	result := TruncateMessage("   ", 10)
	assert.Equal(s.T(), "", result)
}

func (s *UtilTestSuite) TestPluralize_Singular() {
	result := Pluralize(1, "file", "files")
	assert.Equal(s.T(), "file", result)
}

func (s *UtilTestSuite) TestPluralize_Plural() {
	result := Pluralize(2, "file", "files")
	assert.Equal(s.T(), "files", result)
}

func (s *UtilTestSuite) TestPluralize_Zero() {
	result := Pluralize(0, "file", "files")
	assert.Equal(s.T(), "files", result)
}

func (s *UtilTestSuite) TestPluralize_LargeNumber() {
	result := Pluralize(100, "error", "errors")
	assert.Equal(s.T(), "errors", result)
}

func (s *UtilTestSuite) TestPluralizeCount_Singular() {
	result := PluralizeCount(1, "file", "files")
	assert.Equal(s.T(), "1 file", result)
}

func (s *UtilTestSuite) TestPluralizeCount_Plural() {
	result := PluralizeCount(5, "file", "files")
	assert.Equal(s.T(), "5 files", result)
}

func (s *UtilTestSuite) TestPluralizeCount_Zero() {
	result := PluralizeCount(0, "error", "errors")
	assert.Equal(s.T(), "0 errors", result)
}

func (s *UtilTestSuite) TestRunCommand_Success() {
	result := RunCommand(".", "echo", "hello")
	assert.True(s.T(), result.Success())
	assert.Contains(s.T(), result.Stdout, "hello")
	assert.Equal(s.T(), 0, result.ExitCode)
}

func (s *UtilTestSuite) TestRunCommand_Failure() {
	result := RunCommand(".", "ls", "/nonexistent/path/that/does/not/exist")
	assert.False(s.T(), result.Success())
	assert.NotEqual(s.T(), 0, result.ExitCode)
}

func (s *UtilTestSuite) TestRunCommand_NotFound() {
	result := RunCommand(".", "nonexistent_command_xyz123")
	assert.False(s.T(), result.Success())
	assert.Equal(s.T(), -1, result.ExitCode)
}

func (s *UtilTestSuite) TestCommandResult_Output_PreferStderr() {
	result := &CommandResult{
		Stdout: "stdout content",
		Stderr: "stderr content",
	}
	assert.Equal(s.T(), "stderr content", result.Output())
}

func (s *UtilTestSuite) TestCommandResult_Output_FallbackToStdout() {
	result := &CommandResult{
		Stdout: "stdout content",
		Stderr: "",
	}
	assert.Equal(s.T(), "stdout content", result.Output())
}

func (s *UtilTestSuite) TestCommandResult_CombinedOutput() {
	result := &CommandResult{
		Stdout: "stdout",
		Stderr: "stderr",
	}
	assert.Equal(s.T(), "stdoutstderr", result.CombinedOutput())
}

func (s *UtilTestSuite) TestToolAvailable_Exists() {
	// 'echo' should exist on all Unix systems
	assert.True(s.T(), ToolAvailable("echo"))
}

func (s *UtilTestSuite) TestToolAvailable_NotExists() {
	assert.False(s.T(), ToolAvailable("nonexistent_tool_xyz123"))
}

func (s *UtilTestSuite) TestToolNotFoundError_True() {
	result := RunCommand(".", "nonexistent_command_xyz123")
	assert.True(s.T(), ToolNotFoundError(result.Err))
}

func (s *UtilTestSuite) TestToolNotFoundError_False() {
	result := RunCommand(".", "ls", "/nonexistent/path")
	assert.False(s.T(), ToolNotFoundError(result.Err))
}

func (s *UtilTestSuite) TestToolNotFoundError_NilError() {
	assert.False(s.T(), ToolNotFoundError(nil))
}

// ResultBuilder tests

func (s *UtilTestSuite) TestNewResultBuilder() {
	mc := &mockChecker{id: "test:check", name: "Test Check"}
	rb := NewResultBuilder(mc, checker.LangGo)

	assert.NotNil(s.T(), rb)
}

func (s *UtilTestSuite) TestResultBuilder_Pass() {
	mc := &mockChecker{id: "go:build", name: "Go Build"}
	rb := NewResultBuilder(mc, checker.LangGo)

	result := rb.Pass("Build successful")

	assert.Equal(s.T(), "Go Build", result.Name)
	assert.Equal(s.T(), "go:build", result.ID)
	assert.True(s.T(), result.Passed)
	assert.Equal(s.T(), checker.Pass, result.Status)
	assert.Equal(s.T(), "Build successful", result.Message)
	assert.Equal(s.T(), "Build successful", result.Reason)
	assert.Equal(s.T(), checker.LangGo, result.Language)
}

func (s *UtilTestSuite) TestResultBuilder_Fail() {
	mc := &mockChecker{id: "python:tests", name: "Python Tests"}
	rb := NewResultBuilder(mc, checker.LangPython)

	result := rb.Fail("Tests failed: 3 errors")

	assert.Equal(s.T(), "Python Tests", result.Name)
	assert.Equal(s.T(), "python:tests", result.ID)
	assert.False(s.T(), result.Passed)
	assert.Equal(s.T(), checker.Fail, result.Status)
	assert.Equal(s.T(), "Tests failed: 3 errors", result.Message)
	assert.Equal(s.T(), "Tests failed: 3 errors", result.Reason)
	assert.Equal(s.T(), checker.LangPython, result.Language)
}

func (s *UtilTestSuite) TestResultBuilder_Warn() {
	mc := &mockChecker{id: "node:deps", name: "Node Dependencies"}
	rb := NewResultBuilder(mc, checker.LangNode)

	result := rb.Warn("Outdated dependencies found")

	assert.Equal(s.T(), "Node Dependencies", result.Name)
	assert.Equal(s.T(), "node:deps", result.ID)
	assert.False(s.T(), result.Passed)
	assert.Equal(s.T(), checker.Warn, result.Status)
	assert.Equal(s.T(), "Outdated dependencies found", result.Message)
	assert.Equal(s.T(), "Outdated dependencies found", result.Reason)
	assert.Equal(s.T(), checker.LangNode, result.Language)
}

func (s *UtilTestSuite) TestResultBuilder_Info() {
	mc := &mockChecker{id: "common:version", name: "Version Info"}
	rb := NewResultBuilder(mc, checker.LangCommon)

	result := rb.Info("Version: 1.0.0")

	assert.Equal(s.T(), "Version Info", result.Name)
	assert.Equal(s.T(), "common:version", result.ID)
	assert.True(s.T(), result.Passed) // Info doesn't affect pass/fail
	assert.Equal(s.T(), checker.Info, result.Status)
	assert.Equal(s.T(), "Version: 1.0.0", result.Message)
	assert.Equal(s.T(), "Version: 1.0.0", result.Reason)
	assert.Equal(s.T(), checker.LangCommon, result.Language)
}

func (s *UtilTestSuite) TestResultBuilder_ToolNotInstalled_WithHint() {
	mc := &mockChecker{id: "go:deps", name: "Go Dependencies"}
	rb := NewResultBuilder(mc, checker.LangGo)

	result := rb.ToolNotInstalled("govulncheck", "go install golang.org/x/vuln/cmd/govulncheck@latest")

	assert.Equal(s.T(), "Go Dependencies", result.Name)
	assert.Equal(s.T(), "go:deps", result.ID)
	assert.True(s.T(), result.Passed) // Info doesn't affect pass/fail
	assert.Equal(s.T(), checker.Info, result.Status)
	assert.Equal(s.T(), "govulncheck not installed (go install golang.org/x/vuln/cmd/govulncheck@latest)", result.Message)
	assert.Equal(s.T(), "govulncheck not installed (go install golang.org/x/vuln/cmd/govulncheck@latest)", result.Reason)
	assert.Equal(s.T(), checker.LangGo, result.Language)
}

func (s *UtilTestSuite) TestResultBuilder_ToolNotInstalled_WithoutHint() {
	mc := &mockChecker{id: "python:lint", name: "Python Lint"}
	rb := NewResultBuilder(mc, checker.LangPython)

	result := rb.ToolNotInstalled("ruff", "")

	assert.Equal(s.T(), "Python Lint", result.Name)
	assert.Equal(s.T(), "python:lint", result.ID)
	assert.True(s.T(), result.Passed)
	assert.Equal(s.T(), checker.Info, result.Status)
	assert.Equal(s.T(), "ruff not installed", result.Message)
	assert.Equal(s.T(), "ruff not installed", result.Reason)
	assert.Equal(s.T(), checker.LangPython, result.Language)
}

func (s *UtilTestSuite) TestResultBuilder_MultipleResults() {
	mc := &mockChecker{id: "go:vet", name: "Go Vet"}
	rb := NewResultBuilder(mc, checker.LangGo)

	// Builder can be reused to create multiple results
	result1 := rb.Pass("No issues")
	result2 := rb.Fail("Found issues")

	assert.True(s.T(), result1.Passed)
	assert.False(s.T(), result2.Passed)
	assert.Equal(s.T(), result1.Name, result2.Name)
	assert.Equal(s.T(), result1.ID, result2.ID)
}
