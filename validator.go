package nit

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/pkg/errors"
)

type (
	// ImportsValidator defines the type including the rules used for validating
	// the `imports` section.
	ImportsValidator struct {
		LocalPath string
		fsm       *ImportsSectionMachine
	}
)

// Validate validates the token according to the imports rules.
func (i *ImportsValidator) Validate(v *ast.GenDecl, fset *token.FileSet) error {
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
