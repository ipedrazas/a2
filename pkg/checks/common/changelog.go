package common

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// ChangelogCheck verifies that a changelog or release notes exist.
type ChangelogCheck struct{}

func (c *ChangelogCheck) ID() string   { return "common:changelog" }
func (c *ChangelogCheck) Name() string { return "Changelog" }

// Run checks for changelog files and release tooling configuration.
func (c *ChangelogCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	// Check for changelog files
	changelogFiles := []string{
		"CHANGELOG.md",
		"CHANGELOG.txt",
		"CHANGELOG",
		"CHANGES.md",
		"CHANGES.txt",
		"CHANGES",
		"HISTORY.md",
		"HISTORY.txt",
		"HISTORY",
		"NEWS.md",
		"NEWS.txt",
		"NEWS",
		"RELEASES.md",
		"RELEASE_NOTES.md",
	}

	var foundChangelog string
	var changelogContent []byte
	for _, file := range changelogFiles {
		if safepath.Exists(path, file) {
			foundChangelog = file
			// Try to read content for format detection
			content, err := safepath.ReadFile(path, file)
			if err == nil {
				changelogContent = content
			}
			break
		}
	}

	// Check for release tooling (indicates automated changelog/versioning)
	releaseTools := []struct {
		name string
		file string
	}{
		{"GoReleaser", ".goreleaser.yml"},
		{"GoReleaser", ".goreleaser.yaml"},
		{"GoReleaser", "goreleaser.yml"},
		{"GoReleaser", "goreleaser.yaml"},
		{"semantic-release", ".releaserc"},
		{"semantic-release", ".releaserc.json"},
		{"semantic-release", ".releaserc.yaml"},
		{"semantic-release", ".releaserc.yml"},
		{"semantic-release", "release.config.js"},
		{"semantic-release", "release.config.cjs"},
		{"release-please", "release-please-config.json"},
		{"release-please", ".release-please-manifest.json"},
		{"standard-version", ".versionrc"},
		{"standard-version", ".versionrc.json"},
		{"changesets", ".changeset/config.json"},
	}

	var foundReleaseTools []string
	for _, tool := range releaseTools {
		if safepath.Exists(path, tool.file) {
			foundReleaseTools = append(foundReleaseTools, tool.name)
		}
	}

	// Deduplicate release tools
	foundReleaseTools = uniqueStrings(foundReleaseTools)

	// Determine result based on findings
	if foundChangelog != "" {
		// Check if it follows Keep a Changelog format
		format := c.detectChangelogFormat(changelogContent)

		if len(foundReleaseTools) > 0 {
			return rb.Pass(foundChangelog + " found (" + format + "), " + strings.Join(foundReleaseTools, ", ") + " configured"), nil
		}
		return rb.Pass(foundChangelog + " found (" + format + ")"), nil
	}

	// No changelog file, but release tooling is configured
	if len(foundReleaseTools) > 0 {
		return rb.Pass("Release tooling configured: " + strings.Join(foundReleaseTools, ", ")), nil
	}

	// Nothing found
	return rb.Warn("No changelog found (consider adding CHANGELOG.md)"), nil
}

// detectChangelogFormat checks if the changelog follows a known format.
func (c *ChangelogCheck) detectChangelogFormat(content []byte) string {
	if len(content) == 0 {
		return "unknown format"
	}

	contentStr := string(content)

	// Check for Keep a Changelog format
	// https://keepachangelog.com
	keepAChangelogMarkers := []string{
		"## [Unreleased]",
		"## [unreleased]",
		"### Added",
		"### Changed",
		"### Deprecated",
		"### Removed",
		"### Fixed",
		"### Security",
		"keepachangelog.com",
	}

	keepAChangelogScore := 0
	for _, marker := range keepAChangelogMarkers {
		if strings.Contains(contentStr, marker) {
			keepAChangelogScore++
		}
	}

	if keepAChangelogScore >= 2 {
		return "Keep a Changelog format"
	}

	// Check for Conventional Changelog format
	conventionalMarkers := []string{
		"## [",
		"### Features",
		"### Bug Fixes",
		"### BREAKING CHANGES",
		"feat:",
		"fix:",
	}

	conventionalScore := 0
	for _, marker := range conventionalMarkers {
		if strings.Contains(contentStr, marker) {
			conventionalScore++
		}
	}

	if conventionalScore >= 2 {
		return "Conventional Changelog format"
	}

	// Check for basic markdown structure
	if strings.Contains(contentStr, "## ") || strings.Contains(contentStr, "# ") {
		return "markdown format"
	}

	return "plain text"
}
