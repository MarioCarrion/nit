package nit

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/pkg/errors"
)

type (
	// Nitpicker defines the linter.
	Nitpicker struct {
		LocalPath string
		fset      *token.FileSet
		fsm       SectionMachine
	}
)

// Validate nitpicks the filename.
func (v *Nitpicker) Validate(filename string) error {
	v.fset = token.NewFileSet()
	f, err := parser.ParseFile(v.fset, filename, nil, parser.ParseComments)
	if err != nil {
		return errors.Wrap(err, "parsing file failed")
	}

	for _, s := range f.Decls {
		// fmt.Printf("%d == %T - %+v -- %t\n", v.fset.PositionFor(s.Pos(), false).Line, s, s, s.End().IsValid())
		if err := v.validateToken(s); err != nil {
			return err
		}
	}
	return nil
}

func (v *Nitpicker) validateToken(d ast.Decl) error {
	var (
		err       error
		genDecl   *ast.GenDecl
		funcDecl  *ast.FuncDecl
		nextState Section
	)

	switch t := d.(type) {
	case *ast.GenDecl:
		genDecl = t
		nextState, err = NewGenDeclSection(genDecl)
	case *ast.FuncDecl:
		funcDecl = t
		nextState, err = NewFuncDeclSection(funcDecl)
	default:
		return fmt.Errorf("unknown declaration state")
	}
	if err != nil {
		return err
	}

	if err = v.fsm.Transition(nextState); err != nil {
		return errors.Wrap(err, v.fset.PositionFor(d.Pos(), false).String())
	}

	if nextState == SectionImports {
		validator := &ImportsValidator{LocalPath: v.LocalPath}
		if err := validator.Validate(genDecl, v.fset); err != nil {
			return err
		}
	}

	return nil
}
