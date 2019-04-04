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
		comments *BreakComments
		types    []string
	}
)

// NewTypesValidator returns a correctly initialized TypesValidator.
func NewTypesValidator(c *BreakComments) *TypesValidator {
	return &TypesValidator{comments: c}
}

// Types returns the names of all found types.
func (tv *TypesValidator) Types() []string {
	dst := make([]string, len(tv.types))
	copy(dst, tv.types)
	return dst
}

// Validate makes sure the implemented `type` declaration satisfies the
// following rules:
// * Group declaration is parenthesized
// * Sorted exported types are declared first, and
// * Sorted unexported types are declared next
func (tv *TypesValidator) Validate(v *ast.GenDecl, fset *token.FileSet) error { //nolint: gocyclo
	tv.identType = "Type"

	if !v.Lparen.IsValid() {
		return errors.Wrap(errors.New("expected parenthesized declaration"), fset.PositionFor(v.Pos(), false).String())
	}

	for _, t := range v.Specs {
		errPrefix := fset.PositionFor(t.Pos(), false).String()

		s, ok := t.(*ast.TypeSpec)
		if !ok {
			return errors.Wrap(errors.Errorf("invalid token %+v", t), errPrefix)
		}

		if err := tv.validateExported(errPrefix, s.Name); err != nil {
			return err
		}

		next := tv.comments.Next()
		if next != -1 && fset.PositionFor(s.Pos(), false).Line > next {
			tv.last = ""
		}

		if err := tv.validateSortedName(errPrefix, s.Name); err != nil {
			return err
		}

		tv.comments.MoveTo(fset.PositionFor(s.End(), false).Line)

		tv.types = append(tv.types, s.Name.Name)
	}

	return nil
}
