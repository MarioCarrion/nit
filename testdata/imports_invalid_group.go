package testdata

import (
	"fmt"

	"github.com/MarioCarrion/nitpicking"
	"github.com/pkg/errors"
)

func ImportsInvalidGroup() {
	fmt.Println("%s%s", nitpicking.Nitpicker{}, errors.New(""))
}
