package main

import (
	"flag"
	"fmt"
	"go/build"
	"os"
	"path/filepath"

	"github.com/MarioCarrion/nitpicking"
)

func main() {
	//nolint: errcheck
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n%s [packages]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flag.PrintDefaults()
	}

	localPkg := flag.String("pkg", "", "local package")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("missing packages")
		os.Exit(1)
	}

	for _, pkg := range flag.Args() {
		p, err := build.Import(pkg, ".", 0)
		if err != nil {
			fmt.Printf("error importing %s\n", pkg)
			os.Exit(1)
		}

		var failed bool

		for _, f := range p.GoFiles {
			fullpath := filepath.Join(p.Dir, f)
			v := nitpicking.Nitpicker{LocalPath: *localPkg}
			if err := v.Validate(fullpath); err != nil {
				failed = true
				fmt.Println(err)
			}
		}

		if failed {
			os.Exit(1)
		}
	}
}
