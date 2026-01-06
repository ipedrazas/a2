package pythoncheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
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

func (s *LoggingCheckTestSuite) TestIDAndName() {
	check := &LoggingCheck{}
	s.Equal("python:logging", check.ID())
	s.Equal("Python Logging", check.Name())
}

func (s *LoggingCheckTestSuite) TestUsesLoggingNoPrint() {
	code := `import logging

logger = logging.getLogger(__name__)

def main():
    logger.info("Starting application")
`
	err := os.WriteFile(filepath.Join(s.tempDir, "app.py"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Uses logging module")
	s.Contains(result.Message, "no print()")
}

func (s *LoggingCheckTestSuite) TestUsesStructlog() {
	code := `import structlog

logger = structlog.get_logger()

def main():
    logger.info("hello")
`
	err := os.WriteFile(filepath.Join(s.tempDir, "app.py"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
}

func (s *LoggingCheckTestSuite) TestUsesLoguru() {
	code := `from loguru import logger

def main():
    logger.info("hello")
`
	err := os.WriteFile(filepath.Join(s.tempDir, "app.py"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
}

func (s *LoggingCheckTestSuite) TestLoggingWithPrint() {
	code := `import logging

logger = logging.getLogger(__name__)

def main():
    print("Debug message")
    logger.info("Info message")
`
	err := os.WriteFile(filepath.Join(s.tempDir, "app.py"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "Uses logging but found")
	s.Contains(result.Message, "print()")
}

func (s *LoggingCheckTestSuite) TestNoLoggingNoPrint() {
	code := `def main():
    x = 1 + 2
    return x
`
	err := os.WriteFile(filepath.Join(s.tempDir, "app.py"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No logging module detected")
}

func (s *LoggingCheckTestSuite) TestNoLoggingWithPrint() {
	code := `def main():
    print("Hello world")
`
	err := os.WriteFile(filepath.Join(s.tempDir, "app.py"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No logging module")
	s.Contains(result.Message, "print()")
}

func (s *LoggingCheckTestSuite) TestPrintInTestFileIgnored() {
	// Test file with print - should be ignored
	testCode := `def test_something():
    print("Debug output")
    assert True
`
	err := os.WriteFile(filepath.Join(s.tempDir, "test_app.py"), []byte(testCode), 0644)
	s.Require().NoError(err)

	// Main file with logging
	mainCode := `import logging

logger = logging.getLogger(__name__)

def main():
    logger.info("Starting")
`
	err = os.WriteFile(filepath.Join(s.tempDir, "app.py"), []byte(mainCode), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
}

func (s *LoggingCheckTestSuite) TestSkipsVenv() {
	// Create venv with print
	venvDir := filepath.Join(s.tempDir, "venv", "lib")
	err := os.MkdirAll(venvDir, 0755)
	s.Require().NoError(err)

	code := `print("Should be ignored")`
	err = os.WriteFile(filepath.Join(venvDir, "module.py"), []byte(code), 0644)
	s.Require().NoError(err)

	// Main file with logging
	mainCode := `import logging

logger = logging.getLogger(__name__)
`
	err = os.WriteFile(filepath.Join(s.tempDir, "app.py"), []byte(mainCode), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
}

func (s *LoggingCheckTestSuite) TestSkipsPycache() {
	// Create __pycache__ with print
	cacheDir := filepath.Join(s.tempDir, "__pycache__")
	err := os.MkdirAll(cacheDir, 0755)
	s.Require().NoError(err)

	code := `print("Should be ignored")`
	err = os.WriteFile(filepath.Join(cacheDir, "app.cpython-39.py"), []byte(code), 0644)
	s.Require().NoError(err)

	// Main file with logging
	mainCode := `import logging

logger = logging.getLogger(__name__)
`
	err = os.WriteFile(filepath.Join(s.tempDir, "app.py"), []byte(mainCode), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
}

func (s *LoggingCheckTestSuite) TestMultiplePrintStatements() {
	code := `def main():
    print("one")
    print("two")
    print("three")
`
	err := os.WriteFile(filepath.Join(s.tempDir, "app.py"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Contains(result.Message, "3 print()")
}

func (s *LoggingCheckTestSuite) TestFromLoggingImport() {
	code := `from logging import getLogger

logger = getLogger(__name__)

def main():
    logger.info("hello")
`
	err := os.WriteFile(filepath.Join(s.tempDir, "app.py"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
}

func (s *LoggingCheckTestSuite) TestEmptyDirectory() {
	check := &LoggingCheck{}
	result, err := check.Run(s.tempDir)

	// Empty directory means no logging, no prints
	s.NoError(err)
	s.False(result.Passed)
	s.Contains(result.Message, "No logging module detected")
}

func TestLoggingCheckTestSuite(t *testing.T) {
	suite.Run(t, new(LoggingCheckTestSuite))
}
