package nitpicking

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/pkg/errors"
)

type (
	// Section represents a code section.
	Section uint8

	// ImportsSection represents an `imports` section
	ImportsSection uint8

	// SectionValidator FIXME
	SectionValidator interface {
		//		IsValid(d ast.Decl) bool
		// Validate(*ast.GenDecl) error
		Validate(v *ast.GenDecl, fset *token.FileSet) error
	}

	// Imports FIXME
	Imports struct {
		LocalPath string
		fsm       *ImportsSectionMachine
	}
)

const (
	// SectionStart defines the initial State in the validator.
	SectionStart Section = iota
	// SectionImports defines the `import` state.
	SectionImports
	// SectionConsts defines the `const` state.
	SectionConsts
	// SectionTypes defines the `type` state.
	SectionTypes
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

// NewGenDeclState returns a new State that matches the decl type.
func NewGenDeclState(decl *ast.GenDecl) (Section, error) {
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

// NewFuncDeclState returns a new State that matches the decl type.
func NewFuncDeclState(decl *ast.FuncDecl) (Section, error) {
	if decl.Recv == nil {
		return SectionFuncs, nil
	}
	return SectionMethods, nil
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

// Validate validates the token according to the imports rules.
func (i *Imports) Validate(v *ast.GenDecl, fset *token.FileSet) error {
	if !v.Lparen.IsValid() {
		return errors.Wrap(fmt.Errorf("expected parenthesized declaration"), fset.PositionFor(v.Pos(), false).String())
	}

	lastLine := fset.PositionFor(v.Pos(), false).Line

	for _, t := range v.Specs {
		errPrefix := fset.PositionFor(t.Pos(), false).String()

		s, ok := t.(*ast.ImportSpec)
		if !ok {
			return errors.Wrap(fmt.Errorf("invalid token %+v", t), errPrefix)
		}

		section := NewImportsSection(s.Path.Value, i.LocalPath)
		if i.fsm == nil {
			i.fsm = NewImportsSectionMachine(section)
		}
		if err := i.fsm.Transition(section); err != nil {
			return errors.Wrap(err, errPrefix)
		}

		newLine := fset.PositionFor(t.Pos(), false).Line

		if i.fsm.Current() == i.fsm.Previous() {
			if lastLine+1 != newLine {
				return errors.Wrap(errors.New("extra line break in section"), errPrefix)
			}
		} else {
			if lastLine+1 == newLine {
				return errors.Wrap(errors.New("missing line break in section"), errPrefix)
			}
		}

		lastLine = newLine
	}

	return nil

}
