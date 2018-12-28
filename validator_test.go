package nit_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"

	"github.com/MarioCarrion/nit"
)

func TestImportsValidator_Validate(t *testing.T) {
	tests := [...]struct {
		name          string
		filename      string
		expectedError bool
	}{
		{
			"OK",
			"imports_valid.go",
			false,
		},
		{
			"Error: parenthesized declaration",
			"imports_paren.go",
			true,
		},
		{
			"Error: extra line",
			"imports_extra_line.go",
			true,
		},
		{
			"Error: missing line",
			"imports_missing_line.go",
			true,
		},
		{
			"Error: invalid group",
			"imports_invalid_group.go",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, filepath.Join("testdata", tt.filename), nil, parser.ParseComments)
			if err != nil {
				ts.Fatalf("expected no error, got %s", err)
			}

			for _, s := range f.Decls {
				switch g := s.(type) {
				case *ast.GenDecl:
					if g.Tok == token.IMPORT {
						validator := nit.ImportsValidator{LocalPath: "github.com/MarioCarrion"}
						if err := validator.Validate(g, fset); tt.expectedError != (err != nil) {
							ts.Fatalf("expected error %t, got %s", tt.expectedError, err)
						}
						break
					}
				}
			}
		})
	}
}
