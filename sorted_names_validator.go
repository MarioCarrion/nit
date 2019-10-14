package nit

import (
	"go/ast"

	"github.com/pkg/errors"
)

type (
	sortedNamesValidator struct {
		identType string
		exported  *bool
		last      string
	}
)

func (v *sortedNamesValidator) validateName(errPrefix string, name *ast.Ident) error {
	if err := v.validateExported(errPrefix, name); err != nil {
		return err
	}

	if err := v.validateSortedName(errPrefix, name); err != nil {
		return err
	}

	return nil
}

//-

func (v *sortedNamesValidator) validateExported(errPrefix string, name *ast.Ident) error {
	if v.exported == nil || (*v.exported && !name.IsExported()) {
		e := name.IsExported()
		v.exported = &e
	}

	if *v.exported != name.IsExported() {
		return errors.Wrap(errors.Errorf("%s `%s` is not grouped correctly", v.identType, name.Name), errPrefix)
	}

	return nil
}

func (v *sortedNamesValidator) validateSortedName(errPrefix string, name *ast.Ident) error {
	if v.last != "" && v.last > name.Name {
		return errors.Wrap(errors.Errorf("%s `%s` is not sorted", v.identType, name.Name), errPrefix)
	}

	v.last = name.Name

	return nil
}
