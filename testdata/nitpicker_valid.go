package testdata

import (
	"fmt"
)

type (
	Example struct{}
)

const (
	FirstConst = "1"
)

const (
	SecondConst = "2"
)

var (
	varOne = 1
)

func NitpickerValid() {
	fmt.Println("")
}

func (*Example) ExampleMethod() {
	fmt.Println("")
}
