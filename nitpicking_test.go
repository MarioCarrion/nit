package nitpicking_test

import (
	"testing"

	"github.com/MarioCarrion/nitpicking"
)

func TestValidator_Transition(t *testing.T) {
	tests := [...]struct {
		name          string
		factory       func() nitpicking.Validator
		validStates   []nitpicking.State
		invalidStates []nitpicking.State
	}{
		{
			"Start",
			func() nitpicking.Validator {
				return nitpicking.Validator{}
			},
			[]nitpicking.State{nitpicking.StateImports},
			[]nitpicking.State{
				nitpicking.StateStart,
				nitpicking.StateConsts,
				nitpicking.StateTypes,
				nitpicking.StateVars,
				nitpicking.StateFuncs,
				nitpicking.StateMethods,
			},
		},
		{
			"Imports",
			func() nitpicking.Validator {
				v := nitpicking.Validator{}
				v.Transition(nitpicking.StateImports)
				return v
			},
			[]nitpicking.State{nitpicking.StateConsts},
			[]nitpicking.State{
				nitpicking.StateStart,
				nitpicking.StateImports,
				nitpicking.StateTypes,
				nitpicking.StateVars,
				nitpicking.StateFuncs,
				nitpicking.StateMethods,
			},
		},
		{
			"Consts",
			func() nitpicking.Validator {
				v := nitpicking.Validator{}
				v.Transition(nitpicking.StateImports)
				v.Transition(nitpicking.StateConsts)
				return v
			},
			[]nitpicking.State{
				nitpicking.StateTypes,
				nitpicking.StateConsts,
			},
			[]nitpicking.State{
				nitpicking.StateStart,
				nitpicking.StateImports,
				nitpicking.StateVars,
				nitpicking.StateFuncs,
				nitpicking.StateMethods,
			},
		},
		{
			"Types",
			func() nitpicking.Validator {
				v := nitpicking.Validator{}
				v.Transition(nitpicking.StateImports)
				v.Transition(nitpicking.StateConsts)
				v.Transition(nitpicking.StateTypes)
				return v
			},
			[]nitpicking.State{
				nitpicking.StateVars,
			},
			[]nitpicking.State{
				nitpicking.StateStart,
				nitpicking.StateImports,
				nitpicking.StateTypes,
				nitpicking.StateConsts,
				nitpicking.StateFuncs,
				nitpicking.StateMethods,
			},
		},
		{
			"Vars",
			func() nitpicking.Validator {
				v := nitpicking.Validator{}
				v.Transition(nitpicking.StateImports)
				v.Transition(nitpicking.StateConsts)
				v.Transition(nitpicking.StateTypes)
				v.Transition(nitpicking.StateVars)
				return v
			},
			[]nitpicking.State{
				nitpicking.StateFuncs,
			},
			[]nitpicking.State{
				nitpicking.StateStart,
				nitpicking.StateImports,
				nitpicking.StateConsts,
				nitpicking.StateTypes,
				nitpicking.StateVars,
				nitpicking.StateMethods,
			},
		},
		{
			"Funcs",
			func() nitpicking.Validator {
				v := nitpicking.Validator{}
				v.Transition(nitpicking.StateImports)
				v.Transition(nitpicking.StateConsts)
				v.Transition(nitpicking.StateTypes)
				v.Transition(nitpicking.StateVars)
				v.Transition(nitpicking.StateFuncs)
				return v
			},
			[]nitpicking.State{
				nitpicking.StateFuncs,
				nitpicking.StateMethods,
			},
			[]nitpicking.State{
				nitpicking.StateStart,
				nitpicking.StateImports,
				nitpicking.StateConsts,
				nitpicking.StateTypes,
				nitpicking.StateVars,
			},
		},
		{
			"Methods",
			func() nitpicking.Validator {
				v := nitpicking.Validator{}
				v.Transition(nitpicking.StateImports)
				v.Transition(nitpicking.StateConsts)
				v.Transition(nitpicking.StateTypes)
				v.Transition(nitpicking.StateVars)
				v.Transition(nitpicking.StateFuncs)
				v.Transition(nitpicking.StateMethods)
				return v
			},
			[]nitpicking.State{
				nitpicking.StateMethods,
			},
			[]nitpicking.State{
				nitpicking.StateStart,
				nitpicking.StateImports,
				nitpicking.StateConsts,
				nitpicking.StateTypes,
				nitpicking.StateVars,
				nitpicking.StateFuncs,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			for _, s := range tt.validStates {
				validator := tt.factory()
				if err := validator.Transition(s); err != nil {
					t.Fatalf("expected no error, got %s", err)
				}
			}

			for _, s := range tt.invalidStates {
				validator := tt.factory()
				if err := validator.Transition(s); err == nil {
					t.Fatalf("expected error, got nil")
				}
			}
		})
	}
}
