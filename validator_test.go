package nit_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
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
			f, fset, err := newParserFile(ts, tt.filename)
			if err != nil {
				ts.Fatalf("expected no error, got %s", err)
			}

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
			"OK: groupd",
			"funcs_group.go",
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
			validator := nit.FuncsValidator{Comments: &comments}

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

//nolint:dupl
func TestTypesValidator_Validate(t *testing.T) {
	tests := [...]struct {
		name          string
		filename      string
		expectedError bool
	}{
		{
			"OK",
			"types_valid.go",
			false,
		},
		{
			"Error: parenthesized declaration",
			"types_paren.go",
			true,
		},
		{
			"Error: group 1",
			"types_group1.go",
			true,
		},
		{
			"Error: group 2",
			"types_group2.go",
			true,
		},
		{
			"Error: sorted",
			"types_sorted.go",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			f, fset, err := newParserFile(ts, tt.filename)
			if err != nil {
				ts.Fatalf("expected no error, got %s", err)
			}

			for _, s := range f.Decls {
				switch g := s.(type) {
				case *ast.GenDecl:
					if g.Tok == token.TYPE {
						validator := nit.TypesValidator{}
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

//nolint:dupl
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
			f, fset, err := newParserFile(ts, tt.filename)
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
			f, fset, err := newParserFile(ts, tt.filename)
			if err != nil {
				ts.Fatalf("expected no error, got %s", err)
			}

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

//-

func newParserFile(t *testing.T, name string) (*ast.File, *token.FileSet, error) {
	t.Helper()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filepath.Join("testdata", name), nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	return f, fset, nil
}
