package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type EditorconfigCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *EditorconfigCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "editorconfig-test-*")
	s.Require().NoError(err)
}

func (s *EditorconfigCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *EditorconfigCheckTestSuite) TestIDAndName() {
	check := &EditorconfigCheck{}
	s.Equal("common:editorconfig", check.ID())
	s.Equal("Editor Config", check.Name())
}

func (s *EditorconfigCheckTestSuite) TestEditorconfig() {
	content := `root = true

[*]
indent_style = space
indent_size = 2
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".editorconfig"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EditorconfigCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, ".editorconfig")
	s.Contains(result.Reason, "indent")
}

func (s *EditorconfigCheckTestSuite) TestEditorconfigWithEndOfLine() {
	content := `root = true

[*]
end_of_line = lf
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".editorconfig"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EditorconfigCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "line endings")
}

func (s *EditorconfigCheckTestSuite) TestEditorconfigWithCharset() {
	content := `root = true

[*]
charset = utf-8
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".editorconfig"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EditorconfigCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "charset")
}

func (s *EditorconfigCheckTestSuite) TestEditorconfigWithTrimWhitespace() {
	content := `root = true

[*]
trim_trailing_whitespace = true
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".editorconfig"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EditorconfigCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "whitespace")
}

func (s *EditorconfigCheckTestSuite) TestEditorconfigComplete() {
	content := `root = true

[*]
indent_style = space
indent_size = 2
end_of_line = lf
charset = utf-8
trim_trailing_whitespace = true
insert_final_newline = true
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".editorconfig"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EditorconfigCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "indent")
	s.Contains(result.Reason, "line endings")
	s.Contains(result.Reason, "charset")
	s.Contains(result.Reason, "whitespace")
}

func (s *EditorconfigCheckTestSuite) TestVSCodeSettings() {
	vscodeDir := filepath.Join(s.tempDir, ".vscode")
	err := os.MkdirAll(vscodeDir, 0755)
	s.Require().NoError(err)

	content := `{
  "editor.tabSize": 2,
  "editor.formatOnSave": true
}`
	err = os.WriteFile(filepath.Join(vscodeDir, "settings.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EditorconfigCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "VS Code settings")
}

func (s *EditorconfigCheckTestSuite) TestVSCodeExtensions() {
	vscodeDir := filepath.Join(s.tempDir, ".vscode")
	err := os.MkdirAll(vscodeDir, 0755)
	s.Require().NoError(err)

	content := `{
  "recommendations": ["esbenp.prettier-vscode"]
}`
	err = os.WriteFile(filepath.Join(vscodeDir, "extensions.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EditorconfigCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "VS Code extensions")
}

func (s *EditorconfigCheckTestSuite) TestJetBrainsIDE() {
	ideaDir := filepath.Join(s.tempDir, ".idea")
	err := os.MkdirAll(ideaDir, 0755)
	s.Require().NoError(err)

	// Create a placeholder file
	err = os.WriteFile(filepath.Join(ideaDir, "workspace.xml"), []byte(`<project></project>`), 0644)
	s.Require().NoError(err)

	check := &EditorconfigCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "JetBrains IDE")
}

func (s *EditorconfigCheckTestSuite) TestJetBrainsCodeStyles() {
	codeStylesDir := filepath.Join(s.tempDir, ".idea", "codeStyles")
	err := os.MkdirAll(codeStylesDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(codeStylesDir, "Project.xml"), []byte(`<code_scheme></code_scheme>`), 0644)
	s.Require().NoError(err)

	check := &EditorconfigCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "JetBrains code styles")
}

func (s *EditorconfigCheckTestSuite) TestVimConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".vimrc"), []byte(`
set tabstop=4
set shiftwidth=4
`), 0644)
	s.Require().NoError(err)

	check := &EditorconfigCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Vim/Neovim")
}

func (s *EditorconfigCheckTestSuite) TestNvimConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".nvimrc"), []byte(`
vim.opt.tabstop = 4
`), 0644)
	s.Require().NoError(err)

	check := &EditorconfigCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Vim/Neovim")
}

func (s *EditorconfigCheckTestSuite) TestDevContainer() {
	devcontainerDir := filepath.Join(s.tempDir, ".devcontainer")
	err := os.MkdirAll(devcontainerDir, 0755)
	s.Require().NoError(err)

	content := `{
  "name": "My Dev Container",
  "image": "mcr.microsoft.com/devcontainers/go:1.21"
}`
	err = os.WriteFile(filepath.Join(devcontainerDir, "devcontainer.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EditorconfigCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Dev Container")
}

func (s *EditorconfigCheckTestSuite) TestDevContainerAtRoot() {
	content := `{
  "name": "My Dev Container",
  "image": "mcr.microsoft.com/devcontainers/go:1.21"
}`
	err := os.WriteFile(filepath.Join(s.tempDir, ".devcontainer.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EditorconfigCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Dev Container")
}

func (s *EditorconfigCheckTestSuite) TestMultipleConfigs() {
	// Create .editorconfig
	editorContent := `root = true

[*]
indent_style = space
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".editorconfig"), []byte(editorContent), 0644)
	s.Require().NoError(err)

	// Create .vscode/settings.json
	vscodeDir := filepath.Join(s.tempDir, ".vscode")
	err = os.MkdirAll(vscodeDir, 0755)
	s.Require().NoError(err)

	vscodeContent := `{"editor.tabSize": 2}`
	err = os.WriteFile(filepath.Join(vscodeDir, "settings.json"), []byte(vscodeContent), 0644)
	s.Require().NoError(err)

	check := &EditorconfigCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, ".editorconfig")
	s.Contains(result.Reason, "VS Code settings")
}

func (s *EditorconfigCheckTestSuite) TestNoEditorConfigFound() {
	// Create some other files
	err := os.WriteFile(filepath.Join(s.tempDir, "README.md"), []byte(`# Project`), 0644)
	s.Require().NoError(err)

	check := &EditorconfigCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No editor config found")
}

func (s *EditorconfigCheckTestSuite) TestEmptyDirectory() {
	check := &EditorconfigCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No editor config found")
}

func TestEditorconfigCheckTestSuite(t *testing.T) {
	suite.Run(t, new(EditorconfigCheckTestSuite))
}
