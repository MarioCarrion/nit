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

	// FuncsValidator defines the type including the rules used for validating
	// functions.
	FuncsValidator struct {
		sortedNamesValidator
		Comments *BreakComments
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

// Validate makes sure the implemented `const` declaration satisfies the
// following rules:
// * Group declaration is parenthesized
// * Declarations are sorted
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
			if err := c.validateName(errPrefix, name); err != nil {
				return err
			}
		}
	}

	return nil
}

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

// Validate makes sure the implemented `var` declaration satisfies the
// following rules:
// * Group declaration is parenthesized
// * Sorted exported vars are declared first, and
// * Sorted unexported vars are declared next
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
			if err := c.validateName(errPrefix, name); err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *sortedNamesValidator) validateExported(errPrefix string, name *ast.Ident) error {
	if v.exported == nil || (*v.exported && !name.IsExported()) {
		e := name.IsExported()
		v.exported = &e
	}

	if *v.exported != name.IsExported() {
		return errors.Wrap(errors.Errorf("%s is not grouped correctly", name.Name), errPrefix)
	}

	return nil
}

func (v *sortedNamesValidator) validateName(errPrefix string, name *ast.Ident) error {
	if err := v.validateExported(errPrefix, name); err != nil {
		return err
	}

	if err := v.validateSortedName(errPrefix, name); err != nil {
		return err
	}

	return nil
}

func (v *sortedNamesValidator) validateSortedName(errPrefix string, name *ast.Ident) error {
	if v.last != "" && v.last > name.Name {
		return errors.Wrap(errors.Errorf("%s is not sorted", name.Name), errPrefix)
	}

	v.last = name.Name
	return nil
}
