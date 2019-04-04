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
	}
)

const (
	breakComment  = "//-"
	generatedFile = `^// Code generated .* DO NOT EDIT\.$`
)

// NewBreakComments returns all the valid break-like comments.
func NewBreakComments(fset *token.FileSet, comments []*ast.CommentGroup) *BreakComments {
	r := BreakComments{}

	re, _ := regexp.Compile(generatedFile)

	for _, c := range comments {
		for _, c1 := range c.List {
			if re.MatchString(c1.Text) {
				r.generatedFile = true
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
