package nit_test

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/MarioCarrion/nit"
)

//nolint:dupl
func TestConstsValidator_Validate(t *testing.T) {
	tests := [...]struct {
		name          string
		filename      string
		expectedError bool
	}{
		{
			"OK",
			"consts_valid.go",
			false,
		},
		{
			"OK: iota",
			"consts_iota.go",
			false,
		},
		{
			"Error: parenthesized declaration",
			"consts_paren.go",
			true,
		},
		{
			"Error: grouped 1",
			"consts_group1.go",
			true,
		},
		{
			"Error: grouped 2",
			"consts_group2.go",
			true,
		},
		{
			"Error: sorted",
			"consts_sorted.go",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			f, fset := newParserFile(ts, tt.filename)

			for _, s := range f.Decls {
				switch g := s.(type) {
				case *ast.GenDecl:
					if g.Tok == token.CONST {
						validator := nit.ConstsValidator{}
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
