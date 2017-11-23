# go-isatty

[![Godoc Reference](https://godoc.org/github.com/mattn/go-isatty?status.svg)](http://godoc.org/github.com/mattn/go-isatty)
[![Build Status](https://travis-ci.org/mattn/go-isatty.svg?branch=master)](https://travis-ci.org/mattn/go-isatty)
[![Coverage Status](https://coveralls.io/repos/github/mattn/go-isatty/badge.svg?branch=master)](https://coveralls.io/github/mattn/go-isatty?branch=master)
[![Go Report Card](https://goreportcard.com/badge/mattn/go-isatty)](https://goreportcard.com/report/mattn/go-isatty)

isatty for golang

## Usage

```go
package main

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"os"
)

func main() ***REMOVED***
	if isatty.IsTerminal(os.Stdout.Fd()) ***REMOVED***
		fmt.Println("Is Terminal")
	***REMOVED*** else if isatty.IsCygwinTerminal(os.Stdout.Fd()) ***REMOVED***
		fmt.Println("Is Cygwin/MSYS2 Terminal")
	***REMOVED*** else ***REMOVED***
		fmt.Println("Is Not Terminal")
	***REMOVED***
***REMOVED***
```

## Installation

```
$ go get github.com/mattn/go-isatty
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a mattn)

## Thanks

* k-takata: base idea for IsCygwinTerminal

    https://github.com/k-takata/go-iscygpty