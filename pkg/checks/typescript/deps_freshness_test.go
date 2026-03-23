package typescriptcheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type DepsFreshnessCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *DepsFreshnessCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "deps-freshness-ts-test-*")
	s.Require().NoError(err)
}

func (s *DepsFreshnessCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *DepsFreshnessCheckTestSuite) TestIDAndName() {
	check := &DepsFreshnessCheck{}
	s.Equal("typescript:deps_freshness", check.ID())
	s.Equal("TypeScript Dependency Freshness", check.Name())
}

func (s *DepsFreshnessCheckTestSuite) TestNoPackageJSON() {
	check := &DepsFreshnessCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Reason, "package.json")
}

func (s *DepsFreshnessCheckTestSuite) TestCountOutdatedTsPackages_None() {
	count, err := countOutdatedTsPackages("{}")
	s.NoError(err)
	s.Equal(0, count)
}

func (s *DepsFreshnessCheckTestSuite) TestCountOutdatedTsPackages_Some() {
	output := `{"typescript":{"current":"5.0.0","wanted":"5.3.0","latest":"5.3.0"}}`
	count, err := countOutdatedTsPackages(output)
	s.NoError(err)
	s.Equal(1, count)
}

func (s *DepsFreshnessCheckTestSuite) TestCountOutdatedTsPackages_Empty() {
	count, err := countOutdatedTsPackages("")
	s.NoError(err)
	s.Equal(0, count)
}

func (s *DepsFreshnessCheckTestSuite) TestDetectPackageManager() {
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(`{"name":"test"}`), 0644)
	s.Require().NoError(err)

	check := &DepsFreshnessCheck{}
	s.Equal("npm", check.detectPackageManager(s.tempDir))
}

func TestDepsFreshnessCheckTestSuite(t *testing.T) {
	suite.Run(t, new(DepsFreshnessCheckTestSuite))
}
