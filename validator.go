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

	// ImportsValidator defines the type including the rules used for validating
	// the `imports` section.
	ImportsValidator struct {
		LocalPath string
		fsm       *ImportsSectionMachine
	}

	// TypesValidator defines the type including the rules used for validating
	// the `type` sections.
	TypesValidator struct {
		sortedNamesValidator
	}

	// VarsValidator defines the type including the rules used for validating
	// the `var` sections.
	VarsValidator struct {
		sortedNamesValidator
	}

	sortedNamesValidator struct {
		exported *bool
		last     string
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

		for _, name := range s.Names {
			if err := c.validateSortedName(errPrefix, name); err != nil {
				return err
			}
		}
	}

	return nil
}

// Validate valites the token according to the `func` rules.
// func (f *FuncsValidator) Validate(v *ast.FuncDecl, fst *token.FileSet) error {
// 	return nil
// }

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

// Validate validates the token according to the `type` rules.
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

		if err := tv.validateSortedName(errPrefix, s.Name); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates the token according to the `var` rules.
func (c *VarsValidator) Validate(v *ast.GenDecl, fset *token.FileSet) error { //nolint: gocyclo
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
			if err := c.validateSortedName(errPrefix, name); err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *sortedNamesValidator) validateSortedName(errPrefix string, name *ast.Ident) error {
	if v.exported == nil || (*v.exported && !name.IsExported()) {
		e := name.IsExported()
		v.exported = &e
	}

	if *v.exported != name.IsExported() {
		return errors.Wrap(errors.Errorf("%s is not grouped correctly", name.Name), errPrefix)
	}

	if v.last != "" && v.last > name.Name {
		return errors.Wrap(errors.Errorf("%s is not sorted", name.Name), errPrefix)
	}

	v.last = name.Name
	return nil
}
