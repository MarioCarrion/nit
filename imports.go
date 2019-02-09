package nit

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/pkg/errors"
)

type (
	// ImportsSection represents one of the 3 valid `imports` sections.
	ImportsSection uint8

	// ImportsSectionMachine represents the `imports` code organization state
	// machine.
	ImportsSectionMachine struct {
		current  ImportsTransition
		previous ImportsTransition
	}

	// ImportsTransition represents one of the 3 valid `imports` sections, it
	// defines the State Machine transition rules for that concrete section,
	// the order to be expected is: "Standard" -> "External" -> "Local".
	ImportsTransition interface {
		External() (ImportsTransition, error)
		Local() (ImportsTransition, error)
		Standard() (ImportsTransition, error)
	}

	// ImportsValidator defines the type including the rules used for validating
	// the `imports` section as a whole.
	ImportsValidator struct {
		localPath string
		fsm       *ImportsSectionMachine
	}

	externalImportsTransition struct{}
	localImportsTransition    struct{}
	standardImportsTransition struct{}
)

const (
	// ImportsSectionStd represents the Standard Library Packages `imports` section.
	ImportsSectionStd ImportsSection = iota

	// ImportsSectionExternal represents the External Packages `imports` section.
	ImportsSectionExternal

	// ImportsSectionLocal represents the local Packages `imports` section.
	ImportsSectionLocal
)

// NewImportsSection returns the value representing the corresponding Imports
// section.
func NewImportsSection(path, localPathPrefix string) ImportsSection {
	if !strings.Contains(path, ".") {
		return ImportsSectionStd
	}
	if strings.HasPrefix(strings.Replace(path, "\"", "", -1), localPathPrefix) {
		return ImportsSectionLocal
	}
	return ImportsSectionExternal
}

// NewImportsValidator returns a new instalce of ImporstValidator with the
// local path prefix set.
func NewImportsValidator(localPath string) ImportsValidator {
	return ImportsValidator{localPath: localPath}
}

// NewImportsSectionMachine returns a new ImportsSectionMachine with the
// initial state as `start`.
func NewImportsSectionMachine(start ImportsSection) (*ImportsSectionMachine, error) {
	// FIXME: implement tests
	c, err := NewImportsTransition(start)
	if err != nil {
		return nil, err
	}
	return &ImportsSectionMachine{current: c}, nil
}

// NewImportsTransition returns a new transition corresponding to the received
// value.
func NewImportsTransition(s ImportsSection) (ImportsTransition, error) {
	switch s {
	case ImportsSectionStd:
		return standardImportsTransition{}, nil
	case ImportsSectionExternal:
		return externalImportsTransition{}, nil
	case ImportsSectionLocal:
		return localImportsTransition{}, nil
	}
	return nil, errors.New("invalid imports value")
}

//-

// Current returns the current state.
func (s *ImportsSectionMachine) Current() ImportsTransition {
	return s.current
}

// Previous returns the previous state.
func (s *ImportsSectionMachine) Previous() ImportsTransition {
	return s.previous
}

// Transition updates the internal state.
func (s *ImportsSectionMachine) Transition(next ImportsSection) error {
	var (
		res ImportsTransition
		err error
	)

	switch next {
	case ImportsSectionStd:
		res, err = s.current.Standard()
	case ImportsSectionExternal:
		res, err = s.current.External()
	case ImportsSectionLocal:
		res, err = s.current.Local()
	default:
		err = errors.Errorf("invalid imports value: %d", next)
	}
	if err != nil {
		return err
	}

	s.previous = s.current
	s.current = res
	return nil
}

//-

// Validate makes sure the implemented `imports` declaration satisfies the
// following rules:
// * Group declaration is parenthesized
// * Packages are separated by a breaking line like this:
//   * First standard packages,
//   * Next external packages, and
//   * Finally local packages
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

		section := NewImportsSection(s.Path.Value, i.localPath)
		if i.fsm == nil {
			fsm, err := NewImportsSectionMachine(section)
			if err != nil {
				return errors.Wrap(errors.Errorf("invalid imports found: %s", err), errPrefix)
			}
			i.fsm = fsm
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

//-

func (externalImportsTransition) External() (ImportsTransition, error) {
	return externalImportsTransition{}, nil
}

func (externalImportsTransition) Local() (ImportsTransition, error) {
	return localImportsTransition{}, nil
}

func (externalImportsTransition) Standard() (ImportsTransition, error) {
	return nil, errors.New("standard imports is invalid, next one must be external or local.")
}

//-

func (localImportsTransition) External() (ImportsTransition, error) {
	return nil, errors.New("external imports is invalid, next one must be local")
}

func (localImportsTransition) Local() (ImportsTransition, error) {
	return localImportsTransition{}, nil
}

func (localImportsTransition) Standard() (ImportsTransition, error) {
	return nil, errors.New("standard imports is invalid, next one must be local")
}

//-

func (standardImportsTransition) External() (ImportsTransition, error) {
	return externalImportsTransition{}, nil
}

func (standardImportsTransition) Local() (ImportsTransition, error) {
	return localImportsTransition{}, nil
}

func (standardImportsTransition) Standard() (ImportsTransition, error) {
	return standardImportsTransition{}, nil
}
