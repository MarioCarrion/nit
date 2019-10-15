package nit_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/MarioCarrion/nit"
)

func TestBreakComments(t *testing.T) {
	type expected struct {
		result        []int
		generatedCode bool
		noLintNit     bool
	}

	tests := [...]struct {
		name     string
		filename string
		expected expected
	}{
		{
			"OK",
			"break_comments1.go",
			expected{result: []int{8, 13, 19}},
		},
		{
			"OK: code generated",
			"break_comments2.go",
			expected{generatedCode: true},
		},
		{
			"OK: nolint nit",
			"break_comments3.go",
			expected{noLintNit: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			f, fset := newParserFile(ts, tt.filename)

			var actual []int

			bc := nit.NewBreakComments(fset, f.Comments)
			for _, v := range tt.expected.result {
				bc.MoveTo(v)
				v = bc.Next()
				if v == -1 {
					break
				}
				actual = append(actual, v)
			}

			if !cmp.Equal(tt.expected.result, actual) {
				ts.Errorf("expected values do not match: %s", cmp.Diff(tt.expected, actual))
			}

			if tt.expected.generatedCode != bc.HasGeneratedCode() {
				ts.Errorf("expected %t, got %t", tt.expected.generatedCode, bc.HasGeneratedCode())
			}

			if tt.expected.noLintNit != bc.HasNoLintNit() {
				ts.Errorf("expected %t, got %t", tt.expected.noLintNit, bc.HasNoLintNit())
			}
		})
	}
}
