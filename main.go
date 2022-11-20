package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
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
	var err error

	switch len(os.Args) {
	case 2:
		err = outlineCommand(os.Args[1], os.Stdout, func(path string) (io.ReadCloser, error) {
			return os.Open(path)
		})
	default:
		err = helpCommand(os.Stdout)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func outlineCommand(filepath string, w io.Writer, open func(string) (io.ReadCloser, error)) error {
	file, err := open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	tests, err := generateOutline(file)
	if err != nil {
		log.Fatal(err)
	}

	return json.NewEncoder(w).Encode(tests)
}

func helpCommand(w io.Writer) error {
	_, err := fmt.Fprint(w, `gotestoutline is a tool to generate go test outlines

Usage:
	gotestoutline <go-test-file-path>
`)

	return err
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

			test := &Test{
				Path:   []string{},
				Type:   SubtestType,
				LBrace: int(t.Lparen),
				RBrace: int(t.Rparen),
			}

			// Report cases where sub test name is dynamic and not a fixed
			// string, as IDEs might still want to do something with it.
			l, ok := t.Args[0].(*ast.BasicLit)
			if ok {
				test.Name = strings.Trim(l.Value, `"'`)
			} else {
				test.Type = DynamicSubtestType
			}

			tests = append(tests, test)
		}

		return true
	})

	// Assemble the path to each subtest
	lastTestIndex := 0

	for i, test := range tests {
		if test.Type == TestType {
			lastTestIndex = i
			continue
		}

		for j := lastTestIndex; j < i; j++ {
			parent := tests[j]

			if parent.RBrace > test.LBrace {
				test.Path = append(test.Path, parent.Name)
			}
		}
	}

	return tests, nil
}
