[![Go Report Card](https://goreportcard.com/badge/github.com/MarioCarrion/nitpicking)](https://goreportcard.com/report/github.com/MarioCarrion/nitpicking)
[![Circle CI](https://circleci.com/gh/MarioCarrion/nitpicking.svg?style=svg)](https://circleci.com/gh/MarioCarrion/nitpicking)
[![coverage](https://gocover.io/_badge/github.com/MarioCarrion/nitpicking?0 "coverage")](http://gocover.io/github.com/MarioCarrion/nitpicking)

# nitpicking

I'm nitpicking your code.

## What is this?

A really, really nitpicking linter that complains when the code is not organized according to the following very opinionated rules:

1. `imports` is the first section
   * Requires section declaration,
   * One maximum,
   * Separated in 3 blocks: standard, external and same package.
1. `const` is the second section
   * Requires section declaration,
   * Multiple allowed,
   * Section must be sorted: exported first, then unexported.
1. `type` is the third section
   * Requires section declaration,
   * Section must be sorted: exported first, then unexported.
1. `var` is the fourth section
   * Requires section declaration,
   * Section must be sorted: exported first, then unexported.
1. `func` is the fifth section
   * Must be sorted, exported first, then unexported.
1. `func` method, is the sixth section
   * Must be sorted by type, exported first, then unexported.
