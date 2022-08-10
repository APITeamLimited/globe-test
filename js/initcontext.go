package js

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/compiler"
	"go.k6.io/k6/js/eventloop"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/js/modules/k6"
	"go.k6.io/k6/js/modules/k6/crypto"
	"go.k6.io/k6/js/modules/k6/crypto/x509"
	"go.k6.io/k6/js/modules/k6/data"
	"go.k6.io/k6/js/modules/k6/encoding"
	"go.k6.io/k6/js/modules/k6/execution"
	"go.k6.io/k6/js/modules/k6/grpc"
	"go.k6.io/k6/js/modules/k6/html"
	"go.k6.io/k6/js/modules/k6/http"
	"go.k6.io/k6/js/modules/k6/metrics"
	"go.k6.io/k6/js/modules/k6/ws"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/fsext"
	"go.k6.io/k6/loader"
)

type programWithSource struct ***REMOVED***
	pgm     *goja.Program
	src     string
	module  *goja.Object
	exports *goja.Object
***REMOVED***

const openCantBeUsedOutsideInitContextMsg = `The "open()" function is only available in the init stage ` +
	`(i.e. the global scope), see https://k6.io/docs/using-k6/test-life-cycle for more information`

// InitContext provides APIs for use in the init context.
//
// TODO: refactor most/all of this state away, use common.InitEnvironment instead
type InitContext struct ***REMOVED***
	// Bound runtime; used to instantiate objects.
	compiler *compiler.Compiler

	moduleVUImpl *moduleVUImpl

	// Filesystem to load files and scripts from with the map key being the scheme
	filesystems map[string]afero.Fs
	pwd         *url.URL

	// Cache of loaded programs and files.
	programs map[string]programWithSource

	compatibilityMode lib.CompatibilityMode

	logger logrus.FieldLogger

	modules map[string]interface***REMOVED******REMOVED***
***REMOVED***

// NewInitContext creates a new initcontext with the provided arguments
func NewInitContext(
	logger logrus.FieldLogger, rt *goja.Runtime, c *compiler.Compiler, compatMode lib.CompatibilityMode,
	filesystems map[string]afero.Fs, pwd *url.URL,
) *InitContext ***REMOVED***
	return &InitContext***REMOVED***
		compiler:          c,
		filesystems:       filesystems,
		pwd:               pwd,
		programs:          make(map[string]programWithSource),
		compatibilityMode: compatMode,
		logger:            logger,
		modules:           getJSModules(),
		moduleVUImpl: &moduleVUImpl***REMOVED***
			ctx:     context.Background(),
			runtime: rt,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func newBoundInitContext(base *InitContext, vuImpl *moduleVUImpl) *InitContext ***REMOVED***
	// we don't copy the exports as otherwise they will be shared and we don't want this.
	// this means that all the files will be executed again but once again only once per compilation
	// of the main file.
	programs := make(map[string]programWithSource, len(base.programs))
	for key, program := range base.programs ***REMOVED***
		programs[key] = programWithSource***REMOVED***
			src: program.src,
			pgm: program.pgm,
		***REMOVED***
	***REMOVED***
	return &InitContext***REMOVED***
		filesystems: base.filesystems,
		pwd:         base.pwd,
		compiler:    base.compiler,

		programs:          programs,
		compatibilityMode: base.compatibilityMode,
		logger:            base.logger,
		modules:           base.modules,
		moduleVUImpl:      vuImpl,
	***REMOVED***
***REMOVED***

// Require is called when a module/file needs to be loaded by a script
func (i *InitContext) Require(arg string) goja.Value ***REMOVED***
	switch ***REMOVED***
	case arg == "k6", strings.HasPrefix(arg, "k6/"):
		// Builtin or external modules ("k6", "k6/*", or "k6/x/*") are handled
		// specially, as they don't exist on the filesystem. This intentionally
		// shadows attempts to name your own modules this.
		v, err := i.requireModule(arg)
		if err != nil ***REMOVED***
			common.Throw(i.moduleVUImpl.runtime, err)
		***REMOVED***
		return v
	default:
		// Fall back to loading from the filesystem.
		v, err := i.requireFile(arg)
		if err != nil ***REMOVED***
			common.Throw(i.moduleVUImpl.runtime, err)
		***REMOVED***
		return v
	***REMOVED***
***REMOVED***

type moduleVUImpl struct ***REMOVED***
	ctx       context.Context
	initEnv   *common.InitEnvironment
	state     *lib.State
	runtime   *goja.Runtime
	eventLoop *eventloop.EventLoop
***REMOVED***

func (m *moduleVUImpl) Context() context.Context ***REMOVED***
	return m.ctx
***REMOVED***

func (m *moduleVUImpl) InitEnv() *common.InitEnvironment ***REMOVED***
	return m.initEnv
***REMOVED***

func (m *moduleVUImpl) State() *lib.State ***REMOVED***
	return m.state
***REMOVED***

func (m *moduleVUImpl) Runtime() *goja.Runtime ***REMOVED***
	return m.runtime
***REMOVED***

func (m *moduleVUImpl) RegisterCallback() func(func() error) ***REMOVED***
	return m.eventLoop.RegisterCallback()
***REMOVED***

/* This is here to illustrate how to use RegisterCallback to get a promise to work with the event loop
// TODO move this to a common function or remove before merging

// MakeHandledPromise will create and promise and return it's resolve, reject methods as well wrapped in such a way that
// it will block the eventloop from exiting before they are called even if the promise isn't resolved by the time the
// current script ends executing
func (m *moduleVUImpl) MakeHandledPromise() (*goja.Promise, func(interface***REMOVED******REMOVED***), func(interface***REMOVED******REMOVED***)) ***REMOVED***
	callback := m.eventLoop.registerCallback()
	p, resolve, reject := m.runtime.NewPromise()
	return p, func(i interface***REMOVED******REMOVED***) ***REMOVED***
			// more stuff
			callback(func() ***REMOVED*** resolve(i) ***REMOVED***)
		***REMOVED***, func(i interface***REMOVED******REMOVED***) ***REMOVED***
			// more stuff
			callback(func() ***REMOVED*** reject(i) ***REMOVED***)
		***REMOVED***
***REMOVED***
*/

func toESModuleExports(exp modules.Exports) interface***REMOVED******REMOVED*** ***REMOVED***
	if exp.Named == nil ***REMOVED***
		return exp.Default
	***REMOVED***
	if exp.Default == nil ***REMOVED***
		return exp.Named
	***REMOVED***

	result := make(map[string]interface***REMOVED******REMOVED***, len(exp.Named)+2)

	for k, v := range exp.Named ***REMOVED***
		result[k] = v
	***REMOVED***
	// Maybe check that those weren't set
	result["default"] = exp.Default
	// this so babel works with the `default` when it transpiles from ESM to commonjs.
	// This should probably be removed once we have support for ESM directly. So that require doesn't get support for
	// that while ESM has.
	result["__esModule"] = true

	return result
***REMOVED***

func (i *InitContext) requireModule(name string) (goja.Value, error) ***REMOVED***
	mod, ok := i.modules[name]
	if !ok ***REMOVED***
		return nil, fmt.Errorf("unknown module: %s", name)
	***REMOVED***
	if m, ok := mod.(modules.Module); ok ***REMOVED***
		instance := m.NewModuleInstance(i.moduleVUImpl)
		return i.moduleVUImpl.runtime.ToValue(toESModuleExports(instance.Exports())), nil
	***REMOVED***

	return i.moduleVUImpl.runtime.ToValue(mod), nil
***REMOVED***

func (i *InitContext) requireFile(name string) (goja.Value, error) ***REMOVED***
	// Resolve the file path, push the target directory as pwd to make relative imports work.
	pwd := i.pwd
	fileURL, err := loader.Resolve(pwd, name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// First, check if we have a cached program already.
	pgm, ok := i.programs[fileURL.String()]
	if !ok || pgm.module == nil ***REMOVED***
		if filepath.IsAbs(name) && runtime.GOOS == "windows" ***REMOVED***
			i.logger.Warnf("'%s' was imported with an absolute path - this won't be cross-platform and won't work if"+
				" you move the script between machines or run it with `k6 cloud`; if absolute paths are required,"+
				" import them with the `file://` schema for slightly better compatibility",
				name)
		***REMOVED***
		i.pwd = loader.Dir(fileURL)
		defer func() ***REMOVED*** i.pwd = pwd ***REMOVED***()
		exports := i.moduleVUImpl.runtime.NewObject()
		pgm.module = i.moduleVUImpl.runtime.NewObject()
		_ = pgm.module.Set("exports", exports)

		if pgm.pgm == nil ***REMOVED***
			// Load the sources; the loader takes care of remote loading, etc.
			data, err := loader.Load(i.logger, i.filesystems, fileURL, name)
			if err != nil ***REMOVED***
				return goja.Undefined(), err
			***REMOVED***

			pgm.src = string(data.Data)

			// Compile the sources; this handles ES5 vs ES6 automatically.
			pgm.pgm, err = i.compileImport(pgm.src, data.URL.String())
			if err != nil ***REMOVED***
				return goja.Undefined(), err
			***REMOVED***
		***REMOVED***

		i.programs[fileURL.String()] = pgm

		// Run the program.
		f, err := i.moduleVUImpl.runtime.RunProgram(pgm.pgm)
		if err != nil ***REMOVED***
			delete(i.programs, fileURL.String())
			return goja.Undefined(), err
		***REMOVED***
		if call, ok := goja.AssertFunction(f); ok ***REMOVED***
			if _, err = call(exports, pgm.module, exports); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return pgm.module.Get("exports"), nil
***REMOVED***

func (i *InitContext) compileImport(src, filename string) (*goja.Program, error) ***REMOVED***
	pgm, _, err := i.compiler.Compile(src, filename, false)
	return pgm, err
***REMOVED***

// Open implements open() in the init context and will read and return the
// contents of a file. If the second argument is "b" it returns an ArrayBuffer
// instance, otherwise a string representation.
func (i *InitContext) Open(filename string, args ...string) (goja.Value, error) ***REMOVED***
	if i.moduleVUImpl.State() != nil ***REMOVED***
		return nil, errors.New(openCantBeUsedOutsideInitContextMsg)
	***REMOVED***

	if filename == "" ***REMOVED***
		return nil, errors.New("open() can't be used with an empty filename")
	***REMOVED***

	// Here IsAbs should be enough but unfortunately it doesn't handle absolute paths starting from
	// the current drive on windows like `\users\noname\...`. Also it makes it more easy to test and
	// will probably be need for archive execution under windows if always consider '/...' as an
	// absolute path.
	if filename[0] != '/' && filename[0] != '\\' && !filepath.IsAbs(filename) ***REMOVED***
		filename = filepath.Join(i.pwd.Path, filename)
	***REMOVED***
	filename = filepath.Clean(filename)
	fs := i.filesystems["file"]
	if filename[0:1] != afero.FilePathSeparator ***REMOVED***
		filename = afero.FilePathSeparator + filename
	***REMOVED***

	data, err := readFile(fs, filename)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(args) > 0 && args[0] == "b" ***REMOVED***
		ab := i.moduleVUImpl.runtime.NewArrayBuffer(data)
		return i.moduleVUImpl.runtime.ToValue(&ab), nil
	***REMOVED***
	return i.moduleVUImpl.runtime.ToValue(string(data)), nil
***REMOVED***

func readFile(fileSystem afero.Fs, filename string) (data []byte, err error) ***REMOVED***
	defer func() ***REMOVED***
		if errors.Is(err, fsext.ErrPathNeverRequestedBefore) ***REMOVED***
			// loading different files per VU is not supported, so all files should are going
			// to be used inside the scenario should be opened during the init step (without any conditions)
			err = fmt.Errorf(
				"open() can't be used with files that weren't previously opened during initialization (__VU==0), path: %q",
				filename,
			)
		***REMOVED***
	***REMOVED***()

	// Workaround for https://github.com/spf13/afero/issues/201
	if isDir, err := afero.IsDir(fileSystem, filename); err != nil ***REMOVED***
		return nil, err
	***REMOVED*** else if isDir ***REMOVED***
		return nil, fmt.Errorf("open() can't be used with directories, path: %q", filename)
	***REMOVED***

	return afero.ReadFile(fileSystem, filename)
***REMOVED***

// allowOnlyOpenedFiles enables seen only files
func (i *InitContext) allowOnlyOpenedFiles() ***REMOVED***
	fs := i.filesystems["file"]

	alreadyOpenedFS, ok := fs.(fsext.OnlyCachedEnabler)
	if !ok ***REMOVED***
		return
	***REMOVED***

	alreadyOpenedFS.AllowOnlyCached()
***REMOVED***

func getInternalJSModules() map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	return map[string]interface***REMOVED******REMOVED******REMOVED***
		"k6":             k6.New(),
		"k6/crypto":      crypto.New(),
		"k6/crypto/x509": x509.New(),
		"k6/data":        data.New(),
		"k6/encoding":    encoding.New(),
		"k6/execution":   execution.New(),
		"k6/net/grpc":    grpc.New(),
		"k6/html":        html.New(),
		"k6/http":        http.New(),
		"k6/metrics":     metrics.New(),
		"k6/ws":          ws.New(),
	***REMOVED***
***REMOVED***

func getJSModules() map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	result := getInternalJSModules()
	external := modules.GetJSModules()

	// external is always prefixed with `k6/x`
	for k, v := range external ***REMOVED***
		result[k] = v
	***REMOVED***

	return result
***REMOVED***
