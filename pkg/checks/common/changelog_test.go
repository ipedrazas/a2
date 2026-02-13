package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type ChangelogCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *ChangelogCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "changelog-test-*")
	s.Require().NoError(err)
}

func (s *ChangelogCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *ChangelogCheckTestSuite) TestIDAndName() {
	check := &ChangelogCheck{}
	s.Equal("common:changelog", check.ID())
	s.Equal("Changelog", check.Name())
}

func (s *ChangelogCheckTestSuite) TestChangelogMdExists() {
	content := `# Changelog

## [Unreleased]

### Added
- New feature coming soon

## [1.0.0] - 2024-01-01

### Added
- Initial release
`
	err := os.WriteFile(filepath.Join(s.tempDir, "CHANGELOG.md"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ChangelogCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "CHANGELOG.md")
	s.Contains(result.Reason, "Keep a Changelog")
}

func (s *ChangelogCheckTestSuite) TestChangesFileExists() {
	content := `# Changes

v1.0.0 - Initial release
`
	err := os.WriteFile(filepath.Join(s.tempDir, "CHANGES.md"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ChangelogCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "CHANGES.md")
}

func (s *ChangelogCheckTestSuite) TestHistoryFileExists() {
	content := `History
=======

1.0.0 - Initial release
`
	err := os.WriteFile(filepath.Join(s.tempDir, "HISTORY.md"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ChangelogCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "HISTORY.md")
}

func (s *ChangelogCheckTestSuite) TestNewsFileExists() {
	content := `NEWS
====
`
	err := os.WriteFile(filepath.Join(s.tempDir, "NEWS"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ChangelogCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "NEWS")
}

func (s *ChangelogCheckTestSuite) TestGoReleaserConfigured() {
	content := `project_name: myapp
builds:
  - main: .
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".goreleaser.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ChangelogCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "GoReleaser")
}

func (s *ChangelogCheckTestSuite) TestSemanticReleaseConfigured() {
	content := `{
  "branches": ["main"]
}`
	err := os.WriteFile(filepath.Join(s.tempDir, ".releaserc"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ChangelogCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "semantic-release")
}

func (s *ChangelogCheckTestSuite) TestReleasePleaseConfigured() {
	content := `{
  "packages": {}
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "release-please-config.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ChangelogCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "release-please")
}

func (s *ChangelogCheckTestSuite) TestChangesetsConfigured() {
	changesetDir := filepath.Join(s.tempDir, ".changeset")
	err := os.MkdirAll(changesetDir, 0755)
	s.Require().NoError(err)

	content := `{
  "changelog": "@changesets/cli/changelog"
}`
	err = os.WriteFile(filepath.Join(changesetDir, "config.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ChangelogCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "changesets")
}

func (s *ChangelogCheckTestSuite) TestChangelogWithGoReleaser() {
	// Both changelog and release tool
	changelog := `# Changelog

## [1.0.0]
- Initial release
`
	err := os.WriteFile(filepath.Join(s.tempDir, "CHANGELOG.md"), []byte(changelog), 0644)
	s.Require().NoError(err)

	goreleaser := `project_name: myapp`
	err = os.WriteFile(filepath.Join(s.tempDir, ".goreleaser.yaml"), []byte(goreleaser), 0644)
	s.Require().NoError(err)

	check := &ChangelogCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "CHANGELOG.md")
	s.Contains(result.Reason, "GoReleaser")
}

func (s *ChangelogCheckTestSuite) TestKeepAChangelogFormat() {
	content := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- New feature

### Changed
- Updated something

### Fixed
- Bug fix

## [1.0.0] - 2024-01-01

### Added
- Initial release
`
	err := os.WriteFile(filepath.Join(s.tempDir, "CHANGELOG.md"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ChangelogCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Keep a Changelog")
}

func (s *ChangelogCheckTestSuite) TestConventionalChangelogFormat() {
	content := `# Changelog

## [1.1.0] - 2024-02-01

### Features

* feat: add new login page
* feat: add user dashboard

### Bug Fixes

* fix: resolve memory leak
* fix: correct validation error

## [1.0.0] - 2024-01-01

### Features

* Initial release
`
	err := os.WriteFile(filepath.Join(s.tempDir, "CHANGELOG.md"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ChangelogCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Conventional Changelog")
}

func (s *ChangelogCheckTestSuite) TestPlainTextFormat() {
	content := `Version 1.0.0
- Added initial features
- Fixed some bugs
`
	err := os.WriteFile(filepath.Join(s.tempDir, "CHANGELOG.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ChangelogCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "plain text")
}

func (s *ChangelogCheckTestSuite) TestNoChangelog() {
	// Create some other file but no changelog
	err := os.WriteFile(filepath.Join(s.tempDir, "README.md"), []byte("# My Project"), 0644)
	s.Require().NoError(err)

	check := &ChangelogCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No changelog found")
}

func (s *ChangelogCheckTestSuite) TestEmptyDirectory() {
	check := &ChangelogCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No changelog found")
}

func (s *ChangelogCheckTestSuite) TestReleaseNotesFile() {
	content := `# Release Notes

## v1.0.0
Initial release
`
	err := os.WriteFile(filepath.Join(s.tempDir, "RELEASE_NOTES.md"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ChangelogCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "RELEASE_NOTES.md")
}

func (s *ChangelogCheckTestSuite) TestStandardVersionConfigured() {
	content := `{
  "types": [
    {"type": "feat", "section": "Features"},
    {"type": "fix", "section": "Bug Fixes"}
  ]
}`
	err := os.WriteFile(filepath.Join(s.tempDir, ".versionrc.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ChangelogCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "standard-version")
}

func TestChangelogCheckTestSuite(t *testing.T) {
	suite.Run(t, new(ChangelogCheckTestSuite))
}
