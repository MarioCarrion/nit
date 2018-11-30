package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"github.com/MarioCarrion/nitpicking"
)

func main() {
	v := nitpicking.Validator{}

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "example/example.go", nil, 0)
	if err != nil {
		log.Fatalf("parsing error failed %s\n", err)
	}

	for _, s := range f.Decls {
		fmt.Printf("%T - %+v -- %t\n", s, s, s.End().IsValid())
		nextState, err := declToState(s)
		if err != nil {
			log.Fatalf("invalid type %T - %+v", s, s)
		}

		kind := "func"
		if gen, ok := s.(*ast.GenDecl); ok {
			kind = gen.Tok.String()
		}

		if err := v.Transition(nextState); err != nil {
			pos := fset.PositionFor(s.Pos(), false)
			log.Fatalf("%s section `%s` is invalid: %s", pos.String(), kind, err)
		}
	}
}

func declToState(d ast.Decl) (nitpicking.State, error) {
	switch v := d.(type) {
	case *ast.GenDecl:
		switch v.Tok {
		case token.IMPORT:
			return nitpicking.StateImports, nil
		case token.CONST:
			return nitpicking.StateConsts, nil
		case token.TYPE:
			return nitpicking.StateTypes, nil
		case token.VAR:
			return nitpicking.StateVars, nil
		}
	case *ast.FuncDecl:
		if v.Recv == nil {
			return nitpicking.StateFuncs, nil
		}
		return nitpicking.StateMethods, nil
	}
	return nitpicking.State(0), fmt.Errorf("unknown declaration state")
}
