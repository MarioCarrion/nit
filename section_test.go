package nitpicking_test

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/MarioCarrion/nitpicking"
)

func TestNewGenDeclSection(t *testing.T) {
	tests := [...]struct {
		name          string
		input         *ast.GenDecl
		expected      nitpicking.Section
		expectedError bool
	}{
		{
			"Imports",
			&ast.GenDecl{Tok: token.IMPORT},
			nitpicking.SectionImports,
			false,
		},
		{
			"Consts",
			&ast.GenDecl{Tok: token.CONST},
			nitpicking.SectionConsts,
			false,
		},
		{
			"Type",
			&ast.GenDecl{Tok: token.TYPE},
			nitpicking.SectionTypes,
			false,
		},
		{
			"Vars",
			&ast.GenDecl{Tok: token.VAR},
			nitpicking.SectionVars,
			false,
		},
		{
			"Error",
			&ast.GenDecl{Tok: token.COMMENT},
			nitpicking.SectionStart,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			actual, err := nitpicking.NewGenDeclSection(tt.input)
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
		expected nitpicking.Section
	}{
		{
			"Funcs",
			&ast.FuncDecl{},
			nitpicking.SectionFuncs,
		},
		{
			"Methods",
			&ast.FuncDecl{Recv: &ast.FieldList{}},
			nitpicking.SectionMethods,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			actual, err := nitpicking.NewFuncDeclSection(tt.input)
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
		expected            nitpicking.ImportsSection
	}{
		{
			"ImportsSectionStd",
			"fmt",
			"github.com/MarioCarrion/nitpicking",
			nitpicking.ImportsSectionStd,
		},
		{
			"ImportsSectionExternal",
			"github.com/golang/go",
			"github.com/MarioCarrion/nitpicking",
			nitpicking.ImportsSectionExternal,
		},
		{
			"ImportsSectionLocal",
			"github.com/MarioCarrion/nitpicking/something",
			"github.com/MarioCarrion/nitpicking",
			nitpicking.ImportsSectionLocal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			if actual := nitpicking.NewImportsSection(tt.inputPath, tt.inputLocaBasePrefix); actual != tt.expected {
				t.Fatalf("expected %+v, actual %+v", tt.expected, actual)
			}
		})
	}
}
