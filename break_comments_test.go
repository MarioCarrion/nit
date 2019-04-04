package nit_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/MarioCarrion/nit"
)

func TestBreakComments(t *testing.T) {
	tests := [...]struct {
		name                  string
		filename              string
		expected              []int
		expectedCodegenerated bool
	}{
		{
			"OK",
			"break_comments1.go",
			[]int{8, 13, 19},
			false,
		},
		{
			"OK: code generated",
			"break_comments2.go",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			f, fset := newParserFile(ts, tt.filename)

			var actual []int

			bc := nit.NewBreakComments(fset, f.Comments)
			for _, v := range tt.expected {
				bc.MoveTo(v)
				v = bc.Next()
				if v == -1 {
					break
				}
				actual = append(actual, v)
			}

			if !cmp.Equal(tt.expected, actual) {
				ts.Errorf("expected values do not match: %s", cmp.Diff(tt.expected, actual))
			}
			if tt.expectedCodegenerated != bc.HasGeneratedCode() {
				ts.Errorf("expected %t, got %t", tt.expectedCodegenerated, bc.HasGeneratedCode())
			}
		})
	}
}
