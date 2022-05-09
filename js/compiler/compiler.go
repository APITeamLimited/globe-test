/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2017 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package compiler

import (
	_ "embed" // we need this for embedding Babel
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"github.com/go-sourcemap/sourcemap"
	"github.com/sirupsen/logrus"

	"go.k6.io/k6/lib"
)

//go:embed lib/babel.min.js
var babelSrc string //nolint:gochecknoglobals

var (
	DefaultOpts = map[string]interface***REMOVED******REMOVED******REMOVED***
		// "presets": []string***REMOVED***"latest"***REMOVED***,
		"plugins": []interface***REMOVED******REMOVED******REMOVED***
			// es2015 https://github.com/babel/babel/blob/v6.26.0/packages/babel-preset-es2015/src/index.js
			// in goja
			// []interface***REMOVED******REMOVED******REMOVED***"transform-es2015-template-literals", map[string]interface***REMOVED******REMOVED******REMOVED***"loose": false, "spec": false***REMOVED******REMOVED***,
			// "transform-es2015-literals", // in goja
			// "transform-es2015-function-name", // in goja
			// []interface***REMOVED******REMOVED******REMOVED***"transform-es2015-arrow-functions", map[string]interface***REMOVED******REMOVED******REMOVED***"spec": false***REMOVED******REMOVED***, // in goja
			// "transform-es2015-block-scoped-functions", // in goja
			[]interface***REMOVED******REMOVED******REMOVED***"transform-es2015-classes", map[string]interface***REMOVED******REMOVED******REMOVED***"loose": false***REMOVED******REMOVED***,
			"transform-es2015-object-super",
			// "transform-es2015-shorthand-properties", // in goja
			// "transform-es2015-duplicate-keys", // in goja
			// []interface***REMOVED******REMOVED******REMOVED***"transform-es2015-computed-properties", map[string]interface***REMOVED******REMOVED******REMOVED***"loose": false***REMOVED******REMOVED***, // in goja
			// "transform-es2015-for-of", // in goja
			// "transform-es2015-sticky-regex", // in goja
			// "transform-es2015-unicode-regex", // in goja
			// "check-es2015-constants", // in goja
			// []interface***REMOVED******REMOVED******REMOVED***"transform-es2015-spread", map[string]interface***REMOVED******REMOVED******REMOVED***"loose": false***REMOVED******REMOVED***, // in goja
			// "transform-es2015-parameters", // in goja
			// []interface***REMOVED******REMOVED******REMOVED***"transform-es2015-destructuring", map[string]interface***REMOVED******REMOVED******REMOVED***"loose": false***REMOVED******REMOVED***, // in goja
			// "transform-es2015-block-scoping", // in goja
			// "transform-es2015-typeof-symbol", // in goja
			// all the other module plugins are just dropped
			[]interface***REMOVED******REMOVED******REMOVED***"transform-es2015-modules-commonjs", map[string]interface***REMOVED******REMOVED******REMOVED***"loose": false***REMOVED******REMOVED***,
			// "transform-regenerator", // Doesn't really work unless regeneratorRuntime is also added

			// es2016 https://github.com/babel/babel/blob/v6.26.0/packages/babel-preset-es2016/src/index.js
			// "transform-exponentiation-operator",

			// es2017 https://github.com/babel/babel/blob/v6.26.0/packages/babel-preset-es2017/src/index.js
			// "syntax-trailing-function-commas", // in goja
			// "transform-async-to-generator", // Doesn't really work unless regeneratorRuntime is also added
		***REMOVED***,
		"ast":           false,
		"sourceMaps":    false,
		"babelrc":       false,
		"compact":       false,
		"retainLines":   true,
		"highlightCode": false,
	***REMOVED***

	maxSrcLenForBabelSourceMap     = 250 * 1024 //nolint:gochecknoglobals
	maxSrcLenForBabelSourceMapOnce sync.Once    //nolint:gochecknoglobals

	onceBabelCode      sync.Once     // nolint:gochecknoglobals
	globalBabelCode    *goja.Program // nolint:gochecknoglobals
	globalBabelCodeErr error         // nolint:gochecknoglobals
	onceBabel          sync.Once     // nolint:gochecknoglobals
	globalBabel        *babel        // nolint:gochecknoglobals
)

const (
	maxSrcLenForBabelSourceMapVarName = "K6_DEBUG_SOURCEMAP_FILESIZE_LIMIT"
	sourceMapURLFromBabel             = "k6://internal-should-not-leak/file.map"
)

// A Compiler compiles JavaScript source code (ES5.1 or ES6) into a goja.Program
type Compiler struct ***REMOVED***
	logger  logrus.FieldLogger
	babel   *babel
	Options Options
***REMOVED***

// New returns a new Compiler
func New(logger logrus.FieldLogger) *Compiler ***REMOVED***
	return &Compiler***REMOVED***logger: logger***REMOVED***
***REMOVED***

// initializeBabel initializes a separate (non-global) instance of babel specifically for this Compiler.
// An error is returned only if babel itself couldn't be parsed/run which should never be possible.
func (c *Compiler) initializeBabel() error ***REMOVED***
	var err error
	if c.babel == nil ***REMOVED***
		c.babel, err = newBabel()
	***REMOVED***
	return err
***REMOVED***

// Transform the given code into ES5
func (c *Compiler) Transform(src, filename string, inputSrcMap []byte) (code string, srcMap []byte, err error) ***REMOVED***
	if c.babel == nil ***REMOVED***
		onceBabel.Do(func() ***REMOVED***
			globalBabel, err = newBabel()
		***REMOVED***)
		c.babel = globalBabel
	***REMOVED***
	if err != nil ***REMOVED***
		return
	***REMOVED***

	sourceMapEnabled := c.Options.SourceMapLoader != nil
	maxSrcLenForBabelSourceMapOnce.Do(func() ***REMOVED***
		// TODO: drop this code and everything it's connected to when babel is dropped
		v := os.Getenv(maxSrcLenForBabelSourceMapVarName)
		if len(v) > 0 ***REMOVED***
			i, err := strconv.Atoi(v) //nolint:govet // we shadow err on purpose
			if err != nil ***REMOVED***
				c.logger.Warnf("Tried to parse %q from %s as integer but couldn't %s\n",
					v, maxSrcLenForBabelSourceMapVarName, err)
				return
			***REMOVED***
			maxSrcLenForBabelSourceMap = i
		***REMOVED***
	***REMOVED***)
	if sourceMapEnabled && len(src) > maxSrcLenForBabelSourceMap ***REMOVED***
		sourceMapEnabled = false
		c.logger.Warnf("The source for `%s` needs to go through babel but is over %d bytes. "+
			"For performance reasons source map support will be disabled for this particular file.",
			filename, maxSrcLenForBabelSourceMap)
	***REMOVED***

	// check that babel will likely be able to parse the inputSrcMap
	if sourceMapEnabled && len(inputSrcMap) != 0 ***REMOVED***
		if err = verifySourceMapForBabel(inputSrcMap); err != nil ***REMOVED***
			sourceMapEnabled = false
			inputSrcMap = nil
			c.logger.WithError(err).Warnf(
				"The source for `%s` needs to be transpiled by Babel, but its source map will"+
					" not be accepted by Babel, so it was disabled", filename)
		***REMOVED***
	***REMOVED***
	code, srcMap, err = c.babel.transformImpl(c.logger, src, filename, sourceMapEnabled, inputSrcMap)
	return
***REMOVED***

// Options are options to the compiler
type Options struct ***REMOVED***
	CompatibilityMode lib.CompatibilityMode
	SourceMapLoader   func(string) ([]byte, error)
	Strict            bool
***REMOVED***

// compilationState is helper struct to keep the state of a compilation
type compilationState struct ***REMOVED***
	// set when we couldn't load external source map so we can try parsing without loading it
	couldntLoadSourceMap bool
	// srcMap is the current full sourceMap that has been generated read so far
	srcMap      []byte
	srcMapError error
	main        bool

	compiler *Compiler
***REMOVED***

// Compile the program in the given CompatibilityMode, wrapping it between pre and post code
func (c *Compiler) Compile(src, filename string, main bool) (*goja.Program, string, error) ***REMOVED***
	return c.compileImpl(src, filename, main, c.Options.CompatibilityMode, nil)
***REMOVED***

// sourceMapLoader is to be used with goja's WithSourceMapLoader
// it not only gets the file from disk in the simple case, but also returns it if the map was generated from babel
// additioanlly it fixes off by one error in commonjs dependencies due to having to wrap them in a function.
func (c *compilationState) sourceMapLoader(path string) ([]byte, error) ***REMOVED***
	if path == sourceMapURLFromBabel ***REMOVED***
		if !c.main ***REMOVED***
			return c.increaseMappingsByOne(c.srcMap)
		***REMOVED***
		return c.srcMap, nil
	***REMOVED***
	c.srcMap, c.srcMapError = c.compiler.Options.SourceMapLoader(path)
	if c.srcMapError != nil ***REMOVED***
		c.couldntLoadSourceMap = true
		return nil, c.srcMapError
	***REMOVED***
	_, c.srcMapError = sourcemap.Parse(path, c.srcMap)
	if c.srcMapError != nil ***REMOVED***
		c.couldntLoadSourceMap = true
		c.srcMap = nil
		return nil, c.srcMapError
	***REMOVED***
	if !c.main ***REMOVED***
		return c.increaseMappingsByOne(c.srcMap)
	***REMOVED***
	return c.srcMap, nil
***REMOVED***

func (c *Compiler) compileImpl(
	src, filename string, main bool, compatibilityMode lib.CompatibilityMode, srcMap []byte,
) (*goja.Program, string, error) ***REMOVED***
	code := src
	state := compilationState***REMOVED***srcMap: srcMap, compiler: c, main: main***REMOVED***
	if !main ***REMOVED*** // the lines in the sourcemap (if available) will be fixed by increaseMappingsByOne
		code = "(function(module, exports)***REMOVED***\n" + code + "\n***REMOVED***)\n"
	***REMOVED***
	opts := parser.WithDisableSourceMaps
	if c.Options.SourceMapLoader != nil ***REMOVED***
		opts = parser.WithSourceMapLoader(state.sourceMapLoader)
	***REMOVED***
	ast, err := parser.ParseFile(nil, filename, code, 0, opts)

	if state.couldntLoadSourceMap ***REMOVED***
		state.couldntLoadSourceMap = false // reset
		// we probably don't want to abort scripts which have source maps but they can't be found,
		// this also will be a breaking change, so if we couldn't we retry with it disabled
		c.logger.WithError(state.srcMapError).Warnf("Couldn't load source map for %s", filename)
		ast, err = parser.ParseFile(nil, filename, code, 0, parser.WithDisableSourceMaps)
	***REMOVED***
	if err != nil ***REMOVED***
		if compatibilityMode == lib.CompatibilityModeExtended ***REMOVED***
			code, state.srcMap, err = c.Transform(src, filename, state.srcMap)
			if err != nil ***REMOVED***
				return nil, code, err
			***REMOVED***
			// the compatibility mode "decreases" here as we shouldn't transform twice
			return c.compileImpl(code, filename, main, lib.CompatibilityModeBase, state.srcMap)
		***REMOVED***
		return nil, code, err
	***REMOVED***
	pgm, err := goja.CompileAST(ast, c.Options.Strict)
	return pgm, code, err
***REMOVED***

type babel struct ***REMOVED***
	vm        *goja.Runtime
	this      goja.Value
	transform goja.Callable
	m         sync.Mutex
***REMOVED***

func newBabel() (*babel, error) ***REMOVED***
	onceBabelCode.Do(func() ***REMOVED***
		globalBabelCode, globalBabelCodeErr = goja.Compile("<internal/k6/compiler/lib/babel.min.js>", babelSrc, false)
	***REMOVED***)
	if globalBabelCodeErr != nil ***REMOVED***
		return nil, globalBabelCodeErr
	***REMOVED***
	vm := goja.New()
	_, err := vm.RunProgram(globalBabelCode)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	this := vm.Get("Babel")
	bObj := this.ToObject(vm)
	result := &babel***REMOVED***vm: vm, this: this***REMOVED***
	if err = vm.ExportTo(bObj.Get("transform"), &result.transform); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return result, err
***REMOVED***

// increaseMappingsByOne increases the lines in the sourcemap by line so that it fixes the case where we need to wrap a
// required file in a function to support/emulate commonjs
func (c *compilationState) increaseMappingsByOne(sourceMap []byte) ([]byte, error) ***REMOVED***
	var err error
	m := make(map[string]interface***REMOVED******REMOVED***)
	if err = json.Unmarshal(sourceMap, &m); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	mappings, ok := m["mappings"]
	if !ok ***REMOVED***
		// no mappings, no idea what this will do, but just return it as technically we can have sourcemap with sections
		// TODO implement incrementing of `offset` in the sections? to support that case as well
		// see https://sourcemaps.info/spec.html#h.n05z8dfyl3yh
		//
		// TODO (kind of alternatively) drop the newline in the "commonjs" wrapping and have only the first line wrong
		// and drop this whole function
		return sourceMap, nil
	***REMOVED***
	if str, ok := mappings.(string); ok ***REMOVED***
		// ';' is the separator between lines so just adding 1 will make all mappings be for the line after which they were
		// originally
		m["mappings"] = ";" + str
	***REMOVED*** else ***REMOVED***
		// we have mappings but it's not a string - this is some kind of error
		// we still won't abort the test but just not load the sourcemap
		c.couldntLoadSourceMap = true
		return nil, errors.New(`missing "mappings" in sourcemap`)
	***REMOVED***

	return json.Marshal(m)
***REMOVED***

// transformImpl the given code into ES5, while synchronizing to ensure only a single
// bundle instance / Goja VM is in use at a time.
func (b *babel) transformImpl(
	logger logrus.FieldLogger, src, filename string, sourceMapsEnabled bool, inputSrcMap []byte,
) (string, []byte, error) ***REMOVED***
	b.m.Lock()
	defer b.m.Unlock()
	opts := make(map[string]interface***REMOVED******REMOVED***)
	for k, v := range DefaultOpts ***REMOVED***
		opts[k] = v
	***REMOVED***
	if sourceMapsEnabled ***REMOVED***
		// given that the source map should provide accurate lines(and columns), this option isn't needed
		// it also happens to make very long and awkward lines, especially around import/exports and definitely a lot
		// less readable overall. Hopefully it also has some performance improvement not trying to keep the same lines
		opts["retainLines"] = false
		opts["sourceMaps"] = true
		if inputSrcMap != nil ***REMOVED***
			srcMap := new(map[string]interface***REMOVED******REMOVED***)
			if err := json.Unmarshal(inputSrcMap, &srcMap); err != nil ***REMOVED***
				return "", nil, err
			***REMOVED***
			opts["inputSourceMap"] = srcMap
		***REMOVED***
	***REMOVED***
	opts["filename"] = filename

	startTime := time.Now()
	v, err := b.transform(b.this, b.vm.ToValue(src), b.vm.ToValue(opts))
	if err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***
	logger.WithField("t", time.Since(startTime)).Debug("Babel: Transformed")

	vO := v.ToObject(b.vm)
	var code string
	if err = b.vm.ExportTo(vO.Get("code"), &code); err != nil ***REMOVED***
		return code, nil, err
	***REMOVED***
	if !sourceMapsEnabled ***REMOVED***
		return code, nil, nil
	***REMOVED***

	// this is to make goja try to load a sourcemap.
	// it is a special url as it should never leak outside of this code
	// additionally the alternative support from babel is to embed *the whole* sourcemap at the end
	code += "\n//# sourceMappingURL=" + sourceMapURLFromBabel
	stringify, err := b.vm.RunString("(function(m) ***REMOVED*** return JSON.stringify(m)***REMOVED***)")
	if err != nil ***REMOVED***
		return code, nil, err
	***REMOVED***
	c, _ := goja.AssertFunction(stringify)
	mapAsJSON, err := c(goja.Undefined(), vO.Get("map"))
	if err != nil ***REMOVED***
		return code, nil, err
	***REMOVED***
	return code, []byte(mapAsJSON.String()), nil
***REMOVED***

// Pool is a pool of compilers so it can be used easier in parallel tests as they have their own babel.
type Pool struct ***REMOVED***
	c chan *Compiler
***REMOVED***

// NewPool creates a Pool that will be using the provided logger and will preallocate (in parallel)
// the count of compilers each with their own babel.
func NewPool(logger logrus.FieldLogger, count int) *Pool ***REMOVED***
	c := &Pool***REMOVED***
		c: make(chan *Compiler, count),
	***REMOVED***
	go func() ***REMOVED***
		for i := 0; i < count; i++ ***REMOVED***
			go func() ***REMOVED***
				co := New(logger)
				err := co.initializeBabel()
				if err != nil ***REMOVED***
					panic(err)
				***REMOVED***
				c.Put(co)
			***REMOVED***()
		***REMOVED***
	***REMOVED***()

	return c
***REMOVED***

// Get a compiler from the pool.
func (c *Pool) Get() *Compiler ***REMOVED***
	return <-c.c
***REMOVED***

// Put a compiler back in the pool.
func (c *Pool) Put(co *Compiler) ***REMOVED***
	c.c <- co
***REMOVED***

func verifySourceMapForBabel(srcMap []byte) error ***REMOVED***
	// this function exists to do what babel checks in sourcemap before we give it to it.
	m := make(map[string]json.RawMessage)
	err := json.Unmarshal(srcMap, &m)
	if err != nil ***REMOVED***
		return fmt.Errorf("source map is not valid json: %w", err)
	***REMOVED***
	// there are no checks on it's value in babel
	// we technically only support v3 though
	if _, ok := m["version"]; !ok ***REMOVED***
		return fmt.Errorf("source map missing required 'version' field")
	***REMOVED***

	// This actually gets checked by the go implementation so it's not really necessary
	if _, ok := m["mappings"]; !ok ***REMOVED***
		return fmt.Errorf("source map missing required 'mappings' field")
	***REMOVED***
	// the go implementation checks the value even if it doesn't require it exists
	if _, ok := m["sources"]; !ok ***REMOVED***
		return fmt.Errorf("source map missing required 'sources' field")
	***REMOVED***
	return nil
***REMOVED***
