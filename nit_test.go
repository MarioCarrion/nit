package nit_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"
)

func newParserFile(t *testing.T, name string) (*ast.File, *token.FileSet) {
	t.Helper()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filepath.Join("testdata", name), nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	return f, fset
}
