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

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/core"
	"go.k6.io/k6/core/local"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/lib/testutils/minirunner"
)

func TestGetStatus(t *testing.T) ***REMOVED***
	t.Parallel()

	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))
	execScheduler, err := local.NewExecutionScheduler(&minirunner.MiniRunner***REMOVED******REMOVED***, logger)
	require.NoError(t, err)
	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	engine, err := core.NewEngine(execScheduler, lib.Options***REMOVED******REMOVED***, lib.RuntimeOptions***REMOVED******REMOVED***, nil, logger, builtinMetrics)
	require.NoError(t, err)

	rw := httptest.NewRecorder()
	NewHandler().ServeHTTP(rw, newRequestWithEngine(engine, "GET", "/v1/status", nil))
	res := rw.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	t.Run("document", func(t *testing.T) ***REMOVED***
		t.Parallel()

		var doc StatusJSONAPI
		assert.NoError(t, json.Unmarshal(rw.Body.Bytes(), &doc))
		assert.Equal(t, "status", doc.Data.Type)
	***REMOVED***)

	t.Run("status", func(t *testing.T) ***REMOVED***
		t.Parallel()

		var statusEnvelop StatusJSONAPI

		err := json.Unmarshal(rw.Body.Bytes(), &statusEnvelop)
		assert.NoError(t, err)

		status := statusEnvelop.Status()

		assert.True(t, status.Paused.Valid)
		assert.True(t, status.VUs.Valid)
		assert.True(t, status.VUsMax.Valid)
		assert.False(t, status.Stopped)
		assert.False(t, status.Tainted)
	***REMOVED***)
***REMOVED***

func TestPatchStatus(t *testing.T) ***REMOVED***
	t.Parallel()

	testData := map[string]struct ***REMOVED***
		ExpectedStatusCode int
		ExpectedStatus     Status
		Payload            []byte
	***REMOVED******REMOVED***
		"nothing": ***REMOVED***
			ExpectedStatusCode: 200,
			ExpectedStatus:     Status***REMOVED******REMOVED***,
			Payload:            []byte(`***REMOVED***"data":***REMOVED***"type":"status","id":"default","attributes":***REMOVED***"status":0,"paused":null,"vus":null,"vus-max":null,"stopped":false,"running":false,"tainted":false***REMOVED******REMOVED******REMOVED***`),
		***REMOVED***,
		"paused": ***REMOVED***
			ExpectedStatusCode: 200,
			ExpectedStatus:     Status***REMOVED***Paused: null.BoolFrom(true)***REMOVED***,
			Payload:            []byte(`***REMOVED***"data":***REMOVED***"type":"status","id":"default","attributes":***REMOVED***"status":0,"paused":true,"vus":null,"vus-max":null,"stopped":false,"running":false,"tainted":false***REMOVED******REMOVED******REMOVED***`),
		***REMOVED***,
		"max vus": ***REMOVED***
			ExpectedStatusCode: 200,
			ExpectedStatus:     Status***REMOVED***VUsMax: null.IntFrom(20)***REMOVED***,
			Payload:            []byte(`***REMOVED***"data":***REMOVED***"type":"status","id":"default","attributes":***REMOVED***"status":0,"paused":null,"vus":null,"vus-max":20,"stopped":false,"running":false,"tainted":false***REMOVED******REMOVED******REMOVED***`),
		***REMOVED***,
		"max vus below initial": ***REMOVED***
			ExpectedStatusCode: 400,
			ExpectedStatus:     Status***REMOVED***VUsMax: null.IntFrom(5)***REMOVED***,
			Payload:            []byte(`***REMOVED***"data":***REMOVED***"type":"status","id":"default","attributes":***REMOVED***"status":0,"paused":null,"vus":null,"vus-max":5,"stopped":false,"running":false,"tainted":false***REMOVED******REMOVED******REMOVED***`),
		***REMOVED***,
		"too many vus": ***REMOVED***
			ExpectedStatusCode: 400,
			ExpectedStatus:     Status***REMOVED***VUs: null.IntFrom(10), VUsMax: null.IntFrom(0)***REMOVED***,
			Payload:            []byte(`***REMOVED***"data":***REMOVED***"type":"status","id":"default","attributes":***REMOVED***"status":0,"paused":null,"vus":10,"vus-max":0,"stopped":false,"running":false,"tainted":false***REMOVED******REMOVED******REMOVED***`),
		***REMOVED***,
		"vus": ***REMOVED***
			ExpectedStatusCode: 200,
			ExpectedStatus:     Status***REMOVED***VUs: null.IntFrom(10), VUsMax: null.IntFrom(10)***REMOVED***,
			Payload:            []byte(`***REMOVED***"data":***REMOVED***"type":"status","id":"default","attributes":***REMOVED***"status":0,"paused":null,"vus":10,"vus-max":10,"stopped":false,"running":false,"tainted":false***REMOVED******REMOVED******REMOVED***`),
		***REMOVED***,
	***REMOVED***
	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))

	scenarios := lib.ScenarioConfigs***REMOVED******REMOVED***
	err := json.Unmarshal([]byte(`
			***REMOVED***"external": ***REMOVED***"executor": "externally-controlled",
			"vus": 0, "maxVUs": 10, "duration": "1s"***REMOVED******REMOVED***`), &scenarios)
	require.NoError(t, err)
	options := lib.Options***REMOVED***Scenarios: scenarios***REMOVED***
	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)

	for name, testCase := range testData ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			execScheduler, err := local.NewExecutionScheduler(&minirunner.MiniRunner***REMOVED***Options: options***REMOVED***, logger)
			require.NoError(t, err)
			engine, err := core.NewEngine(execScheduler, options, lib.RuntimeOptions***REMOVED******REMOVED***, nil, logger, builtinMetrics)
			require.NoError(t, err)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			run, _, err := engine.Init(ctx, ctx)
			require.NoError(t, err)

			go func() ***REMOVED*** _ = run() ***REMOVED***()
			// wait for the executor to initialize to avoid a potential data race below
			time.Sleep(100 * time.Millisecond)

			rw := httptest.NewRecorder()
			NewHandler().ServeHTTP(rw, newRequestWithEngine(engine, "PATCH", "/v1/status", bytes.NewReader(testCase.Payload)))
			res := rw.Result()

			require.Equal(t, testCase.ExpectedStatusCode, res.StatusCode)

			if testCase.ExpectedStatusCode != 200 ***REMOVED***
				return
			***REMOVED***

			status := NewStatus(engine)
			if testCase.ExpectedStatus.Paused.Valid ***REMOVED***
				assert.Equal(t, testCase.ExpectedStatus.Paused, status.Paused)
			***REMOVED***
			if testCase.ExpectedStatus.VUs.Valid ***REMOVED***
				assert.Equal(t, testCase.ExpectedStatus.VUs, status.VUs)
			***REMOVED***
			if testCase.ExpectedStatus.VUsMax.Valid ***REMOVED***
				assert.Equal(t, testCase.ExpectedStatus.VUsMax, status.VUsMax)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
