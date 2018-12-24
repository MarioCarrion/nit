package testdata

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/MarioCarrion/nitpicking"
)

func ImportsValid() {
	fmt.Println("%s", errors.New(""), nitpicking.Nitpicker{})
}
