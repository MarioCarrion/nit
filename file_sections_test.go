package nit_test

import (
	"testing"

	"github.com/MarioCarrion/nit"
)

func TestNewFileSectionMachine(t *testing.T) {
	_, err := nit.NewFileSectionMachine(nit.FileSection(99))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestNewFileSectionTransition(t *testing.T) {
	_, err := nit.NewFileSectionTransition(nit.FileSection(99))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestFileSectionMachine(t *testing.T) {
	s := nit.FileSectionMachine{}

	if err := s.Transition(nit.FileSection(99)); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

//-

func TestConstsFileSectionTransition(t *testing.T) {
	i, _ := nit.NewFileSectionTransition(nit.FileSectionConsts)

	tests := [...]struct {
		name          string
		transition    func() (nit.FileSectionTransition, error)
		expectedError bool
	}{
		{
			"Consts",
			i.Consts,
			false,
		},
		{
			"Funcs",
			i.Funcs,
			false,
		},
		{
			"Imports",
			i.Imports,
			true,
		},
		{
			"Methods",
			i.Methods,
			false,
		},
		{
			"Types",
			i.Types,
			true,
		},
		{
			"Vars",
			i.Vars,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			tr, err := tt.transition()
			if (err != nil) != tt.expectedError {
				ts.Fatalf("expected error %t, got %s", tt.expectedError, err)
			}
			if err != nil && tr != nil {
				ts.Fatalf("expected nil transition on error")
			}
		})
	}
}

func TestFuncsFileSectionTransition(t *testing.T) {
	i, _ := nit.NewFileSectionTransition(nit.FileSectionFuncs)

	tests := [...]struct {
		name          string
		transition    func() (nit.FileSectionTransition, error)
		expectedError bool
	}{
		{
			"Consts",
			i.Consts,
			true,
		},
		{
			"Funcs",
			i.Funcs,
			false,
		},
		{
			"Imports",
			i.Imports,
			true,
		},
		{
			"Methods",
			i.Methods,
			false,
		},
		{
			"Types",
			i.Types,
			true,
		},
		{
			"Vars",
			i.Vars,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			tr, err := tt.transition()
			if (err != nil) != tt.expectedError {
				ts.Fatalf("expected error %t, got %s", tt.expectedError, err)
			}
			if err != nil && tr != nil {
				ts.Fatalf("expected nil transition on error")
			}
		})
	}
}

func TestImportsFileSectionTransition(t *testing.T) {
	i, _ := nit.NewFileSectionTransition(nit.FileSectionImports)

	tests := [...]struct {
		name          string
		transition    func() (nit.FileSectionTransition, error)
		expectedError bool
	}{
		{
			"Consts",
			i.Consts,
			false,
		},
		{
			"Funcs",
			i.Funcs,
			false,
		},
		{
			"Imports",
			i.Imports,
			false,
		},
		{
			"Methods",
			i.Methods,
			false,
		},
		{
			"Types",
			i.Types,
			false,
		},
		{
			"Vars",
			i.Vars,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			tr, err := tt.transition()
			if (err != nil) != tt.expectedError {
				ts.Fatalf("expected error %t, got %s", tt.expectedError, err)
			}
			if err != nil && tr != nil {
				ts.Fatalf("expected nil transition on error")
			}
		})
	}
}

func TestMethodsFileSectionTransition(t *testing.T) {
	i, _ := nit.NewFileSectionTransition(nit.FileSectionMethods)

	tests := [...]struct {
		name          string
		transition    func() (nit.FileSectionTransition, error)
		expectedError bool
	}{
		{
			"Consts",
			i.Consts,
			true,
		},
		{
			"Funcs",
			i.Funcs,
			true,
		},
		{
			"Imports",
			i.Imports,
			true,
		},
		{
			"Methods",
			i.Methods,
			false,
		},
		{
			"Types",
			i.Types,
			true,
		},
		{
			"Vars",
			i.Vars,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			tr, err := tt.transition()
			if (err != nil) != tt.expectedError {
				ts.Fatalf("expected error %t, got %s", tt.expectedError, err)
			}
			if err != nil && tr != nil {
				ts.Fatalf("expected nil transition on error")
			}
		})
	}
}

func TestTypesFileSectionTransition(t *testing.T) {
	i, _ := nit.NewFileSectionTransition(nit.FileSectionTypes)

	tests := [...]struct {
		name          string
		transition    func() (nit.FileSectionTransition, error)
		expectedError bool
	}{
		{
			"Consts",
			i.Consts,
			false,
		},
		{
			"Funcs",
			i.Funcs,
			false,
		},
		{
			"Imports",
			i.Imports,
			true,
		},
		{
			"Methods",
			i.Methods,
			false,
		},
		{
			"Types",
			i.Types,
			false,
		},
		{
			"Vars",
			i.Vars,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			tr, err := tt.transition()
			if (err != nil) != tt.expectedError {
				ts.Fatalf("expected error %t, got %s", tt.expectedError, err)
			}
			if err != nil && tr != nil {
				ts.Fatalf("expected nil transition on error")
			}
		})
	}
}

func TestVarsFileSectionTransition(t *testing.T) {
	i, _ := nit.NewFileSectionTransition(nit.FileSectionVars)

	tests := [...]struct {
		name          string
		transition    func() (nit.FileSectionTransition, error)
		expectedError bool
	}{
		{
			"Consts",
			i.Consts,
			true,
		},
		{
			"Funcs",
			i.Funcs,
			false,
		},
		{
			"Imports",
			i.Imports,
			true,
		},
		{
			"Methods",
			i.Methods,
			false,
		},
		{
			"Types",
			i.Types,
			true,
		},
		{
			"Vars",
			i.Vars,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			tr, err := tt.transition()
			if (err != nil) != tt.expectedError {
				ts.Fatalf("expected error %t, got %s", tt.expectedError, err)
			}
			if err != nil && tr != nil {
				ts.Fatalf("expected nil transition on error")
			}
		})
	}
}

//-
