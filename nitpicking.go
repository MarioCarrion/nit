package nitpicking

import (
	"fmt"
)

type (
	// State represents a validation state.
	State int64

	// Validator represents the nitpicking validator.
	Validator struct {
		state State
	}
)

const (
	// StateStart defines the initial State in the validator.
	StateStart State = iota
	// StateImports defines the `import` state.
	StateImports
	// StateConsts defines the `const` state.
	StateConsts
	// StateTypes defines the `type` state.
	StateTypes
	// StateVars defines the `var` state.
	StateVars
	// StateFuncs defines the `func` state.
	StateFuncs
	// StateMethods defines the `method` state.
	StateMethods
)

// Transition updates the validator's internal state.
func (v *Validator) Transition(next State) error { //nolint:gocyclo
	switch v.state {
	case StateStart:
		if next != StateImports {
			return fmt.Errorf("next section must be: `imports`")
		}
	case StateImports:
		if next != StateConsts {
			return fmt.Errorf("next section must be: `const`")
		}
	case StateConsts:
		if next != StateTypes && next != StateConsts {
			return fmt.Errorf("next section must be either: `const` or `type`")
		}
	case StateTypes:
		if next != StateVars {
			return fmt.Errorf("next section must be: `vars`")
		}
	case StateVars:
		if next != StateFuncs {
			return fmt.Errorf("next section must be: functions")
		}
	case StateFuncs:
		if next != StateMethods && next != StateFuncs {
			return fmt.Errorf("next section must be: functions or methods")
		}
	case StateMethods:
		if next != StateMethods {
			return fmt.Errorf("next section must be: methods")
		}
	}
	v.state = next
	return nil
}
