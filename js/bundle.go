package js

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"runtime"

	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/compiler"
	"go.k6.io/k6/js/eventloop"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/consts"
	"go.k6.io/k6/loader"
	"go.k6.io/k6/metrics"
)

// A Bundle is a self-contained bundle of scripts and resources.
// You can use this to produce identical BundleInstance objects.
type Bundle struct ***REMOVED***
	Filename *url.URL
	Source   string
	Program  *goja.Program
	Options  lib.Options

	BaseInitContext *InitContext

	RuntimeOptions    lib.RuntimeOptions
	CompatibilityMode lib.CompatibilityMode // parsed value
	registry          *metrics.Registry

	exports map[string]goja.Callable
***REMOVED***

// A BundleInstance is a self-contained instance of a Bundle.
type BundleInstance struct ***REMOVED***
	Runtime *goja.Runtime

	// TODO: maybe just have a reference to the Bundle? or save and pass rtOpts?
	env map[string]string

	exports      map[string]goja.Callable
	moduleVUImpl *moduleVUImpl
	pgm          programWithSource
***REMOVED***

// NewBundle creates a new bundle from a source file and a filesystem.
func NewBundle(
	piState *lib.TestPreInitState, src *loader.SourceData, filesystems map[string]afero.Fs, workerInfo *lib.WorkerInfo,
) (*Bundle, error) ***REMOVED***
	compatMode, err := lib.ValidateCompatibilityMode(piState.RuntimeOptions.CompatibilityMode.String)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Compile sources, both ES5 and ES6 are supported.
	code := string(src.Data)
	c := compiler.New(piState.Logger)
	c.Options = compiler.Options***REMOVED***
		CompatibilityMode: compatMode,
		Strict:            true,
		SourceMapLoader:   generateSourceMapLoader(piState.Logger, filesystems),
	***REMOVED***
	pgm, _, err := c.Compile(code, src.URL.String(), false)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// Make a bundle, instantiate it into a throwaway VM to populate caches.
	rt := goja.New()
	bundle := Bundle***REMOVED***
		Filename:          src.URL,
		Source:            code,
		Program:           pgm,
		BaseInitContext:   NewInitContext(piState.Logger, rt, c, compatMode, filesystems, loader.Dir(src.URL)),
		RuntimeOptions:    piState.RuntimeOptions,
		CompatibilityMode: compatMode,
		exports:           make(map[string]goja.Callable),
		registry:          piState.Registry,
	***REMOVED***
	if err = bundle.instantiate(piState.Logger, rt, bundle.BaseInitContext, 0, workerInfo); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = bundle.getExports(piState.Logger, rt, true)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &bundle, nil
***REMOVED***

// NewBundleFromArchive creates a new bundle from an lib.Archive.
func NewBundleFromArchive(piState *lib.TestPreInitState, arc *lib.Archive, workerInfo *lib.WorkerInfo) (*Bundle, error) ***REMOVED***
	if arc.Type != "js" ***REMOVED***
		return nil, fmt.Errorf("expected bundle type 'js', got '%s'", arc.Type)
	***REMOVED***

	rtOpts := piState.RuntimeOptions // copy the struct from the TestPreInitState
	if !rtOpts.CompatibilityMode.Valid ***REMOVED***
		// `k6 run --compatibility-mode=whatever archive.tar` should override
		// whatever value is in the archive
		rtOpts.CompatibilityMode = null.StringFrom(arc.CompatibilityMode)
	***REMOVED***
	compatMode, err := lib.ValidateCompatibilityMode(rtOpts.CompatibilityMode.String)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c := compiler.New(piState.Logger)
	c.Options = compiler.Options***REMOVED***
		Strict:            true,
		CompatibilityMode: compatMode,
		SourceMapLoader:   generateSourceMapLoader(piState.Logger, arc.Filesystems),
	***REMOVED***
	pgm, _, err := c.Compile(string(arc.Data), arc.FilenameURL.String(), false)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	rt := goja.New()
	initctx := NewInitContext(piState.Logger, rt, c, compatMode, arc.Filesystems, arc.PwdURL)

	env := arc.Env
	if env == nil ***REMOVED***
		// Older archives (<=0.20.0) don't have an "env" property
		env = make(map[string]string)
	***REMOVED***
	for k, v := range rtOpts.Env ***REMOVED***
		env[k] = v
	***REMOVED***
	rtOpts.Env = env

	bundle := &Bundle***REMOVED***
		Filename:          arc.FilenameURL,
		Source:            string(arc.Data),
		Program:           pgm,
		Options:           arc.Options,
		BaseInitContext:   initctx,
		RuntimeOptions:    rtOpts,
		CompatibilityMode: compatMode,
		exports:           make(map[string]goja.Callable),
		registry:          piState.Registry,
	***REMOVED***

	if err = bundle.instantiate(piState.Logger, rt, bundle.BaseInitContext, 0, workerInfo); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Grab exported objects, but avoid overwriting options, which would
	// be initialized from the metadata.json at this point.
	err = bundle.getExports(piState.Logger, rt, false)
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
		Env:               make(map[string]string, len(b.RuntimeOptions.Env)),
		CompatibilityMode: b.CompatibilityMode.String(),
		K6Version:         consts.Version,
		Goos:              runtime.GOOS,
	***REMOVED***
	// Copy env so changes in the archive are not reflected in the source Bundle
	for k, v := range b.RuntimeOptions.Env ***REMOVED***
		arc.Env[k] = v
	***REMOVED***

	return arc
***REMOVED***

// getExports validates and extracts exported objects
func (b *Bundle) getExports(logger logrus.FieldLogger, rt *goja.Runtime, options bool) error ***REMOVED***
	pgm := b.BaseInitContext.programs[b.Filename.String()] // this is the main script and it's always present
	exportsV := pgm.module.Get("exports")
	if goja.IsNull(exportsV) || goja.IsUndefined(exportsV) ***REMOVED***
		return errors.New("exports must be an object")
	***REMOVED***
	exports := exportsV.ToObject(rt)

	for _, k := range exports.Keys() ***REMOVED***
		v := exports.Get(k)
		if fn, ok := goja.AssertFunction(v); ok && k != consts.Options ***REMOVED***
			b.exports[k] = fn
			continue
		***REMOVED***
		switch k ***REMOVED***
		case consts.Options:
			if !options ***REMOVED***
				continue
			***REMOVED***
			data, err := json.Marshal(v.Export())
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			dec := json.NewDecoder(bytes.NewReader(data))
			dec.DisallowUnknownFields()
			if err := dec.Decode(&b.Options); err != nil ***REMOVED***
				if uerr := json.Unmarshal(data, &b.Options); uerr != nil ***REMOVED***
					return uerr
				***REMOVED***
				logger.WithError(err).Warn("There were unknown fields in the options exported in the script")
			***REMOVED***
		case consts.SetupFn:
			return errors.New("exported 'setup' must be a function")
		case consts.TeardownFn:
			return errors.New("exported 'teardown' must be a function")
		***REMOVED***
	***REMOVED***

	if len(b.exports) == 0 ***REMOVED***
		return errors.New("no exported functions in script")
	***REMOVED***

	return nil
***REMOVED***

// Instantiate creates a new runtime from this bundle.
func (b *Bundle) Instantiate(logger logrus.FieldLogger, vuID uint64, workerInfo *lib.WorkerInfo) (*BundleInstance, error) ***REMOVED***
	// Instantiate the bundle into a new VM using a bound init context. This uses a context with a
	// runtime, but no state, to allow module-provided types to function within the init context.
	vuImpl := &moduleVUImpl***REMOVED***runtime: goja.New()***REMOVED***
	init := newBoundInitContext(b.BaseInitContext, vuImpl)
	if err := b.instantiate(logger, vuImpl.runtime, init, vuID, workerInfo); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	rt := vuImpl.runtime
	pgm := init.programs[b.Filename.String()] // this is the main script and it's always present
	bi := &BundleInstance***REMOVED***
		Runtime:      rt,
		exports:      make(map[string]goja.Callable),
		env:          b.RuntimeOptions.Env,
		moduleVUImpl: vuImpl,
		pgm:          pgm,
	***REMOVED***

	// Grab any exported functions that could be executed. These were
	// already pre-validated in cmd.validateScenarioConfig(), just get them here.
	exports := pgm.module.Get("exports").ToObject(rt)
	for k := range b.exports ***REMOVED***
		fn, _ := goja.AssertFunction(exports.Get(k))
		bi.exports[k] = fn
	***REMOVED***

	jsOptions := exports.Get("options")
	var jsOptionsObj *goja.Object
	if jsOptions == nil || goja.IsNull(jsOptions) || goja.IsUndefined(jsOptions) ***REMOVED***
		jsOptionsObj = rt.NewObject()
		err := exports.Set("options", jsOptionsObj)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("couldn't set exported options with merged values: %w", err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		jsOptionsObj = jsOptions.ToObject(rt)
	***REMOVED***

	var instErr error
	b.Options.ForEachSpecified("json", func(key string, val interface***REMOVED******REMOVED***) ***REMOVED***
		if err := jsOptionsObj.Set(key, val); err != nil ***REMOVED***
			instErr = err
		***REMOVED***
	***REMOVED***)

	return bi, instErr
***REMOVED***

// Instantiates the bundle into an existing runtime. Not public because it also messes with a bunch
// of other things, will potentially thrash data and makes a mess in it if the operation fails.

func (b *Bundle) initializeProgramObject(rt *goja.Runtime, init *InitContext) programWithSource ***REMOVED***
	pgm := programWithSource***REMOVED***
		pgm:     b.Program,
		src:     b.Source,
		exports: rt.NewObject(),
		module:  rt.NewObject(),
	***REMOVED***
	_ = pgm.module.Set("exports", pgm.exports)
	init.programs[b.Filename.String()] = pgm
	return pgm
***REMOVED***

func (b *Bundle) instantiate(logger logrus.FieldLogger, rt *goja.Runtime, init *InitContext, vuID uint64, workerInfo *lib.WorkerInfo) (err error) ***REMOVED***
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	rt.SetRandSource(common.NewRandSource())

	env := make(map[string]string, len(b.RuntimeOptions.Env))
	for key, value := range b.RuntimeOptions.Env ***REMOVED***
		env[key] = value
	***REMOVED***
	rt.Set("__ENV", env)
	rt.Set("__VU", vuID)
	_ = rt.Set("console", newConsole(logger))

	if init.compatibilityMode == lib.CompatibilityModeExtended ***REMOVED***
		rt.Set("global", rt.GlobalObject())
	***REMOVED***

	initenv := &common.InitEnvironment***REMOVED***
		Logger:      logger,
		FileSystems: init.filesystems,
		CWD:         init.pwd,
		Registry:    b.registry,
		WorkerInfo:  workerInfo,
	***REMOVED***
	unbindInit := b.setInitGlobals(rt, init)
	init.moduleVUImpl.ctx = context.Background()
	init.moduleVUImpl.initEnv = initenv
	init.moduleVUImpl.eventLoop = eventloop.New(init.moduleVUImpl)
	pgm := b.initializeProgramObject(rt, init)

	err = common.RunWithPanicCatching(logger, rt, func() error ***REMOVED***
		return init.moduleVUImpl.eventLoop.Start(func() error ***REMOVED***
			f, errRun := rt.RunProgram(b.Program)
			if errRun != nil ***REMOVED***
				return errRun
			***REMOVED***
			if call, ok := goja.AssertFunction(f); ok ***REMOVED***
				if _, errRun = call(pgm.exports, pgm.module, pgm.exports); errRun != nil ***REMOVED***
					return errRun
				***REMOVED***
				return nil
			***REMOVED***
			panic("Somehow a commonjs main module is not wrapped in a function")
		***REMOVED***)
	***REMOVED***)

	if err != nil ***REMOVED***
		var exception *goja.Exception
		if errors.As(err, &exception) ***REMOVED***
			err = &scriptException***REMOVED***inner: exception***REMOVED***
		***REMOVED***
		return err
	***REMOVED***
	exportsV := pgm.module.Get("exports")
	if goja.IsNull(exportsV) ***REMOVED***
		return errors.New("exports must be an object")
	***REMOVED***
	pgm.exports = exportsV.ToObject(rt)
	init.programs[b.Filename.String()] = pgm
	unbindInit()
	init.moduleVUImpl.ctx = nil
	init.moduleVUImpl.initEnv = nil

	// If we've already initialized the original VU init context, forbid
	// any subsequent VUs to open new files
	if vuID == 0 ***REMOVED***
		init.allowOnlyOpenedFiles()
	***REMOVED***

	rt.SetRandSource(common.NewRandSource())

	return nil
***REMOVED***

func (b *Bundle) setInitGlobals(rt *goja.Runtime, init *InitContext) (unset func()) ***REMOVED***
	mustSet := func(k string, v interface***REMOVED******REMOVED***) ***REMOVED***
		if err := rt.Set(k, v); err != nil ***REMOVED***
			panic(fmt.Errorf("failed to set '%s' global object: %w", k, err))
		***REMOVED***
	***REMOVED***
	mustSet("require", init.Require)
	mustSet("open", init.Open)
	return func() ***REMOVED***
		mustSet("require", goja.Undefined())
		mustSet("open", goja.Undefined())
	***REMOVED***
***REMOVED***

func generateSourceMapLoader(logger logrus.FieldLogger, filesystems map[string]afero.Fs,
) func(path string) ([]byte, error) ***REMOVED***
	return func(path string) ([]byte, error) ***REMOVED***
		u, err := url.Parse(path)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		data, err := loader.Load(logger, filesystems, u, path)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return data.Data, nil
	***REMOVED***
***REMOVED***
