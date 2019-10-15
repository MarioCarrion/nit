package main

import (
	"flag"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"github.com/MarioCarrion/nit"
)

//nolint: gochecknoglobals
var (
	commit  = "none" //-
	date    = "unknown"
	version = "dev"
)

//nolint: funlen
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
	nolint := flag.Bool("nolint", false, "enable nolint directive")
	includeTests := flag.Bool("include-tests", false, "include test files")
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

	nitpick := func(files []string) {
		for _, f := range files {
			if strings.HasSuffix(f, "_test.go") && !*includeTests {
				continue
			}

			v := nit.Nitpicker{
				LocalPath:         *localPkg,
				SkipGeneratedFile: *skipGenerated,
				NoLint:            *nolint,
			}

			if err := v.Validate(f); err != nil {
				failed = true

				fmt.Println(err)
			}
		}
	}

	for _, pkg := range flag.Args() {
		p, err := build.Import(pkg, ".", 0)
		if err != nil {
			fmt.Printf("error importing %s\n", pkg)
			os.Exit(1)
		}

		gofiles, _ := filepath.Glob(filepath.Join(p.Dir, "*.go"))
		nitpick(gofiles)
	}

	if failed {
		os.Exit(1)
	}
}
