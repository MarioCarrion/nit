package testdata

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/MarioCarrion/nit"
)

func ImportsMissingLine() {
	fmt.Println("%s%s", nitpicking.Nitpicker{}, errors.New(""))
}
