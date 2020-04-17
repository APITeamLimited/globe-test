/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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

package js

import (
	"context"
	"encoding/json"
	"net/url"
	"runtime"

	"github.com/loadimpact/k6/lib/consts"

	"github.com/dop251/goja"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/js/compiler"
	jslib "github.com/loadimpact/k6/js/lib"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/loader"
)

// A Bundle is a self-contained bundle of scripts and resources.
// You can use this to produce identical BundleInstance objects.
type Bundle struct ***REMOVED***
	Filename *url.URL
	Source   string
	Program  *goja.Program
	// exported functions, for validation only
	Exports map[string]struct***REMOVED******REMOVED***
	Options lib.Options

	BaseInitContext *InitContext

	Env               map[string]string
	CompatibilityMode lib.CompatibilityMode
***REMOVED***

// A BundleInstance is a self-contained instance of a Bundle.
type BundleInstance struct ***REMOVED***
	Runtime *goja.Runtime
	Context *context.Context
	// exported functions, ready for execution
	Exports map[string]goja.Callable
***REMOVED***

// NewBundle creates a new bundle from a source file and a filesystem.
func NewBundle(src *loader.SourceData, filesystems map[string]afero.Fs, rtOpts lib.RuntimeOptions) (*Bundle, error) ***REMOVED***
	compatMode, err := lib.ValidateCompatibilityMode(rtOpts.CompatibilityMode.String)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Compile sources, both ES5 and ES6 are supported.
	code := string(src.Data)
	c := compiler.New()
	pgm, _, err := c.Compile(code, src.URL.String(), "", "", true, compatMode)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// Make a bundle, instantiate it into a throwaway VM to populate caches.
	rt := goja.New()
	bundle := Bundle***REMOVED***
		Filename: src.URL,
		Source:   code,
		Program:  pgm,
		Exports:  make(map[string]struct***REMOVED******REMOVED***),
		BaseInitContext: NewInitContext(rt, c, compatMode, new(context.Context),
			filesystems, loader.Dir(src.URL)),
		Env:               rtOpts.Env,
		CompatibilityMode: compatMode,
	***REMOVED***
	if err := bundle.instantiate(rt, bundle.BaseInitContext); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = bundle.getExports(rt)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &bundle, nil
***REMOVED***

// NewBundleFromArchive creates a new bundle from an lib.Archive.
func NewBundleFromArchive(arc *lib.Archive, rtOpts lib.RuntimeOptions) (*Bundle, error) ***REMOVED***
	if arc.Type != "js" ***REMOVED***
		return nil, errors.Errorf("expected bundle type 'js', got '%s'", arc.Type)
	***REMOVED***

	compatModeStr := arc.CompatibilityMode
	if rtOpts.CompatibilityMode.Valid ***REMOVED***
		// `k6 run --compatibility-mode=whatever archive.tar` should  override
		// whatever value is in the archive
		compatModeStr = rtOpts.CompatibilityMode.String
	***REMOVED***
	compatMode, err := lib.ValidateCompatibilityMode(compatModeStr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c := compiler.New()
	pgm, _, err := c.Compile(string(arc.Data), arc.FilenameURL.String(), "", "", true, compatMode)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	rt := goja.New()
	initctx := NewInitContext(rt, c, compatMode,
		new(context.Context), arc.Filesystems, arc.PwdURL)

	env := arc.Env
	if env == nil ***REMOVED***
		// Older archives (<=0.20.0) don't have an "env" property
		env = make(map[string]string)
	***REMOVED***
	for k, v := range rtOpts.Env ***REMOVED***
		env[k] = v
	***REMOVED***

	bundle := &Bundle***REMOVED***
		Filename:          arc.FilenameURL,
		Source:            string(arc.Data),
		Program:           pgm,
		Exports:           make(map[string]struct***REMOVED******REMOVED***),
		Options:           arc.Options,
		BaseInitContext:   initctx,
		Env:               env,
		CompatibilityMode: compatMode,
	***REMOVED***

	if err = bundle.instantiate(rt, bundle.BaseInitContext); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = bundle.getExports(rt)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return bundle, nil
***REMOVED***

func (b *Bundle) makeArchive() *lib.Archive ***REMOVED***
	arc := &lib.Archive***REMOVED***
		Type:              "js",
		Filesystems:       b.BaseInitContext.filesystems,
		Options:           b.Options,
		FilenameURL:       b.Filename,
		Data:              []byte(b.Source),
		PwdURL:            b.BaseInitContext.pwd,
		Env:               make(map[string]string, len(b.Env)),
		CompatibilityMode: b.CompatibilityMode.String(),
		K6Version:         consts.Version,
		Goos:              runtime.GOOS,
	***REMOVED***
	// Copy env so changes in the archive are not reflected in the source Bundle
	for k, v := range b.Env ***REMOVED***
		arc.Env[k] = v
	***REMOVED***

	return arc
***REMOVED***

// getExports validates and extracts exported objects
func (b *Bundle) getExports(rt *goja.Runtime) error ***REMOVED***
	exportsV := rt.Get("exports")
	if goja.IsNull(exportsV) || goja.IsUndefined(exportsV) ***REMOVED***
		return errors.New("exports must be an object")
	***REMOVED***
	exports := exportsV.ToObject(rt)

	for _, k := range exports.Keys() ***REMOVED***
		v := exports.Get(k)
		switch k ***REMOVED***
		case "options":
			data, err := json.Marshal(v.Export())
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := json.Unmarshal(data, &b.Options); err != nil ***REMOVED***
				return err
			***REMOVED***
		case "setup":
			if _, ok := goja.AssertFunction(v); !ok ***REMOVED***
				return errors.New("exported 'setup' must be a function")
			***REMOVED***
		case "teardown":
			if _, ok := goja.AssertFunction(v); !ok ***REMOVED***
				return errors.New("exported 'teardown' must be a function")
			***REMOVED***
		default:
			if _, ok := goja.AssertFunction(v); ok ***REMOVED***
				b.Exports[k] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if len(b.Exports) == 0 ***REMOVED***
		return errors.New("no exported functions in script")
	***REMOVED***

	return nil
***REMOVED***

// Instantiate creates a new runtime from this bundle.
func (b *Bundle) Instantiate() (bi *BundleInstance, instErr error) ***REMOVED***
	// TODO: actually use a real context here, so that the instantiation can be killed
	// Placeholder for a real context.
	ctxPtr := new(context.Context)

	// Instantiate the bundle into a new VM using a bound init context. This uses a context with a
	// runtime, but no state, to allow module-provided types to function within the init context.
	rt := goja.New()
	init := newBoundInitContext(b.BaseInitContext, ctxPtr, rt)
	if err := b.instantiate(rt, init); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	bi = &BundleInstance***REMOVED***
		Runtime: rt,
		Context: ctxPtr,
		Exports: make(map[string]goja.Callable),
	***REMOVED***

	// Grab any exported functions that could be executed. These were
	// already pre-validated in NewBundle(), just get them here.
	exports := rt.Get("exports").ToObject(rt)
	for k := range b.Exports ***REMOVED***
		fn, _ := goja.AssertFunction(exports.Get(k))
		bi.Exports[k] = fn
	***REMOVED***

	jsOptions := rt.Get("options")
	var jsOptionsObj *goja.Object
	if jsOptions == nil || goja.IsNull(jsOptions) || goja.IsUndefined(jsOptions) ***REMOVED***
		jsOptionsObj = rt.NewObject()
		rt.Set("options", jsOptionsObj)
	***REMOVED*** else ***REMOVED***
		jsOptionsObj = jsOptions.ToObject(rt)
	***REMOVED***
	b.Options.ForEachSpecified("json", func(key string, val interface***REMOVED******REMOVED***) ***REMOVED***
		if err := jsOptionsObj.Set(key, val); err != nil ***REMOVED***
			instErr = err
		***REMOVED***
	***REMOVED***)

	return bi, instErr
***REMOVED***

// Instantiates the bundle into an existing runtime. Not public because it also messes with a bunch
// of other things, will potentially thrash data and makes a mess in it if the operation fails.
func (b *Bundle) instantiate(rt *goja.Runtime, init *InitContext) error ***REMOVED***
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	rt.SetRandSource(common.NewRandSource())

	if init.compatibilityMode == lib.CompatibilityModeExtended ***REMOVED***
		if _, err := rt.RunProgram(jslib.GetCoreJS()); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	exports := rt.NewObject()
	rt.Set("exports", exports)
	module := rt.NewObject()
	_ = module.Set("exports", exports)
	rt.Set("module", module)

	env := make(map[string]string, len(b.Env))
	for key, value := range b.Env ***REMOVED***
		env[key] = value
	***REMOVED***
	rt.Set("__ENV", env)
	rt.Set("console", common.Bind(rt, newConsole(), init.ctxPtr))

	*init.ctxPtr = common.WithRuntime(context.Background(), rt)
	unbindInit := common.BindToGlobal(rt, common.Bind(rt, init, init.ctxPtr))
	if _, err := rt.RunProgram(b.Program); err != nil ***REMOVED***
		return err
	***REMOVED***
	unbindInit()
	*init.ctxPtr = nil

	rt.SetRandSource(common.NewRandSource())

	return nil
***REMOVED***
