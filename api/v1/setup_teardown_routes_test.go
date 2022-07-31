/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2018 Load Impact
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

package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/core"
	"go.k6.io/k6/core/local"
	"go.k6.io/k6/js"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/loader"
)

func TestSetupData(t *testing.T) ***REMOVED***
	t.Parallel()
	testCases := []struct ***REMOVED***
		name      string
		script    []byte
		setupRuns [][3]string
	***REMOVED******REMOVED***
		***REMOVED***
			name: "setupReturns",
			script: []byte(`
			export function setup() ***REMOVED***
				return ***REMOVED***"v": 1***REMOVED***;
			***REMOVED***

			export default function(data) ***REMOVED***
				if (data !== undefined) ***REMOVED***
					throw new Error("incorrect data: " + JSON.stringify(data));
				***REMOVED***
			***REMOVED***;

			export function teardown(data) ***REMOVED***
				if (data !== undefined) ***REMOVED***
					throw new Error("incorrect teardown data: " + JSON.stringify(data));
				***REMOVED***
			***REMOVED*** `),
			setupRuns: [][3]string***REMOVED***
				***REMOVED***"GET", "", "***REMOVED******REMOVED***"***REMOVED***,
				***REMOVED***"POST", "", `***REMOVED***"data": ***REMOVED***"v":1***REMOVED******REMOVED***`***REMOVED***,
				***REMOVED***"GET", "", `***REMOVED***"data": ***REMOVED***"v":1***REMOVED******REMOVED***`***REMOVED***,
				***REMOVED***"PUT", `***REMOVED***"v":2, "test":"mest"***REMOVED***`, `***REMOVED***"data": ***REMOVED***"v":2, "test":"mest"***REMOVED******REMOVED***`***REMOVED***,
				***REMOVED***"GET", "", `***REMOVED***"data": ***REMOVED***"v":2, "test":"mest"***REMOVED******REMOVED***`***REMOVED***,
				***REMOVED***"PUT", "", `***REMOVED******REMOVED***`***REMOVED***,
				***REMOVED***"GET", "", `***REMOVED******REMOVED***`***REMOVED***,
			***REMOVED***,
		***REMOVED***, ***REMOVED***

			name: "noSetup",
			script: []byte(`
			export default function(data) ***REMOVED***
				if (!data || data.v != 2) ***REMOVED***
					throw new Error("incorrect data: " + JSON.stringify(data));
				***REMOVED***
			***REMOVED***;

			export function teardown(data) ***REMOVED***
				if (!data || data.v != 2) ***REMOVED***
					throw new Error("incorrect teardown data: " + JSON.stringify(data));
				***REMOVED***
			***REMOVED*** `),
			setupRuns: [][3]string***REMOVED***
				***REMOVED***"GET", "", "***REMOVED******REMOVED***"***REMOVED***,
				***REMOVED***"POST", "", `***REMOVED******REMOVED***`***REMOVED***,
				***REMOVED***"GET", "", `***REMOVED******REMOVED***`***REMOVED***,
				***REMOVED***"PUT", `***REMOVED***"v":2, "test":"mest"***REMOVED***`, `***REMOVED***"data": ***REMOVED***"v":2, "test":"mest"***REMOVED******REMOVED***`***REMOVED***,
				***REMOVED***"GET", "", `***REMOVED***"data": ***REMOVED***"v":2, "test":"mest"***REMOVED******REMOVED***`***REMOVED***,
				***REMOVED***"PUT", "", `***REMOVED******REMOVED***`***REMOVED***,
				***REMOVED***"GET", "", `***REMOVED******REMOVED***`***REMOVED***,
				***REMOVED***"PUT", `***REMOVED***"v":2, "test":"mest"***REMOVED***`, `***REMOVED***"data": ***REMOVED***"v":2, "test":"mest"***REMOVED******REMOVED***`***REMOVED***,
				***REMOVED***"GET", "", `***REMOVED***"data": ***REMOVED***"v":2, "test":"mest"***REMOVED******REMOVED***`***REMOVED***,
			***REMOVED***,
		***REMOVED***, ***REMOVED***
			name: "setupNoReturn",
			script: []byte(`
			export function setup() ***REMOVED***
				let a = ***REMOVED***"v": 1***REMOVED***;
			***REMOVED***
			export default function(data) ***REMOVED***
				if (data === undefined || data !== "") ***REMOVED***
					throw new Error("incorrect data: " + JSON.stringify(data));
				***REMOVED***
			***REMOVED***;

			export function teardown(data) ***REMOVED***
				if (data === undefined || data !== "") ***REMOVED***
					throw new Error("incorrect teardown data: " + JSON.stringify(data));
				***REMOVED***
			***REMOVED*** `),
			setupRuns: [][3]string***REMOVED***
				***REMOVED***"GET", "", "***REMOVED******REMOVED***"***REMOVED***,
				***REMOVED***"POST", "", `***REMOVED******REMOVED***`***REMOVED***,
				***REMOVED***"GET", "", `***REMOVED******REMOVED***`***REMOVED***,
				***REMOVED***"PUT", `***REMOVED***"v":2, "test":"mest"***REMOVED***`, `***REMOVED***"data": ***REMOVED***"v":2, "test":"mest"***REMOVED******REMOVED***`***REMOVED***,
				***REMOVED***"GET", "", `***REMOVED***"data": ***REMOVED***"v":2, "test":"mest"***REMOVED******REMOVED***`***REMOVED***,
				***REMOVED***"PUT", "\"\"", `***REMOVED***"data": ""***REMOVED***`***REMOVED***,
				***REMOVED***"GET", "", `***REMOVED***"data": ""***REMOVED***`***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	runTestCase := func(t *testing.T, tcid int) ***REMOVED***
		testCase := testCases[tcid]
		t.Run(testCase.name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			piState := getTestPreInitState(t)
			runner, err := js.New(
				piState, &loader.SourceData***REMOVED***URL: &url.URL***REMOVED***Path: "/script.js"***REMOVED***, Data: testCase.script***REMOVED***, nil,
			)
			require.NoError(t, err)
			runner.SetOptions(lib.Options***REMOVED***
				Paused:          null.BoolFrom(true),
				VUs:             null.IntFrom(2),
				Iterations:      null.IntFrom(3),
				NoSetup:         null.BoolFrom(true),
				SetupTimeout:    types.NullDurationFrom(5 * time.Second),
				TeardownTimeout: types.NullDurationFrom(5 * time.Second),
			***REMOVED***)
			execScheduler, err := local.NewExecutionScheduler(runner, piState)
			require.NoError(t, err)
			engine, err := core.NewEngine(
				execScheduler, runner.GetOptions(), piState.RuntimeOptions, nil, piState.Logger, piState.Registry,
			)
			require.NoError(t, err)

			require.NoError(t, engine.OutputManager.StartOutputs())
			defer engine.OutputManager.StopOutputs()

			globalCtx, globalCancel := context.WithCancel(context.Background())
			runCtx, runCancel := context.WithCancel(globalCtx)
			run, wait, err := engine.Init(globalCtx, runCtx)
			require.NoError(t, err)

			defer wait()
			defer globalCancel()

			errC := make(chan error)
			go func() ***REMOVED*** errC <- run() ***REMOVED***()

			handler := NewHandler()

			checkSetup := func(method, body, expResult string) ***REMOVED***
				rw := httptest.NewRecorder()
				handler.ServeHTTP(rw, newRequestWithEngine(engine, method, "/v1/setup", bytes.NewBufferString(body)))
				res := rw.Result()
				if !assert.Equal(t, http.StatusOK, res.StatusCode) ***REMOVED***
					t.Logf("body: %s\n", rw.Body.String())
					return
				***REMOVED***

				var doc setUpJSONAPI
				assert.NoError(t, json.Unmarshal(rw.Body.Bytes(), &doc))
				assert.Equal(t, "setupData", doc.Data.Type)

				encoded, err := json.Marshal(doc.Data.Attributes)
				assert.NoError(t, err)
				assert.JSONEq(t, expResult, string(encoded))
			***REMOVED***

			for _, setupRun := range testCase.setupRuns ***REMOVED***
				checkSetup(setupRun[0], setupRun[1], setupRun[2])
			***REMOVED***

			require.NoError(t, engine.ExecutionScheduler.SetPaused(false))

			select ***REMOVED***
			case <-time.After(10 * time.Second):
				runCancel()
				t.Fatal("Test timed out")
			case err := <-errC:
				runCancel()
				require.NoError(t, err)
			***REMOVED***
		***REMOVED***)
	***REMOVED***

	for id := range testCases ***REMOVED***
		id := id
		t.Run(fmt.Sprintf("testcase_%d", id), func(t *testing.T) ***REMOVED***
			t.Parallel()
			runTestCase(t, id)
		***REMOVED***)
	***REMOVED***
***REMOVED***
