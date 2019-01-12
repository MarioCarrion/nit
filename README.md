# nit

[![Go Report Card](https://goreportcard.com/badge/github.com/MarioCarrion/nit)](https://goreportcard.com/report/github.com/MarioCarrion/nit)
[![Circle CI](https://circleci.com/gh/MarioCarrion/nit.svg?style=svg)](https://circleci.com/gh/MarioCarrion/nit)
[![coverage](https://gocover.io/_badge/github.com/MarioCarrion/nit?0 "coverage")](http://gocover.io/github.com/MarioCarrion/nit)

I'm nitpicking your code.

## What is this?

A really, really nitpicking linter that complains when the code is not organized according to the following very opinionated rules:

1. `imports` is the first section
   * Requires section declaration,
   * One maximum,
   * Separated in 3 blocks: standard, external and same package (local).
1. `type` is the second section
   * Requires section declaration,
   * Section must be sorted: exported first, then unexported.
1. `const` is the third section
   * Requires section declaration,
   * Multiple allowed,
   * Section must be sorted: exported first, then unexported.
1. `var` is the fourth section
   * Requires section declaration,
   * Section must be sorted: exported first, then unexported.
1. `func` is the fifth section
   * Must be sorted, exported first, then unexported.
1. `func` method, is the sixth section
   * Must be sorted by type, exported first, then unexported.

![code](code.png "code organization in file")

### Development

Requires [`dep`](https://github.com/golang/dep), you can use [retool](https://github.com/twitchtv/retool) for installing that dependency.
