package nit

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/pkg/errors"
)

type (
	// Nitpicker defines the linter.
	Nitpicker struct {
		LocalPath  string
		fset       *token.FileSet
		fsm        *FileSectionMachine
		fvalidator *FuncsValidator
	}
)

// Validate nitpicks the filename.
func (v *Nitpicker) Validate(filename string) error {
	v.fset = token.NewFileSet()
	f, err := parser.ParseFile(v.fset, filename, nil, parser.ParseComments)
	if err != nil {
		return errors.Wrap(err, "parsing file failed")
	}

	comments := NewBreakComments(v.fset, f.Comments)
	v.fvalidator = &FuncsValidator{Comments: &comments}

	for _, s := range f.Decls {
		// fmt.Printf("%d == %T - %+v -- %t\n", v.fset.PositionFor(s.Pos(), false).Line, s, s, s.End().IsValid())
		if err := v.validateToken(s); err != nil {
			return err
		}
	}

	return nil
}

//nolint:gocyclo
func (v *Nitpicker) validateToken(d ast.Decl) error {
	var (
		err       error
		genDecl   *ast.GenDecl
		funcDecl  *ast.FuncDecl
		nextState FileSection
	)

	switch t := d.(type) {
	case *ast.GenDecl:
		genDecl = t
		nextState, err = NewGenDeclFileSection(genDecl)
	case *ast.FuncDecl:
		funcDecl = t
		nextState, err = NewFuncDeclFileSection(funcDecl)
	default:
		return errors.New("unknown declaration state")
	}
	if err != nil {
		return err
	}

	if v.fsm == nil {
		fsm, err := NewFileSectionMachine(nextState)
		if err != nil {
			return errors.Wrap(err, v.fset.PositionFor(d.Pos(), false).String())
		}
		v.fsm = fsm
	}

	if err = v.fsm.Transition(nextState); err != nil {
		return errors.Wrap(err, v.fset.PositionFor(d.Pos(), false).String())
	}

	switch nextState {
	case FileSectionImports:
		validator := NewImportsValidator(v.LocalPath)
		if err := validator.Validate(genDecl, v.fset); err != nil {
			return err
		}
	case FileSectionTypes:
		validator := &TypesValidator{}
		if err := validator.Validate(genDecl, v.fset); err != nil {
			return err
		}
	case FileSectionConsts:
		validator := &ConstsValidator{}
		if err := validator.Validate(genDecl, v.fset); err != nil {
			return err
		}
	case FileSectionVars:
		validator := &VarsValidator{}
		if err := validator.Validate(genDecl, v.fset); err != nil {
			return err
		}
	case FileSectionFuncs:
		if err := v.fvalidator.Validate(funcDecl, v.fset); err != nil {
			return err
		}
	}

	return nil
}
