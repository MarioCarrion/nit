package nit_test

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/MarioCarrion/nit"
)

func TestNewImportsSection(t *testing.T) {
	tests := [...]struct {
		name                string
		inputPath           string
		inputLocaBasePrefix string
		expected            nit.ImportsSection
	}{
		{
			"ImportsSectionStd",
			"fmt",
			"github.com/MarioCarrion/nit",
			nit.ImportsSectionStd,
		},
		{
			"ImportsSectionExternal",
			"github.com/golang/go",
			"github.com/MarioCarrion/nit",
			nit.ImportsSectionExternal,
		},
		{
			"ImportsSectionLocal",
			"github.com/MarioCarrion/nit/something",
			"github.com/MarioCarrion/nit",
			nit.ImportsSectionLocal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			if actual := nit.NewImportsSection(tt.inputPath, tt.inputLocaBasePrefix); actual != tt.expected {
				ts.Fatalf("expected %+v, actual %+v", tt.expected, actual)
			}
		})
	}
}

//nolint:dupl
func TestImportsValidator(t *testing.T) {
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
			f, fset := newParserFile(ts, tt.filename)

			for _, s := range f.Decls {
				switch g := s.(type) {
				case *ast.GenDecl:
					if g.Tok == token.IMPORT {
						validator := nit.NewImportsValidator("github.com/MarioCarrion")
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

func TestNewImportsTransition(t *testing.T) {
	tests := [...]struct {
		name          string
		input         nit.ImportsSection
		expectedError bool
	}{
		{
			"ImportsSectionStd",
			nit.ImportsSectionStd,
			false,
		},
		{
			"ImportsSectionExternal",
			nit.ImportsSectionExternal,
			false,
		},
		{
			"ImportsSectionLocal",
			nit.ImportsSectionLocal,
			false,
		},
		{
			"ImportsSectionLocal",
			nit.ImportsSection(3),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			_, err := nit.NewImportsTransition(tt.input)
			if (err != nil) != tt.expectedError {
				ts.Fatalf("expected error %t, got %s", tt.expectedError, err)
			}
		})
	}
}

func TestImportsTransitionExternal(t *testing.T) {
	i, _ := nit.NewImportsTransition(nit.ImportsSectionExternal)

	tests := [...]struct {
		name          string
		transition    func() (nit.ImportsTransition, error)
		expectedError bool
	}{
		{
			"External",
			i.External,
			false,
		},
		{
			"Local",
			i.Local,
			false,
		},
		{
			"Standard",
			i.Standard,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			_, err := tt.transition()
			if (err != nil) != tt.expectedError {
				ts.Fatalf("expected error %t, got %s", tt.expectedError, err)
			}
		})
	}
}

func TestImportsTransitionLocal(t *testing.T) {
	i, _ := nit.NewImportsTransition(nit.ImportsSectionLocal)

	tests := [...]struct {
		name          string
		transition    func() (nit.ImportsTransition, error)
		expectedError bool
	}{
		{
			"External",
			i.External,
			true,
		},
		{
			"Local",
			i.Local,
			false,
		},
		{
			"Standard",
			i.Standard,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			_, err := tt.transition()
			if (err != nil) != tt.expectedError {
				ts.Fatalf("expected error %t, got %s", tt.expectedError, err)
			}
		})
	}
}

func TestImportsTransitionStd(t *testing.T) {
	i, _ := nit.NewImportsTransition(nit.ImportsSectionStd)

	tests := [...]struct {
		name          string
		transition    func() (nit.ImportsTransition, error)
		expectedError bool
	}{
		{
			"External",
			i.External,
			false,
		},
		{
			"Local",
			i.Local,
			false,
		},
		{
			"Standard",
			i.Standard,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			_, err := tt.transition()
			if (err != nil) != tt.expectedError {
				ts.Fatalf("expected error %t, got %s", tt.expectedError, err)
			}
		})
	}
}
