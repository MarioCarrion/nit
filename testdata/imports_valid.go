package testdata

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/MarioCarrion/nit"
)

func ImportsValid() {
	fmt.Printf("%s%+v", errors.New(""), nit.Nitpicker{})
}
