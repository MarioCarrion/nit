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
		exported *bool
		last     string
	}

	// ImportsValidator defines the type including the rules used for validating
	// the `imports` section.
	ImportsValidator struct {
		LocalPath string
		fsm       *ImportsSectionMachine
	}
)

// Validate validates the token according to the `const` rules.
func (c *ConstsValidator) Validate(v *ast.GenDecl, fset *token.FileSet) error { //nolint: gocyclo
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

		for _, n := range s.Names {
			if c.exported == nil || (*c.exported && !n.IsExported()) {
				e := n.IsExported()
				c.exported = &e
			}

			if *c.exported != n.IsExported() {
				return errors.Wrap(errors.Errorf("%s is not grouped correctly", n.Name), errPrefix)
			}

			if c.last != "" && c.last > n.Name {
				return errors.Wrap(errors.Errorf("%s is not sorted", n.Name), errPrefix)
			}

			c.last = n.Name
		}
	}

	return nil
}

// Validate validates the token according to the `imports` rules.
func (i *ImportsValidator) Validate(v *ast.GenDecl, fset *token.FileSet) error {
	if !v.Lparen.IsValid() {
		return errors.Wrap(errors.New("expected parenthesized declaration"), fset.PositionFor(v.Pos(), false).String())
	}

	lastLine := fset.PositionFor(v.Pos(), false).Line

	for _, t := range v.Specs {
		errPrefix := fset.PositionFor(t.Pos(), false).String()

		s, ok := t.(*ast.ImportSpec)
		if !ok {
			return errors.Wrap(errors.Errorf("invalid token %+v", t), errPrefix)
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
