package gocheck

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
)

// CyclomaticCheck measures cyclomatic complexity of Go functions.
type CyclomaticCheck struct {
	Threshold int // Default: 15
}

func (c *CyclomaticCheck) ID() string   { return "go:cyclomatic" }
func (c *CyclomaticCheck) Name() string { return "Go Complexity" }

// FunctionComplexity holds complexity info for a single function.
type FunctionComplexity struct {
	Name       string
	File       string
	Line       int
	Complexity int
}

func (c *CyclomaticCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangGo)

	threshold := c.Threshold
	if threshold <= 0 {
		threshold = 15 // Default threshold
	}

	// Check if go.mod exists
	goModPath := filepath.Join(path, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return rb.Fail("go.mod not found"), nil
	}

	var complexFunctions []FunctionComplexity

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}

		// Only check .go files
		if !strings.HasSuffix(filePath, ".go") {
			return nil
		}

		// Skip test files
		if strings.HasSuffix(filePath, "_test.go") {
			return nil
		}

		// Parse the file
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, filePath, nil, 0)
		if err != nil {
			return nil // Skip files that fail to parse
		}

		// Analyze each function
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			complexity := calculateComplexity(fn)
			if complexity > threshold {
				relPath, _ := filepath.Rel(path, filePath)
				if relPath == "" {
					relPath = filePath
				}
				complexFunctions = append(complexFunctions, FunctionComplexity{
					Name:       functionName(fn),
					File:       relPath,
					Line:       fset.Position(fn.Pos()).Line,
					Complexity: complexity,
				})
			}
		}

		return nil
	})

	if err != nil {
		return rb.Warn("Error scanning files: " + err.Error()), nil
	}

	if len(complexFunctions) == 0 {
		return rb.Pass(fmt.Sprintf("No functions exceed complexity threshold (%d)", threshold)), nil
	}

	// Sort by complexity descending
	sort.Slice(complexFunctions, func(i, j int) bool {
		return complexFunctions[i].Complexity > complexFunctions[j].Complexity
	})

	// Build message with top offenders
	msg := checkutil.PluralizeCount(len(complexFunctions), "function exceeds", "functions exceed") +
		fmt.Sprintf(" complexity threshold (%d)", threshold)

	// Show top 3 offenders in message
	showCount := 3
	if len(complexFunctions) < showCount {
		showCount = len(complexFunctions)
	}

	for i := 0; i < showCount; i++ {
		f := complexFunctions[i]
		msg += fmt.Sprintf("\n  • %s (%s:%d) = %d", f.Name, f.File, f.Line, f.Complexity)
	}

	if len(complexFunctions) > showCount {
		msg += fmt.Sprintf("\n  ... and %d more", len(complexFunctions)-showCount)
	}

	// Build raw output with all complex functions
	var rawOutput strings.Builder
	rawOutput.WriteString(fmt.Sprintf("Functions exceeding complexity threshold (%d):\n", threshold))
	for _, f := range complexFunctions {
		rawOutput.WriteString(fmt.Sprintf("  %s:%d %s (complexity: %d)\n", f.File, f.Line, f.Name, f.Complexity))
	}

	return rb.WarnWithOutput(msg, rawOutput.String()), nil
}

// functionName returns the name of a function, including receiver if present.
func functionName(fn *ast.FuncDecl) string {
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		recv := fn.Recv.List[0]
		var typeName string
		switch t := recv.Type.(type) {
		case *ast.StarExpr:
			if ident, ok := t.X.(*ast.Ident); ok {
				typeName = ident.Name
			}
		case *ast.Ident:
			typeName = t.Name
		}
		if typeName != "" {
			return fmt.Sprintf("(%s).%s", typeName, fn.Name.Name)
		}
	}
	return fn.Name.Name
}

// calculateComplexity computes the cyclomatic complexity of a function.
// Cyclomatic complexity = 1 + number of decision points
// Decision points: if, for, case, &&, ||, select case
func calculateComplexity(fn *ast.FuncDecl) int {
	if fn.Body == nil {
		return 1
	}

	complexity := 1 // Base complexity

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.IfStmt:
			complexity++
		case *ast.ForStmt:
			complexity++
		case *ast.RangeStmt:
			complexity++
		case *ast.CaseClause:
			// Each case adds a decision point (except default)
			if node.List != nil {
				complexity++
			}
		case *ast.CommClause:
			// Select case (except default)
			if node.Comm != nil {
				complexity++
			}
		case *ast.BinaryExpr:
			// Logical operators add decision points
			if node.Op == token.LAND || node.Op == token.LOR {
				complexity++
			}
		}
		return true
	})

	return complexity
}
