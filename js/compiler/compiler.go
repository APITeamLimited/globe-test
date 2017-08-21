package compiler

import (
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/dop251/goja"
	"github.com/mitchellh/mapstructure"
	log "github.com/Sirupsen/logrus"
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
