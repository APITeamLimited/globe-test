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

//go:generate rice embed-go

package compiler

import (
	"sync"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

var (
	DefaultOpts = map[string]interface***REMOVED******REMOVED******REMOVED***
		"presets":       []string***REMOVED***"latest"***REMOVED***,
		"ast":           false,
		"sourceMaps":    false,
		"babelrc":       false,
		"compact":       false,
		"retainLines":   true,
		"highlightCode": false,
	***REMOVED***

	once        sync.Once // nolint:gochecknoglobals
	globalBabel *babel    // nolint:gochecknoglobals
)

// CompatibilityMode specifies the JS compatibility mode
// nolint:lll
//go:generate enumer -type=CompatibilityMode -transform=snake -trimprefix CompatibilityMode -output compatibility_mode_gen.go
type CompatibilityMode uint8

const (
	// CompatibilityModeExtended achieves ES6+ compatibility with Babel and core.js
	CompatibilityModeExtended CompatibilityMode = iota + 1
	// CompatibilityModeBase is standard goja ES5.1+
	CompatibilityModeBase
)

// A Compiler compiles JavaScript source code (ES5.1 or ES6) into a goja.Program
type Compiler struct***REMOVED******REMOVED***

// New returns a new Compiler
func New() *Compiler ***REMOVED***
	return &Compiler***REMOVED******REMOVED***
***REMOVED***

// Transform the given code into ES5
func (c *Compiler) Transform(src, filename string) (code string, srcmap *SourceMap, err error) ***REMOVED***
	var b *babel
	if b, err = newBabel(); err != nil ***REMOVED***
		return
	***REMOVED***

	return b.Transform(src, filename)
***REMOVED***

// Compile the program in the given CompatibilityMode, optionally running pre and post code.
func (c *Compiler) Compile(src, filename, pre, post string,
	strict bool, compatMode CompatibilityMode) (*goja.Program, string, error) ***REMOVED***
	code := pre + src + post
	ast, err := parser.ParseFile(nil, filename, code, 0)
	if err != nil ***REMOVED***
		if compatMode == CompatibilityModeExtended ***REMOVED***
			code, _, err = c.Transform(src, filename)
			if err != nil ***REMOVED***
				return nil, code, err
			***REMOVED***
			return c.Compile(code, filename, pre, post, strict, compatMode)
		***REMOVED***
		return nil, code, err
	***REMOVED***
	pgm, err := goja.CompileAST(ast, strict)
	return pgm, code, err
***REMOVED***

type babel struct ***REMOVED***
	vm        *goja.Runtime
	this      goja.Value
	transform goja.Callable
	mutex     sync.Mutex //TODO: cache goja.CompileAST() in an init() function?
***REMOVED***

func newBabel() (*babel, error) ***REMOVED***
	var err error

	once.Do(func() ***REMOVED***
		conf := rice.Config***REMOVED***
			LocateOrder: []rice.LocateMethod***REMOVED***rice.LocateEmbedded***REMOVED***,
		***REMOVED***
		babelSrc := conf.MustFindBox("lib").MustString("babel.min.js")
		vm := goja.New()
		if _, err = vm.RunString(babelSrc); err != nil ***REMOVED***
			return
		***REMOVED***

		this := vm.Get("Babel")
		bObj := this.ToObject(vm)
		globalBabel = &babel***REMOVED***vm: vm, this: this***REMOVED***
		if err = vm.ExportTo(bObj.Get("transform"), &globalBabel.transform); err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***)

	return globalBabel, err
***REMOVED***

// Transform the given code into ES5, while synchronizing to ensure only a single
// bundle instance / Goja VM is in use at a time.
func (b *babel) Transform(src, filename string) (string, *SourceMap, error) ***REMOVED***
	b.mutex.Lock()
	defer b.mutex.Unlock()
	opts := make(map[string]interface***REMOVED******REMOVED***)
	for k, v := range DefaultOpts ***REMOVED***
		opts[k] = v
	***REMOVED***
	opts["filename"] = filename

	startTime := time.Now()
	v, err := b.transform(b.this, b.vm.ToValue(src), b.vm.ToValue(opts))
	if err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***
	logrus.WithField("t", time.Since(startTime)).Debug("Babel: Transformed")

	vO := v.ToObject(b.vm)
	var code string
	if err = b.vm.ExportTo(vO.Get("code"), &code); err != nil ***REMOVED***
		return code, nil, err
	***REMOVED***
	var rawMap map[string]interface***REMOVED******REMOVED***
	if err = b.vm.ExportTo(vO.Get("map"), &rawMap); err != nil ***REMOVED***
		return code, nil, err
	***REMOVED***
	var srcMap SourceMap
	if err = mapstructure.Decode(rawMap, &srcMap); err != nil ***REMOVED***
		return code, &srcMap, err
	***REMOVED***
	return code, &srcMap, err
***REMOVED***
