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
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/core"
	"github.com/loadimpact/k6/core/local"
	"github.com/loadimpact/k6/js"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/testutils"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/loader"
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
	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))
	for _, testCase := range testCases ***REMOVED***
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) ***REMOVED***
			runner, err := js.New(
				logger,
				&loader.SourceData***REMOVED***URL: &url.URL***REMOVED***Path: "/script.js"***REMOVED***, Data: testCase.script***REMOVED***,
				nil,
				lib.RuntimeOptions***REMOVED******REMOVED***,
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
			execScheduler, err := local.NewExecutionScheduler(runner, logger)
			require.NoError(t, err)
			engine, err := core.NewEngine(execScheduler, runner.GetOptions(), lib.RuntimeOptions***REMOVED******REMOVED***, nil, logger)
			require.NoError(t, err)

			globalCtx, globalCancel := context.WithCancel(context.Background())
			runCtx, runCancel := context.WithCancel(globalCtx)
			run, wait, err := engine.Init(globalCtx, runCtx)
			defer wait()
			defer globalCancel()

			require.NoError(t, err)

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

				var doc jsonapi.Document
				assert.NoError(t, json.Unmarshal(rw.Body.Bytes(), &doc))
				require.NotNil(t, doc.Data)
				require.NotNil(t, doc.Data.DataObject)
				assert.Equal(t, "setupData", doc.Data.DataObject.Type)
				assert.JSONEq(t, expResult, string(doc.Data.DataObject.Attributes))
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
***REMOVED***
