package nodecheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type LoggingTestSuite struct {
	suite.Suite
	tempDir string
	check   *LoggingCheck
}

func (s *LoggingTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "node-logging-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &LoggingCheck{}
}

func (s *LoggingTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *LoggingTestSuite) writeFile(name, content string) {
	path := filepath.Join(s.tempDir, name)
	dir := filepath.Dir(path)
	if dir != s.tempDir {
		err := os.MkdirAll(dir, 0755)
		s.NoError(err)
	}
	err := os.WriteFile(path, []byte(content), 0644)
	s.NoError(err)
}

func (s *LoggingTestSuite) TestID() {
	s.Equal("node:logging", s.check.ID())
}

func (s *LoggingTestSuite) TestName() {
	s.Equal("Node Logging", s.check.Name())
}

func (s *LoggingTestSuite) TestRun_NoPackageJson() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Reason, "package.json not found")
}

func (s *LoggingTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangNode, result.Language)
}

func (s *LoggingTestSuite) TestRun_WithStructuredLogging() {
	s.writeFile("package.json", `{
		"name": "test-app",
		"dependencies": {
			"winston": "^3.0.0"
		}
	}`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "structured logging")
}

func (s *LoggingTestSuite) TestRun_WithPino() {
	s.writeFile("package.json", `{
		"name": "test-app",
		"dependencies": {
			"pino": "^8.0.0"
		}
	}`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
}

func (s *LoggingTestSuite) TestRun_WithConsoleStatements() {
	s.writeFile("package.json", `{"name": "test-app"}`)
	s.writeFile("src/index.js", `
		function main() {
			console.log("Hello world");
			console.error("error");
		}
	`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "console.*")
}

func (s *LoggingTestSuite) TestRun_NoLogging() {
	s.writeFile("package.json", `{"name": "test-app"}`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No structured logging")
}

func (s *LoggingTestSuite) TestRun_StructuredWithConsole() {
	s.writeFile("package.json", `{
		"name": "test-app",
		"dependencies": {
			"winston": "^3.0.0"
		}
	}`)
	s.writeFile("src/main.js", `console.log("debug")`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "structured logging but found")
}

func TestLoggingTestSuite(t *testing.T) {
	suite.Run(t, new(LoggingTestSuite))
}
