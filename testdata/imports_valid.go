package testdata

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/MarioCarrion/nit"
)

func ImportsValid() {
	fmt.Println("%s", errors.New(""), nitpicking.Nitpicker{})
}
