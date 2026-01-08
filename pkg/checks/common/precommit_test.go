package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type PrecommitCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *PrecommitCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "precommit-test-*")
	s.Require().NoError(err)
}

func (s *PrecommitCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *PrecommitCheckTestSuite) TestIDAndName() {
	check := &PrecommitCheck{}
	s.Equal("common:precommit", check.ID())
	s.Equal("Pre-commit Hooks", check.Name())
}

func (s *PrecommitCheckTestSuite) TestPreCommitConfigYaml() {
	content := `repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.4.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".pre-commit-config.yaml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "pre-commit")
}

func (s *PrecommitCheckTestSuite) TestPreCommitConfigYml() {
	content := `repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.4.0
    hooks:
      - id: check-yaml
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".pre-commit-config.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "pre-commit")
}

func (s *PrecommitCheckTestSuite) TestHuskyDirectory() {
	huskyDir := filepath.Join(s.tempDir, ".husky")
	err := os.MkdirAll(huskyDir, 0755)
	s.Require().NoError(err)

	// Create a pre-commit hook file
	hookContent := `#!/bin/sh
npm run lint
`
	err = os.WriteFile(filepath.Join(huskyDir, "pre-commit"), []byte(hookContent), 0755)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Husky")
}

func (s *PrecommitCheckTestSuite) TestHuskyInPackageJson() {
	content := `{
  "name": "my-app",
  "version": "1.0.0",
  "devDependencies": {
    "husky": "^8.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Husky")
}

func (s *PrecommitCheckTestSuite) TestHuskyConfigInPackageJson() {
	content := `{
  "name": "my-app",
  "version": "1.0.0",
  "husky": {
    "hooks": {
      "pre-commit": "npm test"
    }
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Husky")
}

func (s *PrecommitCheckTestSuite) TestLefthookYml() {
	content := `pre-commit:
  commands:
    lint:
      run: npm run lint
`
	err := os.WriteFile(filepath.Join(s.tempDir, "lefthook.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Lefthook")
}

func (s *PrecommitCheckTestSuite) TestLefthookYaml() {
	content := `pre-commit:
  commands:
    test:
      run: go test ./...
`
	err := os.WriteFile(filepath.Join(s.tempDir, "lefthook.yaml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Lefthook")
}

func (s *PrecommitCheckTestSuite) TestOvercommit() {
	content := `PreCommit:
  TrailingWhitespace:
    enabled: true
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".overcommit.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Overcommit")
}

func (s *PrecommitCheckTestSuite) TestCommitlintConfigJs() {
	content := `module.exports = {
  extends: ['@commitlint/config-conventional']
};`
	err := os.WriteFile(filepath.Join(s.tempDir, "commitlint.config.js"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "commitlint")
}

func (s *PrecommitCheckTestSuite) TestCommitlintRc() {
	content := `{
  "extends": ["@commitlint/config-conventional"]
}`
	err := os.WriteFile(filepath.Join(s.tempDir, ".commitlintrc.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "commitlint")
}

func (s *PrecommitCheckTestSuite) TestCommitlintInPackageJson() {
	content := `{
  "name": "my-app",
  "version": "1.0.0",
  "commitlint": {
    "extends": ["@commitlint/config-conventional"]
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "commitlint")
}

func (s *PrecommitCheckTestSuite) TestLintStagedConfigJs() {
	content := `module.exports = {
  '*.js': ['eslint --fix', 'prettier --write']
};`
	err := os.WriteFile(filepath.Join(s.tempDir, "lint-staged.config.js"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "lint-staged")
}

func (s *PrecommitCheckTestSuite) TestLintStagedInPackageJson() {
	content := `{
  "name": "my-app",
  "version": "1.0.0",
  "lint-staged": {
    "*.js": ["eslint --fix"]
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "lint-staged")
}

func (s *PrecommitCheckTestSuite) TestGitHooksDirectory() {
	hooksDir := filepath.Join(s.tempDir, ".git", "hooks")
	err := os.MkdirAll(hooksDir, 0755)
	s.Require().NoError(err)

	// Create an executable pre-commit hook
	hookContent := `#!/bin/sh
echo "Running pre-commit hook"
`
	hookPath := filepath.Join(hooksDir, "pre-commit")
	err = os.WriteFile(hookPath, []byte(hookContent), 0755)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "git hooks")
}

func (s *PrecommitCheckTestSuite) TestGitHooksSampleIgnored() {
	hooksDir := filepath.Join(s.tempDir, ".git", "hooks")
	err := os.MkdirAll(hooksDir, 0755)
	s.Require().NoError(err)

	// Create only sample hooks (should be ignored)
	hookContent := `#!/bin/sh
echo "Sample hook"
`
	err = os.WriteFile(filepath.Join(hooksDir, "pre-commit.sample"), []byte(hookContent), 0755)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No pre-commit hooks configured")
}

func (s *PrecommitCheckTestSuite) TestMultipleTools() {
	// Set up pre-commit
	precommitContent := `repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.4.0
    hooks:
      - id: trailing-whitespace
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".pre-commit-config.yaml"), []byte(precommitContent), 0644)
	s.Require().NoError(err)

	// Set up commitlint
	commitlintContent := `module.exports = { extends: ['@commitlint/config-conventional'] };`
	err = os.WriteFile(filepath.Join(s.tempDir, "commitlint.config.js"), []byte(commitlintContent), 0644)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "pre-commit")
	s.Contains(result.Message, "commitlint")
}

func (s *PrecommitCheckTestSuite) TestNoPrecommitHooks() {
	// Create some other files but no pre-commit config
	err := os.WriteFile(filepath.Join(s.tempDir, "README.md"), []byte("# My Project"), 0644)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No pre-commit hooks configured")
}

func (s *PrecommitCheckTestSuite) TestEmptyDirectory() {
	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No pre-commit hooks configured")
}

func (s *PrecommitCheckTestSuite) TestHuskyWithCommitMsg() {
	huskyDir := filepath.Join(s.tempDir, ".husky")
	err := os.MkdirAll(huskyDir, 0755)
	s.Require().NoError(err)

	// Create a commit-msg hook file
	hookContent := `#!/bin/sh
npx --no-install commitlint --edit "$1"
`
	err = os.WriteFile(filepath.Join(huskyDir, "commit-msg"), []byte(hookContent), 0755)
	s.Require().NoError(err)

	check := &PrecommitCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Husky")
}

func TestPrecommitCheckTestSuite(t *testing.T) {
	suite.Run(t, new(PrecommitCheckTestSuite))
}
