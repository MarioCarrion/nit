package nit

import (
	"go/ast"
	"go/token"

	"github.com/pkg/errors"
)

type (
	// TypesValidator defines the type including the rules used for validating
	// the `type` sections.
	TypesValidator struct {
		sortedNamesValidator
	}
)

// Validate makes sure the implemented `type` declaration satisfies the
// following rules:
// * Group declaration is parenthesized
// * Sorted exported types are declared first, and
// * Sorted unexported types are declared next
func (tv *TypesValidator) Validate(v *ast.GenDecl, fset *token.FileSet) error { //nolint: gocyclo
	if !v.Lparen.IsValid() {
		return errors.Wrap(errors.New("expected parenthesized declaration"), fset.PositionFor(v.Pos(), false).String())
	}

	for _, t := range v.Specs {
		errPrefix := fset.PositionFor(t.Pos(), false).String()

		s, ok := t.(*ast.TypeSpec)
		if !ok {
			return errors.Wrap(errors.Errorf("invalid token %+v", t), errPrefix)
		}

		if err := tv.validateName(errPrefix, s.Name); err != nil {
			return err
		}
	}

	return nil
}
