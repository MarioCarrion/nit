package main

import (
	"flag"
	"fmt"
	"go/build"
	"os"
	"path/filepath"

	"github.com/MarioCarrion/nit"
)

//nolint: gochecknoglobals
var (
	commit  = "none" //-
	date    = "unknown"
	version = "dev"
)

func main() {
	//nolint: errcheck
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n%s [packages]\n", os.Args[0])
		flag.PrintDefaults()
	}
	//-

	localPkg := flag.String("pkg", "", "local package")
	skipGenerated := flag.Bool("skip-generated", false, "skip generated files")
	showVersion := flag.Bool("version", false, "prints current version information")

	flag.Parse()

	if *showVersion {
		fmt.Printf("%v, commit %v, built at %v\n", version, commit, date)
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		fmt.Println("missing `pkg` argument.")
		flag.Usage()
		os.Exit(1)
	}

	var failed bool

	for _, pkg := range flag.Args() {
		p, err := build.Import(pkg, ".", 0)
		if err != nil {
			fmt.Printf("error importing %s\n", pkg)
			os.Exit(1)
		}

		for _, f := range p.GoFiles {
			fullpath := filepath.Join(p.Dir, f)
			v := nit.Nitpicker{
				LocalPath:         *localPkg,
				SkipGeneratedFile: *skipGenerated,
			}
			if err := v.Validate(fullpath); err != nil {
				failed = true
				fmt.Println(err)
			}
		}
	}

	if failed {
		os.Exit(1)
	}
}
