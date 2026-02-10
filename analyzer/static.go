package analyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type Issue struct {
	File     string
	Line     int
	Message  string
	Severity string // "error", "warning", "info"
}

func RunStaticAnalysis(rootPath string) ([]Issue, error) {
	var issues []Issue
	fset := token.NewFileSet()

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			fileIssues, err := analyzeFile(path, fset)
			if err != nil {
				return err // Propagate error, let caller log it
			}
			issues = append(issues, fileIssues...)
		}
		return nil
	})

	return issues, err
}

func analyzeFile(path string, fset *token.FileSet) ([]Issue, error) {
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err // Return error directly
	}

	var issues []Issue

	ast.Inspect(node, func(n ast.Node) bool {
		// Example check: Detect usage of "fmt.Println" (discouraged in prod in favor of logging)
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if pkg, ok := sel.X.(*ast.Ident); ok {
					if pkg.Name == "fmt" && (sel.Sel.Name == "Println" || sel.Sel.Name == "Printf") {
						issues = append(issues, Issue{
							File:     path,
							Line:     fset.Position(n.Pos()).Line,
							Message:  "Avoid using fmt.Print* in production code; use a structured logger instead.",
							Severity: "warning",
						})
					}
				}
			}
		}
		// Example check: Detect TODO comments
		return true
	})

	for _, commentGroup := range node.Comments {
		for _, comment := range commentGroup.List {
			if strings.Contains(comment.Text, "TODO") {
				issues = append(issues, Issue{
					File:     path,
					Line:     fset.Position(comment.Pos()).Line,
					Message:  "Found TODO comment: " + comment.Text,
					Severity: "info",
				})
			}
		}
	}

	return issues, nil
}
