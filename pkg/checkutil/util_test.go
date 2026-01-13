package checkutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

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
