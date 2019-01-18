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
	"os"
	"testing"

	"gopkg.in/guregu/null.v3"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	log "github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
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

func TestConsole(t *testing.T) ***REMOVED***
	levels := map[string]log.Level***REMOVED***
		"log":   log.InfoLevel,
		"debug": log.DebugLevel,
		"info":  log.InfoLevel,
		"warn":  log.WarnLevel,
		"error": log.ErrorLevel,
	***REMOVED***
	argsets := map[string]struct ***REMOVED***
		Message string
		Data    log.Fields
	***REMOVED******REMOVED***
		`"string"`:         ***REMOVED***Message: "string"***REMOVED***,
		`"string","a","b"`: ***REMOVED***Message: "string", Data: log.Fields***REMOVED***"0": "a", "1": "b"***REMOVED******REMOVED***,
		`"string",1,2`:     ***REMOVED***Message: "string", Data: log.Fields***REMOVED***"0": "1", "1": "2"***REMOVED******REMOVED***,
		`***REMOVED******REMOVED***`:               ***REMOVED***Message: "[object Object]"***REMOVED***,
	***REMOVED***
	for name, level := range levels ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			for args, result := range argsets ***REMOVED***
				t.Run(args, func(t *testing.T) ***REMOVED***
					r, err := New(&lib.SourceData***REMOVED***
						Filename: "/script",
						Data: []byte(fmt.Sprintf(
							`export default function() ***REMOVED*** console.%s(%s); ***REMOVED***`,
							name, args,
						)),
					***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
					assert.NoError(t, err)

					samples := make(chan stats.SampleContainer, 100)
					vu, err := r.newVU(samples)
					assert.NoError(t, err)

					logger, hook := logtest.NewNullLogger()
					logger.Level = log.DebugLevel
					vu.Console.Logger = logger

					err = vu.RunOnce(context.Background())
					assert.NoError(t, err)

					entry := hook.LastEntry()
					if assert.NotNil(t, entry, "nothing logged") ***REMOVED***
						assert.Equal(t, level, entry.Level)
						assert.Equal(t, result.Message, entry.Message)

						data := result.Data
						if data == nil ***REMOVED***
							data = make(log.Fields)
						***REMOVED***
						assert.Equal(t, data, entry.Data)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestFileConsole(t *testing.T) ***REMOVED***
	logFile := "/tmp/loadtest.log"
	levels := map[string]log.Level***REMOVED***
		"log":   log.InfoLevel,
		"debug": log.DebugLevel,
		"info":  log.InfoLevel,
		"warn":  log.WarnLevel,
		"error": log.ErrorLevel,
	***REMOVED***
	argsets := map[string]struct ***REMOVED***
		Message string
		Data    log.Fields
	***REMOVED******REMOVED***
		`"string"`:         ***REMOVED***Message: "string"***REMOVED***,
		`"string","a","b"`: ***REMOVED***Message: "string", Data: log.Fields***REMOVED***"0": "a", "1": "b"***REMOVED******REMOVED***,
		`"string",1,2`:     ***REMOVED***Message: "string", Data: log.Fields***REMOVED***"0": "1", "1": "2"***REMOVED******REMOVED***,
		`***REMOVED******REMOVED***`:               ***REMOVED***Message: "[object Object]"***REMOVED***,
	***REMOVED***
	for name, level := range levels ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			for args, result := range argsets ***REMOVED***
				t.Run(args, func(t *testing.T) ***REMOVED***
					r, err := New(&lib.SourceData***REMOVED***
						Filename: "/script",
						Data: []byte(fmt.Sprintf(
							`export default function() ***REMOVED*** console.%s(%s); ***REMOVED***`,
							name, args,
						)),
					***REMOVED***, afero.NewMemMapFs(), lib.RuntimeOptions***REMOVED******REMOVED***)
					assert.NoError(t, err)

					err = r.SetOptions(lib.Options***REMOVED***
						ConsoleOutput: null.StringFrom(logFile),
					***REMOVED***)
					assert.NoError(t, err)

					samples := make(chan stats.SampleContainer, 100)
					vu, err := r.newVU(samples)
					assert.NoError(t, err)

					vu.Console.Logger.Level = log.DebugLevel
					hook := logtest.NewLocal(vu.Console.Logger)

					err = vu.RunOnce(context.Background())
					assert.NoError(t, err)

					// Test if the file was created.
					_, err = os.Stat(logFile)
					assert.NoError(t, err)

					entry := hook.LastEntry()
					if assert.NotNil(t, entry, "nothing logged") ***REMOVED***
						assert.Equal(t, level, entry.Level)
						assert.Equal(t, result.Message, entry.Message)

						data := result.Data
						if data == nil ***REMOVED***
							data = make(log.Fields)
						***REMOVED***
						assert.Equal(t, data, entry.Data)

						// Test if what we logged to the hook is the same as what we logged
						// to the file.
						entryStr, err := entry.String()
						assert.NoError(t, err)

						f, err := os.Open(logFile)
						assert.NoError(t, err)

						fileContent, err := ioutil.ReadAll(f)
						assert.NoError(t, err)

						assert.Equal(t, entryStr, string(fileContent))
					***REMOVED***

					os.Remove(logFile)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
