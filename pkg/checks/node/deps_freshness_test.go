package nodecheck

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
	s.tempDir, err = os.MkdirTemp("", "deps-freshness-node-test-*")
	s.Require().NoError(err)
}

func (s *DepsFreshnessCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *DepsFreshnessCheckTestSuite) TestIDAndName() {
	check := &DepsFreshnessCheck{}
	s.Equal("node:deps_freshness", check.ID())
	s.Equal("Node Dependency Freshness", check.Name())
}

func (s *DepsFreshnessCheckTestSuite) TestNoPackageJSON() {
	check := &DepsFreshnessCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Reason, "package.json")
}

func (s *DepsFreshnessCheckTestSuite) TestCountOutdatedNodePackages_None() {
	count, err := countOutdatedNodePackages("{}")
	s.NoError(err)
	s.Equal(0, count)
}

func (s *DepsFreshnessCheckTestSuite) TestCountOutdatedNodePackages_Some() {
	output := `{"lodash":{"current":"4.17.20","wanted":"4.17.21","latest":"4.17.21"},"express":{"current":"4.17.0","wanted":"4.18.2","latest":"5.0.0"}}`
	count, err := countOutdatedNodePackages(output)
	s.NoError(err)
	s.Equal(2, count)
}

func (s *DepsFreshnessCheckTestSuite) TestCountOutdatedNodePackages_Empty() {
	count, err := countOutdatedNodePackages("")
	s.NoError(err)
	s.Equal(0, count)
}

func (s *DepsFreshnessCheckTestSuite) TestDetectPackageManager() {
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(`{"name":"test"}`), 0644)
	s.Require().NoError(err)

	check := &DepsFreshnessCheck{}
	s.Equal("npm", check.detectPackageManager(s.tempDir))

	// With yarn.lock
	err = os.WriteFile(filepath.Join(s.tempDir, "yarn.lock"), []byte(""), 0644)
	s.Require().NoError(err)
	s.Equal("yarn", check.detectPackageManager(s.tempDir))
}

func TestDepsFreshnessCheckTestSuite(t *testing.T) {
	suite.Run(t, new(DepsFreshnessCheckTestSuite))
}
