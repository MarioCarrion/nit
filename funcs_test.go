package nit_test

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/MarioCarrion/nit"
)

//nolint:dupl
func TestFuncsValidator_Validate(t *testing.T) {
	tests := [...]struct {
		name          string
		filename      string
		expectedError bool
	}{
		{
			"OK",
			"funcs_valid.go",
			false,
		},
		{
			"OK: grouped",
			"funcs_group.go",
			false,
		},
		{
			"OK: sorted",
			"funcs_sorted_ok.go",
			false,
		},
		{
			"Error: sorted",
			"funcs_sorted.go",
			true,
		},
		{
			"Error: grouped",
			"funcs_group_error.go",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			var (
				err  error
				f    *ast.File
				fset *token.FileSet
			)

			f, fset, err = newParserFile(ts, tt.filename)
			if err != nil {
				ts.Fatalf("expected no error, got %s", err)
			}

			comments := nit.NewBreakComments(fset, f.Comments)
			validator := nit.NewFuncsValidator(comments)

			for _, s := range f.Decls {
				switch g := s.(type) {
				case *ast.FuncDecl:
					if g.Recv == nil {
						if err = validator.Validate(g, fset); err != nil {
							break
						}
					}
				}
			}
			if tt.expectedError != (err != nil) {
				ts.Fatalf("expected error %t, got %s", tt.expectedError, err)
			}
		})
	}
}
