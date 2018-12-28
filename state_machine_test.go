package nit_test

import (
	"testing"

	"github.com/MarioCarrion/nit"
)

func TestImportsSectionMachine_Current(t *testing.T) {
	i := nit.NewImportsSectionMachine(nit.ImportsSectionStd)
	if i.Current() != nit.ImportsSectionStd {
		t.Fatalf("expected current value does not match")
	}
}

func TestImportsSectionMachine_Previous(t *testing.T) {
	i := nit.NewImportsSectionMachine(nit.ImportsSectionStd)
	if err := i.Transition(nit.ImportsSectionLocal); err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	if i.Previous() != nit.ImportsSectionStd {
		t.Fatalf("expected current value does not match")
	}
}

func TestImportsSectionMachine_Transition(t *testing.T) {
	tests := [...]struct {
		name          string
		startState    nit.ImportsSection
		validStates   []nit.ImportsSection
		invalidStates []nit.ImportsSection
	}{
		{
			"Standard",
			nit.ImportsSectionStd,
			[]nit.ImportsSection{
				nit.ImportsSectionStd,
				nit.ImportsSectionExternal,
				nit.ImportsSectionLocal,
			},
			[]nit.ImportsSection{},
		},
		{
			"External",
			nit.ImportsSectionExternal,
			[]nit.ImportsSection{
				nit.ImportsSectionExternal,
				nit.ImportsSectionLocal,
			},
			[]nit.ImportsSection{nit.ImportsSectionStd},
		},
		{
			"Local",
			nit.ImportsSectionLocal,
			[]nit.ImportsSection{nit.ImportsSectionLocal},
			[]nit.ImportsSection{
				nit.ImportsSectionStd,
				nit.ImportsSectionExternal,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			for _, s := range tt.validStates {
				machine := nit.NewImportsSectionMachine(tt.startState)
				if err := machine.Transition(s); err != nil {
					t.Fatalf("expected no error, to %s", err)
				}
			}

			for _, s := range tt.invalidStates {
				machine := nit.NewImportsSectionMachine(tt.startState)
				if err := machine.Transition(s); err == nil {
					t.Fatalf("expected error, got nil")
				}
			}
		})
	}
}

func TestSectionMachine_Transition(t *testing.T) {
	tests := [...]struct {
		name          string
		factory       func() nit.SectionMachine
		validStates   []nit.Section
		invalidStates []nit.Section
	}{
		{
			"Start",
			func() nit.SectionMachine {
				return nit.SectionMachine{}
			},
			[]nit.Section{
				nit.SectionImports,
				nit.SectionConsts,
				nit.SectionTypes,
				nit.SectionVars,
				nit.SectionFuncs,
				nit.SectionMethods,
			},
			[]nit.Section{
				nit.Section(99),
				nit.SectionStart,
			},
		},
		{
			"Imports",
			func() nit.SectionMachine {
				v := nit.SectionMachine{}
				v.Transition(nit.SectionImports)
				return v
			},
			[]nit.Section{
				nit.SectionConsts,
				nit.SectionTypes,
				nit.SectionVars,
				nit.SectionFuncs,
				nit.SectionMethods,
			},
			[]nit.Section{
				nit.SectionStart,
				nit.SectionImports,
			},
		},
		{
			"Consts",
			func() nit.SectionMachine {
				v := nit.SectionMachine{}
				v.Transition(nit.SectionImports)
				v.Transition(nit.SectionConsts)
				return v
			},
			[]nit.Section{
				nit.SectionConsts,
				nit.SectionTypes,
				nit.SectionVars,
				nit.SectionFuncs,
				nit.SectionMethods,
			},
			[]nit.Section{
				nit.SectionStart,
				nit.SectionImports,
			},
		},
		{
			"Types",
			func() nit.SectionMachine {
				v := nit.SectionMachine{}
				v.Transition(nit.SectionImports)
				v.Transition(nit.SectionConsts)
				v.Transition(nit.SectionTypes)
				return v
			},
			[]nit.Section{
				nit.SectionVars,
				nit.SectionFuncs,
				nit.SectionMethods,
			},
			[]nit.Section{
				nit.SectionStart,
				nit.SectionImports,
				nit.SectionTypes,
				nit.SectionConsts,
			},
		},
		{
			"Vars",
			func() nit.SectionMachine {
				v := nit.SectionMachine{}
				v.Transition(nit.SectionImports)
				v.Transition(nit.SectionConsts)
				v.Transition(nit.SectionTypes)
				v.Transition(nit.SectionVars)
				return v
			},
			[]nit.Section{
				nit.SectionFuncs,
				nit.SectionMethods,
			},
			[]nit.Section{
				nit.SectionStart,
				nit.SectionImports,
				nit.SectionTypes,
				nit.SectionVars,
				nit.SectionConsts,
			},
		},
		{
			"Funcs",
			func() nit.SectionMachine {
				v := nit.SectionMachine{}
				v.Transition(nit.SectionImports)
				v.Transition(nit.SectionConsts)
				v.Transition(nit.SectionTypes)
				v.Transition(nit.SectionVars)
				v.Transition(nit.SectionFuncs)
				return v
			},
			[]nit.Section{
				nit.SectionFuncs,
				nit.SectionMethods,
			},
			[]nit.Section{
				nit.SectionStart,
				nit.SectionImports,
				nit.SectionConsts,
				nit.SectionTypes,
				nit.SectionVars,
			},
		},
		{
			"Methods",
			func() nit.SectionMachine {
				v := nit.SectionMachine{}
				v.Transition(nit.SectionImports)
				v.Transition(nit.SectionConsts)
				v.Transition(nit.SectionTypes)
				v.Transition(nit.SectionVars)
				v.Transition(nit.SectionFuncs)
				v.Transition(nit.SectionMethods)
				return v
			},
			[]nit.Section{
				nit.SectionMethods,
			},
			[]nit.Section{
				nit.SectionStart,
				nit.SectionImports,
				nit.SectionConsts,
				nit.SectionTypes,
				nit.SectionVars,
				nit.SectionFuncs,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			for _, s := range tt.validStates {
				validator := tt.factory()
				if err := validator.Transition(s); err != nil {
					ts.Fatalf("expected no error, got %s", err)
				}
			}

			for _, s := range tt.invalidStates {
				validator := tt.factory()
				if err := validator.Transition(s); err == nil {
					ts.Fatalf("expected error, got nil")
				}
			}
		})
	}
}
