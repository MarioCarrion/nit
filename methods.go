package nit

import (
	"go/ast"
	"go/token"
	"reflect"

	"github.com/pkg/errors"
)

type (
	// MethodsValidator defines the type including the rules used for validating
	// methods.
	MethodsValidator struct {
		comments      *BreakComments
		sortedTypes   sortedNamesValidator
		sortedMethods sortedNamesValidator
		types         map[string]struct{}
		lastType      string
	}

	// TypesFound defines the type returning the types found in the file.
	TypesFound interface {
		Types() []string
	}
)

// NewMethodsValidator returns a correctly initialized MethodsValidator.
func NewMethodsValidator(c *BreakComments, t TypesFound) (*MethodsValidator, error) {
	if t == nil || reflect.ValueOf(t).IsNil() {
		return nil, errors.New("no types found")
	}

	ts := make(map[string]struct{})
	for _, tf := range t.Types() {
		ts[tf] = struct{}{}
	}

	return &MethodsValidator{comments: c, types: ts}, nil
}

// Validate makes sure the implemented methods satisfies the following rules
// considering all previous declared methods:
// * Methods for exported types are declared first, then unexported ones,
// * Sorted exported methods are declared first,
// * Sorted unexported methods are declared next, and
// * Both groups can declare their own sorted subgroups.
func (m *MethodsValidator) Validate(v *ast.FuncDecl, fset *token.FileSet) error {
	var rcvType *ast.Ident
	switch e := v.Recv.List[0].Type.(type) {
	case *ast.Ident:
		rcvType = e
	case *ast.StarExpr:
		rcvType = e.X.(*ast.Ident)
	}

	errPrefix := fset.PositionFor(v.Pos(), false).String()
	if _, ok := m.types[rcvType.Name]; !ok {
		return errors.Wrap(errors.Errorf("Type `%s` is not defined in the file", rcvType.Name), errPrefix)
	}

	validateSorted := func(v *sortedNamesValidator, i *ast.Ident, honorComments bool) error {
		if err := v.validateExported(errPrefix, i); err != nil {
			return err
		}

		if honorComments {
			next := m.comments.Next()

			if next != -1 && fset.PositionFor(i.Pos(), false).Line > next {
				v.last = ""
			}
		}

		if err := v.validateSortedName(errPrefix, i); err != nil {
			return err
		}

		return nil
	}

	if m.lastType != rcvType.Name {
		m.sortedTypes.identType = "Type"
		m.sortedMethods = sortedNamesValidator{}
		if err := validateSorted(&m.sortedTypes, rcvType, false); err != nil {
			return err
		}
		m.lastType = rcvType.Name
	}

	m.sortedMethods.identType = "Method"
	if err := validateSorted(&m.sortedMethods, v.Name, true); err != nil {
		return err
	}

	m.comments.MoveTo(fset.PositionFor(v.End(), false).Line)

	return nil
}
