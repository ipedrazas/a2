package security

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// PatternsTestSuite is test suite for security patterns.
type PatternsTestSuite struct {
	suite.Suite
}

// TestCalculateEntropy tests the entropy calculation function.
func (suite *PatternsTestSuite) TestCalculateEntropy() {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "low entropy (repeating)",
			input:    "aaaaaaaaaaaaaaaaaaaa",
			expected: 0,
		},
		{
			name:     "medium entropy (text)",
			input:    "hello world",
			expected: 3.18, // Approximately
		},
		{
			name:     "high entropy (random)",
			input:    "8f4a2e1c9b7d3a6f",
			expected: 4.0,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := calculateEntropy(tt.input)
			// Allow small tolerance for floating point comparison
			if tt.expected > 0 {
				suite.InDelta(tt.expected, result, 0.5)
			} else {
				suite.Equal(tt.expected, result)
			}
		})
	}
}

// TestIsHighEntropy tests the high entropy detection function.
func (suite *PatternsTestSuite) TestIsHighEntropy() {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "short string",
			input:    "short",
			expected: false,
		},
		{
			name:     "normal text",
			input:    "This is normal text with some variety",
			expected: false,
		},
		{
			name:     "repeating pattern",
			input:    "abcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabc",
			expected: false,
		},
		{
			name: "longer repeating pattern",
			input: "SGVsbG8gV29ybGQgSGVsbG8gV29ybGQgSGVsbG8gV29ybGQgSGVsbG8gV29ybGQg" +
				"SGVsbG8gV29ybGQgSGVsbG8gV29ybGQgSGVsbG8gV29ybGQgSGVsbG8gV29ybGQg",
			expected: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := isHighEntropy(tt.input)
			suite.Equal(tt.expected, result)
		})
	}
}

// TestIsBase64 tests the base64 detection function.
func (suite *PatternsTestSuite) TestIsBase64() {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "not base64",
			input:    "not base64!",
			expected: false,
		},
		{
			name:     "valid base64 with padding",
			input:    "SGVsbG8gV29ybGQ=",
			expected: true,
		},
		{
			name:     "base64 without padding",
			input:    "SGVsbG8gV29ybGQ",
			expected: false, // Our implementation requires proper padding or decode
		},
		{
			name:     "empty string",
			input:    "",
			expected: false, // Empty string should return false
		},
		{
			name:     "properly padded short string",
			input:    "YWJj",
			expected: true, // "abc" in base64
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := isBase64(tt.input)
			suite.Equal(tt.expected, result)
		})
	}
}

// TestIsHex tests the hex detection function.
func (suite *PatternsTestSuite) TestIsHex() {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "not hex",
			input:    "not hex!",
			expected: false,
		},
		{
			name:     "valid hex",
			input:    "48656c6c6f20576f726c64",
			expected: true,
		},
		{
			name:     "odd length",
			input:    "abc",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "short hex",
			input:    "abcd",
			expected: false, // Too short
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := isHex(tt.input)
			suite.Equal(tt.expected, result)
		})
	}
}

// TestHasHexEscapeSequence tests the hex escape sequence detection.
func (suite *PatternsTestSuite) TestHasHexEscapeSequence() {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "no escape sequences",
			input:    "normal string",
			expected: false,
		},
		{
			name:     "one escape sequence",
			input:    "\\x48",
			expected: false, // Less than 3 escapes
		},
		{
			name:     "two escape sequences",
			input:    "\\x48\\x65",
			expected: false, // Less than 3 escapes
		},
		{
			name:     "many escape sequences",
			input:    "\\x48\\x65\\x6c\\x6c\\x6f\\x20\\x57\\x6f\\x72\\x6c\\x64",
			expected: true,
		},
		{
			name:     "unicode escapes",
			input:    "\\u0048\\u0065\\u006c\\u006c",
			expected: true, // 4 unicode escapes
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := hasHexEscapeSequence(tt.input)
			suite.Equal(tt.expected, result)
		})
	}
}

// TestDetectLanguageFromPath tests the language detection from file path.
func (suite *PatternsTestSuite) TestDetectLanguageFromPath() {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "go file",
			input:    "main.go",
			expected: "go",
		},
		{
			name:     "python file",
			input:    "app.py",
			expected: "python",
		},
		{
			name:     "javascript file",
			input:    "index.js",
			expected: "node",
		},
		{
			name:     "typescript file",
			input:    "app.ts",
			expected: "typescript",
		},
		{
			name:     "java file",
			input:    "Main.java",
			expected: "java",
		},
		{
			name:     "unknown file",
			input:    "README.md",
			expected: "",
		},
		{
			name:     "rust file",
			input:    "main.rs",
			expected: "rust",
		},
		{
			name:     "swift file",
			input:    "app.swift",
			expected: "swift",
		},
		{
			name:     "ruby file",
			input:    "app.rb",
			expected: "ruby",
		},
		{
			name:     "php file",
			input:    "index.php",
			expected: "php",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := detectLanguageFromPath(tt.input)
			suite.Equal(tt.expected, result)
		})
	}
}

// TestIsCommentLine tests the comment line detection for various languages.
func (suite *PatternsTestSuite) TestIsCommentLine() {
	tests := []struct {
		name     string
		input    string
		language string
		expected bool
	}{
		{
			name:     "go single line comment",
			input:    "// this is a comment",
			language: "go",
			expected: true,
		},
		{
			name:     "go code",
			input:    "func main() {",
			language: "go",
			expected: false,
		},
		{
			name:     "python comment",
			input:    "# this is a comment",
			language: "python",
			expected: true,
		},
		{
			name:     "python code",
			input:    "def main():",
			language: "python",
			expected: false,
		},
		{
			name:     "javascript comment",
			input:    "// this is a comment",
			language: "node",
			expected: true,
		},
		{
			name:     "javascript block comment",
			input:    "/* this is a comment */",
			language: "node",
			expected: true,
		},
		{
			name:     "ruby comment",
			input:    "# this is a comment",
			language: "ruby",
			expected: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := isCommentLine(tt.input, tt.language)
			suite.Equal(tt.expected, result)
		})
	}
}

// TestIsTestFile tests the test file detection.
func (suite *PatternsTestSuite) TestIsTestFile() {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "go test file",
			input:    "main_test.go",
			expected: true,
		},
		{
			name:     "python test file",
			input:    "test_main.py",
			expected: true,
		},
		{
			name:     "javascript test file",
			input:    "main.test.js",
			expected: true,
		},
		{
			name:     "typescript spec file",
			input:    "main.spec.ts",
			expected: true,
		},
		{
			name:     "normal file",
			input:    "main.go",
			expected: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := isTestFile(tt.input)
			suite.Equal(tt.expected, result)
		})
	}
}

// TestShouldSkipFile tests the file skipping logic.
func (suite *PatternsTestSuite) TestShouldSkipFile() {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "example file",
			input:    "example.go",
			expected: true,
		},
		{
			name:     "sample file",
			input:    "sample.py",
			expected: true,
		},
		{
			name:     "template file",
			input:    "template.js",
			expected: true,
		},
		{
			name:     "normal file",
			input:    "main.go",
			expected: false,
		},
		{
			name:     "README",
			input:    "README.md",
			expected: true, // .md files are skipped (not source files)
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := shouldSkipFile(tt.input)
			suite.Equal(tt.expected, result)
		})
	}
}

// TestIsSourceFile tests the source file detection.
func (suite *PatternsTestSuite) TestIsSourceFile() {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "go file",
			input:    "main.go",
			expected: true,
		},
		{
			name:     "python file",
			input:    "app.py",
			expected: true,
		},
		{
			name:     "javascript file",
			input:    "index.js",
			expected: true,
		},
		{
			name:     "config file",
			input:    "config.yaml",
			expected: true,
		},
		{
			name:     "markdown file",
			input:    "README.md",
			expected: false,
		},
		{
			name:     "text file",
			input:    "notes.txt",
			expected: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := isSourceFile(tt.input)
			suite.Equal(tt.expected, result)
		})
	}
}

// TestFinding tests the Finding struct and its String method.
func (suite *PatternsTestSuite) TestFinding() {
	finding := Finding{
		Type:        "test_type",
		File:        "test.go",
		Line:        42,
		Description: "test finding",
		Severity:    "high",
	}

	suite.Equal("test_type", finding.Type)
	suite.Equal("test.go", finding.File)
	suite.Equal(42, finding.Line)
	suite.Equal("test finding", finding.Description)
	suite.Equal("high", finding.Severity)
	suite.Equal("test finding in test.go:42", finding.String())
}

// TestFormatFindings tests the FormatFindings function.
func (suite *PatternsTestSuite) TestFormatFindings() {
	tests := []struct {
		name     string
		findings []Finding
		maxItems int
		expected string
	}{
		{
			name:     "no findings",
			findings: []Finding{},
			maxItems: 5,
			expected: "No issues found",
		},
		{
			name: "single finding",
			findings: []Finding{
				{Type: "test", File: "test.go", Line: 1, Description: "issue", Severity: "low"},
			},
			maxItems: 5,
			expected: "issue in test.go:1",
		},
		{
			name: "multiple findings within limit",
			findings: []Finding{
				{Type: "test", File: "test.go", Line: 1, Description: "issue1", Severity: "low"},
				{Type: "test", File: "test.go", Line: 2, Description: "issue2", Severity: "low"},
			},
			maxItems: 5,
			expected: "issue1 in test.go:1, issue2 in test.go:2",
		},
		{
			name: "findings exceed limit",
			findings: []Finding{
				{Type: "test", File: "test.go", Line: 1, Description: "issue1", Severity: "low"},
				{Type: "test", File: "test.go", Line: 2, Description: "issue2", Severity: "low"},
				{Type: "test", File: "test.go", Line: 3, Description: "issue3", Severity: "low"},
			},
			maxItems: 2,
			expected: "issue1 in test.go:1, issue2 in test.go:2 (1 more)",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := FormatFindings(tt.findings, tt.maxItems)
			suite.Equal(tt.expected, result)
		})
	}
}

// TestPatternsTestSuite runs all tests in suite.
func TestPatternsTestSuite(t *testing.T) {
	suite.Run(t, new(PatternsTestSuite))
}
