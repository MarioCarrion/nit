package nit_test

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/MarioCarrion/nit"
)

func TestNewGenDeclSection(t *testing.T) {
	tests := [...]struct {
		name          string
		input         *ast.GenDecl
		expected      nit.Section
		expectedError bool
	}{
		{
			"Imports",
			&ast.GenDecl{Tok: token.IMPORT},
			nit.SectionImports,
			false,
		},
		{
			"Consts",
			&ast.GenDecl{Tok: token.CONST},
			nit.SectionConsts,
			false,
		},
		{
			"Type",
			&ast.GenDecl{Tok: token.TYPE},
			nit.SectionTypes,
			false,
		},
		{
			"Vars",
			&ast.GenDecl{Tok: token.VAR},
			nit.SectionVars,
			false,
		},
		{
			"Error",
			&ast.GenDecl{Tok: token.COMMENT},
			nit.SectionStart,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			actual, err := nit.NewGenDeclSection(tt.input)
			if (err != nil) != tt.expectedError {
				ts.Fatalf("expected no error, got %s", err)
			}
			if actual != tt.expected {
				ts.Fatalf("expected %+v, got %+v", tt.expected, actual)
			}
		})
	}
}

func TestNewFuncDeclSection(t *testing.T) {
	tests := [...]struct {
		name     string
		input    *ast.FuncDecl
		expected nit.Section
	}{
		{
			"Funcs",
			&ast.FuncDecl{},
			nit.SectionFuncs,
		},
		{
			"Methods",
			&ast.FuncDecl{Recv: &ast.FieldList{}},
			nit.SectionMethods,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			actual, err := nit.NewFuncDeclSection(tt.input)
			if err != nil {
				ts.Fatalf("expected no error, got %s", err)
			}
			if actual != tt.expected {
				ts.Fatalf("expected %+v, got %+v", tt.expected, actual)
			}
		})
	}
}

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
				t.Fatalf("expected %+v, actual %+v", tt.expected, actual)
			}
		})
	}
}
