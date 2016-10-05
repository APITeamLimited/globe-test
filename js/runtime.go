package js

import (
	"bytes"
	"fmt"
	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/lib"
	"github.com/robertkrimen/otto"
	"gopkg.in/guregu/null.v3"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const wrapper = "(function() ***REMOVED*** var e = ***REMOVED******REMOVED***; (function(exports) ***REMOVED***%s\n***REMOVED***)(e); return e; ***REMOVED***)();"

var (
	libBox      = rice.MustFindBox("lib")
	polyfillBox = rice.MustFindBox("node_modules/babel-polyfill")
)

type Runtime struct ***REMOVED***
	VM      *otto.Otto
	Root    string
	Exports map[string]otto.Value
	Lib     map[string]otto.Value
***REMOVED***

func New() (*Runtime, error) ***REMOVED***
	wd, err := os.Getwd()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	vm := otto.New()

	polyfillJS, err := polyfillBox.String("dist/polyfill.js")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	polyfill, err := vm.Compile("polyfill.js", polyfillJS)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if _, err := vm.Run(polyfill); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Runtime***REMOVED***
		VM:      vm,
		Root:    wd,
		Exports: make(map[string]otto.Value),
		Lib:     make(map[string]otto.Value),
	***REMOVED***, nil
***REMOVED***

func (r *Runtime) Load(filename string) (otto.Value, error) ***REMOVED***
	// You can only require() modules during setup; doing it during runtime would sink your test's
	// performance in potentially nonobvious ways. Only relative imports ("./*", "../*") are
	// allowed for user code; absolute imports are used for system-provided libraries.
	r.VM.Set("require", r.require)
	defer r.VM.Set("require", nil)

	return r.loadFile(filename)
***REMOVED***

func (r *Runtime) require(call otto.FunctionCall) otto.Value ***REMOVED***
	name := call.Argument(0).String()
	if !strings.HasPrefix(name, ".") ***REMOVED***
		exports, err := r.loadLib(name + ".js")
		if err != nil ***REMOVED***
			panic(call.Otto.MakeCustomError("ImportError", err.Error()))
		***REMOVED***
		return exports
	***REMOVED***

	exports, err := r.loadFile(name + ".js")
	if err != nil ***REMOVED***
		panic(call.Otto.MakeCustomError("ImportError", err.Error()))
	***REMOVED***
	return exports
***REMOVED***

func (r *Runtime) ExtractOptions(exports otto.Value, opts *lib.Options) error ***REMOVED***
	expObj := exports.Object()
	if expObj == nil ***REMOVED***
		return nil
	***REMOVED***

	v, err := expObj.Get("options")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	obj := v.Object()
	if obj == nil ***REMOVED***
		return nil
	***REMOVED***

	for _, key := range obj.Keys() ***REMOVED***
		val, err := obj.Get(key)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		switch key ***REMOVED***
		case "vus":
			vus, err := val.ToInteger()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			opts.VUs = null.IntFrom(vus)
		case "vusMax":
			vusMax, err := val.ToInteger()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			opts.VUsMax = null.IntFrom(vusMax)
		case "duration":
			duration, err := val.ToString()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			opts.Duration = null.StringFrom(duration)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *Runtime) loadFile(filename string) (otto.Value, error) ***REMOVED***
	// To protect against directory traversal, prevent loading of files outside the root (pwd) dir
	path, err := filepath.Abs(filename)
	if err != nil ***REMOVED***
		return otto.UndefinedValue(), err
	***REMOVED***
	if !strings.HasPrefix(path, r.Root) ***REMOVED***
		return otto.UndefinedValue(), DirectoryTraversalError***REMOVED***Filename: filename, Root: r.Root***REMOVED***
	***REMOVED***

	// Don't re-compile repeated includes of the same module
	if exports, ok := r.Exports[path]; ok ***REMOVED***
		return exports, nil
	***REMOVED***

	data, err := ioutil.ReadFile(path)
	if err != nil ***REMOVED***
		return otto.UndefinedValue(), err
	***REMOVED***
	exports, err := r.load(path, data)
	if err != nil ***REMOVED***
		return otto.UndefinedValue(), err
	***REMOVED***
	r.Exports[path] = exports

	log.WithField("path", path).Debug("File loaded")

	return exports, nil
***REMOVED***

func (r *Runtime) loadLib(filename string) (otto.Value, error) ***REMOVED***
	if exports, ok := r.Lib[filename]; ok ***REMOVED***
		return exports, nil
	***REMOVED***

	data, err := libBox.Bytes(filename)
	if err != nil ***REMOVED***
		return otto.UndefinedValue(), err
	***REMOVED***
	exports, err := r.load(filename, data)
	if err != nil ***REMOVED***
		return otto.UndefinedValue(), err
	***REMOVED***
	r.Lib[filename] = exports

	log.WithField("filename", filename).Debug("Library loaded")

	return exports, nil
***REMOVED***

func (r *Runtime) load(filename string, data []byte) (otto.Value, error) ***REMOVED***
	// Compile the file with Babel; this subprocess invocation is TEMPORARY:
	// https://github.com/robertkrimen/otto/pull/205
	cmd := exec.Command(babel, "--presets", "latest", "--no-babelrc")
	cmd.Dir = babelDir
	cmd.Stdin = bytes.NewReader(data)
	cmd.Stderr = os.Stderr
	src, err := cmd.Output()
	if err != nil ***REMOVED***
		return otto.UndefinedValue(), err
	***REMOVED***

	// Use a wrapper function to turn the script into an exported module
	s, err := r.VM.Compile(filename, fmt.Sprintf(wrapper, string(src)))
	if err != nil ***REMOVED***
		return otto.UndefinedValue(), err
	***REMOVED***
	exports, err := r.VM.Run(s)
	if err != nil ***REMOVED***
		return otto.UndefinedValue(), err
	***REMOVED***

	return exports, nil
***REMOVED***
