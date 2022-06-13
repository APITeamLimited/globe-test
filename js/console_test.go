package js

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/loader"
	"go.k6.io/k6/metrics"
)

func TestConsoleContext(t *testing.T) ***REMOVED***
	t.Parallel()
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	logger, hook := logtest.NewNullLogger()
	_ = rt.Set("console", &console***REMOVED***logger***REMOVED***)

	_, err := rt.RunString(`console.log("a")`)
	require.NoError(t, err)
	entry := hook.LastEntry()
	require.NotNil(t, entry)
	assert.Equal(t, "a", entry.Message)

	_, err = rt.RunString(`console.log("b")`)
	require.NoError(t, err)
	entry = hook.LastEntry()
	require.NotNil(t, entry)
	require.Equal(t, "b", entry.Message)
***REMOVED***

func getSimpleRunner(tb testing.TB, filename, data string, opts ...interface***REMOVED******REMOVED***) (*Runner, error) ***REMOVED***
	var (
		fs     = afero.NewMemMapFs()
		rtOpts = lib.RuntimeOptions***REMOVED***CompatibilityMode: null.NewString("base", true)***REMOVED***
		logger = testutils.NewLogger(tb)
	)
	for _, o := range opts ***REMOVED***
		switch opt := o.(type) ***REMOVED***
		case afero.Fs:
			fs = opt
		case lib.RuntimeOptions:
			rtOpts = opt
		case *logrus.Logger:
			logger = opt
		***REMOVED***
	***REMOVED***
	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	return New(
		&lib.RuntimeState***REMOVED***
			Logger:         logger,
			RuntimeOptions: rtOpts,
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***,
		&loader.SourceData***REMOVED***
			URL:  &url.URL***REMOVED***Path: filename, Scheme: "file"***REMOVED***,
			Data: []byte(data),
		***REMOVED***,
		map[string]afero.Fs***REMOVED***"file": fs, "https": afero.NewMemMapFs()***REMOVED***,
	)
***REMOVED***

func extractLogger(fl logrus.FieldLogger) *logrus.Logger ***REMOVED***
	switch e := fl.(type) ***REMOVED***
	case *logrus.Entry:
		return e.Logger
	case *logrus.Logger:
		return e
	***REMOVED***
	return nil
***REMOVED***

func TestConsoleLogWithGojaNativeObject(t *testing.T) ***REMOVED***
	t.Parallel()

	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	obj := rt.NewObject()
	err := obj.Set("text", "nativeObject")
	require.NoError(t, err)

	logger := testutils.NewLogger(t)
	hook := logtest.NewLocal(logger)

	c := newConsole(logger)
	c.Log(obj)

	entry := hook.LastEntry()
	require.NotNil(t, entry, "nothing logged")
	require.JSONEq(t, `***REMOVED***"text":"nativeObject"***REMOVED***`, entry.Message)
***REMOVED***

func TestConsoleLogObjectsWithGoTypes(t *testing.T) ***REMOVED***
	t.Parallel()

	type value struct ***REMOVED***
		Text string
	***REMOVED***

	tests := []struct ***REMOVED***
		name string
		in   interface***REMOVED******REMOVED***
		exp  string
	***REMOVED******REMOVED***
		***REMOVED***
			name: "StructLiteral",
			in: value***REMOVED***
				Text: "test1",
			***REMOVED***,
			exp: `***REMOVED***"text":"test1"***REMOVED***`,
		***REMOVED***,
		***REMOVED***
			name: "StructPointer",
			in: &value***REMOVED***
				Text: "test2",
			***REMOVED***,
			exp: `***REMOVED***"text":"test2"***REMOVED***`,
		***REMOVED***,
		***REMOVED***
			name: "Map",
			in: map[string]interface***REMOVED******REMOVED******REMOVED***
				"text": "test3",
			***REMOVED***,
			exp: `***REMOVED***"text":"test3"***REMOVED***`,
		***REMOVED***,
	***REMOVED***

	expFields := logrus.Fields***REMOVED***"source": "console"***REMOVED***
	for _, tt := range tests ***REMOVED***
		tt := tt

		t.Run(tt.name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			rt := goja.New()
			rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
			obj := rt.ToValue(tt.in)

			logger := testutils.NewLogger(t)
			hook := logtest.NewLocal(logger)

			c := newConsole(logger)
			c.Log(obj)

			entry := hook.LastEntry()
			require.NotNil(t, entry, "nothing logged")
			assert.JSONEq(t, tt.exp, entry.Message)
			assert.Equal(t, expFields, entry.Data)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestConsoleLog(t *testing.T) ***REMOVED***
	t.Parallel()

	tests := []struct ***REMOVED***
		in       string
		expected string
	***REMOVED******REMOVED***
		***REMOVED***``, ``***REMOVED***,
		***REMOVED***`""`, ``***REMOVED***,
		***REMOVED***`undefined`, `undefined`***REMOVED***,
		***REMOVED***`null`, `null`***REMOVED***,

		***REMOVED***in: `"string"`, expected: "string"***REMOVED***,
		***REMOVED***in: `"string","a","b"`, expected: "string a b"***REMOVED***,
		***REMOVED***in: `"string",1,2`, expected: "string 1 2"***REMOVED***,

		***REMOVED***in: `["bar", 1, 2]`, expected: `["bar",1,2]`***REMOVED***,
		***REMOVED***in: `"bar", ["bar", 0x01, 2], 1, 2`, expected: `bar ["bar",1,2] 1 2`***REMOVED***,

		***REMOVED***in: `***REMOVED******REMOVED***`, expected: "***REMOVED******REMOVED***"***REMOVED***,
		***REMOVED***in: `***REMOVED***foo:"bar"***REMOVED***`, expected: `***REMOVED***"foo":"bar"***REMOVED***`***REMOVED***,
		***REMOVED***in: `["test1", 2]`, expected: `["test1",2]`***REMOVED***,

		// TODO: the ideal output for a circular object should be like `***REMOVED***a: [Circular]***REMOVED***`
		***REMOVED***in: `function() ***REMOVED***var a = ***REMOVED***foo: ***REMOVED******REMOVED******REMOVED***; a.foo = a; return a***REMOVED***()`, expected: "[object Object]"***REMOVED***,
	***REMOVED***

	for i, tt := range tests ***REMOVED***
		tt := tt
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) ***REMOVED***
			t.Parallel()

			r, err := getSimpleRunner(t, "/script.js", fmt.Sprintf(
				`exports.default = function() ***REMOVED*** console.log(%s); ***REMOVED***`, tt.in))
			require.NoError(t, err)

			samples := make(chan metrics.SampleContainer, 100)
			initVU, err := r.newVU(1, 1, samples)
			require.NoError(t, err)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)

			logger := extractLogger(vu.(*ActiveVU).Console.logger)

			logger.Out = ioutil.Discard
			logger.Level = logrus.DebugLevel
			hook := logtest.NewLocal(logger)

			err = vu.RunOnce()
			require.NoError(t, err)

			entry := hook.LastEntry()

			require.NotNil(t, entry, "nothing logged")
			assert.Equal(t, tt.expected, entry.Message)
			assert.Equal(t, logrus.Fields***REMOVED***"source": "console"***REMOVED***, entry.Data)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestConsoleLevels(t *testing.T) ***REMOVED***
	t.Parallel()
	levels := map[string]logrus.Level***REMOVED***
		"log":   logrus.InfoLevel,
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
		"error": logrus.ErrorLevel,
	***REMOVED***
	argsets := []struct ***REMOVED***
		in  string
		exp string
	***REMOVED******REMOVED***
		***REMOVED***in: `"string"`, exp: "string"***REMOVED***,
		***REMOVED***in: `***REMOVED******REMOVED***`, exp: "***REMOVED******REMOVED***"***REMOVED***,
		***REMOVED***in: `***REMOVED***foo:"bar"***REMOVED***`, exp: `***REMOVED***"foo":"bar"***REMOVED***`***REMOVED***,
	***REMOVED***
	for name, level := range levels ***REMOVED***
		name, level := name, level
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			for _, tt := range argsets ***REMOVED***
				args, result := tt.in, tt.exp
				t.Run(args, func(t *testing.T) ***REMOVED***
					t.Parallel()
					r, err := getSimpleRunner(t, "/script.js", fmt.Sprintf(
						`exports.default = function() ***REMOVED*** console.%s(%s); ***REMOVED***`,
						name, args,
					))
					require.NoError(t, err)

					samples := make(chan metrics.SampleContainer, 100)
					initVU, err := r.newVU(1, 1, samples)
					require.NoError(t, err)

					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)

					logger := extractLogger(vu.(*ActiveVU).Console.logger)

					logger.Out = ioutil.Discard
					logger.Level = logrus.DebugLevel
					hook := logtest.NewLocal(logger)

					err = vu.RunOnce()
					require.NoError(t, err)

					entry := hook.LastEntry()
					require.NotNil(t, entry, "nothing logged")

					assert.Equal(t, level, entry.Level)
					assert.Equal(t, result, entry.Message)
					assert.Equal(t, logrus.Fields***REMOVED***"source": "console"***REMOVED***, entry.Data)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestFileConsole(t *testing.T) ***REMOVED***
	t.Parallel()
	var (
		levels = map[string]logrus.Level***REMOVED***
			"log":   logrus.InfoLevel,
			"debug": logrus.DebugLevel,
			"info":  logrus.InfoLevel,
			"warn":  logrus.WarnLevel,
			"error": logrus.ErrorLevel,
		***REMOVED***
		argsets = map[string]struct ***REMOVED***
			Message string
			Data    logrus.Fields
		***REMOVED******REMOVED***
			`"string"`:         ***REMOVED***Message: "string", Data: logrus.Fields***REMOVED******REMOVED******REMOVED***,
			`"string","a","b"`: ***REMOVED***Message: "string a b", Data: logrus.Fields***REMOVED******REMOVED******REMOVED***,
			`"string",1,2`:     ***REMOVED***Message: "string 1 2", Data: logrus.Fields***REMOVED******REMOVED******REMOVED***,
			`***REMOVED******REMOVED***`:               ***REMOVED***Message: "***REMOVED******REMOVED***", Data: logrus.Fields***REMOVED******REMOVED******REMOVED***,
		***REMOVED***
		preExisting = map[string]bool***REMOVED***
			"log exists":        false,
			"log doesn't exist": true,
		***REMOVED***
		preExistingText = "Prexisting file\n"
	)
	for name, level := range levels ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			for args, result := range argsets ***REMOVED***
				t.Run(args, func(t *testing.T) ***REMOVED***
					t.Parallel()
					// whether the file is existed before logging
					for msg, deleteFile := range preExisting ***REMOVED***
						t.Run(msg, func(t *testing.T) ***REMOVED***
							t.Parallel()
							f, err := ioutil.TempFile("", "")
							if err != nil ***REMOVED***
								t.Fatalf("Couldn't create temporary file for testing: %s", err)
							***REMOVED***
							logFilename := f.Name()
							defer os.Remove(logFilename)
							// close it as we will want to reopen it and maybe remove it
							if deleteFile ***REMOVED***
								f.Close()
								if err := os.Remove(logFilename); err != nil ***REMOVED***
									t.Fatalf("Couldn't remove tempfile: %s", err)
								***REMOVED***
							***REMOVED*** else ***REMOVED***
								// TODO: handle case where the string was no written in full ?
								_, err = f.WriteString(preExistingText)
								_ = f.Close()
								if err != nil ***REMOVED***
									t.Fatalf("Error while writing text to preexisting logfile: %s", err)
								***REMOVED***

							***REMOVED***
							r, err := getSimpleRunner(t, "/script",
								fmt.Sprintf(
									`exports.default = function() ***REMOVED*** console.%s(%s); ***REMOVED***`,
									name, args,
								))
							require.NoError(t, err)

							err = r.SetOptions(lib.Options***REMOVED***
								ConsoleOutput: null.StringFrom(logFilename),
							***REMOVED***)
							require.NoError(t, err)

							samples := make(chan metrics.SampleContainer, 100)
							initVU, err := r.newVU(1, 1, samples)
							require.NoError(t, err)

							ctx, cancel := context.WithCancel(context.Background())
							defer cancel()
							vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
							logger := extractLogger(vu.(*ActiveVU).Console.logger)

							logger.Level = logrus.DebugLevel
							hook := logtest.NewLocal(logger)

							err = vu.RunOnce()
							require.NoError(t, err)

							// Test if the file was created.
							_, err = os.Stat(logFilename)
							require.NoError(t, err)

							entry := hook.LastEntry()
							require.NotNil(t, entry, "nothing logged")
							assert.Equal(t, level, entry.Level)
							assert.Equal(t, result.Message, entry.Message)

							data := result.Data
							if data == nil ***REMOVED***
								data = make(logrus.Fields)
							***REMOVED***
							require.Equal(t, data, entry.Data)

							// Test if what we logged to the hook is the same as what we logged
							// to the file.
							entryStr, err := entry.String()
							require.NoError(t, err)

							f, err = os.Open(logFilename) //nolint:gosec
							require.NoError(t, err)

							fileContent, err := ioutil.ReadAll(f)
							require.NoError(t, err)

							expectedStr := entryStr
							if !deleteFile ***REMOVED***
								expectedStr = preExistingText + expectedStr
							***REMOVED***
							require.Equal(t, expectedStr, string(fileContent))
						***REMOVED***)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
