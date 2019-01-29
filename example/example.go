// Package example is something valid!
//nolint
package example

// import (
// 	"fmt"
// )

import (
	"fmt" // blabla
	ll "log"

	"github.com/MarioCarrion/nit"

	"github.com/pkg/errors"
)

const (
	something = "else"
)

const (
	something1 = "else"
)

var yolo = "hi"
var (
	yolo1 = "hi"
)

type (
	someValue int64
)

func Something() {
	fmt.Println("vim-go")
	ll.Println("hi")

	fmt.Println(nit.FileSection(0))
	fmt.Println(errors.New("bla"))
}
