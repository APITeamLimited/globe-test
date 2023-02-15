package js

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/APITeamLimited/globe-test/js/common"
	"github.com/APITeamLimited/globe-test/js/compiler"
	"github.com/APITeamLimited/globe-test/js/eventloop"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/consts"
	"github.com/APITeamLimited/globe-test/worker/loader"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
)

// A Bundle is a self-contained bundle of scripts and resources.
// You can use this to produce identical BundleInstance objects.
type Bundle struct {
	RootFilename *url.URL

	Filenames []*url.URL
	Sources   []string
	Programs  []*goja.Program

	Options libWorker.Options

	BaseInitContext *InitContext

	registry *workerMetrics.Registry

	exports map[string]map[string]goja.Callable

	// Test data with intialized goja callables
	initializedTestData libWorker.TestData
}

// A BundleInstance is a self-contained instance of a Bundle.
type BundleInstance struct {
	Runtime *goja.Runtime

	// TODO: maybe just have a reference to the Bundle? or save and pass rtOpts?
	env map[string]string

	exports      map[string]map[string]goja.Callable
	moduleVUImpl *moduleVUImpl
	pgms         map[string]programWithSource

	RootFilename *url.URL
}

// NewBundleWorker creates a new bundle from a source file and a filesystem.
func NewBundleWorker(
	piState *libWorker.TestPreInitState, src *[]*loader.SourceData, filesystems map[string]afero.Fs, workerInfo *libWorker.WorkerInfo, testData *libWorker.TestData) (*Bundle, error) {
	return NewBundle(piState, src, filesystems, workerInfo, false, testData)
}

func NewBundle(piState *libWorker.TestPreInitState, src *[]*loader.SourceData, filesystems map[string]afero.Fs, workerInfo *libWorker.WorkerInfo, isOrchestrator bool, testData *libWorker.TestData) (*Bundle, error) {
	c := compiler.New(piState.Logger)

	c.Options = compiler.Options{
		Strict:          true,
		SourceMapLoader: generateSourceMapLoader(piState.Logger, filesystems),
	}

	compiledPrograms := make(map[string]*goja.Program)

	var compileErr error
	compileLockMutex := &sync.Mutex{}

	resultCountChan := make(chan int, len(*src))

	for _, src := range *src {
		go func(src *loader.SourceData) {
			defer func() {
				resultCountChan <- 1
			}()

			srcCompiler := compiler.New(piState.Logger)

			srcCompiler.Options = compiler.Options{
				Strict:          true,
				SourceMapLoader: generateSourceMapLoader(piState.Logger, filesystems),
			}

			// Compile sources, both ES5 and ES6 are supported.
			pgm, _, err := c.Compile(string(src.Data), src.URL.String(), false)

			compileLockMutex.Lock()
			defer compileLockMutex.Unlock()

			if err != nil {
				if compileErr == nil {
					compileErr = err
				}
				return
			}

			compiledPrograms[src.URL.String()] = pgm
		}(src)
	}

	for i := 0; i < len(*src); i++ {
		<-resultCountChan
	}

	if compileErr != nil {
		return nil, compileErr
	}

	filenames := make([]*url.URL, len(*src))
	sources := make([]string, len(*src))
	programs := make([]*goja.Program, len(*src))

	var rootFilename *url.URL

	exports := make(map[string]map[string]goja.Callable)

	for i, srcProgram := range *src {
		if srcProgram.RootSource {
			rootFilename = srcProgram.URL
		}

		filenames[i] = srcProgram.URL
		sources[i] = string(srcProgram.Data)
		programs[i] = compiledPrograms[srcProgram.URL.String()]

		exports[srcProgram.URL.String()] = make(map[string]goja.Callable)
	}

	// Make a bundle, instantiate it into a throwaway VM to populate caches.
	rt := goja.New()
	bundle := Bundle{
		RootFilename:        rootFilename,
		Filenames:           filenames,
		Sources:             sources,
		Programs:            programs,
		BaseInitContext:     NewInitContext(piState.Logger, rt, c, filesystems, &filenames, workerInfo, testData.RootNode),
		exports:             exports,
		registry:            piState.Registry,
		initializedTestData: *testData,
	}
	if err := bundle.instantiate(piState.Logger, rt, bundle.BaseInitContext, 0, workerInfo, isOrchestrator); err != nil {
		return nil, err
	}

	err := bundle.GetExports(piState.Logger, rt, isOrchestrator)
	if err != nil {
		return nil, err
	}

	return &bundle, nil
}

// GetExports validates and extracts exported objects
func (b *Bundle) GetExports(logger logrus.FieldLogger, rt *goja.Runtime, isOrchestrator bool) error {
	rootFilenameString := b.RootFilename.String()

	// Need to validate the exports of all scripts
	for _, filename := range b.Filenames {
		filenameString := filename.String()

		pgm := b.BaseInitContext.programs[filenameString] // this is the main script and it's always present
		exportsV := pgm.module.Get("exports")
		if goja.IsNull(exportsV) || goja.IsUndefined(exportsV) {
			return errors.New("exports must be an object")
		}
		exports := exportsV.ToObject(rt)

		for _, k := range exports.Keys() {
			v := exports.Get(k)
			if fn, ok := goja.AssertFunction(v); ok && k != consts.Options {
				b.exports[filenameString][k] = fn
				continue
			}
			switch k {
			case consts.Options:
				if !isOrchestrator {
					continue
				}

				// Don't need to validate options if not running the main script
				if filenameString != rootFilenameString {
					continue
				}
				data, err := json.Marshal(v.Export())
				if err != nil {
					return err
				}
				dec := json.NewDecoder(bytes.NewReader(data))
				dec.DisallowUnknownFields()

				// Options are being extracted here
				if err := dec.Decode(&b.Options); err != nil {
					if uerr := json.Unmarshal(data, &b.Options); uerr != nil {
						return uerr
					}
					logger.WithError(err).Warn("There were unknown fields in the options exported in the script")
				}
			case consts.SetupFn:
				return errors.New("exported 'setup' must be a function")
			case consts.TeardownFn:
				return errors.New("exported 'teardown' must be a function")
			}
		}

		if len(b.exports) == 0 {
			return errors.New("no exported functions in script")
		}

		if !isOrchestrator {
			innerNode, err := libWorker.GetInnerNode(b.initializedTestData.RootNode, filenameString)
			if err != nil {
				return err
			}

			innerNode.RegisterExports(filenameString, b.exports[filenameString])
		}
	}

	return nil
}

// Instantiate creates a new runtime from this bundle.
func (b *Bundle) Instantiate(logger logrus.FieldLogger, vuID uint64, workerInfo *libWorker.WorkerInfo) (*BundleInstance, error) {
	// Instantiate the bundle into a new VM using a bound init context. This uses a context with a
	// runtime, but no state, to allow module-provided types to function within the init context.
	vuImpl := &moduleVUImpl{runtime: goja.New()}
	init := newBoundInitContext(b.BaseInitContext, vuImpl)
	if err := b.instantiate(logger, vuImpl.runtime, init, vuID, workerInfo, false); err != nil {
		return nil, err
	}

	rt := vuImpl.runtime

	biExports := make(map[string]map[string]goja.Callable)
	for _, filename := range b.Filenames {
		biExports[filename.String()] = make(map[string]goja.Callable)
	}

	bi := &BundleInstance{
		Runtime:      rt,
		exports:      biExports,
		moduleVUImpl: vuImpl,
		pgms:         init.programs,
		RootFilename: b.RootFilename,
		//TestData:     workerInfo.TestData,
	}

	var instErr error

	for filename, bundleFileExports := range b.exports {
		// Grab any exported functions that could be executed. These were
		// already pre-validated in cmd.validateScenarioConfig(), just get them here.
		pgm := init.programs[filename]

		moduleExports := pgm.module.Get("exports").ToObject(rt)

		for k := range bundleFileExports {
			fn, _ := goja.AssertFunction(moduleExports.Get(k))
			bi.exports[filename][k] = fn
		}

		if filename != b.RootFilename.String() {
			continue
		}

		var jsOptionsObj *goja.Object = rt.NewObject()
		err := moduleExports.Set("options", jsOptionsObj)
		if err != nil {
			return nil, fmt.Errorf("couldn't set exported options with merged values: %w", err)
		}

		b.Options.ForEachSpecified("json", func(key string, val interface{}) {
			if err := jsOptionsObj.Set(key, val); err != nil {
				instErr = err
			}
		})
	}

	return bi, instErr
}

// Instantiates the bundle into an existing runtime. Not public because it also messes with a bunch
// of other things, will potentially thrash data and makes a mess in it if the operation fails.
// In isOrchestrator mode, the init stage is limited to 30ms, and the bundle is not allowed to use the console. This is used to make execution safer on orchestrator nodes.
func (b *Bundle) initializeProgramObjects(rt *goja.Runtime, init *InitContext) []programWithSource {
	programs := make([]programWithSource, len(b.Filenames))

	for index, filename := range b.Filenames {
		pgm := programWithSource{
			pgm:     b.Programs[index],
			src:     b.Sources[index],
			srcPwd:  b.Filenames[index],
			exports: rt.NewObject(),
			module:  rt.NewObject(),
		}
		_ = pgm.module.Set("exports", pgm.exports)

		init.programs[filename.String()] = pgm
		programs[index] = pgm
	}

	return programs
}

func (b *Bundle) instantiate(logger logrus.FieldLogger, rt *goja.Runtime, init *InitContext, vuID uint64, workerInfo *libWorker.WorkerInfo, isOrchestrator bool) (err error) {
	rt.SetFieldNameMapper(common.FieldNameMapper{})
	rt.SetRandSource(common.NewRandSource())

	rt.Set("__VU", vuID)

	if !isOrchestrator {
		_ = rt.Set("console", newConsole(logger))

		rt.Set("global", rt.GlobalObject())

	}

	initenv := &common.InitEnvironment{
		Logger:      logger,
		FileSystems: init.filesystems,
		// Believe this is for the root init so use the root CWD
		CWD:        (*init.pwds)[0],
		Registry:   b.registry,
		WorkerInfo: workerInfo,
	}

	unbindInit := b.setInitGlobals(rt, init)
	init.moduleVUImpl.ctx = context.Background()
	init.moduleVUImpl.initEnv = initenv
	init.moduleVUImpl.eventLoop = eventloop.New(init.moduleVUImpl)

	pgms := b.initializeProgramObjects(rt, init)

	if isOrchestrator {
		time.AfterFunc(50*time.Millisecond, func() {
			rt.Interrupt("init stage timeout")
		})
	}

	// This is a temporary VU
	err = common.RunWithPanicCatching(logger, rt, func() error {
		return init.moduleVUImpl.eventLoop.Start(func() error {
			for index, pgm := range pgms {
				// Actually run the compiled code
				f, errRun := rt.RunProgram(b.Programs[index])

				if errRun != nil {
					return errRun
				}
				if call, ok := goja.AssertFunction(f); ok {
					if _, errRun = call(pgm.exports, pgm.module, pgm.exports); errRun != nil {
						return errRun
					}
					continue
				}
				panic("Somehow a commonjs main module is not wrapped in a function")
			}

			return nil
		})
	})

	if err != nil {
		var exception *goja.Exception
		if errors.As(err, &exception) {
			err = &scriptException{inner: exception}
		}
		return err
	}

	for _, pgm := range pgms {
		exportsV := pgm.module.Get("exports")

		if goja.IsNull(exportsV) {
			return errors.New("exports must be an object")
		}
		pgm.exports = exportsV.ToObject(rt)
		init.programs[pgm.srcPwd.String()] = pgm
	}

	unbindInit()
	init.moduleVUImpl.ctx = nil

	rt.SetRandSource(common.NewRandSource())

	return nil
}

func (b *Bundle) setInitGlobals(rt *goja.Runtime, init *InitContext) (unset func()) {
	mustSet := func(k string, v interface{}) {
		if err := rt.Set(k, v); err != nil {
			panic(fmt.Errorf("failed to set '%s' global object: %w", k, err))
		}
	}
	mustSet("require", init.Require)
	return func() {
		mustSet("require", goja.Undefined())
		mustSet("open", goja.Undefined())
	}
}

func generateSourceMapLoader(logger logrus.FieldLogger, filesystems map[string]afero.Fs,
) func(path string) ([]byte, error) {
	return func(path string) ([]byte, error) {
		u, err := url.Parse(path)
		if err != nil {
			return nil, err
		}
		data, err := loader.Load(logger, filesystems, u, path)
		if err != nil {
			return nil, err
		}
		return data.Data, nil
	}
}
