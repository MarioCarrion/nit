package testdata

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/MarioCarrion/nit"
)

func ImportsMissingLine() {
	fmt.Printf("%+v%s", nit.Nitpicker{}, errors.New(""))
}
