package nit_test

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/MarioCarrion/nit"
)

//nolint:dupl
func TestMethodsValidator_Validate(t *testing.T) {
	tests := [...]struct {
		name          string
		filename      string
		expectedError bool
	}{
		{
			"OK",
			"methods_valid.go",
			false,
		},
		{
			"OK: sorted",
			"methods_sorted.go",
			false,
		},
		{
			"OK: sorted type",
			"methods_sorted_type_ok.go",
			false,
		},
		{
			"Error: not defined in file",
			"methods_not_defined.go",
			true,
		},
		{
			"Error: sorted",
			"methods_sorted_error.go",
			true,
		},
		{
			"Error: sorted type",
			"methods_sorted_type_error.go",
			true,
		},
		{
			"Error: sorted type comments",
			"methods_sorted_type_error1.go",
			true,
		},
	}

	t.Run("no types found", func(ts *testing.T) {
		_, err := nit.NewMethodsValidator(nil, nil)
		if err == nil {
			ts.Fatalf("expected error, got nil")
		}
	})

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			var (
				err        error
				f          *ast.File
				fset       *token.FileSet
				tvalidator *nit.TypesValidator
				validator  *nit.MethodsValidator
			)

			f, fset = newParserFile(ts, tt.filename)

			comments := nit.NewBreakComments(fset, f.Comments)

			for _, s := range f.Decls {
				switch g := s.(type) {
				case *ast.GenDecl:
					section, err1 := nit.NewGenDeclFileSection(g)
					if err1 != nil {
						ts.Fatalf("expected no error, got %s", err1)
					}
					if section == nit.FileSectionTypes {
						tvalidator = nit.NewTypesValidator(comments)
						if err1 := tvalidator.Validate(g, fset); err1 != nil {
							ts.Fatalf("expected no error, got %s", err1)
						}
					}
				case *ast.FuncDecl:
					if g.Recv == nil {
						continue
					}

					if validator == nil {
						validator, err = nit.NewMethodsValidator(comments, tvalidator)
						if err != nil {
							break
						}
					}

					if err = validator.Validate(g, fset); err != nil {
						break
					}
				}
			}
			if tt.expectedError != (err != nil) {
				ts.Fatalf("expected error %t, got %s", tt.expectedError, err)
			}
		})
	}
}
