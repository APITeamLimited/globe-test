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

package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/core"
	"go.k6.io/k6/core/local"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/lib/testutils/minirunner"
)

func TestGetStatus(t *testing.T) ***REMOVED***
	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))
	execScheduler, err := local.NewExecutionScheduler(&minirunner.MiniRunner***REMOVED******REMOVED***, logger)
	require.NoError(t, err)
	engine, err := core.NewEngine(execScheduler, lib.Options***REMOVED******REMOVED***, lib.RuntimeOptions***REMOVED******REMOVED***, nil, logger)
	require.NoError(t, err)

	rw := httptest.NewRecorder()
	NewHandler().ServeHTTP(rw, newRequestWithEngine(engine, "GET", "/v1/status", nil))
	res := rw.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	t.Run("document", func(t *testing.T) ***REMOVED***
		var doc jsonapi.Document
		assert.NoError(t, json.Unmarshal(rw.Body.Bytes(), &doc))
		if !assert.NotNil(t, doc.Data.DataObject) ***REMOVED***
			return
		***REMOVED***
		assert.Equal(t, "status", doc.Data.DataObject.Type)
	***REMOVED***)

	t.Run("status", func(t *testing.T) ***REMOVED***
		var status Status
		assert.NoError(t, jsonapi.Unmarshal(rw.Body.Bytes(), &status))
		assert.True(t, status.Paused.Valid)
		assert.True(t, status.VUs.Valid)
		assert.True(t, status.VUsMax.Valid)
		assert.False(t, status.Stopped)
		assert.False(t, status.Tainted)
	***REMOVED***)
***REMOVED***

func TestPatchStatus(t *testing.T) ***REMOVED***
	testdata := map[string]struct ***REMOVED***
		StatusCode int
		Status     Status
	***REMOVED******REMOVED***
		"nothing":               ***REMOVED***200, Status***REMOVED******REMOVED******REMOVED***,
		"paused":                ***REMOVED***200, Status***REMOVED***Paused: null.BoolFrom(true)***REMOVED******REMOVED***,
		"max vus":               ***REMOVED***200, Status***REMOVED***VUsMax: null.IntFrom(20)***REMOVED******REMOVED***,
		"max vus below initial": ***REMOVED***400, Status***REMOVED***VUsMax: null.IntFrom(5)***REMOVED******REMOVED***,
		"too many vus":          ***REMOVED***400, Status***REMOVED***VUs: null.IntFrom(10), VUsMax: null.IntFrom(0)***REMOVED******REMOVED***,
		"vus":                   ***REMOVED***200, Status***REMOVED***VUs: null.IntFrom(10), VUsMax: null.IntFrom(10)***REMOVED******REMOVED***,
	***REMOVED***
	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))

	scenarios := lib.ScenarioConfigs***REMOVED******REMOVED***
	err := json.Unmarshal([]byte(`
			***REMOVED***"external": ***REMOVED***"executor": "externally-controlled",
			"vus": 0, "maxVUs": 10, "duration": "1s"***REMOVED******REMOVED***`), &scenarios)
	require.NoError(t, err)
	options := lib.Options***REMOVED***Scenarios: scenarios***REMOVED***

	for name, indata := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			execScheduler, err := local.NewExecutionScheduler(&minirunner.MiniRunner***REMOVED***Options: options***REMOVED***, logger)
			require.NoError(t, err)
			engine, err := core.NewEngine(execScheduler, options, lib.RuntimeOptions***REMOVED******REMOVED***, nil, logger)
			require.NoError(t, err)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			run, _, err := engine.Init(ctx, ctx)
			require.NoError(t, err)

			go func() ***REMOVED*** _ = run() ***REMOVED***()
			// wait for the executor to initialize to avoid a potential data race below
			time.Sleep(100 * time.Millisecond)

			body, err := jsonapi.Marshal(indata.Status)
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			rw := httptest.NewRecorder()
			NewHandler().ServeHTTP(rw, newRequestWithEngine(engine, "PATCH", "/v1/status", bytes.NewReader(body)))
			res := rw.Result()

			if !assert.Equal(t, indata.StatusCode, res.StatusCode) ***REMOVED***
				return
			***REMOVED***
			if indata.StatusCode != 200 ***REMOVED***
				return
			***REMOVED***

			status := NewStatus(engine)
			if indata.Status.Paused.Valid ***REMOVED***
				assert.Equal(t, indata.Status.Paused, status.Paused)
			***REMOVED***
			if indata.Status.VUs.Valid ***REMOVED***
				assert.Equal(t, indata.Status.VUs, status.VUs)
			***REMOVED***
			if indata.Status.VUsMax.Valid ***REMOVED***
				assert.Equal(t, indata.Status.VUsMax, status.VUsMax)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
