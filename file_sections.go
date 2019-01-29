package nit

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/pkg/errors"
)

type (
	// FileSection represents a code section.
	FileSection uint8

	// FileSectionMachine represents the state machine determining the right
	// order for the sections in file.
	FileSectionMachine struct {
		current FileSectionTransition
	}

	// FileSectionTransition represents one of the 6 valid sections in the
	// current file, it defines the State Machine transition rules for that
	// concrete section, the order to be expected is:
	//	"Imports" -> "Types" -> "Consts" -> "Vars" -> "Funcs" -> "Methods"
	FileSectionTransition interface {
		Imports() (FileSectionTransition, error)
		Types() (FileSectionTransition, error)
		Consts() (FileSectionTransition, error)
		Vars() (FileSectionTransition, error)
		Funcs() (FileSectionTransition, error)
		Methods() (FileSectionTransition, error)
	}

	constsFileSection  struct{}
	funcsFileSection   struct{}
	importsFileSection struct{}
	methodsFileSection struct{}
	typesFileSection   struct{}
	varsFileSection    struct{}
)

const (
	// FileSectionImports defines the `import` state.
	FileSectionImports FileSection = iota // FIXME kill

	// FileSectionTypes defines the `type` state.
	FileSectionTypes

	// FileSectionConsts defines the `const` state.
	FileSectionConsts

	// FileSectionVars defines the `var` state.
	FileSectionVars

	// FileSectionFuncs defines the `func` state.
	FileSectionFuncs

	// SectionMethods defines the `method` state.
	FileSectionMethods
)

// NewFileSectionTransition returns a new transition corresponding to the
// received value.
func NewFileSectionMachine(start FileSection) (*FileSectionMachine, error) {
	// FIXME: implement tests
	c, err := NewFileSectionTransition(start)
	if err != nil {
		return nil, err
	}
	return &FileSectionMachine{current: c}, nil
}

// NewFileSectionTransition returns a new transition corresponding to the
// received value.
func NewFileSectionTransition(s FileSection) (FileSectionTransition, error) {
	switch s {
	case FileSectionConsts:
		return constsFileSection{}, nil
	case FileSectionFuncs:
		return funcsFileSection{}, nil
	case FileSectionImports:
		return importsFileSection{}, nil
	case FileSectionMethods:
		return methodsFileSection{}, nil
	case FileSectionTypes:
		return typesFileSection{}, nil
	case FileSectionVars:
		return varsFileSection{}, nil
	}

	return nil, errors.New("invalid file section value")
}

// NewFuncDeclFileSection returns a new State that matches the decl type.
func NewFuncDeclFileSection(decl *ast.FuncDecl) (FileSection, error) {
	if decl.Recv == nil {
		return FileSectionFuncs, nil
	}
	return FileSectionMethods, nil
}

// NewGenDeclFileSection returns a new State that matches the decl type.
func NewGenDeclFileSection(decl *ast.GenDecl) (FileSection, error) {
	switch decl.Tok {
	case token.IMPORT:
		return FileSectionImports, nil
	case token.CONST:
		return FileSectionConsts, nil
	case token.TYPE:
		return FileSectionTypes, nil
	case token.VAR:
		return FileSectionVars, nil
	}

	return FileSectionImports, fmt.Errorf("unknown generic declaration node")
}

// Transition updates the internal state.
func (v *FileSectionMachine) Transition(next FileSection) error { //nolint:gocyclo
	var (
		res FileSectionTransition
		err error
	)

	switch next {
	case FileSectionImports:
		res, err = v.current.Imports()
	case FileSectionTypes:
		res, err = v.current.Types()
	case FileSectionConsts:
		res, err = v.current.Consts()
	case FileSectionVars:
		res, err = v.current.Vars()
	case FileSectionFuncs:
		res, err = v.current.Funcs()
	case FileSectionMethods:
		res, err = v.current.Methods()
	default:
		err = errors.Errorf("invalid file section value: %d", next)
	}
	if err != nil {
		return err
	}

	v.current = res
	return nil
}

//-

func (constsFileSection) Consts() (FileSectionTransition, error) {
	return constsFileSection{}, nil
}

func (constsFileSection) Funcs() (FileSectionTransition, error) {
	return funcsFileSection{}, nil
}

func (constsFileSection) Imports() (FileSectionTransition, error) {
	return nil, errors.New("imports is invalid, next one must be vars, funcs or methods")
}

func (constsFileSection) Methods() (FileSectionTransition, error) {
	return methodsFileSection{}, nil
}

func (constsFileSection) Types() (FileSectionTransition, error) {
	return nil, errors.New("types is invalid, next one must be vars, funcs or methods")
}

func (constsFileSection) Vars() (FileSectionTransition, error) {
	return varsFileSection{}, nil
}

//-

func (funcsFileSection) Consts() (FileSectionTransition, error) {
	return nil, errors.New("const is invalid, next one must be funcs or methods")
}

func (funcsFileSection) Funcs() (FileSectionTransition, error) {
	return funcsFileSection{}, nil
}

func (funcsFileSection) Imports() (FileSectionTransition, error) {
	return nil, errors.New("imports is invalid, next one must be funcs or methods")
}

func (funcsFileSection) Methods() (FileSectionTransition, error) {
	return methodsFileSection{}, nil
}

func (funcsFileSection) Types() (FileSectionTransition, error) {
	return nil, errors.New("type is invalid, next one must be funcs or methods")
}

func (funcsFileSection) Vars() (FileSectionTransition, error) {
	return nil, errors.New("vars is invalid, next one must be funcs or methods")
}

//-

func (importsFileSection) Consts() (FileSectionTransition, error) {
	return constsFileSection{}, nil
}

func (importsFileSection) Funcs() (FileSectionTransition, error) {
	return funcsFileSection{}, nil
}

func (importsFileSection) Imports() (FileSectionTransition, error) {
	return importsFileSection{}, nil
}

func (importsFileSection) Methods() (FileSectionTransition, error) {
	return methodsFileSection{}, nil
}

func (importsFileSection) Types() (FileSectionTransition, error) {
	return typesFileSection{}, nil
}

func (importsFileSection) Vars() (FileSectionTransition, error) {
	return varsFileSection{}, nil
}

//-

func (methodsFileSection) Consts() (FileSectionTransition, error) {
	return nil, errors.New("consts is invalid, next one must be methods")
}

func (methodsFileSection) Funcs() (FileSectionTransition, error) {
	return nil, errors.New("funcs is invalid, next one must be methods")
}

func (methodsFileSection) Imports() (FileSectionTransition, error) {
	return nil, errors.New("imports is invalid, next one must be methods")
}

func (methodsFileSection) Methods() (FileSectionTransition, error) {
	return methodsFileSection{}, nil
}

func (methodsFileSection) Types() (FileSectionTransition, error) {
	return nil, errors.New("types is invalid, next one must be methods")
}

func (methodsFileSection) Vars() (FileSectionTransition, error) {
	return nil, errors.New("vars is invalid, next one must be methods")
}

//-

func (typesFileSection) Consts() (FileSectionTransition, error) {
	return constsFileSection{}, nil
}

func (typesFileSection) Funcs() (FileSectionTransition, error) {
	return funcsFileSection{}, nil
}

func (typesFileSection) Imports() (FileSectionTransition, error) {
	return nil, errors.New("imports is invalid, next one must be const, vars, funcs or method")
}

func (typesFileSection) Methods() (FileSectionTransition, error) {
	return methodsFileSection{}, nil
}

func (typesFileSection) Types() (FileSectionTransition, error) {
	return typesFileSection{}, nil
}

func (typesFileSection) Vars() (FileSectionTransition, error) {
	return varsFileSection{}, nil
}

//-

func (varsFileSection) Consts() (FileSectionTransition, error) {
	return nil, errors.New("types is invalid, next one must be func or methods")
}

func (varsFileSection) Funcs() (FileSectionTransition, error) {
	return funcsFileSection{}, nil
}

func (varsFileSection) Imports() (FileSectionTransition, error) {
	return nil, errors.New("types is invalid, next one must be func or methods")
}

func (varsFileSection) Methods() (FileSectionTransition, error) {
	return methodsFileSection{}, nil
}

func (varsFileSection) Types() (FileSectionTransition, error) {
	return nil, errors.New("types is invalid, next one must be func or methods")
}

func (varsFileSection) Vars() (FileSectionTransition, error) {
	return varsFileSection{}, nil
}
