package nit_test

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/MarioCarrion/nit"
)

//nolint:dupl
func TestVarsValidator_Validate(t *testing.T) {
	tests := [...]struct {
		name          string
		filename      string
		expectedError bool
	}{
		{
			"OK",
			"vars_valid.go",
			false,
		},
		{
			"Error: parenthesized declaration",
			"vars_paren.go",
			true,
		},
		{
			"Error: grouped 1",
			"vars_group1.go",
			true,
		},
		{
			"Error: grouped 2",
			"vars_group2.go",
			true,
		},
		{
			"Error: sorted",
			"vars_sorted.go",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			f, fset := newParserFile(ts, tt.filename)

			for _, s := range f.Decls {
				switch g := s.(type) {
				case *ast.GenDecl:
					if g.Tok == token.VAR {
						validator := nit.VarsValidator{}
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
