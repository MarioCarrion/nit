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
		LocalPath         string
		SkipGeneratedFile bool
		//-
		fset       *token.FileSet
		fsm        *FileSectionMachine
		comments   *BreakComments
		tvalidator *TypesValidator
		fvalidator *FuncsValidator
		mvalidator *MethodsValidator
	}
)

// Validate nitpicks the filename.
func (v *Nitpicker) Validate(filename string) error {
	v.fset = token.NewFileSet()
	f, err := parser.ParseFile(v.fset, filename, nil, parser.ParseComments)
	if err != nil {
		return errors.Wrap(err, "parsing file failed")
	}

	v.comments = NewBreakComments(v.fset, f.Comments)
	if v.comments.HasGeneratedCode() && v.SkipGeneratedFile {
		return nil
	}

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
		if v.tvalidator != nil {
			return errors.New("only one `type` section block is allowed per file")
		}
		v.tvalidator = &TypesValidator{}
		if err := v.tvalidator.Validate(genDecl, v.fset); err != nil {
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
		if v.fvalidator == nil {
			v.fvalidator = NewFuncsValidator(v.comments)
		}
		if err := v.fvalidator.Validate(funcDecl, v.fset); err != nil {
			return err
		}
	case FileSectionMethods:
		if v.mvalidator == nil {
			mvalidator, err := NewMethodsValidator(v.comments, v.tvalidator)
			if err != nil {
				return err
			}
			v.mvalidator = mvalidator
		}
		if err := v.mvalidator.Validate(funcDecl, v.fset); err != nil {
			return err
		}
	}

	return nil
}
