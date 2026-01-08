package javacheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type CoverageCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *CoverageCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "java-coverage-test-*")
	s.Require().NoError(err)
}

func (s *CoverageCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *CoverageCheckTestSuite) TestIDAndName() {
	check := &CoverageCheck{}
	s.Equal("java:coverage", check.ID())
	s.Equal("Java Coverage", check.Name())
}

func (s *CoverageCheckTestSuite) TestRun_NoJaCoCo() {
	check := &CoverageCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "JaCoCo not configured")
}

func (s *CoverageCheckTestSuite) TestRun_ResultLanguage() {
	check := &CoverageCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangJava, result.Language)
}

func (s *CoverageCheckTestSuite) TestRun_JaCoCoConfigured_NoReport() {
	content := `<project>
  <build>
    <plugins>
      <plugin>
        <groupId>org.jacoco</groupId>
        <artifactId>jacoco-maven-plugin</artifactId>
      </plugin>
    </plugins>
  </build>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &CoverageCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "no coverage report found")
}

func (s *CoverageCheckTestSuite) TestRun_JaCoCoConfigured_Gradle() {
	content := `plugins {
    id 'java'
    id 'jacoco'
}

test {
    finalizedBy jacocoTestReport
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &CoverageCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "no coverage report found")
}

func (s *CoverageCheckTestSuite) TestRun_JaCoCoConfigured_GradleKts() {
	content := `plugins {
    java
    jacoco
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle.kts"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &CoverageCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Contains(result.Message, "no coverage report found")
}

func (s *CoverageCheckTestSuite) TestRun_WithReport_AboveThreshold() {
	// Configure JaCoCo
	pomContent := `<project>
  <build>
    <plugins>
      <plugin>
        <groupId>org.jacoco</groupId>
        <artifactId>jacoco-maven-plugin</artifactId>
      </plugin>
    </plugins>
  </build>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(pomContent), 0644)
	s.Require().NoError(err)

	// Create JaCoCo report with 85% coverage
	reportDir := filepath.Join(s.tempDir, "target", "site", "jacoco")
	err = os.MkdirAll(reportDir, 0755)
	s.Require().NoError(err)

	reportContent := `<?xml version="1.0" encoding="UTF-8"?>
<report name="myapp">
  <counter type="LINE" missed="15" covered="85"/>
  <counter type="INSTRUCTION" missed="20" covered="80"/>
</report>`
	err = os.WriteFile(filepath.Join(reportDir, "jacoco.xml"), []byte(reportContent), 0644)
	s.Require().NoError(err)

	check := &CoverageCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "85.0%")
	s.Contains(result.Message, "80")
}

func (s *CoverageCheckTestSuite) TestRun_WithReport_BelowThreshold() {
	// Configure JaCoCo
	pomContent := `<project>
  <build>
    <plugins>
      <plugin>
        <groupId>org.jacoco</groupId>
        <artifactId>jacoco-maven-plugin</artifactId>
      </plugin>
    </plugins>
  </build>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(pomContent), 0644)
	s.Require().NoError(err)

	// Create JaCoCo report with 60% coverage (below 80% threshold)
	reportDir := filepath.Join(s.tempDir, "target", "site", "jacoco")
	err = os.MkdirAll(reportDir, 0755)
	s.Require().NoError(err)

	reportContent := `<?xml version="1.0" encoding="UTF-8"?>
<report name="myapp">
  <counter type="LINE" missed="40" covered="60"/>
</report>`
	err = os.WriteFile(filepath.Join(reportDir, "jacoco.xml"), []byte(reportContent), 0644)
	s.Require().NoError(err)

	check := &CoverageCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "60.0%")
	s.Contains(result.Message, "below threshold")
}

func (s *CoverageCheckTestSuite) TestRun_CustomThreshold() {
	// Configure JaCoCo
	pomContent := `<project>
  <build>
    <plugins>
      <plugin>
        <groupId>org.jacoco</groupId>
        <artifactId>jacoco-maven-plugin</artifactId>
      </plugin>
    </plugins>
  </build>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(pomContent), 0644)
	s.Require().NoError(err)

	// Create JaCoCo report with 60% coverage
	reportDir := filepath.Join(s.tempDir, "target", "site", "jacoco")
	err = os.MkdirAll(reportDir, 0755)
	s.Require().NoError(err)

	reportContent := `<?xml version="1.0" encoding="UTF-8"?>
<report name="myapp">
  <counter type="LINE" missed="40" covered="60"/>
</report>`
	err = os.WriteFile(filepath.Join(reportDir, "jacoco.xml"), []byte(reportContent), 0644)
	s.Require().NoError(err)

	// Custom threshold of 50%
	check := &CoverageCheck{Threshold: 50.0}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "60.0%")
	s.Contains(result.Message, "50.0%")
}

func (s *CoverageCheckTestSuite) TestRun_GradleReportLocation() {
	// Configure JaCoCo via Gradle
	gradleContent := `plugins {
    id 'java'
    id 'jacoco'
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(gradleContent), 0644)
	s.Require().NoError(err)

	// Create JaCoCo report at Gradle location
	reportDir := filepath.Join(s.tempDir, "build", "reports", "jacoco", "test")
	err = os.MkdirAll(reportDir, 0755)
	s.Require().NoError(err)

	reportContent := `<?xml version="1.0" encoding="UTF-8"?>
<report name="myapp">
  <counter type="INSTRUCTION" missed="10" covered="90"/>
</report>`
	err = os.WriteFile(filepath.Join(reportDir, "jacocoTestReport.xml"), []byte(reportContent), 0644)
	s.Require().NoError(err)

	check := &CoverageCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "90.0%")
}

func (s *CoverageCheckTestSuite) TestParseCoverageReport_InvalidXML() {
	check := &CoverageCheck{}

	// Create invalid XML file
	reportDir := filepath.Join(s.tempDir, "target", "site", "jacoco")
	err := os.MkdirAll(reportDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(reportDir, "jacoco.xml"), []byte("invalid xml"), 0644)
	s.Require().NoError(err)

	coverage, found := check.parseCoverageReport(filepath.Join(reportDir, "jacoco.xml"))
	s.False(found)
	s.Equal(0.0, coverage)
}

func (s *CoverageCheckTestSuite) TestParseCoverageReport_NoCounters() {
	check := &CoverageCheck{}

	// Create XML without LINE or INSTRUCTION counters
	reportDir := filepath.Join(s.tempDir, "target", "site", "jacoco")
	err := os.MkdirAll(reportDir, 0755)
	s.Require().NoError(err)

	reportContent := `<?xml version="1.0" encoding="UTF-8"?>
<report name="myapp">
  <counter type="BRANCH" missed="10" covered="20"/>
</report>`
	err = os.WriteFile(filepath.Join(reportDir, "jacoco.xml"), []byte(reportContent), 0644)
	s.Require().NoError(err)

	coverage, found := check.parseCoverageReport(filepath.Join(reportDir, "jacoco.xml"))
	s.False(found)
	s.Equal(0.0, coverage)
}

func (s *CoverageCheckTestSuite) TestParseCoverageReport_ZeroTotal() {
	check := &CoverageCheck{}

	// Create XML with zero totals
	reportDir := filepath.Join(s.tempDir, "target", "site", "jacoco")
	err := os.MkdirAll(reportDir, 0755)
	s.Require().NoError(err)

	reportContent := `<?xml version="1.0" encoding="UTF-8"?>
<report name="myapp">
  <counter type="LINE" missed="0" covered="0"/>
</report>`
	err = os.WriteFile(filepath.Join(reportDir, "jacoco.xml"), []byte(reportContent), 0644)
	s.Require().NoError(err)

	coverage, found := check.parseCoverageReport(filepath.Join(reportDir, "jacoco.xml"))
	s.False(found)
	s.Equal(0.0, coverage)
}

func TestCoverageCheckTestSuite(t *testing.T) {
	suite.Run(t, new(CoverageCheckTestSuite))
}
