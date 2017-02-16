package compiler

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
