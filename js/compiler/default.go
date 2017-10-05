package compiler

import (
	"github.com/dop251/goja"
)

func init() ***REMOVED***
	c, err := New()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	DefaultCompiler = c
***REMOVED***

var DefaultCompiler *Compiler

func Transform(src, filename string) (code string, srcmap SourceMap, err error) ***REMOVED***
	return DefaultCompiler.Transform(src, filename)
***REMOVED***

func Compile(src, filename string, pre, post string, strict bool) (*goja.Program, string, error) ***REMOVED***
	return DefaultCompiler.Compile(src, filename, pre, post, strict)
***REMOVED***
