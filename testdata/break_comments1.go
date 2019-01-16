package testdata

import "fmt"

type (
	BreakCommentType1 int //- NOT-DETECTED

	//- DETECTED

	BreakCommentType2 int
)

//- DETECTED

func BreakComment1() {
	fmt.Println("example") //- NOT-DETECTED
}

//- DETECTED

func BreakComment2() {}
