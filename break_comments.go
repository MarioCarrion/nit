package nit

import (
	"go/ast"
	"go/token"
	"strings"
)

type (
	// BreakComments defines all the found break-like comments in the file.
	BreakComments struct {
		index    int
		comments []int
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
