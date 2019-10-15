package nit

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

type (
	// BreakComments defines all the found break-like comments in the file.
	BreakComments struct {
		index         int
		comments      []int
		generatedFile bool
		nolintNit     bool
	}
)

const (
	breakComment  = "//-"
	generatedFile = `^// Code generated .* DO NOT EDIT\.$`
	nitLinter     = "^//nolint:nit$"
)

// NewBreakComments returns all the valid break-like comments.
func NewBreakComments(fset *token.FileSet, comments []*ast.CommentGroup) *BreakComments {
	r := BreakComments{}

	regen, _ := regexp.Compile(generatedFile)
	renit, _ := regexp.Compile(nitLinter)

	for _, c := range comments {
		for _, c1 := range c.List {
			if regen.MatchString(c1.Text) {
				r.generatedFile = true
			}

			if renit.MatchString(c1.Text) {
				r.nolintNit = true
			}

			if strings.HasPrefix(c1.Text, breakComment) {
				position := fset.PositionFor(c1.Pos(), false)
				if position.Column == 1 || position.Column == 2 { // left most either nested or not nested group declarations
					r.comments = append(r.comments, position.Line)
				}
			}
		}
	}

	return &r
}

// HasGeneratedCode indicates whether the current file contains a "code
// generated expression".
func (c *BreakComments) HasGeneratedCode() bool {
	return c.generatedFile
}

// HasNoLintNit indicates whether the current file contains a "do not lint nit"
// expression.
func (c *BreakComments) HasNoLintNit() bool {
	return c.nolintNit
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

// Returns the next break line if any, when no more left it returns -1.
func (c *BreakComments) Next() int {
	if c.index >= len(c.comments) {
		return -1
	}

	return c.comments[c.index]
}
