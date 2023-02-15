package js

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/APITeamLimited/globe-test/js/common"
	"github.com/APITeamLimited/globe-test/js/compiler"
	"github.com/APITeamLimited/globe-test/js/eventloop"
	"github.com/APITeamLimited/globe-test/js/modules"
	apiteamContext "github.com/APITeamLimited/globe-test/js/modules/apiteam/context"
	"github.com/APITeamLimited/globe-test/js/modules/k6"
	"github.com/APITeamLimited/globe-test/js/modules/k6/crypto"
	"github.com/APITeamLimited/globe-test/js/modules/k6/crypto/x509"
	"github.com/APITeamLimited/globe-test/js/modules/k6/data"
	"github.com/APITeamLimited/globe-test/js/modules/k6/encoding"
	"github.com/APITeamLimited/globe-test/js/modules/k6/execution"
	"github.com/APITeamLimited/globe-test/js/modules/k6/html"
	"github.com/APITeamLimited/globe-test/js/modules/k6/http"
	"github.com/APITeamLimited/globe-test/js/modules/k6/metrics"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type programWithSource struct {
	pgm     *goja.Program
	srcPwd  *url.URL
	src     string
	module  *goja.Object
	exports *goja.Object
}

const openCantBeUsedOutsideInitContextMsg = `The "open()" function is only available in the init stage ` +
	`(i.e. the global scope), see https://k6.io/docs/using-k6/test-life-cycle for more information`

// InitContext provides APIs for use in the init context.
//
// TODO: refactor most/all of this state away, use common.InitEnvironment instead
type InitContext struct {
	// Bound runtime; used to instantiate objects.
	compiler *compiler.Compiler

	moduleVUImpl *moduleVUImpl

	// Filesystem to load files and scripts from with the map key being the scheme
	filesystems map[string]afero.Fs
	pwds        *[]*url.URL

	// Cache of loaded programs and files.
	programs map[string]programWithSource

	logger logrus.FieldLogger

	modules map[string]interface{}
}

// NewInitContext creates a new initcontext with the provided arguments
func NewInitContext(
	logger logrus.FieldLogger, rt *goja.Runtime, c *compiler.Compiler,
	filesystems map[string]afero.Fs, pwds *[]*url.URL, workerInfo *libWorker.WorkerInfo, rootNode libWorker.Node) *InitContext {
	return &InitContext{
		compiler:    c,
		filesystems: filesystems,
		pwds:        pwds,
		programs:    make(map[string]programWithSource),
		logger:      logger,
		modules:     getJSModules(workerInfo),
		moduleVUImpl: &moduleVUImpl{
			ctx:     context.Background(),
			runtime: rt,
		},
	}
}

func newBoundInitContext(base *InitContext, vuImpl *moduleVUImpl) *InitContext {
	// we don't copy the exports as otherwise they will be shared and we don't want this.
	// this means that all the files will be executed again but once again only once per compilation
	// of the main file.
	programs := make(map[string]programWithSource, len(base.programs))
	for key, program := range base.programs {
		programs[key] = programWithSource{
			src: program.src,
			pgm: program.pgm,
		}
	}
	return &InitContext{
		filesystems: base.filesystems,
		pwds:        base.pwds,
		compiler:    base.compiler,

		programs:     programs,
		logger:       base.logger,
		modules:      base.modules,
		moduleVUImpl: vuImpl,
	}
}

// Require is called when a module/file needs to be loaded by a script
func (i *InitContext) Require(arg string) goja.Value {
	switch {
	case arg == "k6", strings.HasPrefix(arg, "k6/"):
		// Builtin or external modules ("k6", "k6/*", or "k6/x/*") are handled
		// specially, as they don't exist on the filesystem. This intentionally
		// shadows attempts to name your own modules this.
		v, err := i.requireModule(arg)
		if err != nil {
			common.Throw(i.moduleVUImpl.runtime, err)
		}
		return v
	case arg == "apiteam", strings.HasPrefix(arg, "apiteam/"):
		// Handle apiteam modules like k6 modules
		v, err := i.requireModule(arg)
		if err != nil {
			common.Throw(i.moduleVUImpl.runtime, err)
		}
		return v
	default:
		// Fall back to loading from the filesystem.
		v, err := i.requireFile(arg)
		if err != nil {
			common.Throw(i.moduleVUImpl.runtime, err)
		}
		return v
	}
}

type moduleVUImpl struct {
	ctx       context.Context
	initEnv   *common.InitEnvironment
	state     *libWorker.State
	runtime   *goja.Runtime
	eventLoop *eventloop.EventLoop
}

func (m *moduleVUImpl) Context() context.Context {
	return m.ctx
}

func (m *moduleVUImpl) InitEnv() *common.InitEnvironment {
	return m.initEnv
}

func (m *moduleVUImpl) State() *libWorker.State {
	return m.state
}

func (m *moduleVUImpl) Runtime() *goja.Runtime {
	return m.runtime
}

func (m *moduleVUImpl) RegisterCallback() func(func() error) {
	return m.eventLoop.RegisterCallback()
}

func toESModuleExports(exp modules.Exports) interface{} {
	if exp.Named == nil {
		return exp.Default
	}
	if exp.Default == nil {
		return exp.Named
	}

	result := make(map[string]interface{}, len(exp.Named)+2)

	for k, v := range exp.Named {
		result[k] = v
	}
	// Maybe check that those weren't set
	result["default"] = exp.Default
	// this so babel works with the `default` when it transpiles from ESM to commonjs.
	// This should probably be removed once we have support for ESM directly. So that require doesn't get support for
	// that while ESM has.
	result["__esModule"] = true

	return result
}

func (i *InitContext) requireModule(name string) (goja.Value, error) {
	mod, ok := i.modules[name]
	if !ok {
		return nil, fmt.Errorf("unknown module: %s", name)
	}
	if m, ok := mod.(modules.Module); ok {
		instance := m.NewModuleInstance(i.moduleVUImpl)
		return i.moduleVUImpl.runtime.ToValue(toESModuleExports(instance.Exports())), nil
	}

	return i.moduleVUImpl.runtime.ToValue(mod), nil
}

func (i *InitContext) requireFile(name string) (goja.Value, error) {
	return nil, fmt.Errorf("globe test does not support files yet")

	/*
		// Resolve the file path, push the target directory as pwd to make relative imports work.
		pwd := i.pwd
		fileURL, err := loader.Resolve(pwd, name)
		if err != nil {
			return nil, err
		}

		// First, check if we have a cached program already.
		pgm, ok := i.programs[fileURL.String()]
		if !ok || pgm.module == nil {
			if filepath.IsAbs(name) && runtime.GOOS == "windows" {
				i.logger.Warnf("'%s' was imported with an absolute path - this won't be cross-platform and won't work if"+
					" you move the script between machines or run it with `k6 cloud`; if absolute paths are required,"+
					" import them with the `file://` schema for slightly better compatibility",
					name)
			}
			i.pwd = loader.Dir(fileURL)
			defer func() { i.pwd = pwd }()
			exports := i.moduleVUImpl.runtime.NewObject()
			pgm.module = i.moduleVUImpl.runtime.NewObject()
			_ = pgm.module.Set("exports", exports)

			if pgm.pgm == nil {
				// Load the sources; the loader takes care of remote loading, etc.
				data, err := loader.Load(i.logger, i.filesystems, fileURL, name)
				if err != nil {
					return goja.Undefined(), err
				}

				pgm.src = string(data.Data)

				// Compile the sources; this handles ES5 vs ES6 automatically.
				pgm.pgm, err = i.compileImport(pgm.src, data.URL.String())
				if err != nil {
					return goja.Undefined(), err
				}
			}

			i.programs[fileURL.String()] = pgm

			// Run the program.
			f, err := i.moduleVUImpl.runtime.RunProgram(pgm.pgm)
			if err != nil {
				delete(i.programs, fileURL.String())
				return goja.Undefined(), err
			}
			if call, ok := goja.AssertFunction(f); ok {
				if _, err = call(exports, pgm.module, exports); err != nil {
					return nil, err
				}
			}
		}

		return pgm.module.Get("exports"), nil
	*/
}

func getInternalJSModules(workerInfo *libWorker.WorkerInfo) map[string]interface{} {
	k6Module := k6.New()
	kyCryptoModule := crypto.New()
	k6Crypto509Module := x509.New()
	k6DataModule := data.New()
	k6EncodingModule := encoding.New()
	k6ExecutionModule := execution.New()
	k6HTMLModule := html.New()
	k6HTTPModule := http.New(workerInfo)
	k6MetricsModule := metrics.New()

	return map[string]interface{}{
		"k6":             k6Module,
		"k6/crypto":      kyCryptoModule,
		"k6/crypto/x509": k6Crypto509Module,
		"k6/data":        k6DataModule,
		"k6/encoding":    k6EncodingModule,
		"k6/execution":   k6ExecutionModule,
		//"k6/net/grpc":         grpc.New(),
		"k6/html":             k6HTMLModule,
		"k6/http":             k6HTTPModule,
		"k6/metrics":          k6MetricsModule,
		"apiteam":             k6Module,
		"apiteam/crypto":      kyCryptoModule,
		"apiteam/crypto/x509": k6Crypto509Module,
		"apiteam/data":        k6DataModule,
		"apiteam/encoding":    k6EncodingModule,
		"apiteam/execution":   k6ExecutionModule,
		//"k6/net/grpc":         grpc.New(),
		"apiteam/html":    k6HTMLModule,
		"apiteam/http":    k6HTTPModule,
		"apiteam/metrics": k6MetricsModule,
		//"k6/ws":               ws.New(),
		"apiteam/context": apiteamContext.New(workerInfo),
	}
}

func getJSModules(workerInfo *libWorker.WorkerInfo) map[string]interface{} {
	result := getInternalJSModules(workerInfo)
	external := modules.GetJSModules()

	// external is always prefixed with `k6/x`
	for k, v := range external {
		result[k] = v
	}

	return result
}
