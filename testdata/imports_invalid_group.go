package testdata

import (
	"fmt"

	"github.com/MarioCarrion/nit"
	"github.com/pkg/errors"
)

func ImportsInvalidGroup() {
	fmt.Printf("%+v%s", nit.Nitpicker{}, errors.New(""))
}
