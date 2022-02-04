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
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/loader"
	"go.k6.io/k6/stats"
)

func TestConsoleContext(t *testing.T) ***REMOVED***
	t.Parallel()
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	logger, hook := logtest.NewNullLogger()
	_ = rt.Set("console", &console***REMOVED***logger***REMOVED***)

	_, err := rt.RunString(`console.log("a")`)
	assert.NoError(t, err)
	if entry := hook.LastEntry(); assert.NotNil(t, entry) ***REMOVED***
		assert.Equal(t, "a", entry.Message)
	***REMOVED***

	_, err = rt.RunString(`console.log("b")`)
	assert.NoError(t, err)
	if entry := hook.LastEntry(); assert.NotNil(t, entry) ***REMOVED***
		assert.Equal(t, "b", entry.Message)
	***REMOVED***
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
		logger,
		&loader.SourceData***REMOVED***
			URL:  &url.URL***REMOVED***Path: filename, Scheme: "file"***REMOVED***,
			Data: []byte(data),
		***REMOVED***,
		map[string]afero.Fs***REMOVED***"file": fs, "https": afero.NewMemMapFs()***REMOVED***,
		rtOpts,
		builtinMetrics,
		registry,
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

func TestConsole(t *testing.T) ***REMOVED***
	t.Parallel()
	levels := map[string]logrus.Level***REMOVED***
		"log":   logrus.InfoLevel,
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
		"error": logrus.ErrorLevel,
	***REMOVED***
	argsets := map[string]struct ***REMOVED***
		Message string
		Data    logrus.Fields
	***REMOVED******REMOVED***
		`"string"`:         ***REMOVED***Message: "string", Data: logrus.Fields***REMOVED***"source": "console"***REMOVED******REMOVED***,
		`"string","a","b"`: ***REMOVED***Message: "string a b", Data: logrus.Fields***REMOVED***"source": "console"***REMOVED******REMOVED***,
		`"string",1,2`:     ***REMOVED***Message: "string 1 2", Data: logrus.Fields***REMOVED***"source": "console"***REMOVED******REMOVED***,
		`***REMOVED******REMOVED***`:               ***REMOVED***Message: "[object Object]", Data: logrus.Fields***REMOVED***"source": "console"***REMOVED******REMOVED***,
	***REMOVED***
	for name, level := range levels ***REMOVED***
		name, level := name, level
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			for args, result := range argsets ***REMOVED***
				args, result := args, result
				t.Run(args, func(t *testing.T) ***REMOVED***
					t.Parallel()
					r, err := getSimpleRunner(t, "/script.js", fmt.Sprintf(
						`exports.default = function() ***REMOVED*** console.%s(%s); ***REMOVED***`,
						name, args,
					))
					assert.NoError(t, err)

					samples := make(chan stats.SampleContainer, 100)
					initVU, err := r.newVU(1, 1, samples)
					assert.NoError(t, err)

					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)

					logger := extractLogger(vu.(*ActiveVU).Console.logger)

					logger.Out = ioutil.Discard
					logger.Level = logrus.DebugLevel
					hook := logtest.NewLocal(logger)

					err = vu.RunOnce()
					assert.NoError(t, err)

					entry := hook.LastEntry()
					if assert.NotNil(t, entry, "nothing logged") ***REMOVED***
						assert.Equal(t, level, entry.Level)
						assert.Equal(t, result.Message, entry.Message)

						data := result.Data
						if data == nil ***REMOVED***
							data = make(logrus.Fields)
						***REMOVED***
						assert.Equal(t, data, entry.Data)
					***REMOVED***
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
			`***REMOVED******REMOVED***`:               ***REMOVED***Message: "[object Object]", Data: logrus.Fields***REMOVED******REMOVED******REMOVED***,
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
							assert.NoError(t, err)

							err = r.SetOptions(lib.Options***REMOVED***
								ConsoleOutput: null.StringFrom(logFilename),
							***REMOVED***)
							assert.NoError(t, err)

							samples := make(chan stats.SampleContainer, 100)
							initVU, err := r.newVU(1, 1, samples)
							assert.NoError(t, err)

							ctx, cancel := context.WithCancel(context.Background())
							defer cancel()
							vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
							logger := extractLogger(vu.(*ActiveVU).Console.logger)

							logger.Level = logrus.DebugLevel
							hook := logtest.NewLocal(logger)

							err = vu.RunOnce()
							assert.NoError(t, err)

							// Test if the file was created.
							_, err = os.Stat(logFilename)
							assert.NoError(t, err)

							entry := hook.LastEntry()
							if assert.NotNil(t, entry, "nothing logged") ***REMOVED***
								assert.Equal(t, level, entry.Level)
								assert.Equal(t, result.Message, entry.Message)

								data := result.Data
								if data == nil ***REMOVED***
									data = make(logrus.Fields)
								***REMOVED***
								assert.Equal(t, data, entry.Data)

								// Test if what we logged to the hook is the same as what we logged
								// to the file.
								entryStr, err := entry.String()
								assert.NoError(t, err)

								f, err := os.Open(logFilename)
								assert.NoError(t, err)

								fileContent, err := ioutil.ReadAll(f)
								assert.NoError(t, err)

								expectedStr := entryStr
								if !deleteFile ***REMOVED***
									expectedStr = preExistingText + expectedStr
								***REMOVED***
								assert.Equal(t, expectedStr, string(fileContent))
							***REMOVED***
						***REMOVED***)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
