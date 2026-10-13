package gocheck

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DepsFreshnessCheckTestSuite struct {
	suite.Suite
}

func (s *DepsFreshnessCheckTestSuite) TestIDAndName() {
	check := &DepsFreshnessCheck{}
	s.Equal("go:deps_freshness", check.ID())
	s.Equal("Go Dependency Freshness", check.Name())
}

func (s *DepsFreshnessCheckTestSuite) TestCountOutdatedGoModules_NoUpdates() {
	output := `github.com/ipedrazas/a2
github.com/foo/bar v1.0.0
github.com/baz/qux v2.3.1`
	s.Equal(0, countOutdatedGoModules(output))
}

func (s *DepsFreshnessCheckTestSuite) TestCountOutdatedGoModules_WithUpdates() {
	output := `github.com/ipedrazas/a2
github.com/foo/bar v1.0.0 [v1.1.0]
github.com/baz/qux v2.3.1
github.com/some/pkg v0.5.0 [v1.0.0]`
	s.Equal(2, countOutdatedGoModules(output))
}

func (s *DepsFreshnessCheckTestSuite) TestCountOutdatedGoModules_Empty() {
	s.Equal(0, countOutdatedGoModules(""))
}

// TestCountOutdatedGoModules_DirectFormat covers the -f output, which prints
// an empty line for every module that does not match the template.
func (s *DepsFreshnessCheckTestSuite) TestCountOutdatedGoModules_DirectFormat() {
	output := "\n\ngithub.com/foo/bar v1.0.0 [v1.1.0]\n\n\n"
	s.Equal(1, countOutdatedGoModules(output))
}

func TestDepsFreshnessCheckTestSuite(t *testing.T) {
	suite.Run(t, new(DepsFreshnessCheckTestSuite))
}
