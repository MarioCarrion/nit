// Package example is something valid!
//nolint
package example

import (
	"fmt"
)

import "log"

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
	log.Println("hi")
}
