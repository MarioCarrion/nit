package nit_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/MarioCarrion/nit"
)

func TestBreakComments(t *testing.T) {
	tests := [...]struct {
		name     string
		filename string
		expected []int
	}{
		{
			"OK",
			"break_comments1.go",
			[]int{8, 13, 19},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			f, fset, err := newParserFile(ts, tt.filename)
			if err != nil {
				ts.Fatalf("expected no error, got %s", err)
			}

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
		})
	}
}
