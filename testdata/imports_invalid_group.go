package testdata

import (
	"fmt"

	"github.com/MarioCarrion/nit"
	"github.com/pkg/errors"
)

func ImportsInvalidGroup() {
	fmt.Println("%s%s", nit.Nitpicker{}, errors.New(""))
}
