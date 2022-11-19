package main

import (
	"encoding/json"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

type (
	Type string

	Test struct {
		Name   string   `json:"name"`
		Type   Type     `json:"type"`
		Path   []string `json:"path"`
		LBrace int      `json:"lbrace"`
		RBrace int      `json:"rbrace"`
	}
)

const (
	TestType           Type = "test"
	SubtestType        Type = "subtest"
	DynamicSubtestType Type = "dynamicsubtest"
)

func main() {
	flag.Parse()

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	tests, err := generateOutline(file)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(os.Stdout).Encode(tests)
}

func generateOutline(src any) ([]*Test, error) {
	tests := make([]*Test, 0)

	f, err := parser.ParseFile(token.NewFileSet(), "", src, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	testingAlias := "testing"

	// Tranverse the AST and find the tests and sub-tests that exist
	ast.Inspect(f, func(n ast.Node) bool {
		switch t := n.(type) {

		// Figure out what the alias for testing library is
		case *ast.ImportSpec:
			if t.Path.Value == `"testing"` && t.Name != nil {
				testingAlias = t.Name.Name
			}

		// Find all the root test functions
		case *ast.FuncDecl:
			if strings.HasPrefix(t.Name.Name, "Test") {
				tests = append(tests, &Test{
					Name:   t.Name.Name,
					Type:   TestType,
					Path:   []string{},
					LBrace: int(t.Pos()),
					RBrace: int(t.End()),
				})
			}

		// Find all the sub-tests
		case *ast.CallExpr:
			sel, ok := t.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "Run" {
				return true
			}

			id, ok := sel.X.(*ast.Ident)
			if !ok {
				return true
			}

			field, ok := id.Obj.Decl.(*ast.Field)
			if !ok {
				return true
			}

			star, ok := field.Type.(*ast.StarExpr)
			if !ok {
				return true
			}

			sel, ok = star.X.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			if sel.Sel.Name != "T" {
				return true
			}

			id, ok = sel.X.(*ast.Ident)
			if !ok || id.Name != testingAlias {
				return true
			}

			// Subtests whose name is dynamic and not a fixed string cannot
			// be properly named, however we still list them in the outline
			// with an empty name.
			l, ok := t.Args[0].(*ast.BasicLit)
			if !ok {
				tests = append(tests, &Test{
					Name:   "",
					Type:   DynamicSubtestType,
					Path:   []string{},
					LBrace: int(t.Lparen),
					RBrace: int(t.Rparen),
				})

				return true
			}

			tests = append(tests, &Test{
				Name:   strings.Trim(l.Value, `"'`),
				Type:   SubtestType,
				Path:   []string{},
				LBrace: int(t.Lparen),
				RBrace: int(t.Rparen),
			})
		}

		return true
	})

	// Assemble the path to the sub-test
	for i, test := range tests {
		for j := 0; j < i; j++ {
			parent := tests[j]

			if parent.RBrace > test.LBrace {
				test.Path = append(test.Path, parent.Name)
			}
		}
	}

	return tests, nil
}
