package nit

import (
	"go/ast"
	"go/token"
)

type (
	// FuncsValidator defines the type including the rules used for validating
	// functions.
	FuncsValidator struct {
		sortedNamesValidator
		Comments *BreakComments
	}
)

// Validate makes sure the implemented function satisfies the following rules
// considering all previous declared functions:
// * Sorted exported functions are declared first,
// * Sorted unexported functions are declared next, and
// * Both groups can declare their own sorted subgroups,
func (f *FuncsValidator) Validate(v *ast.FuncDecl, fset *token.FileSet) error {
	errPrefix := fset.PositionFor(v.Pos(), false).String()

	if err := f.validateExported(errPrefix, v.Name); err != nil {
		return err
	}

	if f.Comments.Next() > fset.PositionFor(v.Pos(), false).Line {
		f.last = ""
	}

	if err := f.validateSortedName(errPrefix, v.Name); err != nil {
		return err
	}

	f.Comments.MoveTo(fset.PositionFor(v.End(), false).Line)

	return nil
}
