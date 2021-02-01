// Heavily influenced by the fantastic work by @dop251 for https://github.com/dop251/goja

package tc39

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"github.com/loadimpact/k6/js/compiler"
	jslib "github.com/loadimpact/k6/js/lib"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/testutils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

const (
	tc39BASE = "TestTC39/test262"
)

//nolint:gochecknoglobals
var (
	errInvalidFormat = errors.New("invalid file format")

	// ignorableTestError = newSymbol(stringEmpty)

	sabStub = goja.MustCompile("sabStub.js", `
		Object.defineProperty(this, "SharedArrayBuffer", ***REMOVED***
			get: function() ***REMOVED***
				throw IgnorableTestError;
			***REMOVED***
		***REMOVED***);`,
		false)

	esIDPrefixAllowList = []string***REMOVED***
		"sec-array",
		"sec-%typedarray%",
		"sec-string",
		"sec-date",
		"sec-number",
		"sec-math",
		"sec-arraybuffer-length",
		"sec-arraybuffer",
		"sec-regexp",
		"sec-variable-statement",
		"sec-ecmascript-standard-built-in-objects",
	***REMOVED***

	featuresBlockList = []string***REMOVED***
		"BigInt",    // not supported at all
		"IsHTMLDDA", // not supported at all
	***REMOVED***
	skipList       = map[string]bool***REMOVED******REMOVED***
	pathBasedBlock = map[string]bool***REMOVED*** // This completely skips any path matching it without any kind of message
		"test/annexB/built-ins/Date":                          true,
		"test/annexB/built-ins/RegExp/prototype/Symbol.split": true,
		"test/annexB/built-ins/String/prototype/anchor":       true,
		"test/annexB/built-ins/String/prototype/big":          true,
		"test/annexB/built-ins/String/prototype/blink":        true,
		"test/annexB/built-ins/String/prototype/bold":         true,
		"test/annexB/built-ins/String/prototype/fixed":        true,
		"test/annexB/built-ins/String/prototype/fontcolor":    true,
		"test/annexB/built-ins/String/prototype/fontsize":     true,
		"test/annexB/built-ins/String/prototype/italics":      true,
		"test/annexB/built-ins/String/prototype/link":         true,
		"test/annexB/built-ins/String/prototype/small":        true,
		"test/annexB/built-ins/String/prototype/strike":       true,
		"test/annexB/built-ins/String/prototype/sub":          true,
		"test/annexB/built-ins/String/prototype/sup":          true,

		// Async/Promise and other totally unsupported functionality
		"test/built-ins/AsyncArrowFunction":             true,
		"test/built-ins/AsyncFromSyncIteratorPrototype": true,
		"test/built-ins/AsyncFunction":                  true,
		"test/built-ins/AsyncGeneratorFunction":         true,
		"test/built-ins/AsyncGeneratorPrototype":        true,
		"test/built-ins/AsyncIteratorPrototype":         true,
		"test/built-ins/Atomics":                        true,
		"test/built-ins/BigInt":                         true,
		"test/built-ins/Promise":                        true,
		"test/built-ins/SharedArrayBuffer":              true,

		"test/built-ins/Date/parse/without-utc-offset.js": true, // some other reason ?!? depending on local time

		"test/built-ins/Array/prototype/concat/arg-length-exceeding-integer-limit.js": true, // takes forever and is broken
		"test/built-ins/Array/prototype/splice/throws-if-integer-limit-exceeded.js":   true, // takes forever and is broken
		"test/built-ins/Array/prototype/unshift/clamps-to-integer-limit.js":           true, // takes forever and is broken
		"test/built-ins/Array/prototype/unshift/throws-if-integer-limit-exceeded.js":  true, // takes forever and is broken
	***REMOVED***
)

//nolint:unused,structcheck
type tc39Test struct ***REMOVED***
	name string
	f    func(t *testing.T)
***REMOVED***

type tc39BenchmarkItem struct ***REMOVED***
	name     string
	duration time.Duration
***REMOVED***

type tc39BenchmarkData []tc39BenchmarkItem

type tc39TestCtx struct ***REMOVED***
	compiler       *compiler.Compiler
	base           string
	t              *testing.T
	prgCache       map[string]*goja.Program
	prgCacheLock   sync.Mutex
	enableBench    bool
	benchmark      tc39BenchmarkData
	benchLock      sync.Mutex
	testQueue      []tc39Test //nolint:unused,structcheck
	expectedErrors map[string]string

	errorsLock sync.Mutex
	errors     map[string]string
***REMOVED***

type TC39MetaNegative struct ***REMOVED***
	Phase, Type string
***REMOVED***

type tc39Meta struct ***REMOVED***
	Negative TC39MetaNegative
	Includes []string
	Flags    []string
	Features []string
	Es5id    string
	Es6id    string
	Esid     string
***REMOVED***

func (m *tc39Meta) hasFlag(flag string) bool ***REMOVED***
	for _, f := range m.Flags ***REMOVED***
		if f == flag ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func parseTC39File(name string) (*tc39Meta, string, error) ***REMOVED***
	f, err := os.Open(name) //nolint:gosec
	if err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***
	defer f.Close() //nolint:errcheck,gosec

	b, err := ioutil.ReadAll(f)
	if err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***

	metaStart := bytes.Index(b, []byte("/*---"))
	if metaStart == -1 ***REMOVED***
		return nil, "", errInvalidFormat
	***REMOVED***

	metaStart += 5
	metaEnd := bytes.Index(b, []byte("---*/"))
	if metaEnd == -1 || metaEnd <= metaStart ***REMOVED***
		return nil, "", errInvalidFormat
	***REMOVED***

	var meta tc39Meta
	err = yaml.Unmarshal(b[metaStart:metaEnd], &meta)
	if err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***

	if meta.Negative.Type != "" && meta.Negative.Phase == "" ***REMOVED***
		return nil, "", errors.New("negative type is set, but phase isn't")
	***REMOVED***

	return &meta, string(b), nil
***REMOVED***

func (*tc39TestCtx) detachArrayBuffer(call goja.FunctionCall) goja.Value ***REMOVED***
	if obj, ok := call.Argument(0).(*goja.Object); ok ***REMOVED***
		var buf goja.ArrayBuffer
		if goja.New().ExportTo(obj, &buf) == nil ***REMOVED***
			// if buf, ok := obj.Export().(goja.ArrayBuffer); ok ***REMOVED***
			buf.Detach()
			return goja.Undefined()
		***REMOVED***
	***REMOVED***
	panic(goja.New().NewTypeError("detachArrayBuffer() is called with incompatible argument"))
***REMOVED***

func (ctx *tc39TestCtx) fail(t testing.TB, name string, strict bool, errStr string) ***REMOVED***
	nameKey := fmt.Sprintf("%s-strict:%v", name, strict)
	expected, ok := ctx.expectedErrors[nameKey]
	if ok ***REMOVED***
		if !assert.Equal(t, expected, errStr) ***REMOVED***
			ctx.errorsLock.Lock()
			fmt.Println("different")
			fmt.Println(expected)
			fmt.Println(errStr)
			ctx.errors[nameKey] = errStr
			ctx.errorsLock.Unlock()
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		assert.Empty(t, errStr)
		ctx.errorsLock.Lock()
		fmt.Println("no error", name)
		ctx.errors[nameKey] = errStr
		ctx.errorsLock.Unlock()
	***REMOVED***
***REMOVED***

func (ctx *tc39TestCtx) runTC39Test(t testing.TB, name, src string, meta *tc39Meta, strict bool) ***REMOVED***
	if skipList[name] ***REMOVED***
		t.Skip("Excluded")
	***REMOVED***
	failf := func(str string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
		str = fmt.Sprintf(str, args...)
		ctx.fail(t, name, strict, str)
	***REMOVED***
	defer func() ***REMOVED***
		if x := recover(); x != nil ***REMOVED***
			failf("panic while running %s: %v", name, x)
		***REMOVED***
	***REMOVED***()
	vm := goja.New()
	_262 := vm.NewObject()
	ignorableTestError := vm.NewGoError(fmt.Errorf(""))
	vm.Set("IgnorableTestError", ignorableTestError)
	_ = _262.Set("detachArrayBuffer", ctx.detachArrayBuffer)
	_ = _262.Set("createRealm", func(goja.FunctionCall) goja.Value ***REMOVED***
		panic(ignorableTestError)
	***REMOVED***)
	vm.Set("$262", _262)
	vm.Set("print", t.Log)
	if _, err := vm.RunProgram(jslib.GetCoreJS()); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	_, err := vm.RunProgram(sabStub)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	if strict ***REMOVED***
		src = "'use strict';\n" + src
	***REMOVED***
	early, err := ctx.runTC39Script(name, src, meta.Includes, vm)

	if err == nil ***REMOVED***
		if meta.Negative.Type != "" ***REMOVED***
			// vm.vm.prg.dumpCode(t.Logf)
			failf("%s: Expected error: %v", name, err)
			return
		***REMOVED***
		nameKey := fmt.Sprintf("%s-strict:%v", name, strict)
		expected, ok := ctx.expectedErrors[nameKey]
		assert.False(t, ok, "%s passes but and error %q was expected", nameKey, expected)
		return
	***REMOVED***

	if meta.Negative.Type == "" ***REMOVED***
		if err, ok := err.(*goja.Exception); ok ***REMOVED***
			if err.Value() == ignorableTestError ***REMOVED***
				t.Skip("Test threw IgnorableTestError")
			***REMOVED***
		***REMOVED***
		failf("%s: %v", name, err)
		return
	***REMOVED***
	if meta.Negative.Phase == "early" && !early || meta.Negative.Phase == "runtime" && early ***REMOVED***
		failf("%s: error %v happened at the wrong phase (expected %s)", name, err, meta.Negative.Phase)
		return
	***REMOVED***
	var errType string

	switch err := err.(type) ***REMOVED***
	case *goja.Exception:
		if o, ok := err.Value().(*goja.Object); ok ***REMOVED*** //nolint:nestif
			if c := o.Get("constructor"); c != nil ***REMOVED***
				if c, ok := c.(*goja.Object); ok ***REMOVED***
					errType = c.Get("name").String()
				***REMOVED*** else ***REMOVED***
					failf("%s: error constructor is not an object (%v)", name, o)
					return
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				failf("%s: error does not have a constructor (%v)", name, o)
				return
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			failf("%s: error is not an object (%v)", name, err.Value())
			return
		***REMOVED***
	case *goja.CompilerSyntaxError, *parser.Error, parser.ErrorList:
		errType = "SyntaxError"
	case *goja.CompilerReferenceError:
		errType = "ReferenceError"
	default:
		failf("%s: error is not a JS error: %v", name, err)
		return
	***REMOVED***

	_ = errType
	if errType != meta.Negative.Type ***REMOVED***
		// vm.vm.prg.dumpCode(t.Logf)
		failf("%s: unexpected error type (%s), expected (%s)", name, errType, meta.Negative.Type)
		return
	***REMOVED***

	/*
		if vm.vm.sp != 0 ***REMOVED***
			t.Fatalf("sp: %d", vm.vm.sp)
		***REMOVED***

		if l := len(vm.vm.iterStack); l > 0 ***REMOVED***
			t.Fatalf("iter stack is not empty: %d", l)
		***REMOVED***
	*/
***REMOVED***

func shouldBeSkipped(t testing.TB, meta *tc39Meta) ***REMOVED***
	if meta.hasFlag("async") ***REMOVED*** // this is totally not supported
		t.Skipf("Skipping as it has flag async")
	***REMOVED***
	if meta.Es6id == "" && meta.Es5id == "" ***REMOVED*** //nolint:nestif
		skip := true

		if skip ***REMOVED***
			if meta.Esid != "" ***REMOVED***
				for _, prefix := range esIDPrefixAllowList ***REMOVED***
					if strings.HasPrefix(meta.Esid, prefix) ***REMOVED***
						skip = false
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		for _, feature := range meta.Features ***REMOVED***
			for _, bl := range featuresBlockList ***REMOVED***
				if feature == bl ***REMOVED***
					t.Skipf("Blocklisted feature %s", feature)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if skip ***REMOVED***
			t.Skipf("Not ES6 or ES5 esid: %s", meta.Esid)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (ctx *tc39TestCtx) runTC39File(name string, t testing.TB) ***REMOVED***
	p := path.Join(ctx.base, name)
	meta, src, err := parseTC39File(p)
	if err != nil ***REMOVED***
		// t.Fatalf("Could not parse %s: %v", name, err)
		t.Errorf("Could not parse %s: %v", name, err)
		return
	***REMOVED***

	shouldBeSkipped(t, meta)

	var startTime time.Time
	if ctx.enableBench ***REMOVED***
		startTime = time.Now()
	***REMOVED***

	hasRaw := meta.hasFlag("raw")

	/*
		if hasRaw || !meta.hasFlag("onlyStrict") ***REMOVED***
			// log.Printf("Running normal test: %s", name)
			// t.Logf("Running normal test: %s", name)
			ctx.runTC39Test(t, name, src, meta, false)
		***REMOVED***
	*/

	if !hasRaw && !meta.hasFlag("noStrict") ***REMOVED***
		// log.Printf("Running strict test: %s", name)
		// t.Logf("Running strict test: %s", name)
		ctx.runTC39Test(t, name, src, meta, true)
	***REMOVED*** else ***REMOVED*** // Run test in non strict mode only if we won't run them in strict
		// TODO uncomment the if above and delete this else so we run both parts when the tests
		// don't take forever
		ctx.runTC39Test(t, name, src, meta, false)
	***REMOVED***

	if ctx.enableBench ***REMOVED***
		ctx.benchLock.Lock()
		ctx.benchmark = append(ctx.benchmark, tc39BenchmarkItem***REMOVED***
			name:     name,
			duration: time.Since(startTime),
		***REMOVED***)
		ctx.benchLock.Unlock()
	***REMOVED***
***REMOVED***

func (ctx *tc39TestCtx) init() ***REMOVED***
	ctx.prgCache = make(map[string]*goja.Program)
	ctx.errors = make(map[string]string)

	b, err := ioutil.ReadFile("./breaking_test_errors.json")
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	ctx.expectedErrors = make(map[string]string, 1000)
	err = json.Unmarshal(b, &ctx.expectedErrors)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (ctx *tc39TestCtx) compile(base, name string) (*goja.Program, error) ***REMOVED***
	ctx.prgCacheLock.Lock()
	defer ctx.prgCacheLock.Unlock()

	prg := ctx.prgCache[name]
	if prg == nil ***REMOVED***
		fname := path.Join(base, name)
		f, err := os.Open(fname) //nolint:gosec
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer f.Close() //nolint:gosec,errcheck

		b, err := ioutil.ReadAll(f)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		str := string(b)
		prg, _, err = ctx.compiler.Compile(str, name, "", "", false, lib.CompatibilityModeExtended)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		ctx.prgCache[name] = prg
	***REMOVED***

	return prg, nil
***REMOVED***

func (ctx *tc39TestCtx) runFile(base, name string, vm *goja.Runtime) error ***REMOVED***
	prg, err := ctx.compile(base, name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = vm.RunProgram(prg)
	return err
***REMOVED***

func (ctx *tc39TestCtx) runTC39Script(name, src string, includes []string, vm *goja.Runtime) (early bool, err error) ***REMOVED***
	early = true
	err = ctx.runFile(ctx.base, path.Join("harness", "assert.js"), vm)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	err = ctx.runFile(ctx.base, path.Join("harness", "sta.js"), vm)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	for _, include := range includes ***REMOVED***
		err = ctx.runFile(ctx.base, path.Join("harness", include), vm)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	var p *goja.Program
	p, _, err = ctx.compiler.Compile(src, name, "", "", false, lib.CompatibilityModeExtended)

	if err != nil ***REMOVED***
		return
	***REMOVED***

	early = false
	_, err = vm.RunProgram(p)

	return
***REMOVED***

func (ctx *tc39TestCtx) runTC39Tests(name string) ***REMOVED***
	files, err := ioutil.ReadDir(path.Join(ctx.base, name))
	if err != nil ***REMOVED***
		ctx.t.Fatal(err)
	***REMOVED***

	for _, file := range files ***REMOVED***
		if file.Name()[0] == '.' ***REMOVED***
			continue
		***REMOVED***
		newName := path.Join(name, file.Name())
		if pathBasedBlock[newName] ***REMOVED***
			ctx.t.Run(newName, func(t *testing.T) ***REMOVED***
				t.Skipf("Skip %s beause of path based block", newName)
			***REMOVED***)
			continue
		***REMOVED***
		if file.IsDir() ***REMOVED***
			ctx.runTC39Tests(newName)
		***REMOVED*** else if strings.HasSuffix(file.Name(), ".js") && !strings.HasSuffix(file.Name(), "_FIXTURE.js") ***REMOVED***
			ctx.runTest(newName, func(t *testing.T) ***REMOVED***
				ctx.runTC39File(newName, t)
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTC39(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip()
	***REMOVED***

	if _, err := os.Stat(tc39BASE); err != nil ***REMOVED***
		t.Skipf("If you want to run tc39 tests, you need to run the 'checkout.sh` script in the directory to get  https://github.com/tc39/test262 at the correct last tested commit (%v)", err)
	***REMOVED***

	ctx := &tc39TestCtx***REMOVED***
		base:     tc39BASE,
		compiler: compiler.New(testutils.NewLogger(t)),
	***REMOVED***
	ctx.init()
	// ctx.enableBench = true

	t.Run("test262", func(t *testing.T) ***REMOVED***
		ctx.t = t
		ctx.runTC39Tests("test/language")
		ctx.runTC39Tests("test/built-ins")
		ctx.runTC39Tests("test/harness")
		ctx.runTC39Tests("test/annexB/built-ins")

		ctx.flush()
	***REMOVED***)

	if ctx.enableBench ***REMOVED***
		sort.Slice(ctx.benchmark, func(i, j int) bool ***REMOVED***
			return ctx.benchmark[i].duration > ctx.benchmark[j].duration
		***REMOVED***)
		bench := ctx.benchmark
		if len(bench) > 50 ***REMOVED***
			bench = bench[:50]
		***REMOVED***
		for _, item := range bench ***REMOVED***
			fmt.Printf("%s\t%d\n", item.name, item.duration/time.Millisecond)
		***REMOVED***
	***REMOVED***
	if len(ctx.errors) > 0 ***REMOVED***
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(ctx.errors)
	***REMOVED***
***REMOVED***
