package nit

import (
	"go/ast"
	"go/token"

	"github.com/pkg/errors"
)

type (
	// VarsValidator defines the type including the rules used for validating
	// the `var` sections.
	VarsValidator struct {
		sortedNamesValidator
	}
)

// Validate makes sure the implemented `var` declaration satisfies the
// following rules:
// * Group declaration is parenthesized
// * Sorted exported vars are declared first, and
// * Sorted unexported vars are declared next
func (c *VarsValidator) Validate(v *ast.GenDecl, fset *token.FileSet) error { //nolint: gocyclo
	c.identType = "Var"

	if !v.Lparen.IsValid() {
		return errors.Wrap(errors.New("expected parenthesized declaration"), fset.PositionFor(v.Pos(), false).String())
	}

	for _, t := range v.Specs {
		errPrefix := fset.PositionFor(t.Pos(), false).String()

		s, ok := t.(*ast.ValueSpec)
		if !ok {
			return errors.Wrap(errors.Errorf("invalid token %+v", t), errPrefix)
		}

		for _, name := range s.Names {
			if err := c.validateName(errPrefix, name); err != nil {
				return err
			}
		}
	}

	return nil
}
