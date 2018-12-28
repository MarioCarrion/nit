package nit

import (
	"fmt"
)

type (
	// ImportsSectionMachine represents the `imports` code organization state machine.
	ImportsSectionMachine struct {
		current  ImportsSection
		previous ImportsSection
	}

	// SectionMachine represents the default sections code organization state machine.
	SectionMachine struct {
		current Section
	}
)

// NewImportsSectionMachine returns a new ImportsSectionMachine with the initial state as `start`
func NewImportsSectionMachine(start ImportsSection) *ImportsSectionMachine {
	return &ImportsSectionMachine{current: start}
}

// Current returns the current state.
func (s *ImportsSectionMachine) Current() ImportsSection {
	return s.current
}

// Previous returns the previous state.
func (s *ImportsSectionMachine) Previous() ImportsSection {
	return s.previous
}

// Transition updates the internal state.
func (s *ImportsSectionMachine) Transition(next ImportsSection) error {
	switch s.current {
	case ImportsSectionStd:
		// All sections are supported as next ones
	case ImportsSectionExternal:
		if next != ImportsSectionExternal && next != ImportsSectionLocal {
			return fmt.Errorf("next `imports` must be either external or local")
		}
	case ImportsSectionLocal:
		if next != ImportsSectionLocal {
			return fmt.Errorf("next `imports` must be local")
		}
	}
	s.previous = s.current
	s.current = next
	return nil
}

// Transition updates the internal state.
func (v *SectionMachine) Transition(next Section) error { //nolint:gocyclo
	switch v.current {
	case SectionStart:
		if next != SectionImports && next != SectionConsts && next != SectionTypes && next != SectionVars && next != SectionFuncs && next != SectionMethods {
			return fmt.Errorf("next section is invalid")
		}
	case SectionImports, SectionConsts:
		if next != SectionConsts && next != SectionTypes && next != SectionVars && next != SectionFuncs && next != SectionMethods {
			return fmt.Errorf("next section must either: `const`, `type`, `var` or funcs/methods")
		}
	case SectionTypes:
		if next != SectionVars && next != SectionFuncs && next != SectionMethods {
			return fmt.Errorf("next section must either: `var` or funcs/methods")
		}
	case SectionVars, SectionFuncs:
		if next != SectionFuncs && next != SectionMethods {
			return fmt.Errorf("next section must either: funcs or methods")
		}
	case SectionMethods:
		if next != SectionMethods {
			return fmt.Errorf("next section must be: methods")
		}
	}
	v.current = next
	return nil
}
