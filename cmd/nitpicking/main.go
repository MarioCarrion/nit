package main

import (
	"fmt"
	"os"

	"github.com/MarioCarrion/nitpicking"
)

func main() {
	// fset := token.NewFileSet()

	// v := nitpicking.Nitpicker{FileSet: fset}
	v := nitpicking.Nitpicker{LocalPath: "github.com/MarioCarrion/nitpicking"}
	if err := v.Validate("example/example.go"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// f, err := parser.ParseFile(fset, "example/example.go", nil, 0)
	// if err != nil {
	// 	log.Fatalf("parsing error failed %s\n", err)
	// }

	// for _, s := range f.Decls {
	// 	fmt.Printf("%d == %T - %+v -- %t\n", fset.PositionFor(s.Pos(), false).Line, s, s, s.End().IsValid())

	// 	if err := v.Validate(s); err != nil {
	// 		pos := fset.PositionFor(s.Pos(), false)
	// 		log.Fatalf("%s section `%s` is invalid: %s", pos.String(), v.LastTokenKind, err)
	// 	}
	// }
}
