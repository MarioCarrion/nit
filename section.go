package nit

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type (
	// ImportsSection represents an `imports` section
	ImportsSection uint8

	// Section represents a code section.
	Section uint8
)

const (
	// SectionStart defines the initial State in the validator.
	SectionStart Section = iota
	// SectionImports defines the `import` state.
	SectionImports
	// SectionTypes defines the `type` state.
	SectionTypes
	// SectionConsts defines the `const` state.
	SectionConsts
	// SectionVars defines the `var` state.
	SectionVars
	// SectionFuncs defines the `func` state.
	SectionFuncs
	// SectionMethods defines the `method` state.
	SectionMethods
)

const (
	// ImportsSectionStd represents the Standard Library Packages `imports` section.
	ImportsSectionStd ImportsSection = iota
	// ImportsSectionExternal represents the External Packages `imports` section.
	ImportsSectionExternal
	// ImportsSectionLocal represents the local Packages `imports` section.
	ImportsSectionLocal
)

// NewFuncDeclSection returns a new State that matches the decl type.
func NewFuncDeclSection(decl *ast.FuncDecl) (Section, error) {
	if decl.Recv == nil {
		return SectionFuncs, nil
	}
	return SectionMethods, nil
}

// NewGenDeclSection returns a new State that matches the decl type.
func NewGenDeclSection(decl *ast.GenDecl) (Section, error) {
	switch decl.Tok {
	case token.IMPORT:
		return SectionImports, nil
	case token.CONST:
		return SectionConsts, nil
	case token.TYPE:
		return SectionTypes, nil
	case token.VAR:
		return SectionVars, nil
	}

	return SectionStart, fmt.Errorf("unknown generic declaration node")
}

// NewImportsSection returns a new ImportsSection from the path value.
func NewImportsSection(path, localPathPrefix string) ImportsSection {
	if !strings.Contains(path, ".") {
		return ImportsSectionStd
	}
	if strings.HasPrefix(strings.Replace(path, "\"", "", -1), localPathPrefix) {
		return ImportsSectionLocal
	}
	return ImportsSectionExternal
}
