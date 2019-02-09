package nit

import (
	"go/ast"
	"go/token"

	"github.com/pkg/errors"
)

type (
	// ConstsValidator defines the type including the rules used for validating
	// the `const` sections.
	ConstsValidator struct {
		sortedNamesValidator
	}
)

// Validate makes sure the implemented `const` declaration satisfies the
// following rules:
// * Group declaration is parenthesized
// * Declarations are sorted
func (c *ConstsValidator) Validate(v *ast.GenDecl, fset *token.FileSet) error { //nolint: gocyclo
	c.identType = "Const"

	if !v.Lparen.IsValid() {
		return errors.Wrap(errors.New("expected parenthesized declaration"), fset.PositionFor(v.Pos(), false).String())
	}

	for _, t := range v.Specs {
		errPrefix := fset.PositionFor(t.Pos(), false).String()

		s, ok := t.(*ast.ValueSpec)
		if !ok {
			return errors.Wrap(errors.Errorf("invalid token %+v", t), errPrefix)
		}

		for _, vss := range s.Values {
			_, ok := vss.(*ast.Ident)
			if ok {
				return nil // iota
			}
		}

		for _, name := range s.Names {
			if err := c.validateName(errPrefix, name); err != nil {
				return err
			}
		}
	}

	return nil
}
