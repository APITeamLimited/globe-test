package compiler

import (
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

var (
	lib      = rice.MustFindBox("lib")
	babelSrc = lib.MustString("babel-standalone-bower/babel.min.js")

	DefaultOpts = map[string]interface***REMOVED******REMOVED******REMOVED***
		"presets":       []string***REMOVED***"latest"***REMOVED***,
		"ast":           false,
		"sourceMaps":    false,
		"babelrc":       false,
		"compact":       false,
		"retainLines":   true,
		"highlightCode": false,
	***REMOVED***
)

// A Compiler uses Babel to compile ES6 code into something ES5-compatible.
type Compiler struct ***REMOVED***
	vm *goja.Runtime

	// JS pointers.
	this      goja.Value
	transform goja.Callable
***REMOVED***

// Constructs a new compiler.
func New() (*Compiler, error) ***REMOVED***
	c := &Compiler***REMOVED***vm: goja.New()***REMOVED***
	if _, err := c.vm.RunString(babelSrc); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c.this = c.vm.Get("Babel")
	thisObj := c.this.ToObject(c.vm)
	if err := c.vm.ExportTo(thisObj.Get("transform"), &c.transform); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return c, nil
***REMOVED***

// Transform the given code into ES5.
func (c *Compiler) Transform(src, filename string) (code string, srcmap SourceMap, err error) ***REMOVED***
	opts := make(map[string]interface***REMOVED******REMOVED***)
	for k, v := range DefaultOpts ***REMOVED***
		opts[k] = v
	***REMOVED***
	opts["filename"] = filename

	startTime := time.Now()
	v, err := c.transform(c.this, c.vm.ToValue(src), c.vm.ToValue(opts))
	if err != nil ***REMOVED***
		return code, srcmap, err
	***REMOVED***
	log.WithField("t", time.Since(startTime)).Debug("Babel: Transformed")
	vO := v.ToObject(c.vm)

	if err := c.vm.ExportTo(vO.Get("code"), &code); err != nil ***REMOVED***
		return code, srcmap, err
	***REMOVED***

	var rawmap map[string]interface***REMOVED******REMOVED***
	if err := c.vm.ExportTo(vO.Get("map"), &rawmap); err != nil ***REMOVED***
		return code, srcmap, err
	***REMOVED***
	if err := mapstructure.Decode(rawmap, &srcmap); err != nil ***REMOVED***
		return code, srcmap, err
	***REMOVED***

	return code, srcmap, nil
***REMOVED***

// Compiles the program, first trying ES5, then ES6.
func (c *Compiler) Compile(src, filename string, pre, post string, strict bool) (*goja.Program, string, error) ***REMOVED***
	return c.compile(src, filename, pre, post, strict, true)
***REMOVED***

func (c *Compiler) compile(src, filename string, pre, post string, strict, tryBabel bool) (*goja.Program, string, error) ***REMOVED***
	code := pre + src + post
	ast, err := parser.ParseFile(nil, filename, code, 0)
	if err != nil ***REMOVED***
		if tryBabel ***REMOVED***
			code, _, err := c.Transform(src, filename)
			if err != nil ***REMOVED***
				return nil, code, err
			***REMOVED***
			return c.compile(code, filename, pre, post, strict, false)
		***REMOVED***
		return nil, src, err
	***REMOVED***
	pgm, err := goja.CompileAST(ast, strict)
	return pgm, code, err
***REMOVED***
