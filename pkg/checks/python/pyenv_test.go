package pythoncheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PyenvTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *PyenvTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "pyenv-test-*")
	s.Require().NoError(err)
	s.tempDir = dir
}

func (s *PyenvTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *PyenvTestSuite) TestResolvePythonEnv_BareCommand() {
	// No venv indicators → bare command
	name, prefix := resolvePythonEnv(s.tempDir, "pytest")
	s.Equal("pytest", name)
	s.Empty(prefix)
}

func (s *PyenvTestSuite) TestResolvePythonEnv_DotVenv() {
	// .venv/bin/pytest exists → use direct path
	binDir := filepath.Join(s.tempDir, ".venv", "bin")
	s.Require().NoError(os.MkdirAll(binDir, 0755))
	s.Require().NoError(os.WriteFile(filepath.Join(binDir, "pytest"), []byte("#!/bin/sh"), 0755))

	name, prefix := resolvePythonEnv(s.tempDir, "pytest")
	s.Equal(filepath.Join(s.tempDir, ".venv", "bin", "pytest"), name)
	s.Empty(prefix)
}

func (s *PyenvTestSuite) TestResolvePythonEnv_Venv() {
	// venv/bin/pytest exists → use direct path
	binDir := filepath.Join(s.tempDir, "venv", "bin")
	s.Require().NoError(os.MkdirAll(binDir, 0755))
	s.Require().NoError(os.WriteFile(filepath.Join(binDir, "pytest"), []byte("#!/bin/sh"), 0755))

	name, prefix := resolvePythonEnv(s.tempDir, "pytest")
	s.Equal(filepath.Join(s.tempDir, "venv", "bin", "pytest"), name)
	s.Empty(prefix)
}

func (s *PyenvTestSuite) TestResolvePythonEnv_DotVenvPreferredOverVenv() {
	// Both .venv and venv exist → .venv wins
	for _, dir := range []string{".venv", "venv"} {
		binDir := filepath.Join(s.tempDir, dir, "bin")
		s.Require().NoError(os.MkdirAll(binDir, 0755))
		s.Require().NoError(os.WriteFile(filepath.Join(binDir, "pytest"), []byte("#!/bin/sh"), 0755))
	}

	name, prefix := resolvePythonEnv(s.tempDir, "pytest")
	s.Equal(filepath.Join(s.tempDir, ".venv", "bin", "pytest"), name)
	s.Empty(prefix)
}

func (s *PyenvTestSuite) TestResolvePythonEnv_VenvCommandNotPresent() {
	// .venv exists but the specific command is not in it → bare command
	binDir := filepath.Join(s.tempDir, ".venv", "bin")
	s.Require().NoError(os.MkdirAll(binDir, 0755))
	// No pytest binary created

	name, prefix := resolvePythonEnv(s.tempDir, "pytest")
	s.Equal("pytest", name)
	s.Empty(prefix)
}

func (s *PyenvTestSuite) TestPythonToolAvailable_InVenv() {
	binDir := filepath.Join(s.tempDir, ".venv", "bin")
	s.Require().NoError(os.MkdirAll(binDir, 0755))
	s.Require().NoError(os.WriteFile(filepath.Join(binDir, "ruff"), []byte("#!/bin/sh"), 0755))

	s.True(pythonToolAvailable(s.tempDir, "ruff"))
}

func (s *PyenvTestSuite) TestPythonToolAvailable_NotInVenv() {
	// tool not in venv and likely not in PATH (use unlikely name)
	s.False(pythonToolAvailable(s.tempDir, "nonexistent-tool-xyz-12345"))
}

func TestPyenvTestSuite(t *testing.T) {
	suite.Run(t, new(PyenvTestSuite))
}
