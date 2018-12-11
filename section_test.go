package nitpicking_test

import (
	"testing"

	"github.com/MarioCarrion/nitpicking"
)

func TestNewImportsSection(t *testing.T) {
	tests := [...]struct {
		name                string
		inputPath           string
		inputLocaBasePrefix string
		expected            nitpicking.ImportsSection
	}{
		{
			"ImportsSectionStd",
			"fmt",
			"github.com/MarioCarrion/nitpicking",
			nitpicking.ImportsSectionStd,
		},
		{
			"ImportsSectionExternal",
			"github.com/golang/go",
			"github.com/MarioCarrion/nitpicking",
			nitpicking.ImportsSectionExternal,
		},
		{
			"ImportsSectionLocal",
			"github.com/MarioCarrion/nitpicking/something",
			"github.com/MarioCarrion/nitpicking",
			nitpicking.ImportsSectionLocal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			if actual := nitpicking.NewImportsSection(tt.inputPath, tt.inputLocaBasePrefix); actual != tt.expected {
				t.Fatalf("expected %+v, actual %+v", tt.expected, actual)
			}
		})
	}
}
