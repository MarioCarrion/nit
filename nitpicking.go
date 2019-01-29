package nit

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/pkg/errors"
)

type (
	// BreakComments defines all the found break-like comments in the file.
	BreakComments struct {
		index    int
		comments []int
	}

	// Nitpicker defines the linter.
	Nitpicker struct {
		LocalPath  string
		fset       *token.FileSet
		fsm        *FileSectionMachine
		fvalidator *FuncsValidator
	}
)

const (
	breakComment = "//-"
)

// NewBreakComments returns all the valid break-like comments.
func NewBreakComments(fset *token.FileSet, comments []*ast.CommentGroup) BreakComments {
	r := BreakComments{}

	for _, c := range comments {
		for _, c1 := range c.List {
			if strings.HasPrefix(c1.Text, breakComment) {
				position := fset.PositionFor(c1.Pos(), false)
				if position.Column == 1 || position.Column == 2 { // left most either nested or not nested group declarations
					r.comments = append(r.comments, position.Line)
				}
			}
		}
	}
	return r
}

// Validate nitpicks the filename.
func (v *Nitpicker) Validate(filename string) error {
	v.fset = token.NewFileSet()
	f, err := parser.ParseFile(v.fset, filename, nil, parser.ParseComments)
	if err != nil {
		return errors.Wrap(err, "parsing file failed")
	}

	comments := NewBreakComments(v.fset, f.Comments)
	v.fvalidator = &FuncsValidator{Comments: &comments}

	for _, s := range f.Decls {
		// fmt.Printf("%d == %T - %+v -- %t\n", v.fset.PositionFor(s.Pos(), false).Line, s, s, s.End().IsValid())
		if err := v.validateToken(s); err != nil {
			return err
		}
	}

	return nil
}

//nolint:gocyclo
func (v *Nitpicker) validateToken(d ast.Decl) error {
	var (
		err       error
		genDecl   *ast.GenDecl
		funcDecl  *ast.FuncDecl
		nextState FileSection
	)

	switch t := d.(type) {
	case *ast.GenDecl:
		genDecl = t
		nextState, err = NewGenDeclFileSection(genDecl)
	case *ast.FuncDecl:
		funcDecl = t
		nextState, err = NewFuncDeclFileSection(funcDecl)
	default:
		return errors.New("unknown declaration state")
	}
	if err != nil {
		return err
	}

	if v.fsm == nil {
		fsm, err := NewFileSectionMachine(nextState)
		if err != nil {
			return errors.Wrap(err, v.fset.PositionFor(d.Pos(), false).String())
		}
		v.fsm = fsm
	}

	if err = v.fsm.Transition(nextState); err != nil {
		return errors.Wrap(err, v.fset.PositionFor(d.Pos(), false).String())
	}

	switch nextState {
	case FileSectionImports:
		validator := NewImportsValidator(v.LocalPath)
		if err := validator.Validate(genDecl, v.fset); err != nil {
			return err
		}
	case FileSectionTypes:
		validator := &TypesValidator{}
		if err := validator.Validate(genDecl, v.fset); err != nil {
			return err
		}
	case FileSectionConsts:
		validator := &ConstsValidator{}
		if err := validator.Validate(genDecl, v.fset); err != nil {
			return err
		}
	case FileSectionVars:
		validator := &VarsValidator{}
		if err := validator.Validate(genDecl, v.fset); err != nil {
			return err
		}
	case FileSectionFuncs:
		if err := v.fvalidator.Validate(funcDecl, v.fset); err != nil {
			return err
		}
	}

	return nil
}

// Returns the next break line if any, when no more left it returns -1.
func (c *BreakComments) Next() int {
	if c.index >= len(c.comments) {
		return -1
	}
	return c.comments[c.index]
}

// Moves current line cursor to the received line.
func (c *BreakComments) MoveTo(line int) {
	if c.index >= len(c.comments) {
		return
	}

	for _, v := range c.comments[c.index:] {
		if v >= line {
			break
		}
		c.index++
	}
}
