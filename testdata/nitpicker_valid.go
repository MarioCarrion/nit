package testdata

import (
	"fmt"
)

const (
	FirstConst = "1"
)

const (
	SecondConst = "2"
)

type (
	Example struct{}
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
