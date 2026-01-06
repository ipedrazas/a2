package gocheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type CyclomaticCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *CyclomaticCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "cyclomatic-test-*")
	s.Require().NoError(err)

	// Create go.mod
	goMod := `module test
go 1.21
`
	err = os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(goMod), 0644)
	s.Require().NoError(err)
}

func (s *CyclomaticCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *CyclomaticCheckTestSuite) TestNoGoMod() {
	// Remove go.mod
	os.Remove(filepath.Join(s.tempDir, "go.mod"))

	check := &CyclomaticCheck{Threshold: 15}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Message, "go.mod not found")
}

func (s *CyclomaticCheckTestSuite) TestSimpleFunction() {
	// Create a simple function with low complexity
	code := `package main

func simple() int {
	return 42
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "main.go"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &CyclomaticCheck{Threshold: 15}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "No functions exceed complexity threshold")
}

func (s *CyclomaticCheckTestSuite) TestComplexFunction() {
	// Create a function with high complexity
	code := `package main

func complex(x, y, z int) int {
	result := 0
	if x > 0 {
		if y > 0 {
			result++
		} else if y < 0 {
			result--
		}
	} else if x < 0 {
		for i := 0; i < 10; i++ {
			if z > 0 && y > 0 {
				result += i
			} else if z < 0 || y < 0 {
				result -= i
			}
		}
	}
	switch {
	case z > 100:
		result = 100
	case z > 50:
		result = 50
	case z > 10:
		result = 10
	default:
		result = 0
	}
	return result
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "main.go"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &CyclomaticCheck{Threshold: 5}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "exceed complexity threshold")
	s.Contains(result.Message, "complex")
}

func (s *CyclomaticCheckTestSuite) TestSkipsTestFiles() {
	// Create a complex function in a test file
	code := `package main

func complexTest(x, y, z int) int {
	if x > 0 && y > 0 && z > 0 {
		if x > y {
			if y > z {
				return x + y + z
			} else if z > x {
				return z + y + x
			}
		}
	}
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if i > j || j > 5 {
				return i + j
			}
		}
	}
	return 0
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "main_test.go"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &CyclomaticCheck{Threshold: 5}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
}

func (s *CyclomaticCheckTestSuite) TestSkipsVendorDirectory() {
	// Create vendor directory with complex code
	vendorDir := filepath.Join(s.tempDir, "vendor", "some-lib")
	err := os.MkdirAll(vendorDir, 0755)
	s.Require().NoError(err)

	code := `package lib

func veryComplex(x int) int {
	if x > 0 && x < 10 || x > 20 && x < 30 {
		if x%2 == 0 {
			for i := 0; i < x; i++ {
				if i > 5 || i < 2 {
					return i
				}
			}
		}
	}
	return 0
}
`
	err = os.WriteFile(filepath.Join(vendorDir, "lib.go"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &CyclomaticCheck{Threshold: 5}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
}

func (s *CyclomaticCheckTestSuite) TestMethodComplexity() {
	// Test that method receivers are included in function name
	code := `package main

type Handler struct{}

func (h *Handler) Handle(x, y int) int {
	if x > 0 && y > 0 {
		if x > y || y > 100 {
			for i := 0; i < x; i++ {
				if i > 5 && i < 10 {
					return i
				}
			}
		}
	}
	return 0
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "handler.go"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &CyclomaticCheck{Threshold: 3}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "(Handler).Handle")
}

func (s *CyclomaticCheckTestSuite) TestDefaultThreshold() {
	// Test default threshold is used when 0
	check := &CyclomaticCheck{Threshold: 0}

	code := `package main

func simple() int {
	return 42
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "main.go"), []byte(code), 0644)
	s.Require().NoError(err)

	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "15") // Default threshold
}

func TestCyclomaticCheckTestSuite(t *testing.T) {
	suite.Run(t, new(CyclomaticCheckTestSuite))
}
