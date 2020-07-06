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

	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/loader"
	"github.com/loadimpact/k6/stats"
)

func TestConsoleContext(t *testing.T) ***REMOVED***
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	ctxPtr := new(context.Context)
	logger, hook := logtest.NewNullLogger()
	rt.Set("console", common.Bind(rt, &console***REMOVED***logger***REMOVED***, ctxPtr))

	_, err := common.RunString(rt, `console.log("a")`)
	assert.NoError(t, err)
	if entry := hook.LastEntry(); assert.NotNil(t, entry) ***REMOVED***
		assert.Equal(t, "a", entry.Message)
	***REMOVED***

	ctx, cancel := context.WithCancel(context.Background())
	*ctxPtr = ctx
	_, err = common.RunString(rt, `console.log("b")`)
	assert.NoError(t, err)
	if entry := hook.LastEntry(); assert.NotNil(t, entry) ***REMOVED***
		assert.Equal(t, "b", entry.Message)
	***REMOVED***

	cancel()
	_, err = common.RunString(rt, `console.log("c")`)
	assert.NoError(t, err)
	if entry := hook.LastEntry(); assert.NotNil(t, entry) ***REMOVED***
		assert.Equal(t, "b", entry.Message)
	***REMOVED***
***REMOVED***
func getSimpleRunner(filename, data string, opts ...interface***REMOVED******REMOVED***) (*Runner, error) ***REMOVED***
	var (
		fs     = afero.NewMemMapFs()
		rtOpts = lib.RuntimeOptions***REMOVED***CompatibilityMode: null.NewString("base", true)***REMOVED***
	)
	for _, o := range opts ***REMOVED***
		switch opt := o.(type) ***REMOVED***
		case afero.Fs:
			fs = opt
		case lib.RuntimeOptions:
			rtOpts = opt
		***REMOVED***
	***REMOVED***
	return New(
		&loader.SourceData***REMOVED***
			URL:  &url.URL***REMOVED***Path: filename, Scheme: "file"***REMOVED***,
			Data: []byte(data),
		***REMOVED***,
		map[string]afero.Fs***REMOVED***"file": fs, "https": afero.NewMemMapFs()***REMOVED***,
		rtOpts,
	)
***REMOVED***

func TestConsole(t *testing.T) ***REMOVED***
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
		`"string"`:         ***REMOVED***Message: "string"***REMOVED***,
		`"string","a","b"`: ***REMOVED***Message: "string", Data: logrus.Fields***REMOVED***"0": "a", "1": "b"***REMOVED******REMOVED***,
		`"string",1,2`:     ***REMOVED***Message: "string", Data: logrus.Fields***REMOVED***"0": "1", "1": "2"***REMOVED******REMOVED***,
		`***REMOVED******REMOVED***`:               ***REMOVED***Message: "[object Object]"***REMOVED***,
	***REMOVED***
	for name, level := range levels ***REMOVED***
		name, level := name, level
		t.Run(name, func(t *testing.T) ***REMOVED***
			for args, result := range argsets ***REMOVED***
				args, result := args, result
				t.Run(args, func(t *testing.T) ***REMOVED***
					r, err := getSimpleRunner("/script.js", fmt.Sprintf(
						`exports.default = function() ***REMOVED*** console.%s(%s); ***REMOVED***`,
						name, args,
					))
					assert.NoError(t, err)

					samples := make(chan stats.SampleContainer, 100)
					initVU, err := r.newVU(1, samples)
					assert.NoError(t, err)

					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)

					logger, hook := logtest.NewNullLogger()
					logger.Level = logrus.DebugLevel
					jsVU := vu.(*ActiveVU)
					jsVU.Console.Logger = logger

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
			`"string"`:         ***REMOVED***Message: "string"***REMOVED***,
			`"string","a","b"`: ***REMOVED***Message: "string", Data: logrus.Fields***REMOVED***"0": "a", "1": "b"***REMOVED******REMOVED***,
			`"string",1,2`:     ***REMOVED***Message: "string", Data: logrus.Fields***REMOVED***"0": "1", "1": "2"***REMOVED******REMOVED***,
			`***REMOVED******REMOVED***`:               ***REMOVED***Message: "[object Object]"***REMOVED***,
		***REMOVED***
		preExisting = map[string]bool***REMOVED***
			"log exists":        false,
			"log doesn't exist": true,
		***REMOVED***
		preExistingText = "Prexisting file\n"
	)
	for name, level := range levels ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			for args, result := range argsets ***REMOVED***
				t.Run(args, func(t *testing.T) ***REMOVED***
					// whether the file is existed before logging
					for msg, deleteFile := range preExisting ***REMOVED***
						t.Run(msg, func(t *testing.T) ***REMOVED***
							var f, err = ioutil.TempFile("", "")
							if err != nil ***REMOVED***
								t.Fatalf("Couldn't create temporary file for testing: %s", err)
							***REMOVED***
							var logFilename = f.Name()
							defer os.Remove(logFilename)
							// close it as we will want to reopen it and maybe remove it
							if deleteFile ***REMOVED***
								f.Close()
								if err := os.Remove(logFilename); err != nil ***REMOVED***
									t.Fatalf("Couldn't remove tempfile: %s", err)
								***REMOVED***
							***REMOVED*** else ***REMOVED***
								// TODO: handle case where the string was no written in full ?
								_, err := f.WriteString(preExistingText)
								f.Close()
								if err != nil ***REMOVED***
									t.Fatalf("Error while writing text to preexisting logfile: %s", err)
								***REMOVED***

							***REMOVED***
							r, err := getSimpleRunner("/script",
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
							initVU, err := r.newVU(1, samples)
							assert.NoError(t, err)

							ctx, cancel := context.WithCancel(context.Background())
							defer cancel()
							vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
							jsVU := vu.(*ActiveVU)
							jsVU.Console.Logger.Level = logrus.DebugLevel
							hook := logtest.NewLocal(jsVU.Console.Logger)

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

								var expectedStr = entryStr
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
